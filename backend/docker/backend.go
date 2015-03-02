package backend_docker

import (
	"strings"
	"time"

	"github.com/daryl/cash"
	"github.com/samalba/dockerclient"
	"gopkg.in/alecthomas/kingpin.v1"

	"github.com/nickschuch/marco/backend"
	"github.com/nickschuch/marco/handling"
)

var (
	selDockerEndpoint  = kingpin.Flag("docker-endpoint", "The Docker endpoint.").Default("unix:///var/run/docker.sock").OverrideDefaultFromEnvar("DOCKER_HOST").String()
	selDockerPorts     = kingpin.Flag("docker-ports", "The ports you wish to proxy.").Default("80,8080,2368,8983").String()
	selDockerDomainEnv = kingpin.Flag("docker-domain-env", "The container environment variable that is used as a domain identifier.").Default("DOMAIN").String()
)

type BackendDocker struct {
	cache *cash.Cash
}

func init() {
	backend.Register("docker", &BackendDocker{})
}

func (b *BackendDocker) Start() error {
	// Create a brand new cache item that we can use to
	// store all the domains and there associated list of
	// addresses.
	b.cache = cash.New(cash.Conf{
		// Default expiration.
		time.Minute,
		// Clean interval.
		30 * time.Minute,
	})

	return nil
}

func (b *BackendDocker) Addresses(domain string) ([]string, error) {
	var list []string

	if v, ok := b.cache.Get(domain); ok {
		// We found a cached item and we should return it's list of urls.
		list = v.([]string)
		return list, nil
	}

	// This call was already at a cost and given we couldn't filter on the domain for results.
	// We might as well set all the associated URLs as well as the one is missing.
	list, err := getListByDomain(domain)
	handling.Check(err)
	b.cache.Set(domain, list, time.Minute)

	return list, nil
}

func getListByDomain(domain string) ([]string, error) {
	var domainList []string

	list, err := getList()
	handling.Check(err)
	if len(list[domain]) > 0 {
		domainList = list[domain]
	}

	return domainList, nil
}

func getList() (map[string][]string, error) {
	// These are the URL's (keyed by domain) that we will return.
	list := make(map[string][]string)

	dockerClient, error := dockerclient.NewDockerClient(*selDockerEndpoint, nil)
	handling.Check(error)

	containers, error := dockerClient.ListContainers(false, false, "")
	handling.Check(error)

	for _, c := range containers {
		container, _ := dockerClient.InspectContainer(c.Id)

		// We try to find the domain environment variable. If we don't have one
		// then we have nothing left to do with this container.
		envDomain := getContainerEnv(*selDockerDomainEnv, container.Config.Env)
		if len(envDomain) <= 0 {
			continue
		}

		// Here we build the proxy URL based on the exposed values provided
		// by NetworkSettings. If a container has not been exposed, it will
		// not work. We then build a URL based on these exposed values and:
		//   * Add a container reference so we can perform safe operations
		//     in the future.
		//   * Add the built url to the proxy lists for load balancing.
		for portString, portObject := range container.NetworkSettings.Ports {
			port := getPort(portString)
			if strings.Contains(*selDockerPorts, port) {
				list[envDomain] = append(list[envDomain], getProxyUrl(portObject))
			}
		}
	}

	return list, nil
}

func getContainerEnv(key string, envs []string) string {
	for _, env := range envs {
		if strings.Contains(env, key) {
			envValue := strings.Split(env, "=")
			return envValue[1]
		}
	}
	return ""
}

func getPort(exposed string) string {
	port := strings.Split(exposed, "/")
	return port[0]
}

func getProxyUrl(binding []dockerclient.PortBinding) string {
	// Ensure we have PortBinding values to build against.
	if len(binding) <= 0 {
		return ""
	}

	// Handle IP 0.0.0.0 the same way Swarm does. We replace this with an IP
	// that uses a local context.
	port := binding[0].HostPort
	ip := binding[0].HostIp
	if ip == "0.0.0.0" {
		ip = "127.0.0.1"
	}

	return "http://" + ip + ":" + port
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
