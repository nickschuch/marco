package main

import (
	"github.com/pmylund/go-cache"
	"strings"
)

type cacheContainerObject struct {
	Domain string
	Url    string
}

var cacheContainers = cache.New(cache.NoExpiration, cache.NoExpiration)

func getCachedContainer(id string) cacheContainerObject {
	var container cacheContainerObject
	if x, found := cacheContainers.Get(id); found {
		container = x.(cacheContainerObject)
	}
	return container
}

func getContainer(id string) cacheContainerObject {
	container := getCachedContainer(id)
	if len(container.Domain) > 0 {
		return container
	} else {
		dockerContainer := getDockerContainer(id)

		// Get the domain so we can work out which site this container belong to.
		envDomain := getContainerEnv("DOMAIN", dockerContainer.Config.Env)

		// Here we build the proxy URL based on the exposed values provided
		// by NetworkSettings. If a container has not been exposed, it will
		// not work. We then build a URL based on these exposed values and:
		//   * Add a container reference so we can perform safe operations
		//     in the future.
		//   * Add the built url to the proxy lists for load balancing.
		var url string
		for portString, portObject := range dockerContainer.NetworkSettings.Ports {
			port := getPort(portString)
			if strings.Contains(ports, port) {
				url = getProxyUrl(portObject)
			}
		}
		container = addContainer(id, envDomain, url)
	}
	return container
}

func addContainer(id string, domain string, url string) cacheContainerObject {
	container := cacheContainerObject{domain, url}
	cacheContainers.Set(id, container, cache.NoExpiration)
	return container
}

func removeContainer(id string) {
	cacheContainers.Delete(id)
}
