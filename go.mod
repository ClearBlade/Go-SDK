module github.com/clearblade/Go-SDK

go 1.24.0

require (
	github.com/clearblade/go-utils v1.1.5-0.20240513160427-a20563b372a5
	github.com/clearblade/mqtt_parsing v0.0.0-20160301165118-6ae49eac0961
	github.com/eclipse/paho.mqtt.golang v1.5.1
	github.com/mitchellh/mapstructure v1.5.0
	github.com/pkg/errors v0.8.1
)

replace github.com/eclipse/paho.mqtt.golang => github.com/mattwalo32/paho.mqtt.golang v1.5.1-jitter

require (
	github.com/fatih/structs v1.1.0
	golang.org/x/net v0.48.0
)

require (
	github.com/clearblade/paho.mqtt.golang v1.1.1-0.20260421160730-b6329df0f71f // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/sync v0.19.0 // indirect
)
