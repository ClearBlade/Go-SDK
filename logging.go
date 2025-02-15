package GoSDK

import (
	"fmt"
	mqtt "github.com/clearblade/paho.mqtt.golang"
	"strings"
)

type MQTTLogger struct{}

func (MQTTLogger) Println(v ...interface{}) {
	fmt.Println(v...)
}

func (MQTTLogger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func SetLoggers(levelStr string) {
	levels := strings.Split(levelStr, ",")
	for _, oneLevel := range levels {
		oneLevel = strings.ToLower(strings.TrimSpace(oneLevel))
		switch oneLevel {
		case "error":
			mqtt.ERROR = MQTTLogger{}
		case "critical":
			mqtt.CRITICAL = MQTTLogger{}
		case "warn", "warning":
			mqtt.WARN = MQTTLogger{}
		case "debug":
			mqtt.DEBUG = MQTTLogger{}
		case "": // just handle the case where nothing was passed in
		default:
			fmt.Printf("Ignoring bad logging level: \"%s\"\n", oneLevel)
		}
	}
}
