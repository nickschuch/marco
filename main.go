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
var endpoint string
var ports string

var c = cache.New(5*time.Minute, 30*time.Second)

func getProxy(name string, r *http.Request) string {
    if x, found := c.Get(name); found {
        proxy_url := x.(string)

        log.WithFields(log.Fields{
            "container": name,
            "path": r.URL,
            "cache": "HIT",
        }).Info("Connecting to: " + proxy_url)

        return proxy_url
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

    // Query Docker for the IP of the container.
    //
    // Todo:
    //   * What happens when we don't have an IP.
    ip := container.NetworkSettings.IpAddress
    if err != nil {
        log.WithFields(log.Fields{
            "container": name,
            "path": r.URL,
        }).Fatal(err)
    }

    // The first port available via the "ports" argument.
    // eg. When using the default ports flag value.
    //     If the port exposes 8080 and 2368 than it will
    //     be 8080.
    //
    // Todo:
    //   * What happens when we don't have a port?
    port := ""
    for exposed := range container.NetworkSettings.Ports {
        p := strings.Split(exposed, "/")
        if strings.Contains(ports, p[0]) {
            port = p[0]
            break
        }
    }

    // Proxy through to the container.
    proxy_url := "http://" + ip + ":" + port
    c.Set(name, proxy_url, cache.DefaultExpiration)

    log.WithFields(log.Fields{
        "container": name,
        "path": r.URL,
        "cache": "MISS",
    }).Info("Connecting to: " + proxy_url)

    return proxy_url
}

func handler(w http.ResponseWriter, r *http.Request) {
    // Get the name of the container based on the domain sent to this proxy.
    // eg. project1.example.com = project1
    //
    // Todo:
    //   * What happens when we don't have a subdomain provided.
    s := strings.Split(r.Host, ".")
    name := s[0]

    proxy_url := getProxy(name, r)
    remote, err := url.Parse(proxy_url)
    if err != nil {
        panic(err)
    } 
    proxy := httputil.NewSingleHostReverseProxy(remote)
    proxy.ServeHTTP(w, r)
}

func main() {
    flag.StringVar(&bind, "bind", "80", "Server traffic through the following port")
    flag.StringVar(&endpoint, "endpoint", "unix:///var/run/docker.sock", "The Docker API endpoint eg. tcp://localhost:2375")
    flag.StringVar(&ports, "ports", "80,8080,2368,8983", "The ports you wish to proxy. Ordered in preference eg. 80,2368,8983")
    flag.Parse()

    http.HandleFunc("/", handler)
    http.ListenAndServe(":" + bind, nil)
}
