package GoSDK

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	_StatsPreamble  = "/api/v/3/platform/statistics"
	_DBConnPreamble = "/api/v/3/platform/dbconnections"
	_LogsPreamble   = "/api/v/3/platform/logs"
)

func (u *UserClient) GetPlatformStatistics(systemKey string, query *Query) (map[string]interface{}, error) {
	return getStatistics(u, systemKey, query, nil)
}

func (d *DevClient) GetPlatformStatistics(systemKey string, query *Query) (map[string]interface{}, error) {
	return getStatistics(d, systemKey, query, nil)
}

func (u *UserClient) GetConnectedEdgeStatistics(edgeName, systemKey string, query *Query) (map[string]interface{}, error) {
	return getStatistics(u, systemKey, query, map[string][]string{"ClearBlade-Edge": []string{edgeName}})
}

func (d *DevClient) GetConnectedEdgeStatistics(edgeName, systemKey string, query *Query) (map[string]interface{}, error) {
	return getStatistics(d, systemKey, query, map[string][]string{"ClearBlade-Edge": []string{edgeName}})
}

func getStatistics(c cbClient, systemKey string, query *Query, headers map[string][]string) (map[string]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	q, err := dbQueryToReqQuery(query)
	if err != nil {
		return nil, err
	}
	resp, err := get(c, _StatsPreamble+"`"+systemKey, q, creds, headers)
	if err != nil {
		return nil, fmt.Errorf("Unable to communicate with platform, got %s", err)
	}
	mapResp, err := mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	mapBody, ok := mapResp.Body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Expected JSON response, got %#v", resp)
	}
	return mapBody, nil
}

func (u *UserClient) GetPlatformDBConnections(systemKey string, query *Query) (map[string]interface{}, error) {
	return getDBConnections(u, systemKey, query, nil)
}

func (d *DevClient) GetPlatformDBConnections(systemKey string, query *Query) (map[string]interface{}, error) {
	return getDBConnections(d, systemKey, query, nil)
}

func (u *UserClient) GetConnectedEdgeDBConnections(edgeName, systemKey string, query *Query) (map[string]interface{}, error) {
	return getDBConnections(u, systemKey, query, map[string][]string{"ClearBlade-Edge": []string{edgeName}})
}

func (d *DevClient) GetConnectedEdgeDBConnections(edgeName, systemKey string, query *Query) (map[string]interface{}, error) {
	return getDBConnections(d, systemKey, query, map[string][]string{"ClearBlade-Edge": []string{edgeName}})
}

func getDBConnections(c cbClient, systemKey string, query *Query, headers map[string][]string) (map[string]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	q, err := dbQueryToReqQuery(query)
	if err != nil {
		return nil, err
	}
	resp, err := get(c, _DBConnPreamble+"/"+systemKey, q, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to communicate with platform, got %s", err)
	}
	respMap, err := mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	respBody, ok := respMap.Body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Expected JSON response, got %#v", resp)
	}
	return respBody, nil
}

func (u *UserClient) GetPlatformLogs(systemKey string, query *Query) (map[string]interface{}, error) {
	return getLogs(u, systemKey, query, nil)
}

func (d *DevClient) GetPlatformLogs(systemKey string, query *Query) (map[string]interface{}, error) {
	return getLogs(d, systemKey, query, nil)
}

func (u *UserClient) GetConnectedEdgeLogs(edgeName, systemKey string, query *Query) (map[string]interface{}, error) {
	return getLogs(u, systemKey, query, map[string][]string{"ClearBlade-Edge": []string{edgeName}})
}

func (d *DevClient) GetConnectedEdgeLogs(edgeName, systemKey string, query *Query) (map[string]interface{}, error) {
	return getLogs(d, systemKey, query, map[string][]string{"ClearBlade-Edge": []string{edgeName}})
}

func getLogs(c cbClient, systemKey string, query *Query, headers map[string][]string) (map[string]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	q, err := dbQueryToReqQuery(query)
	if err != nil {
		return nil, err
	}
	resp, err := get(c, _LogsPreamble+"/"+systemKey, q, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to communicate with platform, got %s", err)
	}
	respMap, err := mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	respBody, ok := respMap.Body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Expected JSON response, got %#v", resp)
	}
	return respBody, nil
}

func dbQueryToReqQuery(query *Query) (map[string]string, error) {
	var qry map[string]string
	if query != nil {
		queryMap := query.serialize()
		queryBytes, err := json.Marshal(queryMap)
		if err != nil {
			return nil, err
		}
		qry = map[string]string{
			"query": url.QueryEscape(string(queryBytes)),
		}
	} else {
		qry = nil
	}
	return qry, nil
}
