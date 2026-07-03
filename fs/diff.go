package fs

import (
	"context"
	"os"
	"path/filepath"

	"github.com/jpl-au/tether/dev"
)

// Compile-time interface check.
var _ diffStore = (*DiffStore)(nil)

// diffStore mirrors [tether.DiffStore] for compile-time verification.
type diffStore interface {
	Save(ctx context.Context, id string, data []byte) error
	Load(ctx context.Context, id string) ([]byte, error)
	Delete(ctx context.Context, id string) error
}

// DiffStore persists differ snapshots to the filesystem. Each
// session's snapshots are stored as a single file named by its ID.
// This offloads snapshot data from Go memory during the reconnect
// window - a memory optimisation, not a recovery mechanism.
//
// Create with [NewDiffStore] and pass to
// [tether.StatefulConfig].DiffStore:
//
//	DiffStore: fs.NewDiffStore("tmp/diffs"),
type DiffStore struct {
	dir string
}

// NewDiffStore creates a DiffStore that writes snapshot files to dir.
// The directory is created if it does not exist.
func NewDiffStore(dir string) *DiffStore {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		dev.Warn("fs.DiffStore: mkdir failed", "dir", dir, "error", err)
	}
	return &DiffStore{dir: dir}
}

// Save writes differ snapshot data to a file. The write goes to a
// temp file and is renamed into place so a crash mid-write cannot
// leave a truncated file behind.
func (s *DiffStore) Save(_ context.Context, id string, data []byte) error {
	if err := safeName(id); err != nil {
		return err
	}
	dev.Debug("fs.DiffStore: save", "id", short(id), "bytes", len(data))
	return writeAtomic(s.path(id), data)
}

// Load reads differ snapshot data from a file. Returns (nil, nil)
// if the file does not exist.
func (s *DiffStore) Load(_ context.Context, id string) ([]byte, error) {
	if err := safeName(id); err != nil {
		return nil, err
	}
	data, err := os.ReadFile(s.path(id))
	if os.IsNotExist(err) {
		return nil, nil
	}
	return data, err
}

// Delete removes a differ snapshot file. Returns nil if the file
// does not exist.
func (s *DiffStore) Delete(_ context.Context, id string) error {
	if err := safeName(id); err != nil {
		return err
	}
	dev.Debug("fs.DiffStore: delete", "id", short(id))
	err := os.Remove(s.path(id))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (s *DiffStore) path(id string) string {
	return filepath.Join(s.dir, id+".diff")
}
