package config

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Slack  Slack  `yaml:"slack"`
	WeWork WeWork `yaml:"wework"`
}

func Parse(b []byte) (*Config, error) {
	c := new(Config)
	if err := yaml.Unmarshal(b, c); err != nil {
		return nil, fmt.Errorf("yaml.Unmarshal error: %v", err)
	}

	return c, nil
}

func Load(filename string) (*Config, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Load config error: %v", err)
	}

	return Parse(b)
}
