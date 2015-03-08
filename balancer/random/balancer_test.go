package balancer_random

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/nickschuch/marco/balancer"
)

func TestBalancerRandom(t *testing.T) {
	var address string

	random, _ := balancer.New("random")

	// Ensure we handle "no list" gracefully.
	address, _ = random.Next()
	assert.Equal(t, "", address, "Received an empty return.")

	// Assign some addresses so we can start testing against a list.
	addresses := []string{
		"1.2.3.4:81",
		"1.2.3.4:80",
	}
	random.SetAddressList(addresses)

	// Ensure we can get the same addresslist back.
	balancerList, _ := random.GetAddressList()
	assert.Equal(t, balancerList, addresses, "We can get an address list.")

	// Ensure all our attempts result in the random record.
	address, _ = random.Next()
	assert.Contains(t, addresses, address, "The address exists on the list.")
	address, _ = random.Next()
	assert.Contains(t, addresses, address, "The address exists on the list.")
}
