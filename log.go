package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	WarnLogger  = log.New(log.Writer(), "WARN: ", log.Lshortfile)
	ErrorLogger = log.New(log.Writer(), "ERROR: ", log.Lshortfile)
)

func initLog() {
	log.SetFlags(0)
}

func logToFile(name string, content []byte) error {
	logFilePath := filepath.Join(getLogDir(), name)
	logDir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("无法创建日志目录: %s: %w", logDir, err)
	}
	if err := os.WriteFile(logFilePath, content, 0644); err != nil {
		return fmt.Errorf("无法写入日志文件: %s: %w", logFilePath, err)
	}
	return nil
}
