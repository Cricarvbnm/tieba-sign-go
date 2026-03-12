package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

type Config struct {
	BDUSS  string
	STOKEN string
}

func Load(path string) (*Config, error) {
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("无法读取配置文件: %s: %w", path, err)
	}

	cfg := &Config{}
	if err := toml.Unmarshal(configBytes, cfg); err != nil {
		return nil, fmt.Errorf("无法解析配置文件: %s: %w", path, err)
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.BDUSS == "" {
		return fmt.Errorf("配置缺失 BDUSS")
	}
	if c.STOKEN == "" {
		return fmt.Errorf("配置缺失 STOKEN")
	}
	return nil
}

func Home() string {
	home := os.Getenv("XDG_CONFIG_HOME")
	if home == "" {
		home = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return home
}

func DefaultPath() string {
	return filepath.Join(Home(), "tieba-sign", "config.toml")
}
