package health

import (
	"errors"
	"fmt"

	ver "github.com/go-chassis/mesher/adminapi/version"

	metricsink "github.com/ServiceComb/cse-collector"
	"github.com/ServiceComb/go-cc-client/configcenter-client"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"
)

//GetMesherHealth returns health
func GetMesherHealth() *Health {
	serviceName, version, err := getServiceStatus()
	resp := &Health{
		ServiceName: serviceName,
		Version:     version,
		Status:      Green,
		ConnectedConfigCenterClient: isConfigCenterConnected(),
		ConnectedMonitoring:         isMornitorServerConnected(),
		Error:                       "",
	}
	if err != nil {
		lager.Logger.Error("health check failed", err)
		resp.Status = Red
		resp.Error = err.Error()
	}
	return resp
}

func getServiceStatus() (serviceName, version string, err error) {
	appID := config.GlobalDefinition.AppID
	microServiceName := config.SelfServiceName
	version = config.SelfVersion
	if version == "" {
		version = ver.DefaultVersion
	}
	environment := config.MicroserviceDefinition.ServiceDescription.Environment
	serviceID, err := registry.DefaultServiceDiscoveryService.GetMicroServiceID(appID, microServiceName, version, environment)
	if err != nil {
		return microServiceName, version, err
	}
	if len(serviceID) == 0 {
		return microServiceName, version, errors.New("serviceID is empty")
	}
	instances, err := registry.DefaultServiceDiscoveryService.GetMicroServiceInstances(serviceID, serviceID)
	if err != nil {
		return microServiceName, version, err
	}
	if len(instances) == 0 {
		return microServiceName, version, errors.New("no instance found")
	}
	for _, instance := range instances {
		ok, err := registry.DefaultRegistrator.Heartbeat(serviceID, instance.InstanceID)
		if err != nil {
			return microServiceName, version, err
		}
		if !ok {
			e := fmt.Errorf("heartbeat failed, instanceId: %s", instance.InstanceID)
			return microServiceName, version, e
		}
	}
	return microServiceName, version, nil
}

func isConfigCenterConnected() bool {
	if configcenterclient.MemberDiscoveryService == nil {
		return false
	}

	// Getting config center ip's using refresh members handled in GetConfigServer function based on Autodiscovery
	configServerHosts, err := configcenterclient.MemberDiscoveryService.GetConfigServer()
	if err != nil || len(configServerHosts) == 0 {
		return false
	}
	return true
}

func isMornitorServerConnected() bool {
	return metricsink.IsMonitoringConnected
}
