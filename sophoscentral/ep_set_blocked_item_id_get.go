package sophoscentral

import (
	"context"
	"fmt"
)

// GetBlockedItem Get a blocked item by ID.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/blocked-items/{blockedItemId}
func (e *EndpointService) GetBlockedItem(ctx context.Context, tenantID string, tenantURL BaseURL, blockedItemID string) (*EPSettingItem, error) {
	path := fmt.Sprintf("%ssettings/blocked-items/%s", e.basePath, blockedItemID)

	return e.getEPSettingItem(ctx, tenantID, path, tenantURL)

}
