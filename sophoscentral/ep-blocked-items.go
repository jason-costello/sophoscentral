package sophoscentral

import (
	"context"
	"errors"
	"fmt"
	"time"
)


type BlockedItems struct {
	Pages *PageByOffset `json:"pages,omitempty"`
	Items []AllowedItem `json:"items,omitempty"`
}

type BlockedItem struct {
	CreatedAt      *time.Time             `json:"createdAt,omitempty"`
	CreatedBy      *BICreatedBy           `json:"createdBy,omitempty"`
	OriginEndpoint *BIOriginEndpoint      `json:"originEndpoint,omitempty"`
	Comment        *string                `json:"comment,omitempty"`
	ID             *string                `json:"id,omitempty"`
	Type           *string                `json:"type,omitempty"`
	OriginPerson   interface{}            `json:"originPerson,omitempty"`
	Properties     *AllowedItemProperties `json:"properties,omitempty"`
	UpdatedAt      *time.Time             `json:"updatedAt,omitempty"`
}
type BlockedItemProperties struct {
	CertificateSigner *string `json:"certificateSigner,omitempty"`
	Path              *string `json:"path,omitempty"`
	FileName          *string `json:"fileName,omitempty"`
	Sha256            *string `json:"sha256,omitempty"`
}
type BIOriginEndpoint struct {
	ID *string `json:"id,omitempty"`
}
type BICreatedBy struct {
	Name *string `json:"name,omitempty"`
	ID   *string `json:"id,omitempty"`
}


type BlockFromExonerationRequest struct{
	// Type - Property by which an item is allowed.
	// Allowed values:
	// sha256
	Type string `json:"type"`
	Properties []EAIProperty `json:"properties"`
	// Comment indicating why the item should be allowed.
	Comment string `json:"comment"`
}

// ListBlockedItems Get all blocked items.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/blocked-items
func (e *EndpointService) ListBlockedItems(ctx context.Context, tenantID string, tenantURL BaseURL, opts PageByOffsetOptions) (*BlockedItems, error) {
	path := fmt.Sprintf("%ssettings/blocked-items", e.basePath)

	var maxSize = 100
	// GetAllowedItems for an endpoint.
	if opts.PageSize == 0{
		opts.PageSize = maxSize
	}

	path, err := addOptions(path, opts)
	if err != nil {
		return nil, err
	}
	req, err := e.client.NewRequest("GET", path, &tenantURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")
	e.client.Token.SetAuthHeader(req)

	ogAI := new(BlockedItems)
	_, err = e.client.Do(ctx, req, ogAI)
	if err != nil {
		return nil, err
	}

	return ogAI, nil


}

// BlockItemFromExoneration - Block an item from exoneration.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/blocked-items
func (e *EndpointService) BlockItemFromExoneration(ctx context.Context, tenantID string, tenantURL BaseURL, blockRequest BlockFromExonerationRequest) (*ExemptAllowedItemResponse, *Response, error) {
	path := fmt.Sprintf("%sendpoints/settings/allowed-items", e.basePath)

	req, err := e.client.NewRequest("POST", path, &tenantURL, blockRequest)
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

// GetBlockedItem Get a blocked item by ID.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/blocked-items/{blockedItemId}
func (e *EndpointService) GetBlockedItem(ctx context.Context, tenantID string, tenantURL BaseURL, blockedItemID string) (*BlockedItem, error) {
	path := fmt.Sprintf("%ssettings/blocked-items/%s", e.basePath, blockedItemID)


	req, err := e.client.NewRequest("GET", path, &tenantURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")
	e.client.Token.SetAuthHeader(req)

	ai := new(BlockedItem)
	_, err = e.client.Do(ctx, req, ai)
	if err != nil {
		return nil, err
	}

	return ai, nil


}


// DeleteBlockedItem  deletes a blocked item by item
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/blocked-items/{blockedItemId}
func (e *EndpointService) DeleteBlockedItem(ctx context.Context, tenantID string, tenantURL BaseURL, blockedItemID string) (*DeletedResponse,  error) {
	// url path to call
	path := fmt.Sprintf("%ssettings/allowed-items/%s", e.basePath, blockedItemID)

	if blockedItemID == "" {
		return nil,  errors.New("blockedItemID is empty")
	}

	req, err := e.client.NewRequest("DELETE", path, &tenantURL, nil)
	if err != nil {
		return nil,  err
	}

	req.Header.Set("X-Tenant-ID", tenantID)

	dr := new(DeletedResponse)
	_, err = e.client.Do(ctx, req, dr)
	if err != nil {
		return nil,  err
	}
	return dr,  nil

}
