package config

import (
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/server"
	"github.com/ServiceComb/go-chassis/util/fileutil"
	"github.com/go-chassis/mesher/cmd"
	"github.com/go-chassis/mesher/common"
	egressmodel "github.com/go-chassis/mesher/config/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

//Constant for mesher conf file
const (
	ConfFile       = "mesher.yaml"
	EgressConfFile = "egress.yaml"
)

//Mode is of type string which gives mode of mesher deployment
var Mode string
var mesherConfig *MesherConfig
var egressConfig *egressmodel.EgressConfig

//GetConfig returns mesher config
func GetConfig() *MesherConfig {
	return mesherConfig
}

//SetConfig sets new mesher config from input config
func SetConfig(nc *MesherConfig) {
	if mesherConfig == nil {
		mesherConfig = &MesherConfig{}
	}
	*mesherConfig = *nc
}

//GetEgressConfig returns Egress config
func GetEgressConfig() *egressmodel.EgressConfig {
	return egressConfig
}

//SetEgressConfig sets new egress config from input config
func SetEgressConfig(nc *egressmodel.EgressConfig) {
	if egressConfig == nil {
		egressConfig = &egressmodel.EgressConfig{}
	}
	*egressConfig = *nc
}

//GetConfigFilePath returns config file path
func GetConfigFilePath(key string) (string, error) {
	if cmd.Configs.ConfigFile == "" {
		wd, err := fileutil.GetWorkDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(wd, "conf", key), nil
	}
	return cmd.Configs.ConfigFile, nil
}

//InitProtocols initiates protocols
func InitProtocols() error {
	// todo if sdk init failed, do not call the data
	if len(config.GlobalDefinition.Cse.Protocols) == 0 {
		config.GlobalDefinition.Cse.Protocols = map[string]model.Protocol{
			common.HTTPProtocol: {Listen: "127.0.0.1:30101"},
		}

		return server.Init()
	}
	return nil
}

//Init reads config and initiates
func Init() error {
	mesherConfig = &MesherConfig{}
	contents, err := GetConfigContents(ConfFile)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal([]byte(contents), mesherConfig); err != nil {
		return err
	}

	egressConfig = &egressmodel.EgressConfig{}
	egressContents, err := GetConfigContents(EgressConfFile)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal([]byte(egressContents), egressConfig); err != nil {
		return err
	}

	return nil
}

//GetConfigContents returns config contents
func GetConfigContents(key string) (string, error) {
	f, err := GetConfigFilePath(key)
	if err != nil {
		return "", err
	}
	var contents string
	//route rule yaml file's content is value of a key
	//So read from config center first,if it is empty, Try to set file content into memory key value
	contents = archaius.GetString(key, "")
	if contents == "" {
		contents = SetKeyValueByFile(key, f)
	}
	return contents, nil
}

//SetKeyValueByFile reads mesher.yaml and gets key and value
func SetKeyValueByFile(key, f string) string {
	var contents string
	if _, err := os.Stat(f); err != nil {
		lager.Logger.Warn(err.Error(), nil)
		return ""
	}
	b, err := ioutil.ReadFile(f)
	if err != nil {
		lager.Logger.Error("Can not read yaml file", err)
		return ""
	}
	contents = string(b)
	archaius.AddKeyValue(key, contents)
	return contents
}
