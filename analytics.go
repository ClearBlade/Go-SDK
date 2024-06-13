package GoSDK

import (
	"encoding/json"
	"net/url"
)

// TODO Implement all Analytics endpoints
const (
	_ANALYTICS_STORAGE_PREAMBLE  = "/api/v/2/analytics/storage"
)

func (d *DevClient) GetStorage(query *AnalyticsQuery) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	q, err := analyticsQueryToReqQuery(query)
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _ANALYTICS_STORAGE_PREAMBLE, q, creds, nil)
	if err != nil {
		return nil, err
	}
	mapResp, mapErr := mapResponse(resp, err)
	if mapErr != nil {
		return nil, err
	}
	return mapResp.Body.(map[string]interface{}), nil
}

func analyticsQueryToReqQuery(query *AnalyticsQuery) (map[string]string, error) {
	var q map[string]string
	if query != nil {
		query_bytes, err := json.Marshal(query.serialize())
		if err != nil {
			return nil, err
		}
		q = map[string]string{
			"query": url.QueryEscape(string(query_bytes)),
		}
	} else {
		q = nil
	}
	return q, nil
}

