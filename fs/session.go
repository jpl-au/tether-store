package fs

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/jpl-au/tether/dev"
)

// Compile-time interface check.
var _ sessionStore = (*SessionStore)(nil)

// sessionStore mirrors [tether.SessionStore] so this package can
// verify interface satisfaction without importing the full tether
// package in the type signature. The actual satisfaction is checked
// by the caller assigning the store to StatefulConfig.SessionStore.
type sessionStore interface {
	Save(ctx context.Context, id string, data []byte, ttl time.Duration) error
	Load(ctx context.Context, id string) ([]byte, error)
	Delete(ctx context.Context, id string) error
}

// SessionStore persists session state to the filesystem. Each session
// is stored as a single file named by its ID. Suitable for
// development, single-node deployments, and low-traffic applications.
//
// The TTL parameter from Save is logged but not enforced - the
// framework calls Delete on reconnect and destroy, so orphaned files
// are the only case where TTL matters. For automatic expiry, use a
// store backed by Redis or another system with native TTL support.
//
// Create with [NewSessionStore] and pass to
// [tether.StatefulConfig].SessionStore:
//
//	SessionStore: fs.NewSessionStore("tmp/sessions"),
type SessionStore struct {
	dir string
}

// NewSessionStore creates a SessionStore that writes session files
// to dir. The directory is created if it does not exist.
func NewSessionStore(dir string) *SessionStore {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		dev.Warn("fs.SessionStore: mkdir failed", "dir", dir, "error", err)
	}
	return &SessionStore{dir: dir}
}

// Save writes session data to a file. The TTL is logged for
// observability but not enforced by the filesystem store.
func (s *SessionStore) Save(_ context.Context, id string, data []byte, ttl time.Duration) error {
	dev.Debug("fs.SessionStore: save", "id", short(id), "bytes", len(data), "ttl", ttl)
	return os.WriteFile(s.path(id), data, 0o600)
}

// Load reads session data from a file. Returns (nil, nil) if the
// file does not exist - the framework treats this as "no session to
// restore".
func (s *SessionStore) Load(_ context.Context, id string) ([]byte, error) {
	data, err := os.ReadFile(s.path(id))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err == nil {
		dev.Debug("fs.SessionStore: load", "id", short(id), "bytes", len(data))
	}
	return data, err
}

// Delete removes a session file. Returns nil if the file does not
// exist.
func (s *SessionStore) Delete(_ context.Context, id string) error {
	dev.Debug("fs.SessionStore: delete", "id", short(id))
	err := os.Remove(s.path(id))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (s *SessionStore) path(id string) string {
	return filepath.Join(s.dir, id+".session")
}
