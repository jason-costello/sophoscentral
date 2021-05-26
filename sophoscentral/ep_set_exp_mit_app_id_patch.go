package sophoscentral

import (
	"context"
	"fmt"
)

// ExploitMitigationAppUpdate - Update Exploit Mitigation settings for an application.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/exploit-mitigation/applications/{exploitMitigationApplicationId}
func (e *EndpointService) ExploitMitigationAppUpdate(ctx context.Context, tenantID string,
	tenantURL BaseURL, exploitMitigationAppId string, item ExploitMitigationNewAppRequest) (*ExploitMitigationAppItem, error) {
	path := fmt.Sprintf("%ssettings/exploit-mitigation/applications/%s", e.basePath, exploitMitigationAppId)

	if err := verifyExploitMitigationNewAppRequest(item); err != nil{
		return nil, err
	}


	req, err := e.client.NewRequest(ctx, "PATCH", path, &tenantURL, item)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")
	e.client.Token.SetAuthHeader(req)
	sei := new(ExploitMitigationAppItem)
	_, err = e.client.Do(ctx, req, sei)
	if err != nil {
		return nil, err
	}

	return sei, nil

}
