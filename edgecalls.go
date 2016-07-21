package GoSDK

import (
	"os"
	"os/exec"
)

const (
	_EDGES_PREAMBLE        = "/admin/edges/"
	_EDGES_USER_PREAMBLE   = "/api/v/2/edges/"
	_EDGES_SYNC_MANAGEMENT = "/admin/edges/sync/"
)

type EdgeConfig struct {
	EdgeName     string
	EdgeToken    string
	NoviIp       string
	ParentSystem string
	HttpPort     string
	MqttPort     string
	MqttTlsPort  string
	WsPort       string
	WssPort      string
	AuthPort     string
	AuthWsPort   string
	Lean         bool
	Cache        bool
	Stdout       *os.File
	Stderr       *os.File
}

func CreateNewEdge(e EdgeConfig) (*os.Process, error) {
	_, err := exec.LookPath("edge")
	if err != nil {
		println("edge not found in $PATH")
		return nil, err
	}
	cmd := parseEdgeConfig(e)
	return cmd.Process, cmd.Start()
}

func (u *UserClient) GetEdges(systemKey string) ([]interface{}, error) {
	creds, err := u.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(u, _EDGES_USER_PREAMBLE+systemKey, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) GetEdges(systemKey string) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _EDGES_PREAMBLE+systemKey, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) GetEdge(systemKey, name string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _EDGES_PREAMBLE+systemKey+"/"+name, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) CreateEdge(systemKey, name string,
	data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := post(d, _EDGES_PREAMBLE+systemKey+"/"+name, data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) DeleteEdge(systemKey, name string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := delete(d, _EDGES_PREAMBLE+systemKey+"/"+name, nil, creds, nil)
	_, err = mapResponse(resp, err)
	return err
}

func (d *DevClient) UpdateEdge(systemKey, name string, data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := put(d, _EDGES_PREAMBLE+systemKey+"/"+name, data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

type ResourceType string

const (
	ServiceSync ResourceType = "service"
	LibrarySync ResourceType = "library"
	TriggerSync ResourceType = "trigger"
)

func (d *DevClient) GetSyncResourcesForEdge(systemKey string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _EDGES_SYNC_MANAGEMENT+systemKey, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) SyncResourceToEdge(systemKey, edgeName string, add map[ResourceType][]string, remove map[ResourceType][]string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	if add == nil {
		add = map[ResourceType][]string{}
	}
	if remove == nil {
		remove = map[ResourceType][]string{}
	}
	changes := map[string][]map[string]interface{}{
		"add":    mapSyncChanges(add),
		"remove": mapSyncChanges(remove),
	}
	resp, err := put(d, _EDGES_SYNC_MANAGEMENT+systemKey+"/"+edgeName, changes, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}
