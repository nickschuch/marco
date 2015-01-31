package tutum

import (
	"gopkg.in/alecthomas/kingpin.v1"

	backend ".."
)

var (
	selTutumRefresh = kingpin.Flag("tutum-refresh", "How long to wait between querying the Tutum api.").Default("10s").String()
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
	var addresses map[string][]string
	addresses["*"] = []string{}
	return addresses, nil
}
