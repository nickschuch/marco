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

func TestStart(t *testing.T) {
	docker, _ := backend.New("docker")
	docker.Start()
}

func TestGetPort(t *testing.T) {
	port := getPort("3306/tcp")
	assert.Equal(t, "3306", port, "Can slice a port out of a full port string.")
}

func TestStringInSlice(t *testing.T) {
	var exists bool
	list := []string{
		"1.2.3.4",
		"5.6.7.8",
	}

	// If the string exists.
	exists = stringInSlice("1.2.3.4", list)
	assert.Equal(t, true, exists, "Can find string in slice.")

	// If the string does not exist.
	exists = stringInSlice("1.2.3.5", list)
	assert.Equal(t, false, exists, "Cannot find string in slice.")
}
