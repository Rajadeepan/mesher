package pilot

import (
	"fmt"
	"sync"

	"github.com/ServiceComb/go-chassis/core/config"
	egressmodel "github.com/go-chassis/mesher/config/model"
	"github.com/go-chassis/mesher/egress"
)

func init() { egress.InstallEgressService("pilot", newPilotEgress) }

func newPilotEgress() (egress.Egress, error) { return &PilotEgress{}, nil }

//PilotEgress is pilot egress service
type PilotEgress struct{}

//FetchRouteRule return all rules
func (r *PilotEgress) FetchEgressRule() map[string][]*egressmodel.EgressRule {
	return GetEgressRule()
}

//SetRouteRule set rules
func (r *PilotEgress) SetEgressRule(rr map[string][]*egressmodel.EgressRule) {
	SetEgressRule(rr)
}

//FetchRouteRuleByServiceName get rules for service
func (r *PilotEgress) FetchEgressRuleByName(service string) []*egressmodel.EgressRule {
	return GetEgressRuleByKey(service)
}

//Init init router config
func (r *PilotEgress) Init(o egress.Options) error {
	// SetDestinations router destinations
	for target := range config.GetRouterReference() {
		SetEgressRuleByKey(target, nil)
	}
	// the manager use dests to init, so must init after dests
	if err := InitPilotFetcher(o); err != nil {
		return err
	}
	return refresh()
}


// refresh all the router config
func refresh() error {
	configs := pilotfetcher.GetConfigurations()
	d := make(map[string][]*egressmodel.EgressRule)
	for k, v := range configs {
		rules, ok := v.([]*egressmodel.EgressRule)
		if !ok {
			err := fmt.Errorf("Egress rule type assertion fail, key: %s", k)
			return err
		}
		d[k] = rules
	}
	 ok, _ := egress.ValidateEgressRule(d)
	if ok{
		dests = d
	}
	return nil
}

var dests = make(map[string][]*egressmodel.EgressRule)
var lock sync.RWMutex

// SetEgressRuleByKey set route rule by key
func SetEgressRuleByKey(k string, r []*egressmodel.EgressRule) {
	lock.Lock()
	dests[k] = r
	lock.Unlock()
}

// DeleteEgressRuleByKey set route rule by key
func DeleteEgressRuleByKey(k string) {
	lock.Lock()
	delete(dests, k)
	lock.Unlock()
}

// GetEgressRuleByKey get route rule by key
func GetEgressRuleByKey(k string) []*egressmodel.EgressRule {
	lock.RLock()
	defer lock.RUnlock()
	return dests[k]
}

// GetEgressRule get route rule
func GetEgressRule() map[string][]*egressmodel.EgressRule {
	lock.RLock()
	defer lock.RUnlock()
	return dests
}

// SetEgressRule set route rule
func SetEgressRule(rule map[string][]*egressmodel.EgressRule) {
	lock.Lock()
	defer lock.Unlock()
	dests = rule
}
