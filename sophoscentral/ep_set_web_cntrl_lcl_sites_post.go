package sophoscentral

import (
	"context"
	"fmt"
)

// WebControlLocalSitesUpdate - Update a local site definition
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/web-control/local-sites/{localSiteId}
func (e *EndpointService) WebControlLocalSitesAdd(ctx context.Context, tenantID string,  tenantURL BaseURL, localSite WebControlLocalSiteRequest ) (*WebControlLocalSite, error) {
	path := fmt.Sprintf("%ssettings/web-control/local-sites", e.basePath)

	if err := localSite.Verify();err != nil{
		return nil, err
	}


	req, err := e.client.NewRequest(ctx, "POST", path, &tenantURL, localSite)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")
	e.client.Token.SetAuthHeader(req)
	sei := new(WebControlLocalSite)
	_, err = e.client.Do(ctx, req, sei)
	if err != nil {
		return nil, err
	}

	return sei, nil

}
