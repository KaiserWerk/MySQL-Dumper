package main

import (
	"embed"
	"errors"
	"io/ioutil"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

type appConfig struct {
	DSN            string `yaml:"dsn" envconfig:"dsn"`
	BackupInterval uint   `yaml:"backup_interval" envconfig:"backup_interval"`
	BackupPath     string `yaml:"backup_path" envconfig:"backup_path"`
}

//go:embed assets/config
var configAssets embed.FS

var (
	configFile = "app.yaml"
)

func setConfigFile(f string) {
	configFile = f
}

func setupConfig() (*appConfig, bool, error) {
	var (
		conf    appConfig
		created bool
		cont    []byte
		err     error
	)
	if _, err = os.Stat(configFile); err != nil && errors.Is(err, os.ErrNotExist) {
		cont, err = configAssets.ReadFile("assets/config/app.dist.yaml")
		if err != nil {
			return nil, false, err
		}
		created = true
		err = ioutil.WriteFile(configFile, cont, 0644)
		if err != nil {
			return nil, false, err
		}
	} else {
		cont, err = ioutil.ReadFile(configFile)
		if err != nil {
			return nil, false, err
		}
	}

	err = yaml.Unmarshal(cont, &conf)
	if err != nil {
		return nil, false, err
	}
	err = envconfig.Process("mysqldumper", &conf)
	if err != nil {
		return nil, false, err
	}

	return &conf, created, nil
}
