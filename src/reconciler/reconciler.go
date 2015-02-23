package reconciler

import (
	"time"

	backend "../backend"
	_ "../backend/docker"
	_ "../backend/tutum"
	balancer "../balancer"
	_ "../balancer/round"
	handling "../handling"
	logging "../logging"
)

type Reconciler struct {
	backend  backend.Backend
	balancer string
	balancers map[string]balancer.Balancer
	refresh string
}

func (r *Reconciler) AddBackend(t string) error {
	driver, error := backend.New(t)
	handling.Check(error)
	r.backend = driver

	logging.Info("Set backend: " + t)
	return nil
}

func (r *Reconciler) AddBalancer(t string) error {
	r.balancer = t
	logging.Info("Set balancer: " + t)
	return nil
}

func (r *Reconciler) NewBalancer() (balancer.Balancer, error) {
	driver, error := balancer.New(r.balancer)
	handling.Check(error)
	return driver, nil
}

func (r *Reconciler) Start() error {
	go r.Watch()
	return nil
}

func (r *Reconciler) Watch() error {
	duration, error := time.ParseDuration(r.refresh)
	handling.Check(error)	

	// This is an infinite loop that repopulates the load balancers
	// based on the "backend" configuration.
	for {
		domains, error := r.backend.GetAddresses()
		handling.Check(error)

		for domain, addresses := range domains {
			// Setup a new balancer to take the place of the old one.
			newBalancer, error := r.NewBalancer()
			handling.Check(error)

			// Get the addresses based on the latest backend poll.
			newBalancer.SetAddresses(addresses)
			r.balancers = make(map[string]balancer.Balancer)
			r.balancers[domain] = newBalancer
		}
		time.Sleep(duration)
	}

	return nil
}

func (r *Reconciler) GetAddress(domain string) (string, error) {
	// If the backend only supports a single set of containers then
	// it should have returned them with the following delta.
	if ! r.backend.MultipleDomains() {
		domain = "*"
	}

	if _, ok := r.balancers[domain]; ok {
    	address, error := r.balancers[domain].GetAddress()
		handling.Check(error)

		// We log this so we can see that containers we in the mix.
		logging.ContainerFound(domain, address)

		return address, nil
	} else {
		// Make sure we record the fact that we could not find a
		// container for the requested domain.
		logging.ContainerNotFound(domain)

		return "", nil
	}
}

func (r *Reconciler) SetRefresh(t string) error {
	r.refresh = t
	return nil
}
