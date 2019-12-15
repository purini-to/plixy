package memstore

import (
	"sync"

	"github.com/purini-to/plixy/pkg/api"
)

// Store represents a in memory store.
type Store struct {
	sync.RWMutex
	def *api.Definition
}

// New creates a in memory store.
func New() *Store {
	return &Store{}
}

// GetDefinition get a api definition.
func (s *Store) GetDefinition() (*api.Definition, error) {
	s.RLock()
	defer s.RUnlock()
	return s.def, nil
}

// SetDefinition set a api definition.
// definition not set if that is invalid.
func (s *Store) SetDefinition(def *api.Definition) error {
	if _, err := def.Validate(); err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()
	s.def = def
	return nil
}

// Close is nothing func.
func (s *Store) Close() error {
	return nil
}
