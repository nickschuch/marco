package main

import (
	"strings"
	"runtime"
	"net/http"
	"net/http/httputil"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v1"

	reconciler "./reconciler"
	handling "./handling"
	logging "./logging"
)

var (
	selPort     = kingpin.Flag("port", "The port to bind to.").Default("80").String()
	selRefresh  = kingpin.Flag("refresh", "How often to update the load balancer.").Default("10s").String()
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
	reconciled.AddBackend(*selBackend)
	reconciled.AddBalancer(*selBalancer)
	reconciled.SetRefresh(*selRefresh)
	reconciled.Start()

	// Start the webserver.
	logging.Info("Starting on port " + *selPort)
	http.HandleFunc("/", proxyCallback)
	log.Fatal(http.ListenAndServe(":" + *selPort, nil))
}

func proxyCallback(w http.ResponseWriter, r *http.Request) {
	// This returns an address as per the rules of the
	// Load balancers implementation.
	domain := strings.Split(r.Host, ":")
	address, error := reconciled.GetAddress(domain[0])
	handling.Check(error)

	// Proxy the connection through.
	remote, error := url.Parse(address)
	handling.Check(error)
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
}
