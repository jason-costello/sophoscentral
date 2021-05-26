package sophoscentral

import (
	"context"
	"fmt"
)



type WebControlLocalSites struct{

	Items []WebControlLocalSite `json:"items"`
	Pages PagesByOffset `json:"pages"`
}

type WebControlLocalSite struct{

	// ID of local site
	ID string `json:"id"`
	// CategoryID for the local site
	CategoryID int `json:"categoryId,omitempty"`
	// Tags associated with this local site
	Tags []string `json:"tags,omitempty"`
	// URL is a local site url
	URL string `json:"url,omitempty"`
	// Comment indicating why the local site was added.
	// length ≤ 300
	Comment string `json:"comment"`

}

type WebControlLocalSiteRequest struct{
	// CategoryID associated with this local site.
	//Either categoryId or tags must be provided.
	//1 ≤ value ≤ 57
	CategoryID int `json:"categoryId,omitempty"`
	// Tags associated with this local site
	// Array of tags associated with this local site setting.
	// Either categoryId or tags must be provided.
	Tags []string `json:"tags,omitempty"`
	// URL is a local site url
	// 1<= length<= 2048
	URL string `json:"url,omitempty"`
	// Comment indicating why the local site was added.
	// length ≤ 300
	Comment string `json:"comment,omitempty"`

}

// WebControlLocalSitesList - List all local sites.
// Pagination is using page by offset.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/web-control/local-sites
func (e *EndpointService) WebControlLocalSitesList(ctx context.Context, tenantID string,  tenantURL BaseURL, opts *ListByPageOffset) (*WebControlLocalSites, error) {
	path := fmt.Sprintf("%ssettings/web-control/local-sites", e.basePath)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil,  err
	}
	req, err := e.client.NewRequest(ctx, "GET", path, &tenantURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")
	e.client.Token.SetAuthHeader(req)
	sei := new(WebControlLocalSites)
	_, err = e.client.Do(ctx, req, sei)
	if err != nil {
		return nil, err
	}

	return sei, nil

}
