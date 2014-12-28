package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetContainerEnv(t *testing.T)  {
	key := "DOMAIN"
	envs := []string{
		"FOO=bar",
		"DOMAIN=www.example.com",
	}
	result := getContainerEnv(key, envs)
	assert.Equal(t, result, "www.example.com", "Domain should be www.example.com");
}

func TestGetPort(t *testing.T) {
	port := getPort("2365/tcp")
	assert.Equal(t, port, "2365", "Port should be 2365");
}

func TestBuildProxyUrl(t *testing.T) {
	specified := buildProxyUrl("192.168.1.20", "10000")
	assert.Equal(t, specified, "http://192.168.1.20:10000", "Proxy should be http://192.168.1.20:10000");
	open := buildProxyUrl("0.0.0.0", "10000")
	assert.NotEqual(t, open, "http://0.0.0.0:10000", "Proxy should not be http://0.0.0.0:10000");
}
