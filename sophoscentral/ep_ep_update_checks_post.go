package sophoscentral

import (
	"context"
	"fmt"
)

// EndpointCheckForUpdate - Sends a request to the endpoint to check for Sophos management agent software updates
// as well as protection data updates.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/endpoints/{endpointId}/update-checks
// post empty body to endpointID
func (e *EndpointService) EndpointCheckForUpdate(ctx context.Context, endpointID, tenantID string, tenantURL BaseURL) (*EndpointScanResponse, *Response, error) {
	path := fmt.Sprintf("%sendpoints/%s/scans", e.basePath, endpointID)

	req, err := e.client.NewRequest(ctx, "POST", path, &tenantURL, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("X-Tenant-ID", tenantID)

	tps := new(EndpointScanResponse)
	resp, err := e.client.Do(ctx, req, tps)
	if err != nil {
		return nil, resp, err
	}

	return tps, resp, nil

}
