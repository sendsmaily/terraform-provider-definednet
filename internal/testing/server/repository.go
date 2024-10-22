package server

import (
	"fmt"
	"sync"
)

// NewRepository creates a fake API data repository.
func NewRepository[O Object]() *Repository[O] {
	return &Repository[O]{
		data: make(map[string]O),
	}
}

// Repository is a fake API data repository.
type Repository[O Object] struct {
	mu   sync.Mutex
	data map[string]O
}

// Object is the object stored in the repository.
type Object interface {
	Key() string
}

// Add an object to repository.
func (r *Repository[O]) Add(m O) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Primary identifiers must be unique.
	if _, exists := r.data[m.Key()]; exists {
		return fmt.Errorf("object with id %q already exists", m.Key())
	}

	// TODO: implement additional constraints enforcement.
	// For example, Host.Name values must be unique.

	r.data[m.Key()] = m

	return nil
}

// Get an object from the repository.
func (r *Repository[O]) Get(id string) (*O, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	obj, exists := r.data[id]
	if !exists {
		return nil, fmt.Errorf("object with id %q does not exist", id)
	}

	return &obj, nil
}

// Remove an object from the repository.
func (r *Repository[O]) Remove(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[id]; !exists {
		return fmt.Errorf("object with id %q does not exist", id)
	}

	delete(r.data, id)

	return nil
}

// Replace an object in the repository.
func (r *Repository[O]) Replace(m O) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[m.Key()]; !exists {
		return fmt.Errorf("object with id %q does not exist", m.Key())
	}

	r.data[m.Key()] = m

	return nil
}
