package sophoscentral

import (
	"context"
	"errors"
	"fmt"
)

type AllowedItemPatchReq struct {
	Comment string `json:"comment"`
}

// UpdateAllowedItem by id
func (e *EndpointService) UpdateAllowedItem(ctx context.Context, tenantID string, tenantURL BaseURL, allowedItemID string, aip AllowedItemPatchReq) (*EPSettingItem, *Response, error) {
	path := fmt.Sprintf("%ssettings/allowed-items/%s", e.basePath, allowedItemID)
	req, err := e.client.NewRequest(ctx, "PATCH", path, &tenantURL, aip)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("X-Tenant-ID", tenantID)
	updatedItem := new(EPSettingItem)

	resp, err := e.client.Do(ctx, req, updatedItem)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != 200 {
		return nil,resp, errors.New(resp.Status)
	}
	return updatedItem, resp, nil
}
