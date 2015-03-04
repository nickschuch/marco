package balancer_first

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/nickschuch/marco/balancer"
)

func TestBalancerFirst(t *testing.T) {
	var address string

	first, _ := balancer.New("first")

	// Ensure we handle "no list" gracefully.
	address, _ = first.Next()
	assert.Equal(t, "", address, "Received an empty return.")

	// Assign some addresses so we can start testing against a list.
	addresses := []string{
		"1.2.3.4:81",
		"1.2.3.4:80",
	}
	first.SetAddressList(addresses)

	// Ensure we can get the same addresslist back.
	balancerList, _ := first.GetAddressList()
	assert.Equal(t, balancerList, addresses, "We can get an address list.")

	// Ensure all our attempts result in the first record.
	address, _ = first.Next()
	assert.Equal(t, "1.2.3.4:81", address, "We should see the address 1.2.3.4:81")
	address, _ = first.Next()
	assert.Equal(t, "1.2.3.4:81", address, "We should see the address 1.2.3.4:81")
}
