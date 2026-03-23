package memory

import (
	"context"
	"sync"
)

// DiffStore is an in-process implementation of tether.DiffStore.
// All data is held in a map protected by a mutex. Data lives until
// explicitly deleted or the process exits.
//
// Safe for concurrent use. Create with [NewDiffStore].
type DiffStore struct {
	mu   sync.Mutex
	data map[string][]byte
}

// NewDiffStore creates an empty in-memory diff store.
func NewDiffStore() *DiffStore {
	return &DiffStore{data: make(map[string][]byte)}
}

// Save stores differ snapshot data.
func (s *DiffStore) Save(_ context.Context, id string, data []byte) error {
	cp := make([]byte, len(data))
	copy(cp, data)
	s.mu.Lock()
	s.data[id] = cp
	s.mu.Unlock()
	return nil
}

// Load retrieves differ snapshot data. Returns (nil, nil) if the ID
// is not found.
func (s *DiffStore) Load(_ context.Context, id string) ([]byte, error) {
	s.mu.Lock()
	data := s.data[id]
	s.mu.Unlock()
	return data, nil
}

// Delete removes differ snapshot data. Returns nil if the ID is not
// found.
func (s *DiffStore) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	delete(s.data, id)
	s.mu.Unlock()
	return nil
}

// Len returns the number of stored snapshots. Useful in tests.
func (s *DiffStore) Len() int {
	s.mu.Lock()
	n := len(s.data)
	s.mu.Unlock()
	return n
}
