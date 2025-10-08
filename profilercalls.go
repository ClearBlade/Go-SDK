package GoSDK

import (
	"fmt"
)

func (d *DevClient) GetRunningProfilers(systemKey, serviceId string) (map[string]any, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("%s/profiles/%s/%s", _CODE_ADMIN_PREAMBLE_V4, systemKey, serviceId)
	resp, err := get(d, endpoint, nil, creds, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed with status code %d: %v", resp.StatusCode, resp.Body)
	}

	profiles, ok := resp.Body.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected map[string]any, got %T", resp.Body)
	}

	return profiles, nil
}

func (d *DevClient) StartProfiler(systemKey, serviceId, profileType string) (any, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("%s/profiles/%s/%s?type=%s", _CODE_ADMIN_PREAMBLE_V4, systemKey, serviceId, profileType)
	resp, err := get(d, endpoint, nil, creds, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed with status code %d: %v", resp.StatusCode, resp.Body)
	}

	return resp.Body, nil
}

func (d *DevClient) StopProfiler(systemKey, serviceId, profileId string) (string, error) {
	creds, err := d.credentials()
	if err != nil {
		return "", err
	}

	req := &CbReq{
		Method:   "DELETE",
		Endpoint: fmt.Sprintf("%s/profiles/%s/%s/%s", _CODE_ADMIN_PREAMBLE_V4, systemKey, serviceId, profileId),
		NoDecode: true,
	}

	resp, err := do(d, req, creds)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("failed with status code %d: %v", resp.StatusCode, resp.Body)
	}

	profile, ok := resp.Body.(string)
	if !ok {
		return "", fmt.Errorf("expected string, got %T", resp.Body)
	}

	return profile, nil
}

func (d *DevClient) TakeHeapSnapshot(systemKey, serviceId string) (string, error) {
	creds, err := d.credentials()
	if err != nil {
		return "", err
	}

	req := &CbReq{
		Method:   "GET",
		Endpoint: fmt.Sprintf("%s/heapsnapshot/%s/%s/", _CODE_ADMIN_PREAMBLE_V4, systemKey, serviceId),
		NoDecode: true,
	}

	resp, err := do(d, req, creds)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("failed with status code %d: %v", resp.StatusCode, resp.Body)
	}

	profile, ok := resp.Body.(string)
	if !ok {
		return "", fmt.Errorf("expected string, got %T", resp.Body)
	}

	return profile, nil
}
