package GoSDK

import (
	"fmt"
)

const (
	_MQTT_CONNECTOR_PREAMBLE = "/api/v/1/mqtt-connectors"
)

type MqttConnectorUnencrypted struct {
	SystemKey   string                 `json:"system_key"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Credentials map[string]interface{} `json:"credentials"`
	Config      map[string]interface{} `json:"config"`
}

type MqttConnectorEncrypted struct {
	SystemKey   string                 `json:"system_key"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Credentials string                 `json:"credentials"`
	Config      map[string]interface{} `json:"config"`
}

func (d *DevClient) CreateMqttConnector(connector *MqttConnectorUnencrypted) (*MqttConnectorEncrypted, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/%s/%s", _MQTT_CONNECTOR_PREAMBLE, connector.SystemKey, connector.Name)
	resp, err := post(d, url, connector, creds, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad status code in response %v: %d", resp.Body, resp.StatusCode)
	}

	var result *MqttConnectorEncrypted
	if err := decodeMapToStruct(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("could not decode response: %w", err)
	}

	return result, nil
}
