package reconciler

import (
	"github.com/stretchr/testify/assert"
	"testing"

	_ "../backend/mock"
	_ "../balancer/round"
)

func TestSetBackendType(t *testing.T) {
	reconciled := Reconciler{}
	reconciled.SetBackendType("foo")
	actual := reconciled.GetBackendType()
	assert.Equal(t, "foo", actual, "Can set the backend type.")
}

func TestSetBalancerType(t *testing.T) {
	reconciled := Reconciler{}
	reconciled.SetBalancerType("foo")
	actual := reconciled.GetBalancerType()
	assert.Equal(t, "foo", actual, "Can set the balancer type.")
}

func TestAddress(t *testing.T) {
	// Setup a basic balancer and mock backend so we can start to assert
	// addresses against it.
	reconciled := Reconciler{}
	reconciled.SetBalancerType("round")
	reconciled.SetBackendType("mock")
	reconciled.Start()

	// Please see the "mock" backend type for the list of IP addresses.
	var actual string
	domain := "foo.com"
	actual, _ = reconciled.Address(domain)
	assert.Equal(t, "1.2.3.4", actual, "Can get an address.")
	actual, _ = reconciled.Address(domain)
	assert.Equal(t, "5.6.7.8", actual, "Can get the next address.")
	actual, _ = reconciled.Address(domain)
	assert.Equal(t, "1.2.3.4", actual, "Can get the first address.")

	// Ensure we handle having no addresses gracefully.
	actual, _ = reconciled.Address("bar.com")
	assert.Equal(t, "", actual, "Can handle no addresses gracefully.")
}
