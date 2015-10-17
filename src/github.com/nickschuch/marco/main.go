package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/mailgun/oxy/stream"
	"github.com/nickschuch/marco-lib"
	"gopkg.in/alecthomas/kingpin.v1"
)

var (
	cliPort    = kingpin.Flag("port", "The port to bind to for balanced requests.").Default("80").String()
	cliReceive = kingpin.Flag("receive", "The port to bind to for backend data notifications.").Default("81").String()

	// This is a list of all the balancers that we are running.
	// We are running a balancer per domain eg. example.com and www.example.com.
	balancers map[string]*Balancer
)

func main() {
	kingpin.Parse()

	balancers = make(map[string]*Balancer)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		log.Info("Receiving backend data on port " + *cliReceive)
		log.Fatal(http.ListenAndServe(":"+*cliReceive, &Receive{}))
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		log.Info("Balancing connections on port " + *cliPort)
		s := &http.Server{
			Addr:    ":" + *cliPort,
			Handler: http.HandlerFunc(proxy),
		}
		s.ListenAndServe()
		wg.Done()
	}()

	wg.Wait()
}

func proxy(w http.ResponseWriter, r *http.Request) {
	// This is split by ":" so we can get just the domain
	// and not pass the port to the backend.
	domain := strings.Split(r.Host, ":")
	if n, ok := balancers[domain[0]]; ok {
		s, err := stream.New(n, stream.Retry(`IsNetworkError() && Attempts() < 2`))
		if err == nil {
			s.ServeHTTP(w, r)
			return
		}
	}

	w.WriteHeader(http.StatusServiceUnavailable)
	fmt.Fprintf(w, "Service Unavailable - This can indicate the instance is being deployed or has been terminated.")
}

type Receive struct{}

func (r *Receive) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	backends := &[]marco.Backend{}
	err := decoder.Decode(&backends)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Info(err)
		return
	}

	// Was any backend data actually sent to this service.
	if backends == nil {
		return
	}

	// Determine if we need a new balancer, or update an existing one.
	for _, b := range *backends {
		log.WithFields(log.Fields{
			"type":   "received",
			"domain": b.Domain,
			"source": b.Type,
		}).Info(b.List)
		if _, ok := balancers[b.Domain]; ok {
			// Update an existing load balancer.
			balancers[b.Domain].Update(b.Type, b.Weight, b.List)
		} else {
			// Else we just spin up a fresh balancer service.
			n, err := NewBalancer(b.Type, b.Weight, b.List)
			if err != nil {
				continue
			}
			balancers[b.Domain] = n
		}
	}
}
