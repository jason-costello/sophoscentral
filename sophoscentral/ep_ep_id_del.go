package sophoscentral

import (
	"context"
	"errors"
	"fmt"
)
// Delete an endpoint
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/endpoints/{endpointId}
func (e *EndpointService) Delete(ctx context.Context, tenantID string, tenantURL BaseURL, endpointID string) (*Item, *Response, error) {
	// url path to call
	if endpointID == "" {
		return nil, nil, errors.New("endpointID is empty")
	}
	path := fmt.Sprintf("%sendpoints/%s", e.basePath, endpointID)

	req, err := e.client.NewRequest(ctx, "DELETE", path, &tenantURL, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("X-Tenant-ID", tenantID)

	ep := new(Item)
	resp, err := e.client.Do(ctx, req, ep)
	if err != nil {
		return nil, resp, err
	}
	return ep, resp, nil

}
