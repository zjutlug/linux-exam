package cmd

import (
	"os"
	"strings"

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
	conf = &Config{
		APIBaseURL:  "http://127.0.0.1:8080",
		ContainerId: "",
	}
	if viper.IsSet("API_BASE_URL") {
		conf.APIBaseURL = viper.GetString("API_BASE_URL")
	}

	b, err := os.ReadFile("/etc/container_id")
	if err != nil {
		conf.ContainerId = ""
	}
	conf.ContainerId = strings.TrimSpace(string(b))
}
