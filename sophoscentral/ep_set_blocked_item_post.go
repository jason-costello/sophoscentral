package sophoscentral

import (
	"context"
	"errors"
	"fmt"
)

// BlockFromExonerationRequest is the payload sent when making a request to block an item
type BlockFromExonerationRequest struct {
	// Type - Property by which an item is allowed.
	// Allowed values:
	// sha256
	Type       string        `json:"type"`
	Properties struct {
		FileName string `json:"fileName"`
		Path     string `json:"path"`
		Sha256   string `json:"sha256"`
	} `json:"properties"`	// Comment indicating why the item should be allowed.
	Comment string `json:"comment"`
}



// BlockItemFromExoneration - Block an item from exoneration.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/blocked-items
func (e *EndpointService) BlockItemFromExoneration(ctx context.Context, tenantID string, tenantURL BaseURL, blockRequest BlockFromExonerationRequest) (*ExemptAllowedItemResponse, *Response, error) {
	path := fmt.Sprintf("%sendpoints/settings/allowed-items", e.basePath)

	req, err := e.client.NewRequest(ctx, "POST", path, &tenantURL, blockRequest)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("X-Tenant-ID", tenantID)
	blockResp := new(ExemptAllowedItemResponse)
	resp, err := e.client.Do(ctx, req, blockResp)
	if err != nil {
		return nil, resp, err
	}

	if resp.StatusCode != 201 {
		return nil, resp, errors.New(resp.Status)
	}
	return blockResp, resp, nil

}
