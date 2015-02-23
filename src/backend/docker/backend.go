package docker

import (
	"strings"

	"gopkg.in/alecthomas/kingpin.v1"
	"github.com/samalba/dockerclient"

	backend ".."
	handling "../../handling"
)

var (
	selDockerEndpoint  = kingpin.Flag("docker-endpoint", "The Docker endpoint.").Default("unix:///var/run/docker.sock").OverrideDefaultFromEnvar("DOCKER_HOST").String()
	selDockerPorts     = kingpin.Flag("docker-ports", "The ports you wish to proxy.").Default("80,8080,2368,8983").String()
	selDockerDomainEnv = kingpin.Flag("docker-domain-env", "The container environment variable that is used as a domain identifier.").Default("DOMAIN").String()
)

type BackendDocker struct {}

func init() {
	backend.Register("docker", &BackendDocker{})
}

func (b *BackendDocker) MultipleDomains() bool {
	return true
}

func (b *BackendDocker) Start() error {
	return nil
}

func (b *BackendDocker) GetAddresses() (map[string][]string, error) {
	urls, error := getContainerUrls()
	handling.Check(error)
	return urls, nil
}

func getContainerUrls() (map[string][]string, error) {
	// These are the URL's (keyed by domain) that we will return.
	urls := make(map[string][]string)

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
				urls[envDomain] = append(urls[envDomain], getProxyUrl(portObject))				
			}
		}
    }

    return urls, nil
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
