package GoSDK

import (
	"fmt"
	"net/url"
	"strconv"
)

const _MTLS_PREAMBLE = "/api/v/1/mtls-certs/"

type MtlsRevocation struct {
	SerialNumber     string `json:"serial_number"`
	IssuerSpkiSha256 string `json:"issuer_spki_sha256"`
	ReasonCode       *int   `json:"reason_code,omitempty"`
	Description      string `json:"description,omitempty"`
}

func (d *DevClient) AddMTLSSystemRoot(systemKey, name, certificate string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	body := map[string]interface{}{"certificate": certificate}
	resp, err := post(d, mtlsRootsEndpoint(systemKey, name), body, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) GetMTLSSystemRoots(systemKey string) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, mtlsRootsEndpoint(systemKey), nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) DeleteMTLSSystemRoot(systemKey, name string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := delete(d, mtlsRootsEndpoint(systemKey, name), nil, creds, nil)
	_, err = mapResponse(resp, err)
	return err
}

func (d *DevClient) RevokeMTLSCertificate(systemKey string, revocation *MtlsRevocation) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := put(d, mtlsRevokedEndpoint(systemKey), revocation, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

// GetMTLSRevokedCertificates returns one page of the system's revocations. An empty issuerSpkiSha256 returns all issuers.
func (d *DevClient) GetMTLSRevokedCertificates(systemKey, issuerSpkiSha256 string, pageSize, page int) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	query := map[string]string{}
	if issuerSpkiSha256 != "" {
		query["issuer_spki_sha256"] = issuerSpkiSha256
	}
	if pageSize > 0 {
		query["page_size"] = strconv.Itoa(pageSize)
	}
	if page > 0 {
		query["page"] = strconv.Itoa(page)
	}
	resp, err := get(d, mtlsRevokedEndpoint(systemKey), query, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) UnrevokeMTLSCertificate(systemKey, issuerSpkiSha256, serialNumber string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	query := map[string]string{
		"issuer_spki_sha256": issuerSpkiSha256,
		"serial_number":      serialNumber,
	}
	resp, err := delete(d, mtlsRevokedEndpoint(systemKey), query, creds, nil)
	_, err = mapResponse(resp, err)
	return err
}

func mtlsRootsEndpoint(systemKey string, name ...string) string {
	endpoint := _MTLS_PREAMBLE + "roots/" + url.PathEscape(systemKey)
	for _, n := range name {
		endpoint += "/" + url.PathEscape(n)
	}
	return endpoint
}

func mtlsRevokedEndpoint(systemKey string) string {
	return fmt.Sprintf("%srevoked/%s", _MTLS_PREAMBLE, url.PathEscape(systemKey))
}
