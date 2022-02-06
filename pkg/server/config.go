package server

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Signl4 ConfigSignl4 `yaml:"signl4"`
}

type ConfigSignl4 struct {
	StatusKey  string `yaml:"statusKey"`
	TeamSecret string `yaml:"teamSecret"`
	Template   string `yaml:"template"`
}

func ReadConfig(filePath string) (ConfigSignl4, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return ConfigSignl4{}, err
	}

	var c Config
	err = yaml.Unmarshal(content, &c)
	if err != nil {
		return ConfigSignl4{}, nil
	}

	return c.Signl4, nil
}
