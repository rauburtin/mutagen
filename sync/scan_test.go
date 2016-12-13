package sync

import (
	"bytes"
	"crypto/sha1"
	"hash"
	"runtime"
	"testing"
)

// gorootSnapshot is a snapshot of GOROOT. It is shared between multiple tests.
var gorootSnapshot *Entry

// gorootCache is a snapshot cache generated by a GOROOT snapshot. It is shared
// between multiple tests to avoid the cost of rebuilding such a snapshot cache
// multiple times.
var gorootCache *Cache

func init() {
	// Create the GOROOT snapshot and cache.
	if snapshot, cache, err := Scan(runtime.GOROOT(), sha1.New(), nil); err != nil {
		panic("couldn't create GOROOT snapshot: " + err.Error())
	} else if snapshot.Kind != EntryKind_Directory {
		panic("GOROOT snapshot is not a directory")
	} else {
		gorootSnapshot = snapshot
		gorootCache = cache
	}
}

// gorootRebuildHashProxy wraps an instance of and implements hash.Hash, but it
// signals a test error if any hashing occurs.  It is a test fixture for
// TestGorootRebuild.
type gorootRebuildHashProxy struct {
	hash.Hash
	t *testing.T
}

// Sum implements hash.Hash's Sum method, delegating to the underlying hash, but
// signals an error if invoked.
func (p *gorootRebuildHashProxy) Sum(b []byte) []byte {
	p.t.Error("rehashing occurred")
	return p.Hash.Sum(b)
}

// TestEfficientRebuild rebuilds the GOROOT snapshot with the existing cache,
// ensuring that no re-hashing occurs and that results are consistent.
func TestEfficientRebuild(t *testing.T) {
	hasher := &gorootRebuildHashProxy{sha1.New(), t}
	if snapshot, _, err := Scan(runtime.GOROOT(), hasher, gorootCache); err != nil {
		t.Fatal("couldn't rebuild GOROOT snapshot:", err)
	} else if !snapshot.Equal(gorootSnapshot) {
		t.Error("re-snapshotting produced a non-equivalent snapshot")
	}
}

// TestBuilderNonExistent verifies that Scan returns a nil root for paths that
// don't exist.
func TestBuilderNonExistent(t *testing.T) {
	// Create the snapshotter.
	snapshot, cache, err := Scan("THIS/DOES/NOT/EXIST", sha1.New(), nil)

	// Ensure that the snapshot root is nil.
	if snapshot != nil {
		t.Error("snapshot of non-existent path should be nil")
	}

	// Ensure that the cache is non-nil.
	if cache == nil {
		t.Error("snapshot cache of non-existent path should be non-nil")
	}

	// Ensure that the error is nil.
	if err != nil {
		t.Error("snapshot of non-existent path returned error:", err)
	}
}

// TODO: Add verification of change detection.

// TODO: Add verification of reference entries in GOROOT, including files.
