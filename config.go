package main

import (
	"io/ioutil"
	"path/filepath"
	"gopkg.in/yaml.v2"
)

type Config struct {
	DB struct {
		Host     string `yaml:"host,omitempty"`
		Port     int    `yaml:"port,omitempty"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
		Encoding string `yaml:"names,omitempty"`
	} `yaml:"db"`

	Timezone string `yaml:"timezone"`
}

func NewConfig(path string) (*Config, error) {
	config := &Config{}

	filename, err := filepath.Abs(path)

	if err != nil {
		return nil, err
	}

	raw, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(raw, &config)

	if err != nil {
		return nil, err
	}

	return config, nil
}