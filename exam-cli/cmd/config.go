package cmd

import (
	"github.com/spf13/viper"
)

// Config 配置项结构
type Config struct {
	APIBaseURL  string
	ContainerId string
}

var conf *Config

func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("EXAM")
	conf = &Config{
		APIBaseURL:  "http://127.0.0.1",
		ContainerId: "",
	}
	if viper.IsSet("API_BASE_URL") {
		conf.APIBaseURL = viper.GetString("API_BASE_URL")
	}

	if viper.IsSet("CONTAINER_ID") {
		conf.ContainerId = viper.GetString("CONTAINER_ID")
	}
}
