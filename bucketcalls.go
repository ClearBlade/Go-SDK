package GoSDK

import (
//"fmt"
)

const (
	_BUCKET_SETS_PREAMBLE = "/api/v/4/bucket_sets/"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) GetBucketSets(systemKey string) ([]interface{}, error) {
	return getBucketSets(d, _BUCKET_SETS_PREAMBLE+systemKey)
}

func (u *UserClient) GetBucketSets(systemKey string) ([]interface{}, error) {
	return getBucketSets(u, _BUCKET_SETS_PREAMBLE+systemKey)
}

func (dv *DeviceClient) GetBucketSets(systemKey string) ([]interface{}, error) {
	return getBucketSets(dv, _BUCKET_SETS_PREAMBLE+systemKey)
}

func getBucketSets(c cbClient, endpoint string) ([]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(c, endpoint, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.([]interface{}), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) GetBucketSet(systemKey, deploymentName string) (map[string]interface{}, error) {
	return getBucketSet(d, _BUCKET_SETS_PREAMBLE+systemKey+"/"+deploymentName)
}

func (u *UserClient) GetBucketSet(systemKey, deploymentName string) (map[string]interface{}, error) {
	return getBucketSet(u, _BUCKET_SETS_PREAMBLE+systemKey+"/"+deploymentName)
}

func (dv *DeviceClient) GetBucketSet(systemKey, deploymentName string) (map[string]interface{}, error) {
	return getBucketSet(dv, _BUCKET_SETS_PREAMBLE+systemKey+"/"+deploymentName)
}

func getBucketSet(c cbClient, endpoint string) (map[string]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(c, endpoint, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) GetBucketSetFiles(systemKey, deploymentName, box string) (map[string]interface{}, error) {
	return getBucketSetFiles(d, _BUCKET_SETS_PREAMBLE+systemKey+"/"+deploymentName+"/files", box)
}

func (u *UserClient) GetBucketSetFiles(systemKey, deploymentName, box string) (map[string]interface{}, error) {
	return getBucketSetFiles(u, _BUCKET_SETS_PREAMBLE+systemKey+"/"+deploymentName+"/files", box)
}

func (dv *DeviceClient) GetBucketSetFiles(systemKey, deploymentName, box string) (map[string]interface{}, error) {
	return getBucketSetFiles(dv, _BUCKET_SETS_PREAMBLE+systemKey+"/"+deploymentName+"/files", box)
}

func getBucketSetFiles(c cbClient, endpoint, box string) (map[string]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	qs := map[string]string{"box": box}
	resp, err := get(c, endpoint, qs, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) CreateBucketSetFile(systemKey, deploymentName, box, relPath, contents string) (map[string]interface{}, error) {
	return createBucketSetFile(d, _BUCKET_SETS_PREAMBLE+systemKey+"/"+deploymentName+"/file/create", box, relPath, contents)
}

func (u *UserClient) CreateBucketSetFile(systemKey, deploymentName, box, relPath, contents string) (map[string]interface{}, error) {
	return createBucketSetFile(u, _BUCKET_SETS_PREAMBLE+systemKey+"/"+deploymentName+"/file/create", box, relPath, contents)
}

func (dv *DeviceClient) CreateBucketSetFile(systemKey, deploymentName, box, relPath, contents string) (map[string]interface{}, error) {
	return createBucketSetFile(dv, _BUCKET_SETS_PREAMBLE+systemKey+"/"+deploymentName+"/file/create", box, relPath, contents)
}

func createBucketSetFile(c cbClient, endpoint, box, relPath, contents string) (map[string]interface{}, error) {

	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"box":      box,
		"path":     relPath,
		"contents": contents,
	}

	resp, err := post(c, endpoint, data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) GetBucketSetFile(systemKey, deploymentName, box, relPath string) (map[string]interface{}, error) {
	return getBucketSetFile(d, _BUCKET_SETS_PREAMBLE+systemKey+"/"+deploymentName+"/file/meta", box, relPath)
}

func (u *UserClient) GetBucketSetFile(systemKey, deploymentName, box, relPath string) (map[string]interface{}, error) {
	return getBucketSetFile(u, _BUCKET_SETS_PREAMBLE+systemKey+"/"+deploymentName+"/file/meta", box, relPath)
}

func (dv *DeviceClient) GetBucketSetFile(systemKey, deploymentName, box, relPath string) (map[string]interface{}, error) {
	return getBucketSetFile(dv, _BUCKET_SETS_PREAMBLE+systemKey+"/"+deploymentName+"/file/meta", box, relPath)
}

func getBucketSetFile(c cbClient, endpoint, box, relPath string) (map[string]interface{}, error) {
	queries := map[string]string{"box": box, "path": relPath}
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(c, endpoint, queries, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) ReadBucketSetFile(systemKey, deploymentName, box, relPath string) (string, error) {
	return readBucketSetFile(d, _BUCKET_SETS_PREAMBLE+systemKey+"/"+deploymentName+"/file/read", box, relPath)
}

func (u *UserClient) ReadGetBucketSetFile(systemKey, deploymentName, box, relPath string) (string, error) {
	return readBucketSetFile(u, _BUCKET_SETS_PREAMBLE+systemKey+"/"+deploymentName+"/file/read", box, relPath)
}

func (dv *DeviceClient) ReadGetBucketSetFile(systemKey, deploymentName, box, relPath string) (string, error) {
	return readBucketSetFile(dv, _BUCKET_SETS_PREAMBLE+systemKey+"/"+deploymentName+"/file/read", box, relPath)
}

func readBucketSetFile(c cbClient, endpoint, box, relPath string) (string, error) {
	queries := map[string]string{"box": box, "path": relPath}
	creds, err := c.credentials()
	if err != nil {
		return "", err
	}
	resp, err := get(c, endpoint, queries, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return "", err
	}
	return resp.Body.(string), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) CopyBucketSetFile(systemKey, deploymentName, fromBox, fromPath, toBox, toPath string) error {
	return copyBucketSetFile(d, _BUCKET_SETS_PREAMBLE+systemKey+"/"+deploymentName+"/file/copy", fromBox, fromPath, toBox, toPath)
}

func (u *UserClient) CopyBucketSetFile(systemKey, deploymentName, fromBox, fromPath, toBox, toPath string) error {
	return copyBucketSetFile(u, _BUCKET_SETS_PREAMBLE+systemKey+"/"+deploymentName+"/file/copy", fromBox, fromPath, toBox, toPath)
}

func (dv *DeviceClient) CopyBucketSetFile(systemKey, deploymentName, fromBox, fromPath, toBox, toPath string) error {
	return copyBucketSetFile(dv, _BUCKET_SETS_PREAMBLE+systemKey+"/"+deploymentName+"/file/copy", fromBox, fromPath, toBox, toPath)
}

func copyBucketSetFile(c cbClient, endpoint, fromBox, fromPath, toBox, toPath string) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"from_box":  fromBox,
		"from_path": fromPath,
		"to_box":    toBox,
		"to_path":   toPath,
	}
	_, err = post(c, endpoint, data, creds, nil)
	return err
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) MoveBucketSetFile(systemKey, deploymentName, fromBox, fromPath, toBox, toPath string) error {
	return moveBucketSetFile(d, _BUCKET_SETS_PREAMBLE+systemKey+"/"+deploymentName+"/file/move", fromBox, fromPath, toBox, toPath)
}

func (u *UserClient) MoveBucketSetFile(systemKey, deploymentName, fromBox, fromPath, toBox, toPath string) error {
	return moveBucketSetFile(u, _BUCKET_SETS_PREAMBLE+systemKey+"/"+deploymentName+"/file/move", fromBox, fromPath, toBox, toPath)
}

func (dv *DeviceClient) MoveBucketSetFile(systemKey, deploymentName, fromBox, fromPath, toBox, toPath string) error {
	return moveBucketSetFile(dv, _BUCKET_SETS_PREAMBLE+systemKey+"/"+deploymentName+"/file/move", fromBox, fromPath, toBox, toPath)
}

func moveBucketSetFile(c cbClient, endpoint, fromBox, fromPath, toBox, toPath string) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"from_box":  fromBox,
		"from_path": fromPath,
		"to_box":    toBox,
		"to_path":   toPath,
	}
	_, err = post(c, endpoint, data, creds, nil)
	return err
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) DeleteBucketSetFile(systemKey, deploymentName, box, relPath string) error {
	return deleteBucketSetFile(d, _BUCKET_SETS_PREAMBLE+systemKey+"/"+deploymentName+"/file/delete", box, relPath)
}

func (u *UserClient) DeleteBucketSetFile(systemKey, deploymentName, box, relPath string) error {
	return deleteBucketSetFile(u, _BUCKET_SETS_PREAMBLE+systemKey+"/"+deploymentName+"/file/delete", box, relPath)
}

func (dv *DeviceClient) DeleteBucketSetFile(systemKey, deploymentName, box, relPath string) error {
	return deleteBucketSetFile(dv, _BUCKET_SETS_PREAMBLE+systemKey+"/"+deploymentName+"/file/delete", box, relPath)
}

func deleteBucketSetFile(c cbClient, endpoint, box, relPath string) error {

	creds, err := c.credentials()
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"box":  box,
		"path": relPath,
	}

	_, err = post(c, endpoint, data, creds, nil)
	return err
}
