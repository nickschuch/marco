package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAddRemoveContainer(t *testing.T) {
	id := "1234"
	domain := "www.example.com"
	addContainer(id, domain, "http://1.2.3.4:10000")
	containerExists := getContainer(id)
	assert.Equal(t, containerExists.Domain, domain, "Domain should be www.example.com");
	removeContainer(id)
	containerDoesntExist := getCachedContainer(id)
	assert.Equal(t, containerDoesntExist.Domain, "", "Domain should not be www.example.com");
}
