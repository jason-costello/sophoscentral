package sophoscentral

import (
	"context"
	"fmt"
)

// IntrusionPreventionExclusionGet - Get  Intrusion Prevention exclusion by id.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/exclusions/intrusion-prevention/{exclusionId}
func (e *EndpointService) IntrusionPreventionExclusionGet(ctx context.Context, tenantID string,  tenantURL BaseURL, exclusionID string) (*IntrusionPreventionsExclusionItem, error) {
	path := fmt.Sprintf("%ssettings/exclusions/intrusion-prevention/%s", e.basePath, exclusionID)

	req, err := e.client.NewRequest(ctx, "GET", path, &tenantURL,nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")

	e.client.Token.SetAuthHeader(req)
	ipe := new(IntrusionPreventionsExclusionItem)
	_, err = e.client.Do(ctx, req, ipe)
	if err != nil {
		return nil, err
	}

	return ipe, nil

}
