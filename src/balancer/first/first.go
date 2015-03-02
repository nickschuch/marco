package first

import (
	"errors"

	balancer ".."
)

var (
	ErrAdrNotFound = errors.New("Could not find an address.")
)

type BalancerFirst struct {
	list []string
}

func init() {
	balancer.Register("first", &BalancerFirst{})
}

func (b *BalancerFirst) SetAddressList(a []string) error {
	b.list = a
	return nil
}

func (b *BalancerFirst) GetAddressList() ([]string, error) {
	return b.list, nil
}

func (b *BalancerFirst) Next() (string, error) {
	// Check that we don't have an empty address list to
	// begin with.
	if len(b.list) <= 0 {
		return "", ErrAdrNotFound
	}

	return b.list[0], nil
}
