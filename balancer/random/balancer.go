package balancer_random

import (
	"errors"

	"github.com/nickschuch/marco/balancer"
)

var (
	ErrAdrNotFound = errors.New("Could not find an address.")
)

type BalancerRandom struct {
	list []string
}

func init() {
	balancer.Register("random", &BalancerRandom{})
}

func (b *BalancerRandom) SetAddressList(a []string) error {
	b.list = a
	return nil
}

func (b *BalancerRandom) GetAddressList() ([]string, error) {
	return b.list, nil
}

func (b *BalancerRandom) Next() (string, error) {
	// Check that we don't have an empty address list to
	// begin with.
	if len(b.list) <= 0 {
		return "", ErrAdrNotFound
	}

	n := b.list[0]
	return n, nil
}
