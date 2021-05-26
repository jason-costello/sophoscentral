package sophoscentral

import (
	"context"
	"fmt"
)

type TamperProtectionSettings struct {
	Password          *string               `json:"password,omitempty"`
	PreviousPasswords []TPPreviousPasswords `json:"previousPasswords,omitempty"`
	Enabled           bool                  `json:"enabled,omitempty"`
}
type TPPreviousPasswords struct {
	Password      *string `json:"password,omitempty"`
	InvalidatedAt *string `json:"invalidatedAt,omitempty"`
}



// TamperProtection fetches the TamperProtection settings for a specific endpoint
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/endpoints/{endpointId}/tamper-protection
func (e *EndpointService) TamperProtection(ctx context.Context, tenantID string, tenantURL BaseURL, endpointID string) (*TamperProtectionSettings, *Response, error) {

	path := fmt.Sprintf("%sendpoints/%s/tamper-protection", e.basePath, endpointID)

	req, err := e.client.NewRequest(ctx, "GET", path, &tenantURL, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("X-Tenant-ID", tenantID)

	tps := new(TamperProtectionSettings)
	resp, err := e.client.Do(ctx, req, tps)
	if err != nil {
		return nil, resp, err
	}

	return tps, resp, nil

}
