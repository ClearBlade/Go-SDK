package GoSDK

import "fmt"

func composeStorageURL(systemKey string) string {
	return fmt.Sprintf("/api/v/4/message/%s/storage", systemKey)
}

type MessageHistoryStorageEntry struct {
	Edge     bool   `json:"edge"`
	MaxRows  int    `json:"max_rows"`
	MaxTime  int    `json:"max_time"`
	Platform bool   `json:"platform"`
	Topic    string `json:"topic"`
}

type GetMessageHistoryStorageEntry struct {
	MessageHistoryStorageEntry
	SystemKey string `json:"system_key"`
}

func (d *DevClient) UpdateMessageHistoryStorage(systemKey string, data []MessageHistoryStorageEntry) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := put(d, composeStorageURL(systemKey), data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return err
	}
	return nil
}

func (d *DevClient) GetMessageHistoryStorage(systemKey string) ([]GetMessageHistoryStorageEntry, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, composeStorageURL(systemKey), nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}

	switch t := resp.Body.(type) {
	case []interface{}:
		rtn := make([]GetMessageHistoryStorageEntry, 0)
		for i := range t {
			if m, ok := t[i].(map[string]interface{}); ok {
				rtn = append(rtn, GetMessageHistoryStorageEntry{
					MessageHistoryStorageEntry: MessageHistoryStorageEntry{
						Edge:     m["edge"].(bool),
						Platform: m["platform"].(bool),
						MaxRows:  int(m["max_rows"].(float64)),
						MaxTime:  int(m["max_time"].(float64)),
						Topic:    m["topic"].(string),
					},
					SystemKey: m["system_key"].(string),
				})
			} else {
				return nil, fmt.Errorf("Expected map[string]interface{}")
			}

		}
		return rtn, nil
	default:
		return nil, fmt.Errorf("Invalid type returned by GetMessageHistoryStorage. Expected []interface{} but got %s", t)
	}
}
