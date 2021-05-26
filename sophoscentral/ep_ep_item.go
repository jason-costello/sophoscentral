package sophoscentral

import (
	"context"
	"time"
)

type EPSettingItems struct {
	Items []struct {
		Id         string    `json:"id"`
		CreatedAt  *time.Time `json:"createdAt"`
		UpdatedAt  *time.Time `json:"updatedAt"`
		Properties *struct {
			FileName *string `json:"fileName"`
			Path     *string `json:"path"`
			Sha256   *string `json:"sha256"`
		} `json:"properties"`
		Comment        *string `json:"comment"`
		Type           *string `json:"type"`
		OriginEndpoint *struct {
			Id *string `json:"id"`
		} `json:"originEndpoint"`
	} `json:"items"`
	Pages *struct {
		Current *int `json:"current"`
		Size    *int `json:"size"`
		Total   *int `json:"total"`
		Items   *int `json:"items"`
		MaxSize *int `json:"maxSize"`
	} `json:"pages"`
}
type EPSettingItem struct {
	Id         string    `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Properties struct {
		FileName string `json:"fileName"`
		Path     string `json:"path"`
		Sha256   string `json:"sha256"`
	} `json:"properties"`
	Comment        string `json:"comment"`
	Type           string `json:"type"`
	OriginEndpoint struct {
		Id string `json:"id"`
	} `json:"originEndpoint"`
}
type DeletedResponse struct {
	Deleted bool `json:"deleted"`
}






// listEPSettingItems does the legwork for ListAllowedItems and ListBlockedItems functions as
// AllowedItems and BlockedItems are the same thing from a payload perspective.  EPSettingItems is used
// to keep a single struct in place of AllowedItems and BlockedItems
func (e *EndpointService) listEPSettingItems(ctx context.Context, path, tenantID string, tenantURL BaseURL, opts PageByOffsetOptions) (*EPSettingItems, error) {
	var maxSize = 100
	// GetAllowedItems for an endpoint.
	if opts.PageSize == 0 {
		opts.PageSize = maxSize
	}

	path, err := addOptions(path, opts)
	if err != nil {
		return nil, err
	}
	req, err := e.client.NewRequest(ctx, "GET", path, &tenantURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")
	e.client.Token.SetAuthHeader(req)

	epSI := new(EPSettingItems)
	_, err = e.client.Do(ctx, req, epSI)
	if err != nil {
		return nil, err
	}

	return epSI, nil

}
// getEPSettingItem does the legwork for AllowedItem and BlockedItem functions as
// AllowedItems and BlockedItems are the same thing from a payload perspective.
func (e *EndpointService) getEPSettingItem(ctx context.Context, tenantID, path string, tenantURL BaseURL) (*EPSettingItem, error) {
	req, err := e.client.NewRequest(ctx, "GET", path, &tenantURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")
	e.client.Token.SetAuthHeader(req)
	epSI := new(EPSettingItem)
	_, err = e.client.Do(ctx, req, epSI)
	if err != nil {
		return nil, err
	}

	return epSI, nil

}
