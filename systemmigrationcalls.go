package GoSDK

import "fmt"

const (
	_SYSTEM_MIGRATION_DEV_PREAMBLE = "/api/v/1/systemmigration/"
)

type SystemUploadDryRun struct {
	ServicesToCreate            []string              `json:"services_to_create"`
	ServicesToUpdate            []string              `json:"services_to_update"`
	LibrariesToCreate           []string              `json:"libraries_to_create"`
	LibrariesToUpdate           []string              `json:"libraries_to_update"`
	RolesToCreate               []string              `json:"roles_to_create"`
	RolesToUpdate               []string              `json:"roles_to_update"`
	UsersToCreate               []string              `json:"users_to_create"`
	UsersToUpdate               []string              `json:"users_to_update"`
	UserColumnsToAdd            []string              `json:"user_columns_to_add"`
	UserColumnsToDelete         []string              `json:"user_columns_to_delete"`
	DevicesToCreate             []string              `json:"devices_to_create"`
	DevicesToUpdate             []string              `json:"devices_to_update"`
	DeviceColumnsToAdd          []string              `json:"device_columns_to_add"`
	DeviceColumnsToDelete       []string              `json:"device_columns_to_delete"`
	BucketsToCreate             []string              `json:"buckets_to_create"`
	BucketsToUpdate             []string              `json:"buckets_to_update"`
	BucketFilesToCreate         []BucketFileUpdate    `json:"bucket_files_to_create"`
	BucketFilesToUpdate         []BucketFileUpdate    `json:"bucket_files_to_update"`
	ExternalDbsToCreate         []string              `json:"external_databases_to_create"`
	ExternalDbsToUpdate         []string              `json:"external_databases_to_update"`
	WebhooksToCreate            []string              `json:"webhooks_to_create"`
	WebhooksToUpdate            []string              `json:"webhooks_to_update"`
	CollectionsToCreate         []*CollectionUpdate   `json:"collections_to_create"`
	CollectionsToUpdate         []*CollectionUpdate   `json:"collections_to_update"`
	TriggersToCreate            []string              `json:"triggers_to_create"`
	TriggersToUpdate            []string              `json:"triggers_to_update"`
	TimersToCreate              []string              `json:"timers_to_create"`
	TimersToUpdate              []string              `json:"timers_to_update"`
	EdgesToCreate               []string              `json:"edges_to_create"`
	EdgesToUpdate               []string              `json:"edges_to_update"`
	EdgeColumnsToAdd            []string              `json:"edge_columns_to_add"`
	EdgeColumnsToDelete         []string              `json:"edge_columns_to_delete"`
	SecretsToCreate             []string              `json:"secrets_to_create"`
	SecretsToUpdate             []string              `json:"secrets_to_update"`
	PluginsToCreate             []string              `json:"plugins_to_create"`
	PluginsToUpdate             []string              `json:"plugins_to_update"`
	CachesToCreate              []string              `json:"caches_to_create"`
	CachesToUpdate              []string              `json:"caches_to_update"`
	PortalsToCreate             []string              `json:"portals_to_create"`
	PortalsToUpdate             []string              `json:"portals_to_update"`
	DeploymentsToCreate         []string              `json:"deployments_to_create"`
	DeploymentsToUpdate         []string              `json:"deployments_to_update"`
	AdaptorsToCreate            []string              `json:"adaptors_to_create"`
	AdaptorsToUpdate            []string              `json:"adaptors_to_update"`
	AdaptorFilesToCreate        []AdaptorFileUpdate   `json:"adaptor_files_to_create"`
	AdaptorFilesToUpdate        []AdaptorFileUpdate   `json:"adaptor_files_to_update"`
	MessageHistoryStorageTopics []string              `json:"message_history_topics_to_store"`
	MessageTypeTriggers         []*MessageTypeTrigger `json:"message_type_triggers"`
	Errors                      []string              `json:"errors"`
}

type BucketFileUpdate struct {
	BucketName   string `json:"bucket_name"`
	BucketBox    string `json:"bucket_box"`
	RelativePath string `json:"relative_path"`
}

type CollectionUpdate struct {
	Name           string   `json:"name"`
	AddedColumns   []string `json:"added_columns"`
	RemovedColumns []string `json:"removed_columns"`
	AddedIndexes   []string `json:"added_indexes"`
	RemovedIndexes []string `json:"removed_indexes"`
	NumUpserts     int      `json:"num_upserts"`
}

type AdaptorFileUpdate struct {
	AdaptorName string `json:"adaptor_name"`
	FileName    string `json:"file_name"`
}

type MessageTypeTrigger struct {
	MessageType  string `json:"message_type"`
	TopicPattern string `json:"topic_pattern"`
}

func (d *DevClient) UploadToSystemDryRun(systemKey string, zipBuffer []byte) (*SystemUploadDryRun, error) {
	resp, err := d.doSystemUpload(systemKey, zipBuffer, true)
	if err != nil {
		return nil, err
	}

	var dryRun SystemUploadDryRun
	if err = decodeMapToStruct(resp, &dryRun); err != nil {
		return nil, err
	}

	return &dryRun, nil
}

func (d *DevClient) UploadToSystem(systemKey string, zipBuffer []byte) (interface{}, error) {
	return d.doSystemUpload(systemKey, zipBuffer, false)
}

func (d *DevClient) doSystemUpload(systemKey string, zipBuffer []byte, dryRun bool) (interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	zipContents, err := getContentsForFile(zipBuffer)
	if err != nil {
		return nil, err
	}

	body := map[string]interface{}{
		"contents": zipContents,
	}

	url := fmt.Sprintf("%s%s/upload?dryRun=%t", _SYSTEM_MIGRATION_DEV_PREAMBLE, systemKey, dryRun)
	resp, err := post(d, url, body, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (d *DevClient) GetSystemUploadVersion(systemKey string) (interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s%s/upload", _SYSTEM_MIGRATION_DEV_PREAMBLE, systemKey)
	resp, err := get(d, url, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
