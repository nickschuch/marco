package main

import (
	"strings"
    log "github.com/Sirupsen/logrus"
    "github.com/samalba/dockerclient"
)

func getDockerClient() *dockerclient.DockerClient {
    dockerClient, err := dockerclient.NewDockerClient(endpoint, nil)
    if err != nil {
        log.Fatal(err)
    }
    return dockerClient
}

func getDockerContainer(id string) *dockerclient.ContainerInfo {
	dockerClient := getDockerClient()
	container, _ := dockerClient.InspectContainer(id)
	return container
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

func buildProxyUrl(binding []dockerclient.PortBinding) string {
    // Ensure we have PortBinding values to build against.
    if len(binding) <= 0 {
        return ""
    }

    // Handle IP 0.0.0.0 the same way Swarm does. We replace this with an IP
    // that uses a local context.
    ip := binding[0].HostIp
    port := binding[0].HostPort

    if ip == "0.0.0.0" {
        ip = host
    }

    return "http://" + ip + ":" + port
}
