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
	Host       string
	Port       string
	PassSecret string
	PassSHA2   string
	Output     string
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
	var err error
	switch lc.Output {
	case "Stdout":
		lc.setLoggerToStdout()
	case "File":
		err = lc.setLoggerToFile()
		if err != nil {
			fmt.Printf("error opening file: %v", err)
		}
	case "Graylog":
		err = lc.setLoggerToGraylog()
		if err != nil {
			fmt.Printf("error assigning logger to Graylog: %v", err)
		}
	default:
		err = fmt.Errorf("Error logger configure output destination <%v> < should be Graylog, Stdout or File", lc.Output)
	}
	log.SetLevel(log.PanicLevel)
	return err
}

func (lc *LogConfig) setLoggerToFile() error {
	f, err := os.OpenFile("user_manager_api.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(f)
	return err
}

func (lc *LogConfig) setLoggerToStdout() error {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	return nil
}

func (lc *LogConfig) setLoggerToGraylog() error {
	graylogAdr := fmt.Sprintf("%v:%v", lc.Host, lc.Port)
	hook := graylog.NewGraylogHook(graylogAdr, map[string]interface{}{"API": "User management service"})
	log.AddHook(hook)
	log.SetFormatter(new(NullFormatter))
	return nil
}
