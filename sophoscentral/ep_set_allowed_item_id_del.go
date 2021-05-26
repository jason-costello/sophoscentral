package sophoscentral

import (
	"context"
	"errors"
	"fmt"
)





// DeleteAllowedItem an allowed item
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/allowed-items/{allowedItemId}
func (e *EndpointService) DeleteAllowedItem(ctx context.Context, tenantID string, tenantURL BaseURL, allowedItemID string) (*DeletedResponse, error) {
	// url path to call
	path := fmt.Sprintf("%ssettings/allowed-items/%s", e.basePath, allowedItemID)

	if allowedItemID == "" {
		return nil, errors.New("allowedItemID is empty")
	}

	req, err := e.client.NewRequest(ctx, "DELETE", path, &tenantURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Tenant-ID", tenantID)

	dr := new(DeletedResponse)
	_, err = e.client.Do(ctx, req, dr)
	if err != nil {
		return nil, err
	}
	return dr, nil

}
