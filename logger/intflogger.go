// Package logger provides implementation of writing Log messages to Graylog, file, and stdout.
// Supports log levels and destination.
package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Log is a package level variable that access logging function through "Log"
var LogUM Logger

// NullFormatter structure for creating null formatter for logger
type NullFormatter struct {
}

// StorageConfig contains fields used in Connect for DSN
type LogConfig struct {
	Host       string
	Port       int
	PassSecret string
	PassSHA2   string
	Output     string
	Level      string
	Type       string
}

// Logger represent interface for logging function
type Logger interface {
	Panicf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Error(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatal(args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Warnf(format string, args ...interface{})
	Warn(args ...interface{})
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
}

// SetLogger is the setter for Log variable
func SetLogger(cfg *LogConfig) error {
	log := logrus.New()

	err := configLogger(log, cfg)
	if err != nil {
		return fmt.Errorf("setlogger error: %w", err)
	}

	LogUM = log

	return nil
}
