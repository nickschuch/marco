package round

import (
	"github.com/stretchr/testify/assert"
	"testing"

	balancer ".."
)

func TestBalancerRound(t *testing.T) {
	// These are the addresses we want to rotate around.
	var address string
	addresses := []string {
		"1.2.3.4:80",
		"1.2.3.4:81",
		"1.2.3.4:82",
	}

	robin, _ := balancer.New("round")
	robin.SetAddresses(addresses)

	// Run a full circle of the round robin.
	address, _ = robin.GetAddress()
	assert.Equal(t, "1.2.3.4:80", address, "We should see the address 1.2.3.4:80")
	address, _ = robin.GetAddress()
	assert.Equal(t, "1.2.3.4:81", address, "We should see the address 1.2.3.4:81")
	address, _ = robin.GetAddress()
	assert.Equal(t, "1.2.3.4:82", address, "We should see the address 1.2.3.4:82")

	// This tests that we can go a full circle and come back to the beginning of the round.
	address, _ = robin.GetAddress()
	assert.Equal(t, "1.2.3.4:80", address, "We should see the address 1.2.3.4:80")

	// We now remove a record from the list of addresses.
	robin.RemoveAddress("1.2.3.4:81")
	address, _ = robin.GetAddress()
	assert.Equal(t, "1.2.3.4:82", address, "We should see the address 1.2.3.4:82")

	// Add a new address back into the mix and do a full run.
	robin.AddAddress("1.2.3.4:81")
	address, _ = robin.GetAddress()
	assert.Equal(t, "1.2.3.4:81", address, "We should see the address 1.2.3.4:81")
	address, _ = robin.GetAddress()
	assert.Equal(t, "1.2.3.4:80", address, "We should see the address 1.2.3.4:80")
	address, _ = robin.GetAddress()
	assert.Equal(t, "1.2.3.4:82", address, "We should see the address 1.2.3.4:82")
}
