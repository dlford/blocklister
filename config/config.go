package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Lists    []Blocklist `yaml:"lists"`
	Schedule string      `yaml:"schedule"`
}

type Blocklist struct {
	Title string `yaml:"title"`
	URL   string `yaml:"url"`
}

func (c *Config) GetConf() *Config {
	var path string
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	if path == "" {
		// TODO: Change to /etc/blocklister.yml
		path = "blocklister.yml"
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
