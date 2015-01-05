package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringInSlice(t *testing.T) {
	list := []string{
		"foo",
		"bar",
	}
	found := stringInSlice("foo", list)
	assert.Equal(t, found, true, "List should contain foo")
	notFound := stringInSlice("baz", list)
	assert.Equal(t, notFound, false, "List should not contain baz")
}
