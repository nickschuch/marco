package tutum

type PortInfo struct {
	Container   string `json:"container"`
	EndpointUri string `json:"endpoint_uri"`
	InnerPort   int    `json:"inner_port"`
	OuterPort   int    `json:"outer_port"`
	Pubished    bool   `json:"published"`
	Protocol    string `json:"protocol"`
}

type EnvVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
