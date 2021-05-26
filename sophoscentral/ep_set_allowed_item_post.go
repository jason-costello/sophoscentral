package sophoscentral

import (
	"context"
	"errors"
	"fmt"
)

type ExemptAllowedItemRequest struct {
	// Type - Property by which an item is allowed.
	// Allowed values:
	// path, sha256, certificateSigner
	Type       string        `json:"type"`
	Properties struct {
		FileName string `json:"fileName"`
		Path     string `json:"path"`
		Sha256   string `json:"sha256"`
	} `json:"properties"`
	// Comment indicating why the item should be allowed.
	Comment string `json:"comment"`
	// OriginPersonID - Person associated with the endpoint where the item to be allowed was last seen.
	OriginPersonID string `json:"originPersonId"`
	// OriginEndpointID - the endpoint where the item to be allowed was last seen.
	OriginEndpointID string `json:"originEndpointId"`
}

type ExemptAllowedItemResponse struct {
	EPSettingItem

	// CreatedBy - User
	CreatedBy struct {
		// ID -  Unique ID for the user.
		ID string `json:"id"`
		// Name - Person's name.
		Name string `json:"name"`
	} `json:"createdBy:omitempty"`

	// OriginPerson - User
	OriginPerson struct {
		// ID - Unique ID for the user.
		ID string `json:"id"`
		//Name - Person's name.
		Name string `json:"name"`
	} `json:"originPerson:omitempty"`

}

// ExemptAllowedItem - Exempt an item from conviction.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/allowed-items
func (e *EndpointService) ExemptAllowedItem(ctx context.Context, tenantID string, tenantURL BaseURL, exemptionRequest ExemptAllowedItemRequest) (*ExemptAllowedItemResponse, *Response, error) {
	path := fmt.Sprintf("%sendpoints/settings/allowed-items", e.basePath)

	req, err := e.client.NewRequest(ctx, "POST", path, &tenantURL, exemptionRequest)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("X-Tenant-ID", tenantID)
	exemptionResp := new(ExemptAllowedItemResponse)
	resp, err := e.client.Do(ctx, req, exemptionResp)
	if err != nil {
		return nil, resp, err
	}

	if resp.StatusCode != 201 {
		return nil, resp, errors.New(resp.Status)
	}
	return exemptionResp, resp, nil

}
