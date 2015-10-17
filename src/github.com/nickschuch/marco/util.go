package main

import (
	"net/url"
)

func contains(s []*url.URL, e *url.URL) bool {
	for _, a := range s {
		if a.String() == e.String() {
			return true
		}
	}
	return false
}
