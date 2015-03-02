package balancer_round

import (
	"errors"
	"sort"

	"github.com/nickschuch/marco/balancer"
)

var (
	ErrAdrNotFound = errors.New("Could not find an address.")
)

type BalancerRound struct {
	list      []string
	lastDelta int
}

func init() {
	balancer.Register("round", &BalancerRound{})
}

func (b *BalancerRound) SetAddressList(a []string) error {
	// We want to make sure the list is ordered at all times.
	// This will ensure that all instances get a fare shot at
	// being a part of the round robin.
	sort.Strings(a)

	b.list = a
	return nil
}

func (b *BalancerRound) GetAddressList() ([]string, error) {
	return b.list, nil
}

func (b *BalancerRound) Next() (string, error) {
	var address string

	// Check that we don't have an empty address list to
	// begin with.
	if len(b.list) <= 0 {
		return "", ErrAdrNotFound
	}

	// First check if we are going to exceed the length of
	// the addresses. If we are, then we go back to the beginning.
	if b.lastDelta >= len(b.list) {
		b.lastDelta = 0
	}

	// Check if we have a value to return.
	if b.list[b.lastDelta] == "" {
		return "", ErrAdrNotFound
	}

	// Increment for the next run.
	address = b.list[b.lastDelta]
	b.lastDelta++

	return address, nil
}
