package sophoscentral

import (
	"context"
	"fmt"
)

// IsolationSettings something
type IsolationSettings struct {
	LastDisabledBy interface{} `json:"lastDisabledBy,omitempty"`
	LastEnabledBy  struct {
		AccountID   *string `json:"accountId,omitempty"`
		AccountType *string `json:"accountType,omitempty"`
		Name        *string `json:"name,omitempty"`
		ID          *string `json:"id,omitempty"`
		Type        *string `json:"type,omitempty"`
	} `json:"lastEnabledBy,omitempty"`
	Comment        *string `json:"comment,omitempty"`
	LastEnabledAt  *string `json:"lastEnabledAt,omitempty"`
	LastDisabledAt *string `json:"lastDisabledAt,omitempty"`
	Enabled        *bool   `json:"enabled,omitempty"`
}
type ToggleIsolations struct {
	Enabled *bool    `json:"enabled,omitempty"`
	Comment *string  `json:"comment,omitempty"`
	Ids     []string `json:"ids,omitempty"`
}

type ToggleIsolationSettings struct {
	Items []struct {
		Isolation *IsolationSettings `json:"isolation"`
		ID        *string            `json:"id"`
	} `json:"items"`
}
type IsolationSettingsUpdate struct {
	Enabled *bool   `json:"enabled,omitempty"`
	Comment *string `json:"comment,omitempty"`
}

// GetIsolationSettings for an endpoint.
func (e *EndpointService) GetIsolationSettings(ctx context.Context, tenantID string, tenantURL BaseURL, endpointID string) (*IsolationSettings, *Response, error) {

	path := fmt.Sprintf("%sendpoints/%s/isolation", e.basePath, endpointID)

	req, err := e.client.NewRequest("GET", path, &tenantURL, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("X-Tenant-ID", tenantID)

	tps := new(IsolationSettings)
	resp, err := e.client.Do(ctx, req, tps)
	if err != nil {
		return nil, resp, err
	}

	return tps, resp, nil

}

// UpdateIsolationSetting for an endpoint.
func (e *EndpointService) UpdateIsolationSetting(ctx context.Context, tenantID string, tenantURL BaseURL, endpointID string, update IsolationSettingsUpdate) (*IsolationSettings, *Response, error) {

	path := fmt.Sprintf("%sendpoints/%s/isolation", e.basePath, endpointID)

	req, err := e.client.NewRequest("PATCH", path, &tenantURL, update)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("X-Tenant-ID", tenantID)

	tps := new(IsolationSettings)
	resp, err := e.client.Do(ctx, req, tps)
	if err != nil {
		return nil, resp, err
	}

	return tps, resp, nil

}

// ToggleIsolation Turn on or off endpoint isolation for multiple endpoints.
// GetIsolationSettings for an endpoint.
func (e *EndpointService) ToggleIsolation(ctx context.Context, tenantID string, tenantURL BaseURL, ti ToggleIsolations) (*ToggleIsolationSettings, *Response, error) {

	path := "endpoints/isolation"

	req, err := e.client.NewRequest("POST", path, &tenantURL, ti)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("X-Tenant-ID", tenantID)

	tps := new(ToggleIsolationSettings)
	resp, err := e.client.Do(ctx, req, tps)
	if err != nil {
		return nil, resp, err
	}

	return tps, resp, nil

}
