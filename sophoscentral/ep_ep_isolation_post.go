package sophoscentral

import (
	"context"
)

// ToggleIsolation Turn on or off endpoint isolation for multiple endpoints.
// GetIsolationSettings for an endpoint.
func (e *EndpointService) ToggleIsolation(ctx context.Context, tenantID string, tenantURL BaseURL, ti ToggleIsolations) (*ToggleIsolationSettings, *Response, error) {

	path := "endpoints/isolation"

	req, err := e.client.NewRequest(ctx, "POST", path, &tenantURL, ti)
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
