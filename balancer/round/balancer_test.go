package balancer_round

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/nickschuch/marco/balancer"
)

func TestBalancerRound(t *testing.T) {
	var address string

	robin, _ := balancer.New("round")

	// Ensure we handle "no list" gracefully.
	address, _ = robin.Next()
	assert.Equal(t, "", address, "Received an empty return.")

	// Assign some addresses so we can start testing against a list.
	addresses := []string{
		"1.2.3.4:81",
		"1.2.3.4:80",
		"1.2.3.4:82",
	}
	robin.SetAddressList(addresses)

	// Ensure we can get the same addresslist back.
	balancerList, _ := robin.GetAddressList()
	assert.Equal(t, balancerList, addresses, "We can get an address list.")

	// Run a full circle of the round robin.
	address, _ = robin.Next()
	assert.Equal(t, "1.2.3.4:80", address, "We should see the address 1.2.3.4:80")
	address, _ = robin.Next()
	assert.Equal(t, "1.2.3.4:81", address, "We should see the address 1.2.3.4:81")
	address, _ = robin.Next()
	assert.Equal(t, "1.2.3.4:82", address, "We should see the address 1.2.3.4:82")

	// This tests that we can go a full circle and come back to the beginning of the round.
	address, _ = robin.Next()
	assert.Equal(t, "1.2.3.4:80", address, "We should see the address 1.2.3.4:80")

	// We now remove a record from the list of addresses.
	addresses = []string{
		"1.2.3.4:80",
		"1.2.3.4:82",
	}
	robin.SetAddressList(addresses)
	address, _ = robin.Next()
	assert.Equal(t, "1.2.3.4:82", address, "We should see the address 1.2.3.4:82")

	// Add a new address back into the mix and do a full run.
	addresses = []string{
		"1.2.3.4:80",
		"1.2.3.4:81",
		"1.2.3.4:82",
	}
	robin.SetAddressList(addresses)
	address, _ = robin.Next()
	assert.Equal(t, "1.2.3.4:82", address, "We should see the address 1.2.3.4:81")
	address, _ = robin.Next()
	assert.Equal(t, "1.2.3.4:80", address, "We should see the address 1.2.3.4:80")
	address, _ = robin.Next()
	assert.Equal(t, "1.2.3.4:81", address, "We should see the address 1.2.3.4:82")

	// Handle empty values gracefully.
	addresses = []string{
		"",
	}
	robin.SetAddressList(addresses)
	address, _ = robin.Next()
	assert.Equal(t, "", address, "Handle empty values gracefully.")
}
