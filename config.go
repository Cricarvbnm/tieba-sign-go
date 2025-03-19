package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

type Config struct {
	Cookie    string `toml:"cookie"`
	UserAgent string `toml:"user-agent"`
	Request   struct {
		Concurrency int `toml:"concurrency"`
	} `toml:"request"`
}

func loadConfig() (*Config, error) {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, _ := os.UserHomeDir()
		configHome = filepath.Join(home, ".config")
	}

	configDir := filepath.Join(configHome, "tieba-sign")

	configFile := filepath.Join(configDir, "config.toml")
	tree, err := toml.LoadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("加载配置文件失败: %w", err)
	}

	var config Config
	if err := tree.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// default values
	if config.Request.Concurrency <= 0 {
		config.Request.Concurrency = 5
	}

	return &config, nil
}
