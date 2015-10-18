package main

import (
	"net/http"
	"net/url"
	"os"

	"github.com/mailgun/oxy/forward"
	"github.com/mailgun/oxy/trace"
	"github.com/nickschuch/oxy/roundrobin"
)

type Balancer struct {
	Handler *roundrobin.RoundRobin
}

func NewBalancer(tag string, weight int, addresses []*url.URL) (*Balancer, error) {
	var b *Balancer

	// Forwards requests to remote location and rewrites headers.
	fwd, err := forward.New()
	if err != nil {
		return b, err
	}

	// Structured request and response logger.
	t, err := trace.New(fwd, os.Stdout)
	if err != nil {
		return b, err
	}

	// Round robin load balancer for distributing traffic.
	lb, err := roundrobin.New(t)
	if err != nil {
		return b, err
	}

	// Populate the balancer with a list of backends.
	for _, a := range addresses {
		lb.UpsertServer(a, roundrobin.Weight(weight), roundrobin.Tag(tag))
	}

	b = &Balancer{
		Handler: lb,
	}
	return b, nil
}

func (b *Balancer) Update(t string, w int, l []*url.URL) {
	// Get a list of servers from the balancer.
	existing := b.Handler.FindServersByTag(t)

	// Determine if new items are in the list, if not remove them.
	for _, e := range existing {
		if !Contains(l, e) {
			b.Handler.RemoveServer(e)
		}
	}

	// Add the remaining items to the servers list.
	for _, n := range l {
		if !Contains(existing, n) {
			b.Handler.UpsertServer(n, roundrobin.Weight(w), roundrobin.Tag(t))
		}
	}
}

func (b *Balancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.Handler.ServeHTTP(w, r)
}
