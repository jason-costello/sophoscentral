package sophoscentral

import (
	"context"
	"errors"
	"fmt"
)
// TamperProtectionToggle Turns Tamper Protection on or off on an
// endpoint. Or generates a new Tamper Protection password.
// Note that Tamper Protection can be turned on for an endpoint
// only if it has also been turned on globally.

func (e *EndpointService) TamperProtectionToggle(ctx context.Context, tenantID string, tenantURL BaseURL, endpointID string) (*TamperProtectionSettings, *Response, error) {

	path := fmt.Sprintf("%sendpoints/%s/tamper-protection", e.basePath, endpointID)

	req, err := e.client.NewRequest(ctx, "POST", path, &tenantURL, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("X-Tenant-ID", tenantID)

	tps := new(TamperProtectionSettings)
	resp, err := e.client.Do(ctx, req, tps)
	if err != nil {
		return nil, resp, err
	}

	if resp.StatusCode != 201 {
		return nil, resp, errors.New(resp.Status)
	}
	return tps, resp, nil

}
