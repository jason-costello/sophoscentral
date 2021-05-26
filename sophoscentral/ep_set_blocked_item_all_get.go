package sophoscentral

import "fmt"
import "context"


// ListBlockedItems Get all blocked items.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/blocked-items
func (e *EndpointService) ListBlockedItems(ctx context.Context, tenantID string, tenantURL BaseURL, opts PageByOffsetOptions) (*EPSettingItems, error) {
path := fmt.Sprintf("%ssettings/blocked-items", e.basePath)
return e.listEPSettingItems(ctx, path, tenantID, tenantURL, opts)
}
