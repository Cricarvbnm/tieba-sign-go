package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	EMERG   = "EMERG: "
	ALERT   = "ALERT: "
	CRIT    = "CRIT: "
	ERR     = "ERR: "
	WARNING = "WARNING: "
	NOTICE  = "NOTICE: "
	INFO    = "INFO: "
	DEBUG   = "DEBUG: "
)

var (
	Notice  = log.New(os.Stdout, NOTICE, 0)
	Info    = log.New(os.Stdout, INFO, 0)
	Warning = log.New(os.Stdout, WARNING, log.Lshortfile)
	Err     = log.New(os.Stderr, ERR, log.Lshortfile)
	Crit    = log.New(os.Stderr, CRIT, log.Lshortfile)
	Fatal   = log.New(os.Stderr, CRIT, log.Lshortfile)
	Debug   = log.New(os.Stdout, DEBUG, log.Lshortfile)
)

func Home() string {
	home := os.Getenv("XDG_STATE_HOME")
	if home == "" {
		home = filepath.Join(os.Getenv("HOME"), ".local", "state")
	}
	return home
}

func Dir() string {
	return filepath.Join(Home(), "tieba-sign", "log")
}

func ToFile(name string, content []byte) error {
	logFilePath := filepath.Join(Dir(), name)
	logDir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("无法创建日志目录: %s: %w", logDir, err)
	}
	if err := os.WriteFile(logFilePath, content, 0644); err != nil {
		return fmt.Errorf("无法写入日志文件: %s: %w", logFilePath, err)
	}
	return nil
}
