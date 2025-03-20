package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

var (
	configPath string
	config     *Config
)

type Config struct {
	Cookie    string `toml:"cookie"`
	BDUSS     string `toml:"BDUSS"`
	STOKEN    string `toml:"STOKEN"`
	UserAgent string `toml:"user-agent"`
}

func initConfig() {
	// get config path
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		configHome = filepath.Join(os.Getenv("HOME"), ".config")
	}

	configPath = filepath.Join(configHome, "tieba-sign", "config.toml")

	// read
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("failed to read config file: %s: %s\n", configPath, err)
	}

	// unmarshal
	config = &Config{}
	err = toml.Unmarshal(configBytes, config)
	if err != nil {
		log.Fatalf("failed to unmarshal config file: %s: %s\n", configPath, err)
	}
}
