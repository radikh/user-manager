// Package logger provides implementetion of writing Log messages to Graylog, file, and stdout.
// Supports log levels and destination.
package logger

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"
	graylog "gopkg.in/gemnasium/logrus-graylog-hook.v2"
)

// ErrFailedToConfigureLog is the error returned when configuring failed for some reasons
var ErrFailedToConfigureLog = errors.New("Failed to init log: failed to configure ")

// StorageConfig contains fileds used in Connect for DSN
type LogConfig struct {
	Host       string
	Port       string
	PassSecret string
	PassSHA2   string
	Output     string
	Level      string
	Type       string
}
type loggerKeyType int

const loggerKey loggerKeyType = iota

// NullFormatter structure for creating null formatter for logger
type NullFormatter struct {
}

// New returns a context that has a logrus logger
func New(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerKey, WithContext(ctx))
}

// WithContext returns a logrus logger from the context
func WithContext(ctx context.Context) *log.Logger {
	if ctx == nil {
		return log.StandardLogger()
	}

	if ctxLogger, ok := ctx.Value(loggerKey).(*log.Logger); ok {
		return ctxLogger
	}

	return log.StandardLogger()
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
			return err
		}
	case "Graylog":
		lc.setLoggerToGraylog()
	default:
		err = ErrFailedToConfigureLog
		return err
	}
	lev, err := log.ParseLevel(lc.Level)
	if err != nil {
		return err
	}
	log.SetLevel(lev)
	return err
}

// setLoggerToFile initialize logger for writing to file
func (lc *LogConfig) setLoggerToFile() error {
	f, err := os.OpenFile("user_manager_api.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(f)
	return err
}

// setLoggerToStdout initialize logger for writing to stdout
func (lc *LogConfig) setLoggerToStdout() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
}

// setLoggerToGraylog initialize logger for writing to Graylog
func (lc *LogConfig) setLoggerToGraylog() {
	var hook log.Hook
	graylogAdr := fmt.Sprintf("%v:%v", lc.Host, lc.Port)
	if lc.Type == "async" {
		hook = graylog.NewAsyncGraylogHook(graylogAdr, map[string]interface{}{"API": "User management service"})
	} else {
		hook = graylog.NewGraylogHook(graylogAdr, map[string]interface{}{"API": "User management service"})
	}
	log.AddHook(hook)
	log.SetFormatter(new(NullFormatter))
}

// LoggerHandler creates a new middleware
func LoggerHandler() func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(New(r.Context()))

			logg := WithContext(r.Context())
			logg.WithFields(log.Fields{"method": r.Method, "path": r.URL.Path}).Debug("Started request")

			originalURL := &url.URL{}
			*originalURL = *r.URL

			fields := log.Fields{
				"method":      r.Method,
				"host":        r.Host,
				"request":     r.RequestURI,
				"remote-addr": r.RemoteAddr,
				"referer":     r.Referer(),
				"user-agent":  r.UserAgent(),
			}
			if originalURL.String() != r.URL.String() {
				fields["upstream-host"] = r.URL.Host
				fields["upstream-request"] = r.URL.RequestURI()
			}

			logg.WithFields(fields).Info("Completed handling request")
		})
	}
}
