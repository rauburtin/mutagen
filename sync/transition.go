package sync

import (
	"bytes"
	"os"
	pathpkg "path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/havoc-io/mutagen/filesystem"
	"github.com/havoc-io/mutagen/timestamp"
)

func ensureRouteWithProperCase(root, path string, skipLast bool) error {
	// If the path is empty, then there's nothing to check.
	if path == "" {
		return nil
	}

	// Set the initial parent.
	parent := root

	// Decompose the path.
	components := strings.Split(path, "/")

	// If we were requested to not check the last component, then remove it from
	// the list.
	if skipLast && len(components) > 0 {
		components = components[:len(components)-1]
	}

	// While components remain, read the contents of the current parent and
	// ensure that a child with the correct cased entry exists.
	for len(components) > 0 {
		// Grab the contents for this location.
		contents, err := filesystem.DirectoryContents(parent)
		if err != nil {
			return errors.Wrap(err, "unable to read directory contents")
		}

		// Ensure that this path exists in contents. We can do a binary search
		// since contents will be sorted. If there's not a match, we're done.
		index := sort.SearchStrings(contents, components[0])
		if index == len(contents) || contents[index] != components[0] {
			return errors.New("unable to find matching entry")
		}

		// Update the parent.
		parent = filepath.Join(parent, components[0])

		// Reduce the component list.
		components = components[1:]
	}

	// Success.
	return nil
}

func ensureExpected(fullPath, path string, target *Entry, cache *Cache) error {
	// Grab cache information for this path. If we can't find it, we treat this
	// as an immediate fail. This is a bit of a heuristic/hack, because we could
	// recompute the digest of what's on disk, but for our use case this is very
	// expensive and we SHOULD already have this information cached from the
	// last scan.
	cacheEntry, ok := cache.Entries[path]
	if !ok {
		return errors.New("unable to find cache information for path")
	}

	// Grab stat information for this path.
	info, err := os.Lstat(fullPath)
	if err != nil {
		return errors.Wrap(err, "unable to grab file statistics")
	}

	// Convert the modification time to Protocol Buffers format.
	modificationTime, err := timestamp.Convert(info.ModTime())
	if err != nil {
		return errors.Wrap(err, "unable to convert modification timestamp")
	}

	// If stat information doesn't match, don't bother re-hashing, just abort.
	// Note that we don't really have to check executability here (and we
	// shouldn't since it's not preserved on all systems) - we just need to
	// check that it hasn't changed from the perspective of the filesystem, and
	// that is accomplished as part of the mode check. This is why we don't
	// restrict the mode comparison to the type bits.
	match := os.FileMode(cacheEntry.Mode) == info.Mode() &&
		timestamp.Equal(cacheEntry.ModificationTime, modificationTime) &&
		cacheEntry.Size == uint64(info.Size()) &&
		bytes.Equal(cacheEntry.Digest, target.Digest)
	if !match {
		return errors.New("modification detected")
	}

	// Success.
	return nil
}

func ensureNotExists(fullPath string) error {
	// Attempt to grab stat information for the path.
	_, err := os.Lstat(fullPath)

	// Handle error cases (which may indicate success).
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return errors.Wrap(err, "unable to determine path existence")
	}

	// Failure.
	return errors.New("path exists")
}

func removeFile(root, path string, target *Entry, cache *Cache) error {
	// Compute the full path to this file.
	fullPath := filepath.Join(root, path)

	// Ensure that the existing entry hasn't been modified from what we're
	// expecting.
	if err := ensureExpected(fullPath, path, target, cache); err != nil {
		return errors.Wrap(err, "unable to validate existing file")
	}

	// Remove the file.
	return os.Remove(fullPath)
}

