package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBalancer(t *testing.T) {
	l := []*url.URL{
		&url.URL{
			Host:   "127.0.0.1:8080",
			Scheme: "http",
		},
	}
	e := []*url.URL{
		&url.URL{
			Host:   "127.0.0.1:8080",
			Scheme: "http",
		},
	}

	// Build a new balancer with the first backend.
	b, err := NewBalancer("backend1", 1, l)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, e, b.Handler.Servers(), "Added records from backend1")

	// Update addresses from "backend1".
	l = []*url.URL{
		&url.URL{
			Host:   "127.0.0.1:8081",
			Scheme: "http",
		},
	}
	b.Update("backend1", 1, l)
	e = []*url.URL{
		&url.URL{
			Host:   "127.0.0.1:8081",
			Scheme: "http",
		},
	}
	assert.Equal(t, e, b.Handler.Servers(), "Update records from backend1")

	// Add new records from "backend2".
	l = []*url.URL{
		&url.URL{
			Host:   "127.0.0.1:8082",
			Scheme: "http",
		},
	}
	b.Update("backend2", 2, l)
	e = []*url.URL{
		&url.URL{
			Host:   "127.0.0.1:8081",
			Scheme: "http",
		},
		&url.URL{
			Host:   "127.0.0.1:8082",
			Scheme: "http",
		},
	}
	assert.Equal(t, e, b.Handler.Servers(), "Added records from backend2")

	// Update addresses from "backend2"
	l = []*url.URL{
		&url.URL{
			Host:   "127.0.0.1:8082",
			Scheme: "http",
		},
		&url.URL{
			Host:   "127.0.0.1:8083",
			Scheme: "http",
		},
	}
	b.Update("backend2", 2, l)
	e = []*url.URL{
		&url.URL{
			Host:   "127.0.0.1:8081",
			Scheme: "http",
		},
		&url.URL{
			Host:   "127.0.0.1:8082",
			Scheme: "http",
		},
		&url.URL{
			Host:   "127.0.0.1:8083",
			Scheme: "http",
		},
	}
	assert.Equal(t, e, b.Handler.Servers(), "Added records from backend2")
}
