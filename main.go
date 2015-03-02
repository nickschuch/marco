package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"runtime"
	"strings"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v1"

	"github.com/nickschuch/marco/handling"
	"github.com/nickschuch/marco/logging"
	"github.com/nickschuch/marco/reconciler"
)

var (
	selPort     = kingpin.Flag("port", "The port to bind to.").Default("80").String()
	selBackend  = kingpin.Flag("backend", "The name of the backend driver.").Default("docker").String()
	selBalancer = kingpin.Flag("balancer", "The name of the balancer driver.").Default("round").String()

	// We set these as globals so proxyCallback() can access them.
	// @todo, Find a better way to handle this.
	reconciled reconciler.Reconciler
)

func main() {
	// This allows us to serve more than a single request at a time.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Instanciate the command.
	kingpin.Version("0.0.1")
	kingpin.CommandLine.Help = "Marco - Proxy for multiple backends."
	kingpin.Parse()

	// This sets up our "reconciled" object that handles backend connections
	// and load balancing.
	reconciled.SetBackendType(*selBackend)
	reconciled.SetBalancerType(*selBalancer)
	reconciled.Start()

	// Start the webserver.
	logging.Info("Starting on port " + *selPort)
	http.HandleFunc("/", proxyCallback)
	log.Fatal(http.ListenAndServe(":"+*selPort, nil))
}

func proxyCallback(w http.ResponseWriter, r *http.Request) {
	// This returns an address as per the rules of the
	// Load balancers implementation.
	// This is split by ":" so we can get just the domain
	// and not pass the port to the backend.
	domain := strings.Split(r.Host, ":")
	address, error := reconciled.Address(domain[0])
	handling.Check(error)

	// Proxy the connection through.
	remote, error := url.Parse(address)
	handling.Check(error)
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
}
