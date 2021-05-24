package sophoscentral

import (
	"context"
	"errors"
	"fmt"
)

type TamperProtectionSettings struct {
	Password          *string               `json:"password,omitempty"`
	PreviousPasswords []TPPreviousPasswords `json:"previousPasswords,omitempty"`
	Enabled           bool                  `json:"enabled,omitempty"`
}
type TPPreviousPasswords struct {
	Password      *string `json:"password,omitempty"`
	InvalidatedAt *string `json:"invalidatedAt,omitempty"`
}

// TamperProtection fetches the TamperProtection settings for a specific endpoint
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/endpoints/{endpointId}/tamper-protection
func (e *EndpointService) TamperProtection(ctx context.Context, tenantID string, tenantURL BaseURL, endpointID string) (*TamperProtectionSettings, *Response, error) {

	path := fmt.Sprintf("%sendpoints/%s/tamper-protection", e.basePath, endpointID)

	req, err := e.client.NewRequest("GET", path, &tenantURL, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("X-Tenant-ID", tenantID)

	tps := new(TamperProtectionSettings)
	resp, err := e.client.Do(ctx, req, tps)
	if err != nil {
		return nil, resp, err
	}

	return tps, resp, nil

}

// SetTamperProtection Turns Tamper Protection on or off on an
// endpoint. Or generates a new Tamper Protection password.
// Note that Tamper Protection can be turned on for an endpoint
// only if it has also been turned on globally.

func (e *EndpointService) SetTamperProtection(ctx context.Context, tenantID string, tenantURL BaseURL, endpointID string) (*TamperProtectionSettings, *Response, error) {

	path := fmt.Sprintf("%sendpoints/%s/tamper-protection", e.basePath, endpointID)

	req, err := e.client.NewRequest("POST", path, &tenantURL, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("X-Tenant-ID", tenantID)

	tps := new(TamperProtectionSettings)
	resp, err := e.client.Do(ctx, req, tps)
	if err != nil {
		return nil, resp, err
	}

	if resp.StatusCode != 201 {
		return nil, resp, errors.New(resp.Status)
	}
	return tps, resp, nil

}

// GlobalTamperProtectionEnabled is the response object for the call to GlobalTamperProtectionStatus
type GlobalTamperProtectionEnabled struct {
	// indicates that status of global tamper protection
	Enabled bool `json:"enabled,omitempty"`
}

// GlobalTamperProtectionStatus checks whether Tamper Protection is turned on globally.
// only returns 200 or 500
func (e *EndpointService) GlobalTamperProtectionStatus(ctx context.Context, tenantID string, tenantURL BaseURL) (*GlobalTamperProtectionEnabled, *Response, error) {
	path := fmt.Sprintf("%sendpoints/settings/tamper-protection", e.basePath)
	req, err := e.client.NewRequest("GET", path, &tenantURL, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("X-Tenant-ID", tenantID)

	tps := new(GlobalTamperProtectionEnabled)

	resp, err := e.client.Do(ctx, req, tps)
	if err != nil {
		return nil, resp, err
	}

	if resp.StatusCode != 200 {
		return nil, resp, errors.New(resp.Status)
	}
	return tps, resp, nil
}
