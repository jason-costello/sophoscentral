package sophoscentral

import (
	"context"
	"fmt"
)

type IsolationExclusionItems struct{
	Item []IsolationExclusionItem `json:"items"`
	Pages PagesByOffset `json:"pages"`
}

type IsolationExclusionItem struct{
	ID *string `json:"idomitempty"`
	// LocalPorts that are allowed
	// Type refers to exclusion type.
	// this value is always
	// isolation
	Type *string `json:"type,omitempty"`
	LocalPorts []uint16 `json:"localPorts,omitempty"`
	// RemotePorts that are allowed
	RemotePorts []uint16 `json:"remotePorts,omitempty"`
	// Direction - allowed values are
	// inbound, outbound, both
	Direction *ExclusionDirections `json:"direction,omitempty"`
	// RemoteAddress - Remote addresses to exempt from Intrusion Prevention checks.
	RemoteAddresses []string `json:"remoteAddresses,omitempty"`
	// Comment - exclusion comment
	Comment *string `json:"comment,omitempty"`

}
// IsolationExclusionsList - Get all Intrusion Prevention exclusions.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/exclusions/isolation
func (e *EndpointService) IsolationExclusionsList(ctx context.Context, tenantID string,  tenantURL BaseURL, opts *ListByPageOffset) (*IsolationExclusionItems, error) {
	path := fmt.Sprintf("%ssettings/exclusions/isolation", e.basePath)
	var err error
	if path, err = addOptions(path, opts); err != nil{
		return nil, err
	}
	req, err := e.client.NewRequest(ctx, "GET", path, &tenantURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")
	e.client.Token.SetAuthHeader(req)
	ipe := new(IsolationExclusionItems)
	_, err = e.client.Do(ctx, req, ipe)
	if err != nil {
		return nil, err
	}

	return ipe, nil

}
