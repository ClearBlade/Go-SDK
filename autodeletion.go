package GoSDK

import (
	"fmt"
	"net/http"
	"path"

	"github.com/pkg/errors"
)

const (
	autodeletion_preamble = "api/v/4/message"
)

type AutodeletionSettings struct {
	Enabled       bool  `json:"enabled"`
	MaxSizeKb     int   `json:"expiration_age_seconds"`
	MaxRows       int   `json:"max_size_kb"`
	MaxAgeSeconds int64 `json:"max_rows"`
}

func unpackMapToAutodeletionSettings(m interface{}) (*AutodeletionSettings, error) {
	switch m := m.(type) {
	case map[string]interface{}:
		return &AutodeletionSettings{
			Enabled:       m["enabled"].(bool),
			MaxSizeKb:     m["max_size_kb"].(int),
			MaxRows:       m["max_rows"].(int),
			MaxAgeSeconds: m["expiration_age_seconds"].(int64),
		}, nil
	default:
		return nil, fmt.Errorf("Unexpected response from Autodeletion Settings Endpoint: %+v", m)
	}

}

func GetAutodeletionSettings(c cbClient, systemkey string) (*AutodeletionSettings, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, errors.Wrap(err, "Error with credentials")
	}

	resp, err := get(c, path.Join(autodeletion_preamble, systemkey, "autodeletion"), nil, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting autodeletion settings: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting autodeletion settings: %s", resp.Body)
	}

	fmt.Println(resp.Body)

	return unpackMapToAutodeletionSettings(resp.Body)
}

func SetAutodeletionSettings(c cbClient, systemkey string, newSettings AutodeletionSettings) (*AutodeletionSettings, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, errors.Wrap(err, "Error with credentials")
	}

	resp, err := post(c, path.Join(autodeletion_preamble, systemkey, "autodeletion"), newSettings, creds, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error setting autodeletion settings")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error setting autodeletion settings: %d: %s", resp.StatusCode, resp.Body)
	}
	return unpackMapToAutodeletionSettings(resp.Body)
}

func SetAllAutodeletionSettings(c cbClient, newSettings AutodeletionSettings) (systems []string, err error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, errors.Wrap(err, "Error with credentials")
	}

	resp, err := post(c, path.Join(autodeletion_preamble, "bulksetautodelete"), newSettings, creds, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error Setting All Autodeletion Settings")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error Setting All Autodeletion Settings: %d: %s", resp.StatusCode, resp.Body)
	}
	switch body := resp.Body.(type) {
	case []interface{}:
		for _, sys := range body {
			systems = append(systems, sys.(string))
		}
	default:
		return nil, fmt.Errorf("Unexpected response from Autodeletion Settings Endpoint: %+v", body)
	}

	return systems, nil
}
