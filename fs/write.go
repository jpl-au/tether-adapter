package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

// safeName rejects session IDs that could escape the store directory
// when joined into a file path. Tether session IDs are ASCII
// alphanumerics (crypto/rand.Text), with hyphen and underscore also
// accepted for custom tooling; anything else - path separators, dots,
// empty strings - is not an ID this framework issued. The framework
// validates client-supplied IDs before they reach a store; this check
// is defence in depth for the store itself.
func safeName(id string) error {
	if id == "" {
		return fmt.Errorf("fs store: empty session ID")
	}
	for i := 0; i < len(id); i++ {
		c := id[i]
		switch {
		case c >= 'A' && c <= 'Z':
		case c >= 'a' && c <= 'z':
		case c >= '0' && c <= '9':
		case c == '-' || c == '_':
		default:
			return fmt.Errorf("fs store: session ID contains invalid character %q", c)
		}
	}
	return nil
}

// writeAtomic writes data to path via a temporary file in the same
// directory followed by a rename. Save runs during the disconnect and
// shutdown windows - exactly when a deploy is most likely to SIGKILL
// the process - and a torn write would leave a truncated file that
// silently fails to decode on the next reconnect, losing the very
// session the store exists to preserve.
func writeAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".tether-*.tmp")
	if err != nil {
		return err
	}
	name := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(name)
		return err
	}
	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		os.Remove(name)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(name)
		return err
	}
	if err := os.Rename(name, path); err != nil {
		os.Remove(name)
		return err
	}
	return nil
}
