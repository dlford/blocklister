package config

type Config struct {
	ListConfigs []ListConfig `yaml:"lists"`
	Schedule    string       `yaml:"schedule"`
}

type ListConfig struct {
	Title  string   `yaml:"title"`
	URL    string   `yaml:"url"`
	Chains []string `yaml:"chains"`
}
