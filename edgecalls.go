package GoSDK

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	_EDGES_PREAMBLE          = "/admin/edges/"
	_EDGES_USER_PREAMBLE     = "/api/v/2/edges/"
	_EDGES_SYNC_MANAGEMENT   = "/admin/edges/sync/"
	_EDGES_DEPLOY_MANAGEMENT = "/admin/edges/resources/{systemKey}/deploy"
	_EDGES_USER_V3           = "/api/v/3/edges/"
	_EDGES_V1                = "/api/v/1/edges/"
)

type EdgeConfig struct {
	EdgeName       string
	EdgeToken      string
	PlatformIP     string
	PlatformPort   string
	ParentSystem   string
	HttpPort       string
	MqttPort       string
	MqttTlsPort    string
	WsPort         string
	WssPort        string
	AuthPort       string
	AuthWsPort     string
	AdapterRootDir string
	Lean           bool
	Cache          bool
	LogLevel       string
	Insecure       bool
	DevMode        bool
	Stdout         *os.File
	Stderr         *os.File
}

func CreateNewEdgeWithCmd(e EdgeConfig) (*exec.Cmd, *os.Process, error) {
	_, err := exec.LookPath("edge")
	if err != nil {
		println("edge not found in $PATH")
		return nil, nil, err
	}
	cmd := parseEdgeConfig(e)
	return cmd, cmd.Process, cmd.Start()
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
	return u.GetEdgesWithQuery(systemKey, nil)
}

func (u *UserClient) GetEdgesWithQuery(systemKey string, query *Query) ([]interface{}, error) {
	creds, err := u.credentials()
	if err != nil {
		return nil, err
	}

	qry, err := createQueryMap(query)
	if err != nil {
		return nil, err
	}

	resp, err := get(u, _EDGES_USER_PREAMBLE+systemKey, qry, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) GetEdges(systemKey string) ([]interface{}, error) {
	return d.GetEdgesWithQuery(systemKey, nil)
}

func (d *DevClient) GetEdgesWithQuery(systemKey string, query *Query) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	qry, err := createQueryMap(query)
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _EDGES_PREAMBLE+systemKey, qry, creds, nil)
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

