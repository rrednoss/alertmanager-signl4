package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	AllowInsecureTLSConfig bool   `yaml:"allowInsecureTLSConfig"` // TODO (rednoss): Evaluate!
	StatusKey              string `yaml:"statusKey"`
	TeamSecret             string `yaml:"teamSecret"`
	Template               string `yaml:"template"`
}

func NewAppConfig(path string) AppConfig {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var appConfig AppConfig
	err = yaml.Unmarshal(content, &appConfig)
	if err != nil {
		panic(err)
	}
	return appConfig
}
