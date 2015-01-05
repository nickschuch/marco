package main

import (
	log "github.com/Sirupsen/logrus"
)

func populateCache() {
	builtUrls := make(map[string][]string)
	dockerClient := getDockerClient()
	containers, err := dockerClient.ListContainers(false, false, "")
	if err != nil {
		log.Fatal(err)
	}
	for _, c := range containers {
		container := getContainer(c.Id)

		// We need a domain to be able to work out what to proxy.
		if container.Domain == "" {
			continue
		}
		// We need the URL to be able to work out where to proxy.
		if container.Url == "" {
			continue
		}

		addProxy(container.Domain, container.Url)
	}

	// Cache the domains for later.
	for domain, urls := range builtUrls {
		setProxies(domain, urls)
	}
}
