package config

type Config struct {
	ListConfigs []ListConfig `yaml:"lists"`
	Schedule    string       `yaml:"schedule"`
}

type ListConfig struct {
	Title    string   `yaml:"title"`
	URL      string   `yaml:"url"`
	MaxElem  int      `yaml:"max_elem"`
	Chains   []string `yaml:"chains"`
	Schedule string   `yaml:"schedule"`
}
