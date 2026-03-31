package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	EMERG   = 0
	ALERT   = 1
	CRIT    = 2
	ERR     = 3
	WARNING = 4
	NOTICE  = 5
	INFO    = 6
	DEBUG   = 7
)

const (
	EMERG_PREFIX   = "EMERG: "
	ALERT_PREFIX   = "ALERT: "
	CRIT_PREFIX    = "CRIT: "
	ERR_PREFIX     = "ERR: "
	WARNING_PREFIX = "WARNING: "
	NOTICE_PREFIX  = "NOTICE: "
	INFO_PREFIX    = "INFO: "
	DEBUG_PREFIX   = "DEBUG: "
)

var level int = NOTICE

func SetLevel(levelStr string) {
	switch strings.ToLower(levelStr) {
	case "debug":
		level = DEBUG
	case "info":
		level = INFO
	case "notice":
		level = NOTICE
	case "warning":
		level = WARNING
	case "error", "err":
		level = ERR
	case "crit", "critical":
		level = CRIT
	case "alert":
		level = ALERT
	case "emerg":
		level = EMERG
	default:
		level = NOTICE
	}
}

func Enabled(lvl int) bool {
	return lvl <= level
}

type Logger struct {
	stdLogger *log.Logger
	prefix    string
	level     int
}

func (l *Logger) Print(v ...interface{}) {
	if !Enabled(l.level) {
		return
	}
	l.stdLogger.Output(2, fmt.Sprint(v...))
}

func (l *Logger) Printf(format string, v ...interface{}) {
	if !Enabled(l.level) {
		return
	}
	l.stdLogger.Output(2, fmt.Sprintf(format, v...))
}

func (l *Logger) Println(v ...interface{}) {
	if !Enabled(l.level) {
		return
	}
	l.stdLogger.Output(2, fmt.Sprint(v...))
}

func (l *Logger) Fatalln(v ...interface{}) {
	if Enabled(l.level) {
		l.stdLogger.Output(2, fmt.Sprint(v...))
	}
	os.Exit(1)
}

func NewLogger(prefix string, lvl int, flag int) *Logger {
	var output *os.File
	if lvl <= ERR {
		output = os.Stderr
	} else {
		output = os.Stdout
	}
	return &Logger{
		stdLogger: log.New(output, prefix, flag),
		prefix:    prefix,
		level:     lvl,
	}
}

var (
	Notice  = NewLogger(NOTICE_PREFIX, NOTICE, 0)
	Info    = NewLogger(INFO_PREFIX, INFO, 0)
	Warning = NewLogger(WARNING_PREFIX, WARNING, log.Lshortfile)
	Err     = NewLogger(ERR_PREFIX, ERR, log.Lshortfile)
	Crit    = NewLogger(CRIT_PREFIX, CRIT, log.Lshortfile)
	Fatal   = NewLogger(CRIT_PREFIX, CRIT, log.Lshortfile)
	Debug   = NewLogger(DEBUG_PREFIX, DEBUG, log.Lshortfile)
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
