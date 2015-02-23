package balancer

import (
	"errors"
)

type Balancer interface {
	SetAddresses([]string) error
	GetAddresses() ([]string, error)
	AddAddress(string) error
	RemoveAddress(string) error
	GetAddress() (string, error)
}

var (
	Balancers     map[string]Balancer
	ErrNotFound   = errors.New("Could not find the balancer.")
	ErrAlreadyReg = errors.New("Balancer is already defined.")
)

func init() {
	Balancers = make(map[string]Balancer)
}

func Register(name string, balancer Balancer) error {
	if _, exists := Balancers[name]; exists {
		return ErrAlreadyReg
	}
	Balancers[name] = balancer

	return nil
}

func New(name string) (Balancer, error) {
	if p, exists := Balancers[name]; exists {
		return p, nil
	}

	return nil, ErrNotFound
}

