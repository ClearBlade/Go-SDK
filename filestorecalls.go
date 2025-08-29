package GoSDK

import (
	"fmt"
	"strconv"
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

type Filestore struct {
	Name            string         `json:"name"`
	StorageType     string         `json:"storage_type"`
	StorageConfig   map[string]any `json:"storage_config"`
	ReadAuthMethod  string         `json:"read_auth_method"`
	WriteAuthMethod string         `json:"write_auth_method"`
	SystemKey       string         `json:"system_key"`
}

type EncryptedFilestore struct {
	Name            string `json:"name"`
	StorageType     string `json:"storage_type"`
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
	body, err := makeFilestoreRequest(c, systemKey, filestoreName, false)
	if err != nil {
		return nil, err
	}
	filestore := &EncryptedFilestore{}
	err = decodeMapToStruct(body, filestore)
	return filestore, err
}

func makeFilestoreRequest(c cbClient, systemKey, filestoreName string, decrypt bool) (any, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("%s%s/%s?decrypt=%t", _FILESTORES_PREAMBLE, systemKey, filestoreName, decrypt)
	resp, err := get(c, endpoint, nil, creds, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed with status code %d: %v", resp.StatusCode, resp.Body)
	}

	return resp.Body, nil
}

func (d *DevClient) GetDecryptedFilestore(systemKey, name string) (*Filestore, error) {
	return getDecryptedFilestore(d, systemKey, name)
}

func (u *UserClient) GetDecryptedFilestore(systemKey, name string) (*Filestore, error) {
	return getDecryptedFilestore(u, systemKey, name)
}

func getDecryptedFilestore(c cbClient, systemKey, filestoreName string) (*Filestore, error) {
	body, err := makeFilestoreRequest(c, systemKey, filestoreName, true)
	if err != nil {
		return nil, err
	}
	filestore := &Filestore{}
	err = decodeMapToStruct(body, filestore)
	return filestore, err
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) GetFilestores(systemKey string) ([]*EncryptedFilestore, error) {
	return getFilestores(d, systemKey)
}

func (u *UserClient) GetFilestores(systemKey string) ([]*EncryptedFilestore, error) {
	return getFilestores(u, systemKey)
}

func getFilestores(c cbClient, systemKey string) ([]*EncryptedFilestore, error) {
	body, err := makeFilestoresRequest(c, systemKey, false)
	if err != nil {
		return nil, err
	}

	filestores := []*EncryptedFilestore{}
	err = decodeMapToStruct(body, &filestores)
	return filestores, err
}

func (d *DevClient) GetDecryptedFilestores(systemKey string) ([]*Filestore, error) {
	return getDecryptedFilestores(d, systemKey)
}

func (u *UserClient) GetDecryptedFilestores(systemKey string) ([]*Filestore, error) {
	return getDecryptedFilestores(u, systemKey)
}

func getDecryptedFilestores(c cbClient, systemKey string) ([]*Filestore, error) {
	body, err := makeFilestoresRequest(c, systemKey, true)
	if err != nil {
		return nil, err
	}

	filestores := []*Filestore{}
	err = decodeMapToStruct(body, &filestores)
	return filestores, err
}

func makeFilestoresRequest(c cbClient, systemKey string, decrypt bool) (any, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("%s%s?decrypt=%t", _FILESTORES_PREAMBLE, systemKey, decrypt)
	resp, err := get(c, endpoint, nil, creds, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed with status code %d: %v", resp.StatusCode, resp.Body)
	}

	return resp.Body, nil
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

	return []byte(bodyStr), nil
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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type ListOptions struct {
	Depth             int    `json:"depth"`
	Limit             int    `json:"limit"`
	Prefix            string `json:"prefix"`
	ContinuationToken string `json:"continuation_token"`
}

type FileList struct {
	Files             []FileMeta `json:"files"`
	ContinuationToken string     `json:"continuation_token"`
}
type FileMeta struct {
	FullPath      string `json:"full_path"`
	SizeBytes     int64  `json:"size_bytes"`
	Permissions   string `json:"permissions"`
	UpdatedAtUnix int64  `json:"updated_at"`
	IsDir         bool   `json:"is_dir"`
}

func (d *DevClient) ListFilestoreFiles(systemKey, filestore string, opts *ListOptions) (*FileList, error) {
	return listFilestoreFiles(d, systemKey, filestore, opts)
}

func (u *UserClient) ListFilestoreFiles(systemKey, filestore string, opts *ListOptions) (*FileList, error) {
	return listFilestoreFiles(u, systemKey, filestore, opts)
}

func listFilestoreFiles(c cbClient, systemKey, filestore string, opts *ListOptions) (*FileList, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s%s/%s/list/%s", _FILESTORES_PREAMBLE, systemKey, filestore, opts.Prefix)
	query := map[string]string{}
	if opts.Depth != 0 {
		query["depth"] = strconv.Itoa(opts.Depth)
	}

	if opts.Limit != 0 {
		query["limit"] = strconv.Itoa(opts.Limit)
	}

	if opts.ContinuationToken != "" {
		query["continuation_token"] = opts.ContinuationToken
	}

	resp, err := get(c, endpoint, query, creds, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed with status code %d: %v", resp.StatusCode, resp.Body)
	}

	result := &FileList{}
	if err := decodeMapToStruct(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode file list: %w", err)
	}

	return result, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) CopyFilestoreFile(systemKey, filestore, src, dst string) error {
	return copyFilestoreFile(d, systemKey, filestore, src, dst)
}

func (u *UserClient) CopyFilestoreFile(systemKey, filestore, src, dst string) error {
	return copyFilestoreFile(u, systemKey, filestore, src, dst)
}

func copyFilestoreFile(c cbClient, systemKey, filestore, src, dst string) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	qryStr := query_to_string(map[string]string{"destination": dst})
	endpoint := fmt.Sprintf("%s%s/%s/copy/%s?%s", _FILESTORES_PREAMBLE, systemKey, filestore, src, qryStr)
	resp, err := put(c, endpoint, nil, creds, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed with status code %d: %v", resp.StatusCode, resp.Body)
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) MoveFilestoreFile(systemKey, filestore, src, dst string) error {
	return moveFilestoreFile(d, systemKey, filestore, src, dst)
}

func (u *UserClient) MoveFilestoreFile(systemKey, filestore, src, dst string) error {
	return moveFilestoreFile(u, systemKey, filestore, src, dst)
}

func moveFilestoreFile(c cbClient, systemKey, filestore, src, dst string) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	qryStr := query_to_string(map[string]string{"destination": dst})
	endpoint := fmt.Sprintf("%s%s/%s/move/%s?%s", _FILESTORES_PREAMBLE, systemKey, filestore, src, qryStr)
	resp, err := put(c, endpoint, nil, creds, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed with status code %d: %v", resp.StatusCode, resp.Body)
	}

	return nil
}
