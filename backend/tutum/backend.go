package tutum

import (
	"strings"

	"gopkg.in/alecthomas/kingpin.v1"
	"github.com/nickschuch/go-tutum/tutum"

	backend ".."
	handling "../../handling"
)

var (
	selTutumUser      = kingpin.Flag("tutum-user", "The Tutum API user.").OverrideDefaultFromEnvar("TUTUM_USER").String()
	selTutumAPI       = kingpin.Flag("tutum-key", "The Tutum API key.").OverrideDefaultFromEnvar("TUTUM_KEY").String()
	selTutumDomainEnv = kingpin.Flag("tutum-domain-env", "The container environment variable that is used as a domain identifier.").Default("DOMAIN").String()
)

type BackendTutum struct {}

func init() {
	backend.Register("tutum", &BackendTutum{})
}

func (b *BackendTutum) MultipleDomains() bool {
	return true
}

func (b *BackendTutum) Start() error {
	return nil
}

func (b *BackendTutum) GetAddresses() (map[string][]string, error) {
	addresses := make(map[string][]string)

	tutum.User = *selTutumUser
    tutum.ApiKey = *selTutumAPI

	// Get the list of services.
	// https://support.tutum.co/support/solutions/articles/5000525024-your-first-service
	services, error := tutum.ListServices()
    handling.Check(error)

    for _, s := range services {
    	// These are the things we are looking for.
    	//   * Domain delta so we can load balancer more than one.
    	//   * A list of endpoints to route to.
    	var domain string
    	var endpoints []string


	    service, error := tutum.GetService(s.Uuid)
        handling.Check(error)

        // Loop through and look for the domain variable.
        var found bool
        for _, env := range service.ContainerEnvVars {
            if env.Key == *selTutumDomainEnv {
            	domain = env.Value
                found = true
                break
            }
        }

        // If we didn't find any variable identifiers then we should
        // continue onto the next Tutum service.
        if ! found {
        	continue
        }

        // Now we grab our exposes URLs (if they exist).
        for _, uri := range service.ContainerUris {
            uuidSlice := strings.Split(uri, "/")
            uuid := uuidSlice[4]
            container, error := tutum.GetContainer(uuid)
            handling.Check(error)

            // Now that we have the container we need to find out
            // if tutum have exposed an "endpoint" for out container.
            for _, port := range container.ContainerPorts {
                if port.EndpointUri != "" {
                    endpoints = append(endpoints, port.EndpointUri)
                }
            }
        }

        // If we didn't find any endpoints we should continue onto the next
        // service to see if we can find some there.
		if len(endpoints) <= 0 {
			continue
		}

    	// Now we add our endpoints to the map/slice so the reconciler can
    	// assign a balancer.
		addresses[domain] = endpoints
    }

	return addresses, nil
}
