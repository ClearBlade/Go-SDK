package GoSDK

import (
	"encoding/json"
	"net/url"
)

const (
	_STATS_PREAMBLE  = "/api/v/3/platform/statistics"
	_DBCONN_PREAMBLE = "/api/v/3/platform/dbconnections"
	_LOGS_PREAMBLE   = "/api/v/3/platform/logs"
)

func (u *UserClient) GetConnectedEdgeStatistics(systemKey string, query *Query) (map[string]interface{}, error) {
	creds, err := u.credentials()
	if err != nil {
		return nil, err
	}
	q, err := dbQueryToReqQuery(query)
	if err != nil {
		return nil, err
	}
	resp, err := get(u, _STATS_PREAMBLE+"/"+systemKey, q, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) GetConnectedEdgeStatistics(systemKey string, query *Query) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	q, err := dbQueryToReqQuery(query)
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _STATS_PREAMBLE+"/"+systemKey, q, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (u *UserClient) GetConnectedEdgeDBConnections(systemKey string, query *Query) (map[string]interface{}, error) {
	creds, err := u.credentials()
	if err != nil {
		return nil, err
	}
	q, err := dbQueryToReqQuery(query)
	if err != nil {
		return nil, err
	}
	resp, err := get(u, _DBCONN_PREAMBLE+"/"+systemKey, q, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) GetConnectedEdgeDBConnections(systemKey string, query *Query) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	q, err := dbQueryToReqQuery(query)
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _DBCONN_PREAMBLE+"/"+systemKey, q, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (u *UserClient) GetConnectedEdgeLogs(systemKey string, query *Query) (map[string]interface{}, error) {
	creds, err := u.credentials()
	if err != nil {
		return nil, err
	}
	q, err := dbQueryToReqQuery(query)
	if err != nil {
		return nil, err
	}
	resp, err := get(u, _LOGS_PREAMBLE+"/"+systemKey, q, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) GetConnectedEdgeLogs(systemKey string, query *Query) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	q, err := dbQueryToReqQuery(query)
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _LOGS_PREAMBLE+"/"+systemKey, q, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func dbQueryToReqQuery(query *Query) (map[string]string, error) {
	var qry map[string]string
	if query != nil {
		query_map := query.serialize()
		query_bytes, err := json.Marshal(query_map)
		if err != nil {
			return nil, err
		}
		qry = map[string]string{
			"query": url.QueryEscape(string(query_bytes)),
		}
	} else {
		qry = nil
	}
	return qry, nil
}
