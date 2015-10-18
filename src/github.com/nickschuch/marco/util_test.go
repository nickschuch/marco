package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	l := []*url.URL{
		&url.URL{
			Host:   "127.0.0.1:8082",
			Scheme: "http",
		},
		&url.URL{
			Host:   "127.0.0.1:8083",
			Scheme: "http",
		},
	}
	c := &url.URL{
		Host:   "127.0.0.1:8082",
		Scheme: "http",
	}
	assert.True(t, Contains(l, c), "URL list contains URL.")
}
