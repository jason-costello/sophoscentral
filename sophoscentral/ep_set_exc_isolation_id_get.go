package sophoscentral

import (
	"context"
	"fmt"
)




// IsolationExclusionGet - Get a single Isolation exclusion by ID.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/exclusions/isolation/{exclusionId}
func (e *EndpointService) IsolationExclusionGet(ctx context.Context, tenantID string,  tenantURL BaseURL, exclusionID string) (*IsolationExclusionItem, error) {
	path := fmt.Sprintf("%ssettings/exclusions/isolation/%s", e.basePath, exclusionID)

	req, err := e.client.NewRequest(ctx, "GET", path, &tenantURL,nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")

	e.client.Token.SetAuthHeader(req)
	ipe := new(IsolationExclusionItem)
	_, err = e.client.Do(ctx, req, ipe)
	if err != nil {
		return nil, err
	}

	return ipe, nil

}
