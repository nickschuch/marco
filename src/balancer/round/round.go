package round

import (
	"errors"

	balancer ".."
)

var (
	ErrAdrNotFound = errors.New("Could not find an address.")
)

type BalancerRound struct {
	addresses []string
	lastDelta int
}

func init() {
	balancer.Register("round", &BalancerRound{})
}

func (b *BalancerRound) SetAddresses(a []string) error {
	b.addresses = a
	return nil
}

func (b *BalancerRound) GetAddresses() ([]string, error) {
	return b.addresses, nil
}

func (b *BalancerRound) AddAddress(a string) error {
	b.addresses = append(b.addresses, a)
	return nil
}

func (b *BalancerRound) RemoveAddress(a string) error {
	var newAddresses []string
	for _, address := range b.addresses {
        if address != a {
        	newAddresses = append(newAddresses, address)
        }
    }
    b.addresses = newAddresses
	return nil
}

func (b *BalancerRound) GetAddress() (string, error) {
	var address string

	// Check that we don't have an empty address list to
	// begin with.
	if len(b.addresses) <= 0 {
		return "", ErrAdrNotFound
	}

	// First check if we are going to exceed the length of
	// the addresses. If we are, then we go back to the beginning.
	if b.lastDelta >= len(b.addresses) {
		b.lastDelta = 0
	}

	// Check if we have a value to return.
	if b.addresses[b.lastDelta] == "" {
		return "", ErrAdrNotFound
	}

	// Increment for the next run.
	address = b.addresses[b.lastDelta]
	b.lastDelta++

	return address, nil
}
