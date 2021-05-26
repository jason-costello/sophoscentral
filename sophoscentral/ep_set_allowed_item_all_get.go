package sophoscentral

import (
	"context"
	"fmt"
)

// ListAllowedItems gathers all endpoints for a tenant ID
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/allowed-items
func (e *EndpointService) ListAllowedItems(ctx context.Context, tenantID string, tenantURL BaseURL, opts PageByOffsetOptions) (*EPSettingItems, error) {
	path := fmt.Sprintf("%ssettings/allowed-items", e.basePath)

	return e.listEPSettingItems(ctx, path, tenantID, tenantURL, opts)
}
