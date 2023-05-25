package framework

import (
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port int `yaml:"port"`

	path string
}

func ReadConfig(path string) *Config {
	if strings.HasPrefix(path, "/") {
		path = root() + path
	}
	c := &Config{
		path: path + "/" + Env() + ".yaml",
	}
	c.readYaml()
	return c
}

func (c *Config) readYaml() {
	content, err := os.ReadFile(c.path)
	check(err)
	check(yaml.Unmarshal(content, c))
}
