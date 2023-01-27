package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func (c *Config) Load() *Config {
	var path string
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	if path == "" {
		path = "/etc/blocklister.yml"
	}

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Cannot read config file: %s\n#%v", path, err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Cannot unmarshal config file: %v", err)
	}

	return c
}
