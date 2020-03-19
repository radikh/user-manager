// Package logger provides implementetion of writing Log messages to Graylog, file, and stdout.
// Supports log levels and destination.
package logger

import (
	"os"
	"fmt"
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


type NullFormatter struct {
}

// Don't spend time formatting logs
func (NullFormatter) Format(e *log.Entry) ([]byte, error) {
    return []byte{}, nil
    }
}

func NewLogger(lc *LogConfig) {
	switch lc.Output {
	case "Stdout":
		// Log as default ASCII formatter.
		log.SetFormatter(&log.TextFormatter{})
		log.SetOutput(os.Stdout)
	case "File":
		   // open a file
		   f, err := os.OpenFile("user_manager_api.log", os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
		   if err != nil {
			   fmt.Printf("error opening file: %v", err)
		   }
		   // Log as JSON formatter.
		   log.SetFormatter(&log.JSONFormatter{})
		   log.SetOutput(f)	
	default:
		graylog_adr := fmt.Sprinf("%v:%v", lc.Host, lc.Port)
		hook := graylog.NewGraylogHook(graylog_adr, map[string]interface{}{"API": "User management service"})
		log.Hooks.Add(hook)
		log.SetFormatter(new(NullFormatter)) 
	}
 log.SetLevel(log.PanicLevel)	   
}

func Message(m *logrus.Message) {
	log.
}
