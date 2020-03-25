// Package logger provides implementetion of writing Log messages to Graylog, file, and stdout.
// Supports log levels and destination.
package logger

import (
	"errors"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	graylog "gopkg.in/gemnasium/logrus-graylog-hook.v2"
)

// ErrFailedToConfigureLog is the error returned when configuring failed for some reasons
var ErrFailedToConfigureLog = errors.New("Failed to init log: failed to configure ")

// Format config loger with nullformatter, that is onlly log to Graylog, with out ouput to file or stdout
func (NullFormatter) Format(e *logrus.Entry) ([]byte, error) {
	return []byte{}, nil
}

// NewLogger initialized logger according to configuration
func ConfigLogger(log *logrus.Logger, lc *LogConfig) error {
	var err error
	switch lc.Output {
	case "Stdout":
		lc.setLoggerToStdout(log)
	case "File":
		err = lc.setLoggerToFile(log)
		if err != nil {
			return err
		}
	case "Graylog":
		lc.setLoggerToGraylog(log)
	default:
		return ErrFailedToConfigureLog
	}
	lev, err := logrus.ParseLevel(lc.Level)
	if err != nil {
		return err
	}
	log.SetLevel(lev)
	return err
}

// setLoggerToFile initialize logger for writing to file
func (lc *LogConfig) setLoggerToFile(log *logrus.Logger) error {
	f, err := os.OpenFile("../user_manager_api.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(f)
	return err
}

// setLoggerToStdout initialize logger for writing to stdout
func (c *LogConfig) setLoggerToStdout(log *logrus.Logger) {
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

// setLoggerToGraylog initialize logger for writing to Graylog
func (lc *LogConfig) setLoggerToGraylog(log *logrus.Logger) {
	var hook logrus.Hook
	graylogAdr := fmt.Sprintf("%v:%v", lc.Host, lc.Port)
	if lc.Type == "async" {
		hook = graylog.NewAsyncGraylogHook(graylogAdr, map[string]interface{}{"API": "User management service"})
	} else {
		hook = graylog.NewGraylogHook(graylogAdr, map[string]interface{}{"API": "User management service"})
	}
	log.AddHook(hook)
	log.SetFormatter(new(NullFormatter))
}
