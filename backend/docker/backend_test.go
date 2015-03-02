package backend_docker

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/nickschuch/marco/backend"
)

func TestInit(t *testing.T) {
	keys := backend.List()
	assert.Contains(t, keys, "docker", "The Docker backend is registered.")
}
