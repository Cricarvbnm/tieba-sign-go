package main

import (
	"os"
	"path/filepath"
)

func getConfigHome() string {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		configHome = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return configHome
}

func getDataHome() string {
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		dataHome = filepath.Join(os.Getenv("HOME"), ".local", "share")
	}
	return dataHome
}

func getConfigDir() string {
	configDir := filepath.Join(getConfigHome(), "tieba-sign")
	return configDir
}

func getDataDir() string {
	dataDir := filepath.Join(getDataHome(), "tieba-sign")
	return dataDir
}

func getLogDir() string {
	logDir := filepath.Join(getDataDir(), "log")
	return logDir
}
