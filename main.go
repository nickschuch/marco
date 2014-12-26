package main

import (
    "flag"
    "time"
    log "github.com/Sirupsen/logrus"
    "strings"
    "net/url"
    "net/http"
    "net/http/httputil"
    "github.com/samalba/dockerclient"
    "github.com/pmylund/go-cache"
)

var bind string
var host string
var endpoint string
var ports string

var c = cache.New(5*time.Minute, 30*time.Second)

func getPort(exposed string) string {
    port := strings.Split(exposed, "/")
    return port[0]
}

func buildProxyUrl(binding []dockerclient.PortBinding) string {
    // Ensure we have PortBinding values to build against.
    if len(binding) <= 0 {
        return ""
    }

    // Handle IP 0.0.0.0 the same way Swarm does. We replace this with an IP
    // that uses a local context.
    // @todo, Add the logic.
    ip := binding[0].HostIp
    port := binding[0].HostPort

    if ip == "0.0.0.0" {
        ip = host
    }

    return "http://" + ip + ":" + port
}

func getProxy(name string, r *http.Request) string {
    var builtUrl string

    if x, found := c.Get(name); found {
        builtUrl = x.(string)

        log.WithFields(log.Fields{
            "container": name,
            "path": r.URL,
            "cache": "HIT",
        }).Info("Connecting to: " + builtUrl)

        return builtUrl
    }

    // Connect to the Docker daemon with the flag.
    docker, _ := dockerclient.NewDockerClient(endpoint, nil)
    container, err := docker.InspectContainer(name)
    if err != nil {
        log.WithFields(log.Fields{
            "container": name,
            "path": r.URL,
        }).Fatal(err)
    }

    // Here we build the proxy URL based on the exposed values provided
    // by NetworkSettings. If a container has not been exposed, it will
    // not work.
    for portString, portObject := range container.NetworkSettings.Ports {
        port := getPort(portString)
        if strings.Contains(ports, port) {
            builtUrl = buildProxyUrl(portObject)
            break
        }
    }

    // Cache the value for later. This ensures that we don't have to
    // query the Docker daemon on every page request.
    c.Set(name, builtUrl, cache.DefaultExpiration)

    // Ensure we can debug this later on.
    log.WithFields(log.Fields{
        "container": name,
        "path": r.URL,
        "cache": "MISS",
    }).Info("Connecting to: " + builtUrl)

    return builtUrl
}

func handler(w http.ResponseWriter, r *http.Request) {
    // Get the name of the container based on the domain sent to this proxy.
    // eg. project1.example.com = project1
    //
    // Todo:
    //   * What happens when we don't have a subdomain provided.
    s := strings.Split(r.Host, ".")
    name := s[0]

    proxyUrl := getProxy(name, r)
    remote, err := url.Parse(proxyUrl)
    if err != nil {
        log.WithFields(log.Fields{
            "container": name,
            "path": r.URL,
        }).Fatal(err)
    } 
    proxy := httputil.NewSingleHostReverseProxy(remote)
    proxy.ServeHTTP(w, r)
}

func main() {
    flag.StringVar(&bind, "bind", "80", "Server traffic through the following port")
    flag.StringVar(&host, "host", "172.17.42.1", "The IP or DNS of the host exposing ports.")
    flag.StringVar(&endpoint, "endpoint", "unix:///var/run/docker.sock", "The Docker API endpoint eg. tcp://localhost:2375")
    flag.StringVar(&ports, "ports", "80,8080,2368,8983", "The ports you wish to proxy. Ordered in preference eg. 80,2368,8983")
    flag.Parse()

    http.HandleFunc("/", handler)
    http.ListenAndServe(":" + bind, nil)
}
