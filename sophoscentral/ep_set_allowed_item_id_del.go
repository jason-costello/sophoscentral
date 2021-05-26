package sophoscentral

import (
	"context"
	"errors"
	"fmt"
)


// AllowedItemDelete an allowed item
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/allowed-items/{allowedItemId}
func (e *EndpointService) AllowedItemDelete(ctx context.Context, tenantID string, tenantURL BaseURL, allowedItemID string) (*DeletedResponse, *Response, error) {
	// url path to call
	path := fmt.Sprintf("%ssettings/allowed-items/%s", e.basePath, allowedItemID)

	if allowedItemID == "" {
		return nil, nil, errors.New("allowedItemID is empty")
	}

	req, err := e.client.NewRequest(ctx, "DELETE", path, &tenantURL, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("X-Tenant-ID", tenantID)

	dr := new(DeletedResponse)
	resp, err := e.client.Do(ctx, req, dr)
	if err != nil {
		return nil, resp, err
	}
	return dr, resp, nil

}
