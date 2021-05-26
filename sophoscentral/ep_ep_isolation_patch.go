package sophoscentral

import (
	"context"
	"fmt"
)

// UpdateIsolationSetting for an endpoint.
func (e *EndpointService) UpdateIsolationSetting(ctx context.Context, tenantID string, tenantURL BaseURL, endpointID string, update IsolationSettingsUpdate) (*IsolationSettings, *Response, error) {

path := fmt.Sprintf("%sendpoints/%s/isolation", e.basePath, endpointID)

req, err := e.client.NewRequest(ctx, "PATCH", path, &tenantURL, update)
if err != nil {
return nil, nil, err
}
req.Header.Set("X-Tenant-ID", tenantID)

tps := new(IsolationSettings)
resp, err := e.client.Do(ctx, req, tps)
if err != nil {
return nil, resp, err
}

return tps, resp, nil

}
