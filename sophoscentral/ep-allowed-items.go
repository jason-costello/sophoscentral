package sophoscentral

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type PageByOffsetOptions struct {
	ListByPageOffset
}
type AllowedItems struct {
	Pages *PageByOffset `json:"pages,omitempty"`
	Items []AllowedItem `json:"items,omitempty"`
}
type PageByOffset struct {
	Current *int `json:"current,omitempty"`
	Total   *int `json:"total,omitempty"`
	Size    *int `json:"size,omitempty"`
	Maxsize *int `json:"maxSize,omitempty"`
	Items   *int `json:"items,omitempty"`
}
type AllowedItem struct {
	CreatedAt      *time.Time             `json:"createdAt,omitempty"`
	CreatedBy      *AICreatedBy           `json:"createdBy,omitempty"`
	OriginEndpoint *AIOriginEndpoint      `json:"originEndpoint,omitempty"`
	Comment        *string                `json:"comment,omitempty"`
	ID             *string                `json:"id,omitempty"`
	Type           *string                `json:"type,omitempty"`
	OriginPerson   interface{}            `json:"originPerson,omitempty"`
	Properties     *AllowedItemProperties `json:"properties,omitempty"`
	UpdatedAt      *time.Time             `json:"updatedAt,omitempty"`
}
type AllowedItemProperties struct {
	CertificateSigner *string `json:"certificateSigner,omitempty"`
	Path              *string `json:"path,omitempty"`
	FileName          *string `json:"fileName,omitempty"`
	Sha256            *string `json:"sha256,omitempty"`
}
type AIOriginEndpoint struct {
	ID *string `json:"id,omitempty"`
}
type AICreatedBy struct {
	Name *string `json:"name,omitempty"`
	ID   *string `json:"id,omitempty"`
}

type ExemptAllowedItemRequest struct{
	// Type - Property by which an item is allowed.
	// Allowed values:
	// path, sha256, certificateSigner
	Type string `json:"type"`
	Properties []EAIProperty `json:"properties"`
	// Comment indicating why the item should be allowed.
	Comment string `json:"comment"`
	// OriginPersonID - Person associated with the endpoint where the item to be allowed was last seen.
	OriginPersonID string `json:"originPersonId"`
	// OriginEndpointID - the endpoint where the item to be allowed was last seen.
	OriginEndpointID string `json:"originEndpointId"`
}

type EAIProperty struct{
	// FileName for application
	FileName string `json:"fileName"`
	// Path for the application.
	Path string `json:"path"`
	// SHA256 value for the application.
	SHA256 string `json:"sha256"`
	// CertificateSigner - Value saved for the certificateSigner.
	CertificateSigner string `json:"certificateSigner"`

}

type ExemptAllowedItemResponse struct{
	// ID  - Unique ID for the allowed application.
	ID *string `json:"id,omitempty"`
	// CreatedAt - Date and time (UTC) when the allowed application was created.
	CreatedAt *string `json:"createdAt,omitempty"`
	// UpdatedAt - Date and time (UTC) when the allowed application was updated.
	UpdatedAt *string `json:"updatedAt,omitempty"`
	// Properties - Allowed item properties.
	Properties *ExemptAllowedItemRequest `json:"properties,omitempty"`
	// Comment indicating why the item was allowed.
	Comment string  `json:"comment:omitempty"`
	// Type - Property by which an item is allowed.
	// Allowed values:
	// path, sha256, certificateSigner
	Type string  `json:"type:omitempty"`

	// CreatedBy - User
	CreatedBy struct{
		// ID -  Unique ID for the user.
		ID string `json:"id"`
		// Name - Person's name.
		Name string `json:"name"`
	} `json:"createdBy:omitempty"`

	// OriginPerson - User
	OriginPerson struct{
		// ID - Unique ID for the user.
		ID string `json:"id"`
		//Name - Person's name.
		Name string `json:"name"`
	} `json:"originPerson:omitempty"`

	// OriginEndpoint - Represents a referenced object.
	OriginEndpoint struct {
		// ID of the referenced object.
		ID string `json:"id"`
	} `json:"originEndpoint"`
}
type DeletedResponse struct{
	Deleted bool `json:"deleted"`
}

// ListAllowedItems gathers all endpoints for a tenant ID
// https://api-{region}.central.sophos.com/endpoint/v1/endpoints
func (e *EndpointService) ListAllowedItems(ctx context.Context, tenantID string, tenantURL BaseURL, opts PageByOffsetOptions) (*AllowedItems, error) {
	path := fmt.Sprintf("%ssettings/allowed-items", e.basePath)

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

	ogAI := new(AllowedItems)
	_, err = e.client.Do(ctx, req, ogAI)
	if err != nil {
		return nil, err
	}

return ogAI, nil


}

// ExemptAllowedItem - Exempt an item from conviction.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/allowed-items
func (e *EndpointService) ExemptAllowedItem(ctx context.Context, tenantID string, tenantURL BaseURL, exemptionRequest ExemptAllowedItemRequest) (*ExemptAllowedItemResponse, *Response, error) {
	path := fmt.Sprintf("%sendpoints/settings/allowed-items", e.basePath)

	req, err := e.client.NewRequest("POST", path, &tenantURL, exemptionRequest)
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

// GetAllowedItem returns one AllowedItem for the requested id
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/allowed-items/{allowedItemId}
func (e *EndpointService) GetAllowedItem(ctx context.Context, tenantID string, tenantURL BaseURL, allowedItemID string) (*AllowedItem, error) {
	path := fmt.Sprintf("%ssettings/allowed-items/%s", e.basePath, allowedItemID)


	req, err := e.client.NewRequest("GET", path, &tenantURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")
	e.client.Token.SetAuthHeader(req)

	ai := new(AllowedItem)
	_, err = e.client.Do(ctx, req, ai)
	if err != nil {
		return nil, err
	}

	return ai, nil


}

type AllowedItemPatchReq struct{
	Comment string `json:"comment"`
}

// UpdateAllowedItem by id
func (e *EndpointService) UpdateAllowedItem(ctx context.Context, tenantID string, tenantURL BaseURL, allowedItemID string) (*AllowedItem, error) {
	path := fmt.Sprintf("%ssettings/allowed-items/%s", e.basePath, allowedItemID)
	aip := new(AllowedItemPatchReq)
	req, err := e.client.NewRequest("PATCH", path, &tenantURL, aip)
	if err != nil {
		return nil,  err
	}

	req.Header.Set("X-Tenant-ID", tenantID)
	updatedItem := new(AllowedItem)

	resp, err := e.client.Do(ctx, req,updatedItem)
	if err != nil {
		return  nil, err
	}

	if resp.StatusCode != 200 {
		return  nil, errors.New(resp.Status)
	}
	return updatedItem,  nil
}


// DeleteAllowedItem an allowed item
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/allowed-items/{allowedItemId}
func (e *EndpointService) DeleteAllowedItem(ctx context.Context, tenantID string, tenantURL BaseURL, allowedItemID string) (*DeletedResponse,  error) {
	// url path to call
	path := fmt.Sprintf("%ssettings/allowed-items/%s", e.basePath, allowedItemID)

	if allowedItemID == "" {
		return nil,  errors.New("allowedItemID is empty")
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
