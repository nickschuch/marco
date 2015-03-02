package backend_tutum

import (
	"strings"
	"time"

	"github.com/daryl/cash"
	"github.com/nickschuch/go-tutum/tutum"
	"gopkg.in/alecthomas/kingpin.v1"

	"github.com/nickschuch/marco/backend"
	"github.com/nickschuch/marco/handling"
)

var (
	selTutumUser      = kingpin.Flag("tutum-user", "The Tutum API user.").OverrideDefaultFromEnvar("TUTUM_USER").String()
	selTutumAPI       = kingpin.Flag("tutum-key", "The Tutum API key.").OverrideDefaultFromEnvar("TUTUM_KEY").String()
	selTutumDomainEnv = kingpin.Flag("tutum-domain-env", "The container environment variable that is used as a domain identifier.").Default("DOMAIN").String()
)

type BackendTutum struct {
	cache *cash.Cash
}

func init() {
	backend.Register("tutum", &BackendTutum{})
}

func (b *BackendTutum) Start() error {
	// Create a brand new cache item that we can use to
	// store all the domains and there associated list of
	// addresses.
	b.cache = cash.New(cash.Conf{
		// Default expiration.
		5 * time.Minute,
		// Clean interval.
		30 * time.Minute,
	})

	return nil
}

func (b *BackendTutum) Addresses(domain string) ([]string, error) {
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
	b.cache.Set(domain, list, 5*time.Minute)

	return list, nil
}

func getListByDomain(domain string) ([]string, error) {
	var endpoints []string

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
		var endpoints []string

		service, error := tutum.GetService(s.Uuid)
		handling.Check(error)

		// Loop through and look for the domain variable.
		var found bool
		for _, env := range service.ContainerEnvVars {
			if (env.Key == *selTutumDomainEnv) && (env.Value == domain) {
				found = true
				break
			}
		}

		// If we didn't find any variable identifiers then we should
		// continue onto the next Tutum service.
		if !found {
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
		if len(endpoints) > 0 {
			return endpoints, nil
		}
	}

	return endpoints, nil
}
