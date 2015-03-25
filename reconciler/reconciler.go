package reconciler

import (
	"github.com/nickschuch/marco/backend"
	_ "github.com/nickschuch/marco/backend/docker"
	_ "github.com/nickschuch/marco/backend/mock"
	_ "github.com/nickschuch/marco/backend/tutum"
	"github.com/nickschuch/marco/balancer"
	_ "github.com/nickschuch/marco/balancer/first"
	_ "github.com/nickschuch/marco/balancer/random"
	_ "github.com/nickschuch/marco/balancer/round"
	"github.com/nickschuch/marco/handling"
	"github.com/nickschuch/marco/logging"
)

type Reconciler struct {
	backendType  string
	balancerType string
	backend      backend.Backend
	balancers    map[string]balancer.Balancer
}

func (r *Reconciler) SetBackendType(t string) {
	r.backendType = t
}

func (r *Reconciler) GetBackendType() string {
	return r.backendType
}

func (r *Reconciler) SetBalancerType(t string) {
	r.balancerType = t
}

func (r *Reconciler) GetBalancerType() string {
	return r.balancerType
}

func (r *Reconciler) GetBalancer(domain string) (balancer.Balancer, error) {
	driver, error := balancer.New(r.balancerType)
	handling.Check(error)
	return driver, nil
}

func (r *Reconciler) Start() error {
	// Kick off the Backend watch process.
	// This will allow the backend to have the opitunity to catch
	// new hosts / environments.
	backend, err := backend.New(r.backendType)
	handling.Check(err)

	// Set and backend and start the process.
	r.backend = backend
	r.backend.Start()
	return nil
}

func (r *Reconciler) Address(domain string) (string, error) {
	// Check if the list of IP's have changed.
	addresses, _ := r.backend.Addresses(domain)

	// Apply the list of addresses to the appropriate balancer.
	b, _ := r.GetBalancer(domain)
	b.SetAddressList(addresses)

	// Get the new address from the appropriate balancer.
	address, _ := b.Next()

	// We need to make sure we have an address to hand back for this domain.
	// If we don't then we should fail and log it so we can follow up.
	if len(address) > 0 {
		// We log this so we can see that containers we in the mix.
		logging.ContainerFound(domain, address)
		return address, nil
	} else {
		// Make sure we record the fact that we could not find a
		// container for the requested domain.
		logging.ContainerNotFound(domain)
	}

	return "", nil
}
