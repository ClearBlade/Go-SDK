package GoSDK

import (
	"fmt"
	mqtt "github.com/clearblade/paho.mqtt.golang"
	"strings"
)

type MQTTLogger struct {
	level string
}

func (l MQTTLogger) Println(v ...interface{}) {
	v = []interface{}{l.level, v}
	fmt.Println(v...)
}

func (l MQTTLogger) Printf(format string, v ...interface{}) {
	fmt.Printf(l.level+": "+format, v...)
}

func SetPahoLoggers(levelStr string) {
	levels := strings.Split(levelStr, ",")
	for _, oneLevel := range levels {
		oneLevel = strings.ToLower(strings.TrimSpace(oneLevel))
		switch oneLevel {
		case "critical":
			mqtt.CRITICAL = MQTTLogger{level: "CRITICAL"}
		case "error":
			mqtt.ERROR = MQTTLogger{level: "ERROR"}
		case "warn", "warning":
			mqtt.WARN = MQTTLogger{level: "WARN"}
		case "debug":
			mqtt.DEBUG = MQTTLogger{level: "DEBUG"}
		case "": // just handle the case where nothing was passed in
		default:
			fmt.Printf("Ignoring bad logging level: \"%s\"\n", oneLevel)
		}
	}
}
