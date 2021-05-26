package sophoscentral

import (
	"context"
	"fmt"
)

// ScanningExclusionGet - Get  scanning exclusions by exclusion ID.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/exclusions/scanning/{exclusionId}
func (e *EndpointService) ScanningExclusionGet(ctx context.Context, tenantID string,  tenantURL BaseURL, exclusionID string) (*ScanningExclusionItem, error) {
	path := fmt.Sprintf("%ssettings/exclusions/scanning/%s", e.basePath, exclusionID)

	req, err := e.client.NewRequest(ctx, "GET", path, &tenantURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")
	e.client.Token.SetAuthHeader(req)
	sei := new(ScanningExclusionItem)
	_, err = e.client.Do(ctx, req, sei)
	if err != nil {
		return nil, err
	}

	return sei, nil

}
