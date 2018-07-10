package common

//Constants for default app and version
const (
	DefaultApp     = "default"
	DefaultVersion = "0.0.1"
)

//Constants for buildtag app and version
const (
	BuildInTagApp     = "app"
	BuildInTagVersion = "version"
)

//ComponentName is contant for component name
const ComponentName = "mesher"

//ModeSidecar is constant for side car mode
const ModeSidecar = "sidecar"

//ModePerHost is constant for side car mode
const ModePerHost = "per-host"

//Constants for env specific addr and service ports
const (
	//EnvSpecificAddr Deprecated
	EnvSpecificAddr = "SPECIFIC_ADDR"
	EnvServicePorts = "SERVICE_PORTS"
)

//HTTPProtocol is constant for protocol
const HTTPProtocol = "http"

//Constants for provider and consumer handlers
const (
	ChainConsumerOutgoing = "outgoing"
	ChainProviderIncoming = "incoming"
	ChainConsumerEgress   = "egress"
)

//ConsumerEgress constant for egress
const ConsumerEgress = "ConsumerEgress"