func (u *UserClient) GetEdge(systemKey, name string) (map[string]interface{}, error) {
	creds, err := u.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(u, _EDGES_PREAMBLE+systemKey+"/"+name, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) CreateEdge(systemKey, name string, data map[string]interface{}) (map[string]interface{}, error) {
	return createEdge(d, systemKey, _EDGES_PREAMBLE, name, data)
}

func (u *UserClient) CreateEdge(systemKey, name string, data map[string]interface{}) (map[string]interface{}, error) {
	return createEdge(u, systemKey, _EDGES_USER_V3, name, data)
}

func createEdge(client cbClient, systemKey, preamble string, name string, data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := client.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := post(client, preamble+systemKey+"/"+name, data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) DeleteEdge(systemKey, name string) error {
	return deleteEdge(d, systemKey, _EDGES_PREAMBLE, name)
}

func (u *UserClient) DeleteEdge(systemKey, name string) error {
	return deleteEdge(u, systemKey, _EDGES_USER_V3, name)
}

func deleteEdge(client cbClient, systemKey, preamble string, name string) error {
	creds, err := client.credentials()
	if err != nil {
		return err
	}
	resp, err := delete(client, preamble+systemKey+"/"+name, nil, creds, nil)
	_, err = mapResponse(resp, err)
	return err
}

func (d *DevClient) UpdateEdge(systemKey string, name string, changes map[string]interface{}) (map[string]interface{}, error) {
	return updateEdge(d, systemKey, _EDGES_PREAMBLE, name, changes)
}

func (u *UserClient) UpdateEdge(systemKey string, name string, changes map[string]interface{}) (map[string]interface{}, error) {
	return updateEdge(u, systemKey, _EDGES_USER_V3, name, changes)
}

func updateEdge(client cbClient, systemKey, preamble string, name string, data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := client.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := put(client, preamble+systemKey+"/"+name, data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

const (
	ServiceSync = "service"
	LibrarySync = "library"
	TriggerSync = "trigger"
	TimerSync   = "timer"
)

func (d *DevClient) GetDeployResourcesForSystem(systemKey string) ([]map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, strings.Replace(_EDGES_DEPLOY_MANAGEMENT, "{systemKey}", systemKey, 1), nil, creds, nil)
	if err != nil {
		return nil, err
	}
	return makeSliceOfMaps(resp.Body)
}

func (d *DevClient) serializeQuery(qIF interface{}) (string, error) {
	switch qIF.(type) {
	case string:
		return qIF.(string), nil
	case *Query:
		q := qIF.(*Query)
		qs, err := json.Marshal(q.serialize())
		if err != nil {
			return "", err
		}
		return string(qs), nil
	default:
		return "", fmt.Errorf("Bad query type: %T", qIF)
	}
}

func (d *DevClient) CreateDeployResourcesForSystem(systemKey, resourceName, resourceType string, platform bool, edgeQueryInfo interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	queryString, err := d.serializeQuery(edgeQueryInfo)
	//queryString, err := json.Marshal(edgeQuery.serialize())
	if err != nil {
		return nil, err
	}
	deploySpec := map[string]interface{}{
		"edge":                string(queryString[:]),
		"platform":            platform,
		"resource_identifier": resourceName,
		"resource_type":       resourceType,
	}
	resp, err := post(d, strings.Replace(_EDGES_DEPLOY_MANAGEMENT, "{systemKey}", systemKey, 1), deploySpec, creds, nil)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) UpdateDeployResourcesForSystem(systemKey, resourceName, resourceType string, platform bool, edgeQuery *Query) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	queryString, err := json.Marshal(edgeQuery.serialize())
	if err != nil {
		return nil, err
	}
	updatedDeploySpec := map[string]interface{}{
		"edge":                queryString,
		"platform":            platform,
		"resource_identifier": resourceName,
		"resource_type":       resourceType,
	}
	resp, err := put(d, strings.Replace(_EDGES_DEPLOY_MANAGEMENT, "{systemKey}", systemKey, 1), updatedDeploySpec, creds, nil)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) DeleteDeployResourcesForSystem(systemKey, resourceName, resourceType string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	urlString := strings.Replace(_EDGES_DEPLOY_MANAGEMENT, "{systemKey}", systemKey, 1)
	urlString += "?resource_type=" + resourceType + "&resource_identifier=" + resourceName
	_, err = put(d, urlString, nil, creds, nil)
	return err
}

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

func (d *DevClient) SyncResourceToEdge(systemKey, edgeName string, add map[string][]string, remove map[string][]string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	if add == nil {
		add = map[string][]string{}
	}
	if remove == nil {
		remove = map[string][]string{}
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

func (d *DevClient) CreateEdgeColumn(systemKey, colName, colType string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"column_name": colName,
		"type":        colType,
	}
	resp, err := post(d, _EDGES_PREAMBLE+systemKey+"/columns", data, creds, nil)
	_, err = mapResponse(resp, err)
	return err
}

func (d *DevClient) DeleteEdgeColumn(systemKey, colName string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := delete(d, _EDGES_PREAMBLE+systemKey+"/columns", map[string]string{"column": colName}, creds, nil)
	_, err = mapResponse(resp, err)
	return err
}

func (d *DevClient) GetEdgeColumns(systemKey string) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _EDGES_PREAMBLE+systemKey+"/columns", nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) GetEdgesCountWithQuery(systemKey string, query *Query) (CountResp, error) {
	return getEdgesCount(d, systemKey, _EDGES_USER_V3, query)
}

func (d *DevClient) MonitorEdges(daysInReport int) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	q := map[string]string{"report_days": strconv.Itoa(daysInReport)}
	resp, err := get(d, "/admin/edge_report", q, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func getEdgesCount(client cbClient, systemKey string, preamble string, query *Query) (CountResp, error) {
	creds, err := client.credentials()
	if err != nil {
		return CountResp{Count: 0}, err
	}

	qry, err := createQueryMap(query)
	if err != nil {
		return CountResp{Count: 0}, err
	}

	resp, err := get(client, preamble+systemKey+"/count", qry, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return CountResp{Count: 0}, err
	}
	rval, ok := resp.Body.(map[string]interface{})
	if !ok {
		return CountResp{Count: 0}, fmt.Errorf("Bad type returned by getEdgesCount: %T, %s", resp.Body, resp.Body.(string))
	}

	return CountResp{
		Count: rval["count"].(float64),
	}, nil
}

func (d *DevClient) GetEdgePublicKeys(systemKey, edgeName string) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _EDGES_PREAMBLE+"public_key/"+systemKey+"/"+edgeName, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) AddEdgePublicKey(systemKey, edgeName, publicKey, expirationTime string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	body := map[string]interface{}{
		"public_key": publicKey,
	}

	if expirationTime != "" {
		body["expiration_time"] = expirationTime
	}

	resp, err := post(d, _EDGES_PREAMBLE+"public_key/"+systemKey+"/"+edgeName, body, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) UpdateEdgePublicKey(systemKey, edgeName, publicKey, id, expirationTime string) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	body := map[string]interface{}{
		"public_key": publicKey,
	}
	if expirationTime != "" {
		body["expiration_time"] = expirationTime
	}
	resp, err := put(d, _EDGES_PREAMBLE+"public_key/"+systemKey+"/"+edgeName, map[string]interface{}{
		"id":   id,
		"$set": body,
	}, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) DeleteEdgePublicKey(systemKey, edgeName string, query *Query) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	query_map := query.serialize()
	query_bytes, err := json.Marshal(query_map)
	if err != nil {
		return nil, err
	}
	qry := map[string]string{
		"query": string(query_bytes),
	}
	_, err = delete(d, _EDGES_PREAMBLE+"public_key/"+systemKey+"/"+edgeName, qry, creds, nil)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (d *DevClient) RemoteCommandExecEdge(systemKey, edgeName string, cmd map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := post(d, _EDGES_V1+systemKey+"/"+edgeName+"/command-exec", cmd, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (u *UserClient) RemoteCommandExecEdge(systemKey, edgeName string, cmd map[string]interface{}) (map[string]interface{}, error) {
	creds, err := u.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := post(u, _EDGES_V1+systemKey+"/"+edgeName+"/command-exec", cmd, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) RemoteRestartEdge(systemKey, edgeName string) (string, error) {
	creds, err := d.credentials()
	if err != nil {
		return "", err
	}
	resp, err := post(d, _EDGES_V1+systemKey+"/"+edgeName+"/restart", nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return "", err
	}
	return resp.Body.(string), nil
}

func (u *UserClient) RemoteRestartEdge(systemKey, edgeName string) (string, error) {
	creds, err := u.credentials()
	if err != nil {
		return "", err
	}
	resp, err := post(u, _EDGES_V1+systemKey+"/"+edgeName+"/restart", nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return "", err
	}
	return resp.Body.(string), nil
}

func (d *DevClient) GetEdgeConfig(systemKey, edgeName string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _EDGES_V1+systemKey+"/"+edgeName+"/config", nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (u *UserClient) GetEdgeConfig(systemKey, edgeName string) (map[string]interface{}, error) {
	creds, err := u.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(u, _EDGES_V1+systemKey+"/"+edgeName+"/config", nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) UpdateEdgeConfig(systemKey, edgeName string, changes map[string]interface{}) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := put(d, _EDGES_V1+systemKey+"/"+edgeName+"/config", changes, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserClient) UpdateEdgeConfig(systemKey, edgeName string, changes map[string]interface{}) error {
	creds, err := u.credentials()
	if err != nil {
		return err
	}
	resp, err := put(u, _EDGES_V1+systemKey+"/"+edgeName+"/config", changes, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return err
	}
	return nil
}

func (d *DevClient) RemoteDBWipeEdge(systemKey, edgeName string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := post(d, _EDGES_V1+systemKey+"/"+edgeName+"/db-wipe", nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserClient) RemoteDBWipeEdge(systemKey, edgeName string) error {
	creds, err := u.credentials()
	if err != nil {
		return err
	}
	resp, err := post(u, _EDGES_V1+systemKey+"/"+edgeName+"/db-wipe", nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return err
	}
	return nil
}
