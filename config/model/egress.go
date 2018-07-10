package model

//EgressConfig is the struct having info about egress rule destinations
type EgressConfig struct {
	Destinations map[string][]*EgressRule `yaml:"egressRule"`
}

//EgressRule has hosts and ports information
type EgressRule struct {
	Hosts []string      `yaml:"hosts"`
	Ports []*EgressPort `yaml:"ports"`
}

//EgressPort protocol and the corresponding port
type EgressPort struct {
	Port     int32  `yaml:"port"`
	Protocol string `yaml:"protocol"`
}
