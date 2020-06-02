package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Project struct {
	Name      string   `yaml:"name"`
	Dir       string   `yaml:"dir"`
	ErrLog    string   `yaml:"errLog"`
	AppLog    string   `yaml:"appLog"`
	Commands  []string `yaml:"commands"`
	OnFailure string   `yaml:"onFailure"`
}

type Config struct {
	Projects []*Project `yaml:"projects,omitempty"`
}

func NewConfig(path string) (*Config, error) {
	c := &Config{}
	y, err := ioutil.ReadFile(path)
	if err != nil {
		return c, err
	}
	yaml.Unmarshal(y, &c)
	return c, nil
}

func (c *Config) getProject(name string) *Project {
	for _, p := range c.Projects {
		if name == p.Name {
			return p
		}
	}
	return &Project{}
}
