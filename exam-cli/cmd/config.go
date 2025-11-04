package cmd

import (
	"os"
)

// Config 配置项结构
type Config struct {
	APIBaseURL string
}

// readConfig 从环境变量读取配置
func readConfig() map[string]string {
	config := make(map[string]string)

	// 从环境变量读取 API 地址
	if apiURL := os.Getenv("EXAM_API_URL"); apiURL != "" {
		config["api_base_url"] = apiURL
	} else {
		config["api_base_url"] = "http://127.0.0.1:8080" // 默认值
	}

	return config
}
