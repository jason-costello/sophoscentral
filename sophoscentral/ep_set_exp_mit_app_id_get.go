package sophoscentral

import (
	"context"
	"fmt"
)

// ExploitMitigationApplicationGet - Get Exploit Mitigation settings for an application.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/exploit-mitigation/applications/{exploitMitigationApplicationId}
func (e *EndpointService) ExploitMitigationApplicationGet(ctx context.Context, tenantID, exploitMitigationAppId string,  tenantURL BaseURL) (*ExploitMitigationAppItem, error) {
	path := fmt.Sprintf("%ssettings/exploit-mitigation/applications/%s", e.basePath, exploitMitigationAppId)

	req, err := e.client.NewRequest(ctx, "GET", path, &tenantURL, nil)
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
