package sophoscentral


import (
	"context"
	"fmt"
)

// IsolationExclusionDelete - Deletes an Isolation exclusion.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/exclusions/isolation/{exclusionId}
func (e *EndpointService) IsolationExclusionDelete(ctx context.Context, tenantID string, tenantURL BaseURL, exclusionID string) (*DeletedResponse, *Response, error) {
	path := fmt.Sprintf("%ssettings/exclusions/isolation/%s", e.basePath, exclusionID)


	req, err := e.client.NewRequest(ctx, "DELETE", path, &tenantURL, nil)
	if err != nil {
		return nil, nil,err
	}

	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")

	ipe := new(DeletedResponse)
	resp, err := e.client.Do(ctx, req, ipe)
	if err != nil {
		return nil, resp,err
	}

	return ipe, resp, nil


}
