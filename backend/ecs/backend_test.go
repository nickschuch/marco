package backend_ecs

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/nickschuch/marco/backend"
)

func TestInit(t *testing.T) {
	keys := backend.List()
	assert.Contains(t, keys, "ecs", "The ECS backend is registered.")
}
