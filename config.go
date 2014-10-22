package main

import (
	"fmt"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"os"
)

type Config struct {
	Host              string
	User              string
	Db                string
	Pass              string
	Assets            string
	BackupStoragePath string
	Log               string
}

func getConfigFromYaml(yamlPath string) (Config, error) {
	config := Config{}

	if len(yamlPath) == 0 {
		return config, fmt.Errorf("config path not supplied")
	}

	if _, err := os.Stat(yamlPath); err != nil {
		return config, fmt.Errorf("config path " + yamlPath + " not valid")
	}

	ymlData, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal([]byte(ymlData), &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func getConfig(flagCfg Config, yamlPath string) (Config, error) {
	config, err := getConfigFromYaml(yamlPath)

	if err != nil {
		return config, err
	}

	// Cli overrides

	if len(flagCfg.Host) > 0 {
		config.Host = flagCfg.Host
	}

	if len(flagCfg.Assets) > 0 {
		config.Assets = flagCfg.Assets
	}

	if len(flagCfg.User) > 0 {
		config.User = flagCfg.User
	}

	if len(flagCfg.Db) > 0 {
		config.Db = flagCfg.Db
	}

	if len(flagCfg.Pass) > 0 {
		config.Pass = flagCfg.Pass
	}

	if len(flagCfg.Log) > 0 {
		config.Log = flagCfg.Log
	}

	if len(flagCfg.BackupStoragePath) > 0 {
		config.BackupStoragePath = flagCfg.BackupStoragePath
	}

	return config, nil
}
