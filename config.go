package aegateway

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type (
	GatewayRoute struct {
		Pattern string `yaml:"pattern"`
		Dest string `yaml:"dest"`
	}

	GatewayConfig struct {
		Routes []GatewayRoute
	}
)

func LoadConfig(name string) *GatewayConfig {
	yamlStr, err := ioutil.ReadFile("./gateway.yaml")
	if err != nil {
		panic(err)
	}
	
	config := GatewayConfig{}
	if err := yaml.Unmarshal(yamlStr, &config); err != nil {
		panic(err)
	}

	return &config
}


