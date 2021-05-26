package sophoscentral

import (
	"context"
	"errors"
	"fmt"
)





// GlobalTamperProtectionEnabled is the response object for the call to GlobalTamperProtectionStatus
type GlobalTamperProtectionEnabled struct {
	// indicates that status of global tamper protection
	Enabled bool `json:"enabled,omitempty"`
}

// GlobalTamperProtectionStatus checks whether Tamper Protection is turned on globally.
// only returns 200 or 500
func (e *EndpointService) GlobalTamperProtectionStatus(ctx context.Context, tenantID string, tenantURL BaseURL) (*GlobalTamperProtectionEnabled, *Response, error) {
	path := fmt.Sprintf("%sendpoints/settings/tamper-protection", e.basePath)
	req, err := e.client.NewRequest(ctx, "GET", path, &tenantURL, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("X-Tenant-ID", tenantID)

	tps := new(GlobalTamperProtectionEnabled)

	resp, err := e.client.Do(ctx, req, tps)
	if err != nil {
		return nil, resp, err
	}

	if resp.StatusCode != 200 {
		return nil, resp, errors.New(resp.Status)
	}
	return tps, resp, nil
}
