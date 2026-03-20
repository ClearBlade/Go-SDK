package GoSDK

import (
	"fmt"
	"time"
)

//
// DevClient must already be platform admin to use these endpoints
//

func (d *DevClient) PromoteDevToPlatformAdmin(email string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"email": email,
	}

	resp, err := post(d, d.preamble()+"/promotedev", data, creds, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error promoting %s to admin: %v", email, resp.Body)
	}
	return nil
}

func (d *DevClient) DemoteDevFromPlatformAdmin(email string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"email": email,
	}

	resp, err := post(d, d.preamble()+"/demotedev", data, creds, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error demoting %s to admin: %v", email, resp.Body)
	}
	return nil
}

func (d *DevClient) ResetDevelopersPassword(email, newPassword string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"email":        email,
		"new_password": newPassword,
	}

	resp, err := post(d, d.preamble()+"/resetpassword", data, creds, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error resetting %s's password: %v", email, resp.Body)
	}
	return nil
}

func (d *DevClient) GetSystemAnalytics(systemKey string) (interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	analyticsEndpoint := d.preamble() + "/platform/system/" + systemKey
	resp, err := get(d, analyticsEndpoint, nil, creds, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting %s's analytics: %v", systemKey, resp.Body)
	}
	return resp.Body, nil
}

func (d *DevClient) GetAllSystemsAnalytics(query string) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	analyticsEndpoint := d.preamble() + "/platform/systems"
	q := make(map[string]string)
	if len(query) != 0 {
		q["query"] = query
	} else {
		q = nil
	}
	resp, err := get(d, analyticsEndpoint, q, creds, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting analytics: %v", resp.Body)
	}
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) DisableSystem(systemKey string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	disableSystemEndpoint := d.preamble() + "/platform/" + systemKey
	resp, err := delete(d, disableSystemEndpoint, nil, creds, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error disabling system %s: %v", systemKey, resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) EnableSystem(systemKey string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	disableSystemEndpoint := d.preamble() + "/platform/" + systemKey
	resp, err := post(d, disableSystemEndpoint, nil, creds, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error disabling system %s: %v", systemKey, resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) GetDeveloper(devEmail string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	q := make(map[string]string)
	q["developer"] = devEmail

	developerEndpoint := d.preamble() + "/platform/developer"
	resp, err := get(d, developerEndpoint, q, creds, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting developer %s: %v", devEmail, resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) GetAllDevelopers() (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	allDevelopersEndpoint := d.preamble() + "/platform/developers"
	resp, err := get(d, allDevelopersEndpoint, nil, creds, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting all developers: %v", resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) SetDeveloper(email string, admin, disabled bool) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{
		"email":    email,
		"admin":    admin,
		"disabled": disabled,
	}

	developerEndpoint := d.preamble() + "/platform/developer"
	resp, err := post(d, developerEndpoint, data, creds, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error setting developer %s: %v", email, resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) GetMetrics(metricType string) (interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	metricsEndpoint := d.preamble() + "/platform/" + metricType
	resp, err := get(d, metricsEndpoint, nil, creds, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting metric %s: %v", metricType, resp.Body)
	}
	return resp.Body, nil
}

func (d *DevClient) KillClient(systemKey, clientID string) (interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	path := "/admin/" + systemKey + "/killclient?clientid=" + clientID
	resp, err := post(d, path, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (d *DevClient) DeleteDeveloper(email string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := delete(d, "/admin/platform/developer", map[string]string{"email": email}, creds, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("%+v", resp.Body)
	}
	return nil
}

type DeveloperUpdate struct {
	Disabled    bool   `json:"disabled"`
	Admin       bool   `json:"admin"`
	Email       string `json:"email"`
	OidcEnabled bool   `json:"oidc_enabled"`
}

func (d *DevClient) UpdatePlatformDeveloper(update DeveloperUpdate) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	bodyMap := structToMap(update)
	resp, err := post(d, "/admin/platform/developer", bodyMap, creds, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("%+v", resp.Body)
	}
	return nil
}

func (d *DevClient) DeleteSelf() error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := delete(d, "/admin/userinfo", nil, creds, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("%+v", resp.Body)
	}
	return nil
}

type SystemAlias struct {
	Alias          string `json:"alias"`
	SystemKey      string `json:"system_key"`
	IsPrimaryAlias bool   `json:"is_primary_alias"`
	Disabled       bool   `json:"disabled"`
}

func (d *DevClient) GetSystemAliases(systemKey string) ([]SystemAlias, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("/admin/aliases/%s", systemKey)
	resp, err := get(d, uri, nil, creds, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%+v", resp.Body)
	}

	var aliases []SystemAlias
	if err := decodeMapToStruct(resp.Body, &aliases); err != nil {
		return nil, fmt.Errorf("could not decode aliases: %w", err)
	}

	return aliases, nil
}

func (d *DevClient) CreateSystemAlias(alias SystemAlias) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	uri := fmt.Sprintf("/admin/aliases/%s", alias.SystemKey)
	resp, err := post(d, uri, alias, creds, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("%+v", resp.Body)
	}

	return nil
}

func (d *DevClient) UpdateSystemAlias(systemKey, alias string, changes map[string]any) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	uri := fmt.Sprintf("/admin/aliases/%s/%s", systemKey, alias)
	resp, err := put(d, uri, changes, creds, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("%+v", resp.Body)
	}

	return nil
}

func (d *DevClient) DeleteSystemAlias(systemKey, alias string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	uri := fmt.Sprintf("/admin/aliases/%s/%s", systemKey, alias)
	resp, err := delete(d, uri, nil, creds, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("%+v", resp.Body)
	}

	return nil
}

type ProfileFile struct {
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"last_modified"`
	ProfileType  string    `json:"profile_type"`
	NodeId       string    `json:"node_id"`
}

func (d *DevClient) ListPlatformProfiles() ([]ProfileFile, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	uri := "/admin/platform/profiles"
	resp, err := get(d, uri, nil, creds, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%+v", resp.Body)
	}

	var profiles []ProfileFile
	if err := decodeMapToStruct(resp.Body, &profiles); err != nil {
		return nil, fmt.Errorf("could not decode profiles: %w", err)
	}

	return profiles, nil
}

func (d *DevClient) GetPlatformProfile(nodeId, profileType, name string) ([]byte, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("/admin/platform/profiles/%s/%s", profileType, name)
	queryParams := map[string]string{
		"node": nodeId,
	}

	req := &CbReq{
		Method:      "GET",
		Endpoint:    uri,
		QueryString: query_to_string(queryParams),
		NoDecode:    true,
	}

	resp, err := do(d, req, creds)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%+v", resp.Body)
	}

	respStr, ok := resp.Body.(string)
	if !ok {
		return nil, fmt.Errorf("expected string, got %T", resp.Body)
	}

	return []byte(respStr), nil
}
