package config

import (
	"encoding/json"
	"idp/Jobs/api"
	"idp/Jobs/usecases"
	"io/ioutil"
	"os"
)

type MicroserviceConfig struct {
	confFile string
}

func NewMicroserviceConfig(confFile string) *MicroserviceConfig {
	return &MicroserviceConfig{
		confFile: confFile,
	}
}

func (mc *MicroserviceConfig) GetConfig() (*api.ServerConfig, error) {
	file, err := os.Open(mc.confFile)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var jsonResult map[string]interface{}
	if err = json.Unmarshal([]byte(content), &jsonResult); err != nil {
		return nil, err
	}

	serverConfig := &api.ServerConfig{
		JobManagerConf: usecases.JobManagerConfig{
			DBURL: jsonResult["DBURL"].(string),
			DBUser: jsonResult["DBUser"].(string),
			DBPass: jsonResult["DBPass"].(string),
		},
		Transport: jsonResult["transport"].(string),
		Port: jsonResult["port"].(string),
	}

	return serverConfig, nil
}