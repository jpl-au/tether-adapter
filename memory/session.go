package memory

import (
	"context"
	"sync"
	"time"
)

// SessionStore is an in-process implementation of tether.SessionStore.
// All data is held in a map protected by a mutex. The TTL parameter
// from Save is ignored - data lives until explicitly deleted or the
// process exits. For automatic expiry, use a store backed by Redis
// or another system with native TTL support.
//
// Safe for concurrent use. Create with [NewSessionStore].
type SessionStore struct {
	mu   sync.Mutex
	data map[string][]byte
}

// NewSessionStore creates an empty in-memory session store.
func NewSessionStore() *SessionStore {
	return &SessionStore{data: make(map[string][]byte)}
}

// Save stores session data. The TTL hint is accepted but not
// enforced - the framework calls Delete on reconnect and destroy.
func (s *SessionStore) Save(_ context.Context, id string, data []byte, _ time.Duration) error {
	cp := make([]byte, len(data))
	copy(cp, data)
	s.mu.Lock()
	s.data[id] = cp
	s.mu.Unlock()
	return nil
}

// Load retrieves session data. Returns (nil, nil) if the ID is not
// found - the framework treats this as "no session to restore". The
// returned slice is a copy, matching Save's semantics, so the caller
// cannot alias the store's internal buffer.
func (s *SessionStore) Load(_ context.Context, id string) ([]byte, error) {
	s.mu.Lock()
	data := s.data[id]
	s.mu.Unlock()
	if data == nil {
		return nil, nil
	}
	cp := make([]byte, len(data))
	copy(cp, data)
	return cp, nil
}

// Delete removes session data. Returns nil if the ID is not found.
func (s *SessionStore) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	delete(s.data, id)
	s.mu.Unlock()
	return nil
}

// Len returns the number of stored sessions. Useful in tests.
func (s *SessionStore) Len() int {
	s.mu.Lock()
	n := len(s.data)
	s.mu.Unlock()
	return n
}
