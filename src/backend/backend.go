package backend

import (
	"errors"
)

type Backend interface {
	MultipleDomains() bool
	Start() error
	GetAddresses() (map[string][]string, error)
}

var (
	backends      map[string]Backend
	ErrNotFound   = errors.New("Could not find the backend.")
	ErrAlreadyReg = errors.New("Backend is already defined.")
)

func init() {
	backends = make(map[string]Backend)
}

func Register(name string, backend Backend) error {
	if _, exists := backends[name]; exists {
		return ErrAlreadyReg
	}
	backends[name] = backend

	return nil
}

func New(name string) (Backend, error) {
	if b, exists := backends[name]; exists {
		return b, nil
	}

	return nil, ErrNotFound
}

