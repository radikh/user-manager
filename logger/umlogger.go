// Package logger provides implementetion of writing Log messages to Graylog, file, and stdout.
// Supports log levels and destination.
package logger

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	graylog "gopkg.in/gemnasium/logrus-graylog-hook.v2"
)

// StorageConfig contains fileds used in Connect for DSN
type LogConfig struct {
	Host        string
	Port        string
	Pass_Secret string
	Pass_SHA2   string
	Output      string
}

// NullFormatter structure for creating null formatter for logger
type NullFormatter struct {
}

// Format config loger with nullformatter, that is onlly log to Graylog, with out ouput to file or stdout
func (NullFormatter) Format(e *log.Entry) ([]byte, error) {
	return []byte{}, nil
}

// NewLogger initialized logger according to configuration
func NewLogger(lc *LogConfig) error {
	switch lc.Output {
	case "Stdout":
		log.SetFormatter(&log.TextFormatter{})
		log.SetOutput(os.Stdout)
	case "File":
		f, err := os.OpenFile("user_manager_api.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			fmt.Printf("error opening file: %v", err)
		}
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(f)
	case "Graylog":
		graylog_adr := fmt.Sprintf("%v:%v", lc.Host, lc.Port)
		hook := graylog.NewGraylogHook(graylog_adr, map[string]interface{}{"API": "User management service"})
		log.AddHook(hook)
		log.SetFormatter(new(NullFormatter))
	default:
		return fmt.Errorf("Error logger configure output destination <%v> < should be Graylog, Stdout or File", lc.Output)
	}
	log.SetLevel(log.PanicLevel)
	return nil
}

// Message log message depending to log level
func Message(level log.Level, m string) {
	switch level {
	case log.DebugLevel:
		log.Debug(m)
	case log.InfoLevel:
		log.Info(m)
	case log.WarnLevel:
		log.Warn(m)
	case log.ErrorLevel:
		log.Error(m)
	case log.FatalLevel:
		log.Fatal(m)
	case log.PanicLevel:
		log.Panic(m)
	}
}
