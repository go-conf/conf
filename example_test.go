package conf_test

import (
	"gopkg.in/conf.v0"
)

type Config struct {
	JsonFile   string // Allow env/flag config to control json file
	AppName    string `conf:"name" json:"app_name"`
	ListenPort int    `conf:"port"`
	Listen     string
}

func Example() {
	// config stores defaults and also receivesresultant config
	config := &Config{}

	c := conf.New("env", "flag", "json")
	c.Load(config)
}
