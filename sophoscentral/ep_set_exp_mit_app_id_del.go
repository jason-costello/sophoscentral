package sophoscentral

import (
	"context"
	"fmt"
)

// ExploitMitigationApplicationDelete -Get detected exploits and the number of each detected exploit.
// page by offset pagination
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/exploit-mitigation/detected-exploits
func (e *EndpointService) ExploitMitigationApplicationDelete(ctx context.Context, tenantID, exploitMitigationAppId string,  tenantURL BaseURL) (*DeletedResponse, error) {
	path := fmt.Sprintf("%ssettings/exploit-mitigation/applications/%s", e.basePath, exploitMitigationAppId)

	req, err := e.client.NewRequest(ctx, "DELETE", path, &tenantURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")
	e.client.Token.SetAuthHeader(req)
	sei := new(DeletedResponse)
	_, err = e.client.Do(ctx, req, sei)
	if err != nil {
		return nil, err
	}

	return sei, nil

}