func removeDirectory(root, path string, target *Entry, cache *Cache) []Problem {
	// Compute the full path to this directory.
	fullPath := filepath.Join(root, path)

	// List the contents for this directory.
	contentNames, err := filesystem.DirectoryContents(fullPath)
	if err != nil {
		return []Problem{newProblem(path, errors.Wrap(err, "unable to read directory contents"))}
	}

	// Loop through contents and remove them. We do this to ensure that what
	// we're removing has the proper case. If we were to just pass the OS what
	// exists in our content map and it were case insensitive, we could delete
	// a file that had been unmodified but renamed. There is, of course, a race
	// condition here between the time we grab the directory contents and the
	// time we remove, but it is very small and we also compare file contents,
	// so the chance of deleting something we shouldn't is very small.
	//
	// Note that we don't need to check that we've removed all entries listed in
	// the target. If they aren't in the directory contents, then they must have
	// already been deleted.
	var problems []Problem
	for _, name := range contentNames {
		// Compute the content path.
		contentPath := pathpkg.Join(path, name)

		// Grab the corresponding entry. If we don't know anything about this
		// entry, then mark that as a problem and ignore for now.
		entry, ok := target.Contents[name]
		if !ok {
			problems = append(problems, newProblem(
				contentPath,
				errors.New("unknown content encountered on disk"),
			))
			continue
		}

		// Handle its removal accordingly.
		var contentProblems []Problem
		if entry.Kind == EntryKind_Directory {
			contentProblems = removeDirectory(root, contentPath, entry, cache)
		} else if entry.Kind == EntryKind_File {
			if err = removeFile(root, contentPath, entry, cache); err != nil {
				contentProblems = append(contentProblems, newProblem(
					contentPath,
					errors.Wrap(err, "unable to remove file"),
				))
			}
		} else {
			contentProblems = append(contentProblems, newProblem(
				contentPath,
				errors.New("unknown entry type found in removal target"),
			))
		}

		// If there weren't any problems, than removal succeeded, so remove this
		// entry from the target. Otherwise add the problems to the complete
		// list.
		if len(contentProblems) == 0 {
			delete(target.Contents, name)
		} else {
			problems = append(problems, contentProblems...)
		}
	}

	// Attempt to remove the directory. If this succeeds, then clear any prior
	// problems, because clearly they no longer matter. This isn't a recursive
	// removal, so if something below failed to delete, this will still fail.
	if err := os.Remove(fullPath); err != nil {
		problems = append(problems, newProblem(
			path,
			errors.Wrap(err, "unable to remove directory"),
		))
	} else {
		problems = nil
	}

	// Done.
	return problems
}

func remove(root, path string, target *Entry, cache *Cache) (*Entry, []Problem) {
	// If the target is nil, we're done.
	if target == nil {
		return nil, nil
	}

	// Ensure that the path of the target exists (relative to the root) with the
	// specificed casing.
	if err := ensureRouteWithProperCase(root, path, false); err != nil {
		return target, []Problem{newProblem(
			path,
			errors.Wrap(err, "unable to verify path to target"),
		)}
	}

	// Create a copy of target for mutation.
	targetCopy := target.copy()

	// Check the target type and handle accordingly.
	var problems []Problem
	if target.Kind == EntryKind_Directory {
		problems = removeDirectory(root, path, targetCopy, cache)
	} else if target.Kind == EntryKind_File {
		if err := removeFile(root, path, targetCopy, cache); err != nil {
			problems = []Problem{newProblem(
				path,
				errors.Wrap(err, "unable to remove file"),
			)}
		}
	} else {
		problems = []Problem{newProblem(
			path,
			errors.New("removal requested for unknown entry type"),
		)}
	}

	// If there were any problems, then at least the root of the target will
	// have failed to remove, so return the reduced target.
	if len(problems) > 0 {
		return targetCopy, problems
	}

	// Success.
	return nil, nil
}

func swap(root, path string, oldEntry, newEntry *Entry, cache *Cache, provider StagingProvider) error {
	// Compute the full path to this file.
	fullPath := filepath.Join(root, path)

	// Ensure that the path of the target exists (relative to the root) with the
	// specificed casing.
	if err := ensureRouteWithProperCase(root, path, false); err != nil {
		return errors.Wrap(err, "unable to verify path to target")
	}

	// Ensure that the existing entry hasn't been modified from what we're
	// expecting.
	if err := ensureExpected(fullPath, path, oldEntry, cache); err != nil {
		return errors.Wrap(err, "unable to validate existing file")
	}

	// Compute the path to the staged file.
	stagedPath, err := provider(path, newEntry)
	if err != nil {
		return errors.Wrap(err, "unable to locate staged file")
	}

	// Rename the staged file.
	if err := filesystem.RenameFileAtomic(stagedPath, fullPath); err != nil {
		return errors.Wrap(err, "unable to relocate staged file")
	}

	// Success.
	return nil
}

func createFile(root, path string, target *Entry, provider StagingProvider) (*Entry, error) {
	// Compute the full path to the target.
	fullPath := filepath.Join(root, path)

	// Compute the path to the staged file.
	stagedPath, err := provider(path, target)
	if err != nil {
		return nil, errors.Wrap(err, "unable to locate staged file")
	}

	// Rename the staged file.
	if err := filesystem.RenameFileAtomic(stagedPath, fullPath); err != nil {
		return nil, errors.Wrap(err, "unable to relocate staged file")
	}

	// Success.
	return target, nil
}

