package config

import (
	"io/ioutil"
	"log"
	"sync"

	"gopkg.in/yaml.v2"
)

var (
	once   sync.Once
	Signl4 Signl4Config
)

func NewSignl4Config() Signl4Config {
	once.Do(func() {
		config, err := readConfig("/conf/signl4.yaml")
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		Signl4 = config
	})
	return Signl4
}

type AppConfig struct {
	Signl4 Signl4Config `yaml:"signl4"`
}

type Signl4Config struct {
	// AllowInsecureTLSConfig bool   `yaml:"allowInsecureTLSConfig"`
	StatusKey  string `yaml:"statusKey"`
	TeamSecret string `yaml:"teamSecret"`
	Template   string `yaml:"template"`
}

func readConfig(filePath string) (Signl4Config, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Signl4Config{}, err
	}

	var c AppConfig
	err = yaml.Unmarshal(content, &c)
	if err != nil {
		return Signl4Config{}, nil
	}

	return c.Signl4, nil
}
