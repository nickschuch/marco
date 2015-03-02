package first

import (
	"github.com/stretchr/testify/assert"
	"testing"

	balancer ".."
)

func TestBalancerFirst(t *testing.T) {
	// These are the addresses we want to rotate around.
	var address string
	addresses := []string{
		"1.2.3.4:81",
		"1.2.3.4:80",
	}

	first, _ := balancer.New("first")
	first.SetAddressList(addresses)

	// Ensure all our attempts result in the first record.
	address, _ = first.Next()
	assert.Equal(t, "1.2.3.4:81", address, "We should see the address 1.2.3.4:81")
	address, _ = first.Next()
	assert.Equal(t, "1.2.3.4:81", address, "We should see the address 1.2.3.4:81")
}
