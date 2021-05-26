package sophoscentral

import (
	"context"
	"fmt"
)



type WebControlCategories []struct{
	ID int `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Label string `json:"label,omitempty"`
}


// WebControlCategoriesList - List all Web Control categories.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/web-control/categories
func (e *EndpointService) WebControlCategoriesList(ctx context.Context, tenantID string,  tenantURL BaseURL) (*WebControlCategories, error) {
	path := fmt.Sprintf("%ssettings/web-control/categories", e.basePath)

	req, err := e.client.NewRequest(ctx, "GET", path, &tenantURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")
	e.client.Token.SetAuthHeader(req)
	sei := new(WebControlCategories)
	_, err = e.client.Do(ctx, req, sei)
	if err != nil {
		return nil, err
	}

	return sei, nil

}
