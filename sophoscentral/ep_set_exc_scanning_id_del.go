package sophoscentral

import (
	"context"
	"fmt"
)

// ScanningExclusionDelete - Deletes a scanning exclusion. by exclusion ID.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/exclusions/scanning/{exclusionId}
func (e *EndpointService) ScanningExclusionDelete(ctx context.Context, tenantID string,  tenantURL BaseURL, exclusionID string) (*DeletedResponse, *Response, error) {
	path := fmt.Sprintf("%ssettings/exclusions/scanning/%s", e.basePath, exclusionID)

	req, err := e.client.NewRequest(ctx, "DELETE", path, &tenantURL, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")

	dr := new(DeletedResponse)

	resp, err := e.client.Do(ctx, req, dr)
	if err != nil {
		return nil,resp, err
	}

	return dr, resp, nil

}
