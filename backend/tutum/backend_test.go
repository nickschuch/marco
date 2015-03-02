package backend_tutum

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/nickschuch/marco/backend"
)

func TestInit(t *testing.T) {
	keys := backend.List()
	assert.Contains(t, keys, "tutum", "The Tutum backend is registered.")
}
