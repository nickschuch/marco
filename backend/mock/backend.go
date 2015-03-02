package backend_mock

import (
	"github.com/nickschuch/marco/backend"
)

type BackendMock struct{}

func init() {
	backend.Register("mock", &BackendMock{})
}

func (b *BackendMock) Addresses(domain string) ([]string, error) {
	addresses := []string{
		"1.2.3.4",
		"5.6.7.8",
	}
	return addresses, nil
}

func (b *BackendMock) Start() error {
	// Nothing to see here.
	return nil
}
