package sophoscentral

import (
	"context"
	"errors"
	"fmt"
)

// DeleteBlockedItem  deletes a blocked item by item
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/blocked-items/{blockedItemId}
func (e *EndpointService) DeleteBlockedItem(ctx context.Context, tenantID string, tenantURL BaseURL, blockedItemID string) (*DeletedResponse, *Response, error) {
	// url path to call
	path := fmt.Sprintf("%ssettings/blocked-items/%s", e.basePath, blockedItemID)

	if blockedItemID == "" {
		return nil, nil,errors.New("blockedItemID is empty")
	}

	req, err := e.client.NewRequest(ctx, "DELETE", path, &tenantURL, nil)
	if err != nil {
		return nil,nil, err
	}

	req.Header.Set("X-Tenant-ID", tenantID)

	dr := new(DeletedResponse)
	resp, err := e.client.Do(ctx, req, dr)
	if err != nil {
		return nil,resp, err
	}
	return dr, resp, nil

}
