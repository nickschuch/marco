package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v1"
)

var (
	cliPort    = kingpin.Flag("port", "The port to bind to for balanced requests.").Default("80").String()
	cliReceive = kingpin.Flag("receive", "The port to bind to for backend data notifications.").Default("81").String()

	// This is a list of all the balancers that we are running.
	// We are running a balancer per domain eg. example.com and www.example.com.
	balancers map[string]Balancer

	// Messages.
	msgUnavailable = "Service Unavailable - This can indicate the instance is being deployed or has been terminated."
)

func main() {
	kingpin.Parse()

	balancers = make(map[string]Balancer)

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
		n.ServeHTTP(w, r)
		return
	}

	// At this point we assume that we cannot find the
	w.WriteHeader(http.StatusServiceUnavailable)
	fmt.Fprintf(w, msgUnavailable)
}

type Receive struct{}

func (r *Receive) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	backends := &[]Backend{}
	err := decoder.Decode(&backends)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, b := range *backends {
		n, err := NewBalancer(b.List)
		if err != nil {
			continue
		}
		balancers[b.Domain] = n

		log.WithFields(log.Fields{
			"type":   "received",
			"domain": b.Domain,
			"source": b.Type,
		}).Info(strings.Join(b.List, ","))
	}
}
