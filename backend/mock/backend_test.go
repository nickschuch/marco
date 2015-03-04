package backend_mock

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/nickschuch/marco/backend"
)

func TestInit(t *testing.T) {
	keys := backend.List()
	assert.Contains(t, keys, "mock", "The mock backend is registered.")
}

func TestAddresses(t *testing.T) {
	var list []string
	var expected []string

	mock, _ := backend.New("mock")

	// Get an empty list.
	list, _ = mock.Addresses("bar.com")
	assert.Equal(t, list, expected, "The backend has empty data for bar.com")

	// Get the addresses for all the other domains.
	list, _ = mock.Addresses("foobar.com")
	expected = []string{
		"1.2.3.4",
		"5.6.7.8",
	}
	assert.Equal(t, list, expected, "The mock backend has sample data for foobar.com")
}

func TestStart(t *testing.T) {
	// This is for both:
	//  * Code coverage.
	//  * Bootstrap for tests that related to starting.
	mock, _ := backend.New("mock")
	mock.Start()
}
