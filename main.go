package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"runtime"
)

var bind string
var host string
var endpoint string
var ports string

func eventCallback(event *dockerclient.Event, args ...interface{}) {
	statesRemove := []string{
		"die",
		"stop",
		"destroy",
	}
	statesAdd := []string{
		"start",
	}

	if stringInSlice(event.Status, statesRemove) {
		container := getContainer(event.Id)
		removeProxy(container.Domain, container.Url)
		log.WithFields(log.Fields{
			"event": event.Status,
		}).Info("Removed container " + container.Domain + " (" + event.Id + ") out of rotation.")

		// Only remove the container when we destroy it.
		if event.Status == "destroy" {
			removeContainer(event.Id)
		}
	}
	if stringInSlice(event.Status, statesAdd) {
		container := getContainer(event.Id)
		addProxy(container.Domain, container.Url)
		log.WithFields(log.Fields{
			"event": event.Status,
		}).Info("Added container " + container.Domain + " (" + event.Id + ") into rotation.")
	}
}

func proxyCallback(w http.ResponseWriter, r *http.Request) {
	// Get a list of URL's that we can proxy this connection through to.
	//
	// Todo:
	//   * Make this pluggable so we can have different type of container discovery.
	proxyUrls := getProxies(r.Host)
	if len(proxyUrls) <= 0 {
		return
	}

	// Here is implement a basic random load balancer.
	//
	// Todo:
	//   * Make this pluggable (custom load balancer).
	proxyUrl := proxyUrls[rand.Intn(len(proxyUrls))]

	// Ensure we keep a log of the connection so we can go back and
	// debug if anything goes wrong.
	log.WithFields(log.Fields{
		"host":   r.Host,
		"uri":    r.URL,
		"method": r.Method,
	}).Info("Proxy to: " + proxyUrl)

	// Proxy the connection through.
	remote, err := url.Parse(proxyUrl)
	if err != nil {
		log.Fatal(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
}

func main() {
	// This allows us to serve more than a single request at a time.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// These are all the flags that can be passed to this application.
	flag.StringVar(&bind, "bind", "80", "Server traffic through the following port")
	flag.StringVar(&host, "host", "172.17.42.1", "The IP or DNS of the host exposing ports.")
	flag.StringVar(&endpoint, "endpoint", "unix:///var/run/docker.sock", "The Docker API endpoint eg. tcp://localhost:2375")
	flag.StringVar(&ports, "ports", "80,8080,2368,8983", "The ports you wish to proxy. Ordered in preference eg. 80,2368,8983")
	flag.Parse()

	// Get the Docker client that we can resuse for building a Proxy URL list
	// as well as monitor events.
	dockerClient := getDockerClient()
	dockerClient.StartMonitorEvents(eventCallback)

	// Build a cached copy of the containers and the urls that they
	// can be proxied to.
	populateCache()

	// Register and run.
	http.HandleFunc("/", proxyCallback)
	http.ListenAndServe(":"+bind, nil)
}
