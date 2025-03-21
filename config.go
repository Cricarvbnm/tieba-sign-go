package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

var (
	config *Config
)

type Config struct {
	BDUSS  string
	STOKEN string
}

func loadConfig() error {
	// get path
	configPath := filepath.Join(getConfigDir(), "config.toml")
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("无法读取配置: %s: %w", configPath, err)
	}

	// parse config
	config = &Config{}
	if err := toml.Unmarshal(configBytes, config); err != nil {
		return fmt.Errorf("无法解析配置: %s: %w", configPath, err)
	} else if err := checkConfig(config); err != nil {
		return fmt.Errorf("错误配置: %s: %w", configPath, err)
	}

	return nil
}

func checkConfig(config *Config) error {
	if config.BDUSS == "" {
		return fmt.Errorf("缺失BDUSS")
	}

	if config.STOKEN == "" {
		return fmt.Errorf("缺失STOKEN")
	}

	return nil
}
