package main

import (
    "flag"
    "time"
    log "github.com/Sirupsen/logrus"
    "strings"
    "net/url"
    "net/http"
    "net/http/httputil"
    "runtime"
    "math/rand"
    "github.com/samalba/dockerclient"
    "github.com/pmylund/go-cache"
)

var bind string
var host string
var endpoint string
var ports string

var c = cache.New(5*time.Minute, 30*time.Second)

// Helper function to get a value of a Docker container ENV.
func getContainerEnv(key string, envs []string) string {
    for _, env := range envs {
        if strings.Contains(env, key) {
            envValue := strings.Split(env, "=")
            return envValue[1]
        }
    }
    return ""
}

// Helper function to convert "2365/tcp" into "2365".
func getPort(exposed string) string {
    port := strings.Split(exposed, "/")
    return port[0]
}

// Converts the Docker HostIP and HostPort details into a single
// string that we can use for proxy connections.
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

// Build up a list of proxy urls so we can load balance against them.
func getProxyUrls(domain string, r *http.Request) []string {
    var builtUrls []string

    // Return the cached response if we already have one.
    if x, found := c.Get(domain); found {
        builtUrls = x.([]string)
        return builtUrls
    }

    // Connect to the Docker daemon with the flag.
    docker, _ := dockerclient.NewDockerClient(endpoint, nil)
    containers, err := docker.ListContainers(false, false, "")
    if err != nil {
        log.Fatal(err)
    }
    for _, c := range containers {
        // We compare the name of the container against the subdomain of the request.
        // eg. project = all the containers with the name "project".
        container, _ := docker.InspectContainer(c.Id)

        // We need to check if a Domain has been set in the environment variables of the container.
        envDomain := getContainerEnv("DOMAIN", container.Config.Env)
        if envDomain == "" {
            log.WithFields(log.Fields{
              "host": r.Host,
              "uri": r.URL,
              "method": r.Method,
            }).Info("Could not find domain ENV value assigned to: " + container.Name)
            continue
        }

        // Don't include in the pool if this container's domain ENV value does not match.
        if envDomain != domain {
            continue
        }

        // Here we build the proxy URL based on the exposed values provided
        // by NetworkSettings. If a container has not been exposed, it will
        // not work.
        for portString, portObject := range container.NetworkSettings.Ports {
            port := getPort(portString)
            if strings.Contains(ports, port) {
                builtUrl := buildProxyUrl(portObject)
                if builtUrl != "" {
                    builtUrls = append(builtUrls, builtUrl)
                }
            }
        }
    }

    // Cache the value for later. This ensures that we don't have to
    // query the Docker daemon on every page request.
    c.Set(domain, builtUrls, cache.DefaultExpiration)

    return builtUrls
}

// This is the callback for the HTTP server.
func handler(w http.ResponseWriter, r *http.Request) {
    // Get a list of URL's that we can proxy this connection through to.
    //
    // Todo:
    //   * Make this pluggable so we can have different type of container discovery.
    proxyUrls := getProxyUrls(r.Host, r)
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
      "host": r.Host,
      "uri": r.URL,
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

// Run at the time of execution.
func main() {
    // This allows us to serve more than a single request at a time.
    runtime.GOMAXPROCS(runtime.NumCPU())

    // These are all the flags that can be passed to this application.
    flag.StringVar(&bind, "bind", "80", "Server traffic through the following port")
    flag.StringVar(&host, "host", "172.17.42.1", "The IP or DNS of the host exposing ports.")
    flag.StringVar(&endpoint, "endpoint", "unix:///var/run/docker.sock", "The Docker API endpoint eg. tcp://localhost:2375")
    flag.StringVar(&ports, "ports", "80,8080,2368,8983", "The ports you wish to proxy. Ordered in preference eg. 80,2368,8983")
    flag.Parse()

    // Register and run.
    http.HandleFunc("/", handler)
    http.ListenAndServe(":" + bind, nil)
}