func createDirectory(root, path string, target *Entry, provider StagingProvider) (*Entry, []Problem) {
	// Compute the full path to the target.
	fullPath := filepath.Join(root, path)

	// Attempt to create the directory.
	if err := os.Mkdir(fullPath, 0700); err != nil {
		return nil, []Problem{newProblem(
			path,
			errors.Wrap(err, "unable to create directory"),
		)}
	}

	// Create a shallow copy of the target that we'll populate as we create its
	// contents.
	created := target.copyShallow()

	// If there are contents in the target, allocate a map for created, because
	// we'll need to populate it.
	if len(target.Contents) > 0 {
		created.Contents = make(map[string]*Entry)
	}

	// Attempt to create the target contents. Track problems as we go.
	var problems []Problem
	for name, entry := range target.Contents {
		// Compute the content path.
		contentPath := pathpkg.Join(path, name)

		// Handle content creation based on type.
		var createdContent *Entry
		var contentProblems []Problem
		if entry.Kind == EntryKind_Directory {
			createdContent, contentProblems = createDirectory(root, contentPath, entry, provider)
		} else if entry.Kind == EntryKind_File {
			var err error
			createdContent, err = createFile(root, contentPath, entry, provider)
			if err != nil {
				contentProblems = append(contentProblems, newProblem(
					contentPath,
					errors.Wrap(err, "unable to create file"),
				))
			}
		} else {
			contentProblems = append(contentProblems, newProblem(
				contentPath,
				errors.New("creation requested for unknown entry type"),
			))
		}

		// If the created content is non-nil, then at least some portion of it
		// was created successfully, so record that.
		if createdContent != nil {
			created.Contents[name] = createdContent
		}

		// Record any problems that occurred when attempting to create the
		// content.
		problems = append(problems, contentProblems...)
	}

	// Return the portion of the target that was created and any problems that
	// occurred.
	return created, problems
}

func create(root, path string, target *Entry, provider StagingProvider) (*Entry, []Problem) {
	// If the target is nil, we're done.
	if target == nil {
		return nil, nil
	}

	// Ensure that the parent of the target path exists with the proper casing.
	if err := ensureRouteWithProperCase(root, path, true); err != nil {
		return nil, []Problem{newProblem(
			path,
			errors.Wrap(err, "unable to verify path to target"),
		)}
	}

	// Compute the full path to this file.
	fullPath := filepath.Join(root, path)

	// Ensure that the target path doesn't exist.
	if err := ensureNotExists(fullPath); err != nil {
		return nil, []Problem{newProblem(
			path,
			errors.Wrap(err, "unable to ensure path does not exist"),
		)}
	}

	// Check the target type and handle accordingly.
	if target.Kind == EntryKind_Directory {
		return createDirectory(root, path, target, provider)
	} else if target.Kind == EntryKind_File {
		if created, err := createFile(root, path, target, provider); err != nil {
			return created, []Problem{newProblem(
				path,
				errors.Wrap(err, "unable to create file"),
			)}
		} else {
			return created, nil
		}
	}
	return nil, []Problem{newProblem(
		path,
		errors.New("creation requested for unknown entry type"),
	)}
}

func Transition(root string, transitions []Change, cache *Cache, provider StagingProvider) ([]Change, []Problem) {
	// Set up results.
	var results []Change
	var problems []Problem

	// Iterate through transitions.
	for _, t := range transitions {
		// TODO: Should we check for transitions here that don't make any sense
		// but which aren't really logic errors? E.g. it doesn't make sense to
		// have a nil-to-nil transition. Likewise, it doesn't make sense to have
		// a directory-to-directory transition (although it might later on if
		// we're implementing permissions changes). If we see these types of
		// transitions, then there is an error in reconciliation or somewhere
		// else in the synchronization pipeline.

		// Handle the special case where both old and new are a file. In this
		// case we can do a simple swap. It makes sense to handle this specially
		// because it is a very common case and doing it with a swap will remove
		// any window where the path is empty on the filesystem.
		fileToFile := t.Old != nil && t.New != nil &&
			t.Old.Kind == EntryKind_File &&
			t.New.Kind == EntryKind_File
		if fileToFile {
			if err := swap(root, t.Path, t.Old, t.New, cache, provider); err != nil {
				results = append(results, Change{Path: t.Path, New: t.Old})
				problems = append(problems, newProblem(
					t.Path,
					errors.Wrap(err, "unable to swap file"),
				))
			} else {
				results = append(results, Change{Path: t.Path, New: t.New})
			}
			continue
		}

		// Reduce whatever we expect to see on disk to nil (remove it). If we
		// don't expect to see anything (t.Old == nil), this is a no-op. If this
		// fails, record the reduced entry as well as any problems preventing
		// full removal and continue to the next transition.
		if r, p := remove(root, t.Path, t.Old, cache); r != nil {
			results = append(results, Change{Path: t.Path, New: r})
			problems = append(problems, p...)
			continue
		}

		// At this point, we should have nil on disk. Transition to whatever the
		// new entry is. If the new-entry is nil, this is a no-op. Record
		// whatever portion of the target we create as well as any problems
		// preventing full creation.
		c, p := create(root, t.Path, t.New, provider)
		results = append(results, Change{Path: t.Path, New: c})
		problems = append(problems, p...)
	}

	// Done.
	return results, problems
}
