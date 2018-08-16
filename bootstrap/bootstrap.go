/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package bootstrap

import (
	"log"
	"strings"

	"github.com/go-chassis/mesher/adminapi"
	"github.com/go-chassis/mesher/adminapi/version"
	"github.com/go-chassis/mesher/cmd"
	"github.com/go-chassis/mesher/common"
	"github.com/go-chassis/mesher/config"
	"github.com/go-chassis/mesher/register"
	"github.com/go-chassis/mesher/resolver"

	"github.com/go-chassis/go-chassis"
	chassisHandler "github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/metadata"
)

// Start initialize configs and components
func Start() error {
	if err := config.InitProtocols(); err != nil {
		return err
	}
	if err := config.Init(); err != nil {
		return err
	}
	if err := resolver.Init(); err != nil {
		return err
	}
	if err := DecideMode(); err != nil {
		return err
	}
	if err := adminapi.Init(); err != nil {
		log.Println("Error occurred in starting admin server", err)
	}
	register.AdaptEndpoints()
	if cmd.Configs.LocalServicePorts == "" {
		lager.Logger.Warnf("local service ports is missing, service can not be called by mesher")
	} else {
		lager.Logger.Infof("local service ports is [%v]", cmd.Configs.PortsMap)
	}

	return nil

}

//DecideMode get config mode
func DecideMode() error {
	config.Mode = cmd.Configs.Mode
	lager.Logger.Info("Running as "+config.Mode, nil)
	return nil
}

//RegisterFramework registers framework
func RegisterFramework() {
	if framework := metadata.NewFramework(); cmd.Configs.Mode == common.ModeSidecar {
		version := GetVersion()
		framework.SetName("Mesher")
		framework.SetVersion(version)
		framework.SetRegister("SIDECAR")
	} else if cmd.Configs.Mode == common.ModePerHost {
		framework.SetName("Mesher")
	}
}

//GetVersion returns version
func GetVersion() string {
	versionID := version.Ver().Version
	if len(versionID) == 0 {
		return version.DefaultVersion
	}
	return versionID
}

//SetHandlers leverage go-chassis API to set default handlers if there is no define in chassis.yaml
func SetHandlers() {
	consumerChain := strings.Join([]string{
		chassisHandler.Router,
		chassisHandler.RatelimiterConsumer,
		chassisHandler.BizkeeperConsumer,
		chassisHandler.Loadbalance,
		chassisHandler.Transport,
	}, ",")
	providerChain := strings.Join([]string{
		chassisHandler.RatelimiterProvider,
		chassisHandler.Transport,
	}, ",")
	consumerChainMap := map[string]string{
		common.ChainConsumerOutgoing: consumerChain,
	}
	providerChainMap := map[string]string{
		common.ChainProviderIncoming: providerChain,
	}
	chassis.SetDefaultConsumerChains(consumerChainMap)
	chassis.SetDefaultProviderChains(providerChainMap)
}
