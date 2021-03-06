package session

import (
	"crypto/sha1"
	"sort"

	"github.com/golang/protobuf/proto"

	"github.com/havoc-io/mutagen/sync"
)

func stableCopy(entry *sync.Entry) *StableEntry {
	// If the entry is nil, then the copy is nil.
	if entry == nil {
		return nil
	}

	// Create the result.
	result := &StableEntry{
		Kind:       entry.Kind,
		Executable: entry.Executable,
		Digest:     entry.Digest,
	}

	// Copy contents.
	for name, entry := range entry.Contents {
		result.Contents = append(result.Contents, &StableEntryContent{
			Name:  name,
			Entry: stableCopy(entry),
		})
	}

	// Sort contents by name.
	sort.Slice(result.Contents, func(i, j int) bool {
		return result.Contents[i].Name < result.Contents[j].Name
	})

	// Done.
	return result
}

func stableMarshal(entry *sync.Entry) ([]byte, error) {
	// Convert the entry to a stable copy.
	stableEntry := stableCopy(entry)

	// Wrap it in an archive in case it's nil.
	stableArchive := &StableArchive{Root: stableEntry}

	// Attempt to marshal.
	return proto.Marshal(stableArchive)
}

func stableUnmarshal(encoded []byte) (*sync.Entry, error) {
	// We can unmarshal directly into a normal archive since they are
	// byte-compatible with stable archives.
	archive := &Archive{}
	if err := proto.Unmarshal(encoded, archive); err != nil {
		return nil, err
	}

	// Success.
	return archive.Root, nil
}

func checksum(snapshotBytes []byte) []byte {
	result := sha1.Sum(snapshotBytes)
	return result[:]
}
