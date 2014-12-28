package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAddRemoveProxy(t *testing.T) {
	domain := "www.example.com"
	url := "http://1.2.3.4:10000"
	addProxy(domain, url)
	proxyExists := getProxies(domain)
	assert.Equal(t, stringInSlice(url, proxyExists), true, "Domain should be www.example.com");
	removeProxy(domain, url)
	proxyDoesntExist := getProxies(domain)
	assert.Equal(t, stringInSlice(url, proxyDoesntExist), false, "Domain should not be www.example.com");
}
