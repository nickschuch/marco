package main

import (
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/mailgun/oxy/forward"
	"github.com/mailgun/oxy/roundrobin"
)

type Balancer struct {
	Robin *roundrobin.RoundRobin
}

func NewBalancer(addresses []string) (Balancer, error) {
	var b Balancer

	fwd, err := forward.New()
	if err != nil {
		return b, err
	}

	lb, err := roundrobin.New(fwd)
	if err != nil {
		return b, err
	}

	for _, a := range addresses {
		u, err := url.Parse(a)
		if err != nil {
			log.Printf("Cannot parse url: %s", a)
			continue
		}
		lb.UpsertServer(u)
	}

	b = Balancer{
		Robin: lb,
	}

	return b, nil
}

func (b *Balancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.Robin.ServeHTTP(w, r)
}
