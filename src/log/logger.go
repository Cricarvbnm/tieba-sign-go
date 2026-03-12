package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	Info  = log.New(os.Stdout, "", 0)
	Warn  = log.New(os.Stdout, "WARN: ", log.Lshortfile)
	Error = log.New(os.Stderr, "ERROR: ", log.Lshortfile)
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
