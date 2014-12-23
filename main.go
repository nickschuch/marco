package main

import (
    "flag"
    "log"
    "strings"
    "net/url"
    "net/http"
    "net/http/httputil"
    "github.com/samalba/dockerclient"
)

var endpoint string
var ports string

func handler(w http.ResponseWriter, r *http.Request) {
    // Get the name of the container based on the domain sent to this proxy.
    // eg. project1.example.com = project1
    s := strings.Split(r.Host, ".")
    name := s[0]

    docker, _ := dockerclient.NewDockerClient(endpoint, nil)
    container, err := docker.InspectContainer(name)
    if err != nil {
        log.Fatal(err)
    }

    ip := container.NetworkSettings.IpAddress
    if err != nil {
        log.Fatal(err)
    }

    port := ""
    for exposed := range container.NetworkSettings.Ports {
        p := strings.Split(exposed, "/")
        if strings.Contains(ports, p[0]) {
            port = p[0]
            break
        }
    }

    remote, err := url.Parse("http://" + ip + ":" + port)
    if err != nil {
       panic(err)
    } 

    proxy := httputil.NewSingleHostReverseProxy(remote)
    proxy.ServeHTTP(w, r)
}

func main() {
    flag.StringVar(&endpoint, "endpoint", "unix:///var/run/docker.sock", "The Docker API endpoint eg. tcp://localhost:2375")
    flag.StringVar(&ports, "ports", "80,8080,2368,8983", "The ports you wish to proxy. Ordered in preference eg. 80,2368,8983")
    flag.Parse()

    http.HandleFunc("/", handler)
    // @todo, Make this an option.
    http.ListenAndServe(":80", nil)
}
