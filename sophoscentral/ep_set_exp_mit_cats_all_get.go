package sophoscentral

import (
	"context"
	"fmt"
)


type ExploitMitigationCategories struct{
	Items []ExploitMitigationCategory `json:"items,omitempty"`

}

type ExploitMitigationCategory struct{
	ID string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Settings []struct{
		Setting string
		Protected bool
	} `json:"settings,omitempty"`
}

type ExploitMitigationListOptions struct{

	ListByPageOffset

	// Type - Exploit mitigation application type
	// The following values are allowed:
	// detected, custom
	Type ExploitMitigationAppType `url:"type,omitempty"`
	// Whether or not Exploit Mitigation Application has been customized.
	Modified bool `url:"modified,omitempty"`
}

// ExploitMitigationCategoriesList - List Exploit Mitigation categories.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/exploit-mitigation/categories
func (e *EndpointService) ExploitMitigationCategoriesList(ctx context.Context, tenantID string,  tenantURL BaseURL, opts *ExploitMitigationListOptions) (*ExploitMitigationCategory, error) {
	path := fmt.Sprintf("%ssettings/exploit-mitigation/categories", e.basePath)
	var err error
	if path, err = addOptions(path, opts); err != nil{
		return nil, err
	}
	req, err := e.client.NewRequest(ctx, "GET", path, &tenantURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")
	e.client.Token.SetAuthHeader(req)
	sei := new(ExploitMitigationCategory)
	_, err = e.client.Do(ctx, req, sei)
	if err != nil {
		return nil, err
	}

	return sei, nil

}
