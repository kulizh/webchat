package utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ListenPort string
	Routes     struct {
		Ws      string `yaml:"ws"`
		Counter string `yaml:"counter"`
	}
}

func ParseConfig(filename string) Config {
	yamlFile, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	checkStructure(config)

	return config
}

func checkStructure(config Config) {
	if len(config.ListenPort) == 0 {
		panic("Missing ListenPort")
	}

	if len(config.Routes.Ws) == 0 || len(config.Routes.Counter) == 0 {
		panic("Empty routes")
	}
}
