package sophoscentral

import (
	"context"
	"errors"
	"fmt"
)

// WebControlLocalSitesUpdate - Update a local site definition
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/web-control/local-sites/{localSiteId}
func (e *EndpointService) WebControlLocalSitesUpdate(ctx context.Context, tenantID string,  tenantURL BaseURL, localSiteID string, localSite WebControlLocalSiteRequest) (*WebControlLocalSite, error) {
	path := fmt.Sprintf("%ssettings/web-control/local-sites/%s", e.basePath, localSiteID)
	if err := localSite.Verify();err != nil{
		return nil, err
	}


	req, err := e.client.NewRequest(ctx, "PATCH", path, &tenantURL, localSite)
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


func (w WebControlLocalSiteRequest) Verify()error{

	if w.CategoryID == 0 && len(w.Tags) == 0{
		return errors.New("must have either categoryID or tags")
	}

	if w.CategoryID > 57 && len(w.Tags) == 0{
		return errors.New("category ID must be between 1 and 57 inclusive")
	}

	for x := range w.Tags{
		if len(w.Tags[x]) < 1 || len(w.Tags[x]) > 50{
			return errors.New(fmt.Sprintf("tag index %d must be <=1 length <= 50", x))
		}

	}
	if len(w.URL) > 2048{
		return errors.New("url must be 1 <= length <= 2048")
	}

	if len(w.Comment) > 300{
		return errors.New("comment must be less than 301 chars")
	}

	return nil

}
