package sophoscentral

import (
	"context"
	"errors"
	"fmt"
)

// ScanningExclusionAdd - Add a new scanning exclusion.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/exclusions/scanning
func (e *EndpointService) ScanningExclusionAdd(ctx context.Context, tenantID string,  tenantURL BaseURL, exclusionItem ScanningExclusionItem) (*ScanningExclusionItem, error) {
	path := fmt.Sprintf("%ssettings/exclusions/scanning", e.basePath)


	if len(*exclusionItem.Comment) > 100{
		return nil, errors.New("comment field is limited to 100 chars or less")
	}
	req, err := e.client.NewRequest(ctx, "POST", path, &tenantURL, exclusionItem)
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
