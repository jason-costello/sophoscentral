package sophoscentral

import (
	"context"
	"fmt"
)

// WebControlLocalSitesDelete - Delete  local site by id.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/web-control/local-sites/{localSiteId}
func (e *EndpointService) WebControlLocalSitesDelete(ctx context.Context, tenantID string,  tenantURL BaseURL, localSiteID string) (*DeletedResponse, error) {
	path := fmt.Sprintf("%ssettings/web-control/local-sites/%s", e.basePath, localSiteID)

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
