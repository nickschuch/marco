package backend_mock

import (
	"github.com/stretchr/testify/assert"
	"testing"

	backend ".."
)

func TestInit(t *testing.T) {
	keys := backend.List()
	assert.Contains(t, keys, "mock", "The mock backend is registered.")
}

func TestAddresses(t *testing.T) {
	mock, _ := backend.New("mock")
	list, _ := mock.Addresses("foobar.com")
	expected := []string{
		"1.2.3.4",
		"5.6.7.8",
	}
	assert.Equal(t, list, expected, "The mock backend has sample data.")
}
