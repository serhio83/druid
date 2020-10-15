package config

import "os"

// configure rergistry connect & credentials
// configure listen port
// configure bbolt storage path

// Config ...
type Config struct {
	RegHost    string
	RegPort    string
	RegUser    string
	RegPass    string
	ListenPort string
	BboltPath  string
}

// New returns new Config
func New() *Config {
	c := Config{}
	c.RegHost = os.Getenv("DRUID_REGISTRY_HOST")
	c.RegPort = os.Getenv("DRUID_REGISTRY_PORT")
	c.RegUser = os.Getenv("DRUID_REGISTRY_USER")
	c.RegPass = os.Getenv("DRUID_REGISTRY_PASSWORD")
	c.ListenPort = os.Getenv("DRUID_LISTEN_PORT")
	c.BboltPath = os.Getenv("DRUID_DATA_PATH")

	return &c
}
