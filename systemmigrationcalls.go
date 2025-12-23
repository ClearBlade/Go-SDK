package GoSDK

import (
	"errors"
	"fmt"
)

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
	FilestoresToCreate          []string              `json:"file_stores_to_create"`
	FilestoresToUpdate          []string              `json:"file_stores_to_update"`
	FilestoreFilesToCreate      []FileStoreFileUpdate `json:"file_store_files_to_create"`
	FilestoreFilesToUpdate      []FileStoreFileUpdate `json:"file_store_files_to_update"`
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
	MessageTypeTriggers         []*TriggeredMsgType   `json:"message_type_triggers"`
	Warnings                    []string              `json:"warnings"`
	Errors                      []string              `json:"errors"`
}

type BucketFileUpdate struct {
	BucketName   string `json:"bucket_name"`
	BucketBox    string `json:"bucket_box"`
	RelativePath string `json:"relative_path"`
}

type FileStoreFileUpdate struct {
	FileStoreName string `json:"file_store_name"`
	Path          string `json:"path"`
}

type CollectionUpdate struct {
	Name                        string   `json:"name"`
	AddedColumns                []string `json:"added_columns"`
	RemovedColumns              []string `json:"removed_columns"`
	AddedIndexes                []string `json:"added_indexes"`
	RemovedIndexes              []string `json:"removed_indexes"`
	NumUpserts                  int      `json:"num_upserts"`
	UpdatedHypertableProperties []string `json:"updated_hypertable_properties"`
}

type AdaptorFileUpdate struct {
	AdaptorName string `json:"adaptor_name"`
	FileName    string `json:"file_name"`
}

type TriggeredMsgType struct {
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

type SystemUploadChanges struct {
	UpdatedLibraries            []string              `json:"updated_libraries"`
	UpdatedServices             []string              `json:"updated_services"`
	UpdatedRoles                []string              `json:"updated_roles"`
	CreatedServices             []string              `json:"created_services"`
	CreatedLibraries            []string              `json:"created_libraries"`
	CreatedRoles                []string              `json:"created_roles"`
	CreatedUsers                []string              `json:"created_users"`
	UpdatedUsers                []string              `json:"updated_users"`
	AddedUserColumns            []string              `json:"added_user_columns"`
	DroppedUserColumns          []string              `json:"dropped_user_columns"`
	CreatedDevices              []string              `json:"created_devices"`
	UpdatedDevices              []string              `json:"updated_devices"`
	AddedDeviceColumns          []string              `json:"added_device_columns"`
	DroppedDeviceColumns        []string              `json:"dropped_device_columns"`
	CreatedBuckets              []string              `json:"created_buckets"`
	UpdatedBuckets              []string              `json:"updated_buckets"`
	CreatedBucketFiles          []BucketFileUpdate    `json:"created_bucket_files"`
	UpdatedBucketFiles          []BucketFileUpdate    `json:"updated_bucket_files"`
	CreatedFileStores           []string              `json:"created_file_stores"`
	UpdatedFileStores           []string              `json:"updated_file_stores"`
	CreatedFileStoreFiles       []FileStoreFileUpdate `json:"created_file_store_files"`
	UpdatedFileStoreFiles       []FileStoreFileUpdate `json:"updated_file_store_files"`
	CreatedExternalDbs          []string              `json:"created_external_databases"`
	UpdatedExternalDbs          []string              `json:"updated_external_databases"`
	CreatedCollections          []*CollectionUpdate   `json:"created_collections"`
	UpdatedCollections          []*CollectionUpdate   `json:"updated_collections"`
	CreatedWebhooks             []string              `json:"created_webhooks"`
	UpdatedWebhooks             []string              `json:"updated_webhooks"`
	CreatedTriggers             []string              `json:"created_triggers"`
	UpdatedTriggers             []string              `json:"updated_triggers"`
	CreatedTimers               []string              `json:"created_timers"`
	UpdatedTimers               []string              `json:"updated_timers"`
	CreatedEdges                []string              `json:"created_edges"`
	UpdatedEdges                []string              `json:"updated_edges"`
	AddedEdgeColumns            []string              `json:"added_edge_columns"`
	DroppedEdgeColumns          []string              `json:"dropped_edge_columns"`
	CreatedSecrets              []string              `json:"created_secrets"`
	UpdatedSecrets              []string              `json:"updated_secrets"`
	CreatedPlugins              []string              `json:"created_plugins"`
	UpdatedPlugins              []string              `json:"updated_plugins"`
	CreatedCaches               []string              `json:"created_caches"`
	UpdatedCaches               []string              `json:"updated_caches"`
	CreatedPortals              []string              `json:"created_portals"`
	UpdatedPortals              []string              `json:"updated_portals"`
	CreatedDeployments          []string              `json:"created_deployments"`
	UpdatedDeployments          []string              `json:"updated_deployments"`
	CreatedAdaptors             []string              `json:"created_adaptors"`
	UpdatedAdaptors             []string              `json:"updated_adaptors"`
	CreatedAdaptorFiles         []AdaptorFileUpdate   `json:"created_adaptor_files"`
	UpdatedAdaptorFiles         []AdaptorFileUpdate   `json:"updated_adaptor_files"`
	CollectionNameToId          map[string]string     `json:"collection_name_to_id"`
	RoleNameToId                map[string]string     `json:"role_name_to_id"`
	UserEmailToId               map[string]string     `json:"user_email_to_id"`
	MessageHistoryStorageTopics []string              `json:"message_history_storage_topics"`
	MessageTypeTriggers         []*TriggeredMsgType   `json:"message_type_triggers"`
	Errors                      []string              `json:"errors"`
}

func (r *SystemUploadChanges) Error() error {
	if len(r.Errors) == 0 {
		return nil
	}

	errs := make([]error, len(r.Errors))
	for i, err := range r.Errors {
		errs[i] = errors.New(err)
	}

	return fmt.Errorf("encountered the following errors while pushing: %w", errors.Join(errs...))
}

func (d *DevClient) UploadToSystem(systemKey string, zipBuffer []byte) (*SystemUploadChanges, error) {
	resp, err := d.doSystemUpload(systemKey, zipBuffer, false)
	if err != nil {
		return nil, err
	}

	var changes SystemUploadChanges
	if err = decodeMapToStruct(resp, &changes); err != nil {
		return nil, err
	}

	return &changes, nil
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
