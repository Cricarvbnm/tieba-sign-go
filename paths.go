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

func getStateHome() string {
	stateHome := os.Getenv("XDG_STATE_HOME")
	if stateHome == "" {
		stateHome = filepath.Join(os.Getenv("HOME"), ".local", "state")
	}
	return stateHome
}

func getConfigDir() string {
	configDir := filepath.Join(getConfigHome(), "tieba-sign")
	return configDir
}

func getStateDir() string {
	stateDir := filepath.Join(getStateHome(), "tieba-sign")
	return stateDir
}

func getLogDir() string {
	logDir := filepath.Join(getStateDir(), "log")
	return logDir
}
