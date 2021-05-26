package sophoscentral

import (
	"context"
	"fmt"
)

// GetAllowedItem returns one EPSettingItem for the requested id
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/allowed-items/{allowedItemId}
func (e *EndpointService) GetAllowedItem(ctx context.Context, tenantID string, tenantURL BaseURL, allowedItemID string) (*EPSettingItem, error) {
path := fmt.Sprintf("%ssettings/allowed-items/%s", e.basePath, allowedItemID)
return e.getEPSettingItem(ctx, tenantID, path, tenantURL)

}
