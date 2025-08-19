package GoSDK

import (
	"encoding/base64"
	"fmt"
)

const (
	_FILESTORES_PREAMBLE = "/api/v/1/filestore/"
)

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type FilestoreConfig struct {
	Name            string         `json:"name"`
	StorageType     string         `json:"storage_type"`
	StorageConfig   map[string]any `json:"storage_config"`
	ReadAuthMethod  string         `json:"read_auth_method"`
	WriteAuthMethod string         `json:"write_auth_method"`
}

func (d *DevClient) CreateFilestore(systemKey string, config *FilestoreConfig) error {
	return createFilestore(d, systemKey, config)
}

func (u *UserClient) CreateFilestore(systemKey string, config *FilestoreConfig) error {
	return createFilestore(u, systemKey, config)
}

func createFilestore(c cbClient, systemKey string, config *FilestoreConfig) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	endpoint := _FILESTORES_PREAMBLE + systemKey
	body := structToMap(config)
	resp, err := post(c, endpoint, body, creds, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed with status code %d: %v", resp.StatusCode, resp.Body)
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type EncryptedFilestore struct {
	Name            string `json:"name"`
	StorageType     string `json:"storage_type"`
	StorageConfig   string `json:"storage_config"`
	ReadAuthMethod  string `json:"read_auth_method"`
	WriteAuthMethod string `json:"write_auth_method"`
}

func (d *DevClient) GetFilestore(systemKey, name string) (*EncryptedFilestore, error) {
	return getFilestore(d, systemKey, name)
}

func (u *UserClient) GetFilestore(systemKey, name string) (*EncryptedFilestore, error) {
	return getFilestore(u, systemKey, name)
}

func getFilestore(c cbClient, systemKey, filestoreName string) (*EncryptedFilestore, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	endpoint := _FILESTORES_PREAMBLE + systemKey
	resp, err := get(c, endpoint, map[string]string{"name": filestoreName}, creds, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed with status code %d: %v", resp.StatusCode, resp.Body)
	}

	filestore := &EncryptedFilestore{}
	err = decodeMapToStruct(resp.Body, filestore)
	return filestore, err
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) DeleteFilestore(systemKey, name string) error {
	return deleteFilestore(d, systemKey, name)
}

func (u *UserClient) DeleteFilestore(systemKey, name string) error {
	return deleteFilestore(u, systemKey, name)
}

func deleteFilestore(c cbClient, systemKey, filestoreName string) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	endpoint := _FILESTORES_PREAMBLE + systemKey
	_, err = delete(c, endpoint, map[string]string{"name": filestoreName}, creds, nil)
	if err != nil {
		return err
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) ReadFilestoreFile(systemKey, filestore, path string) ([]byte, error) {
	return readFilestoreFile(d, systemKey, filestore, path)
}

func (u *UserClient) ReadFilestoreFile(systemKey, filestore, path string) ([]byte, error) {
	return readFilestoreFile(u, systemKey, filestore, path)
}

func readFilestoreFile(c cbClient, systemKey, filestore, path string) ([]byte, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("%s%s/%s/file/%s", _FILESTORES_PREAMBLE, systemKey, filestore, path)
	resp, err := get(c, endpoint, nil, creds, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed with status code %d: %v", resp.StatusCode, resp.Body)
	}

	bodyStr, ok := resp.Body.(string)
	if !ok {
		return nil, fmt.Errorf("expected response body to be string, got: %T", resp.Body)
	}

	decoded, err := base64.StdEncoding.DecodeString(bodyStr)
	if err != nil {
		return nil, fmt.Errorf("could not decode response: %w", err)
	}

	return []byte(decoded), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) WriteFilestoreFile(systemKey, filestore, path string, contents []byte) error {
	return writeFilestoreFile(d, systemKey, filestore, path, contents)
}

func (u *UserClient) WriteFilestoreFile(systemKey, filestore, path string, contents []byte) error {
	return writeFilestoreFile(u, systemKey, filestore, path, contents)
}

func writeFilestoreFile(c cbClient, systemKey, filestore, path string, contents []byte) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	endpoint := fmt.Sprintf("%s%s/%s/file/%s", _FILESTORES_PREAMBLE, systemKey, filestore, path)
	resp, err := put(c, endpoint, contents, creds, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed with status code %d: %v", resp.StatusCode, resp.Body)
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) DeleteFilestoreFile(systemKey, filestore, path string) error {
	return deleteFilestoreFile(d, systemKey, filestore, path)
}

func (u *UserClient) DeleteFilestoreFile(systemKey, filestore, path string) error {
	return deleteFilestoreFile(u, systemKey, filestore, path)
}

func deleteFilestoreFile(c cbClient, systemKey, filestore, path string) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	endpoint := fmt.Sprintf("%s%s/%s/file/%s", _FILESTORES_PREAMBLE, systemKey, filestore, path)
	resp, err := delete(c, endpoint, nil, creds, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed with status code %d: %v", resp.StatusCode, resp.Body)
	}

	return nil
}
