package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v1"
)

type Config struct {
	Controller struct {
		URL     string `yaml:"url"`
		AuthKey string `yaml:"authkey"`
	} `yaml:"controller"`
}

func loadConfig(filepath string) (*Config, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	config := &Config{}

	err = yaml.Unmarshal(data, config)

	return config, nil
}

func saveConfig(filepath string, conf *Config) error {
	data, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath, data, 0600)
}
