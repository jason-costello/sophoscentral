package sophoscentral

import (
	"context"
	"fmt"
)




type ExclusionDirections string
const Inbound ExclusionDirections = "inbound"
const Outbound ExclusionDirections = "outbound"
const Both ExclusionDirections = "both"
type IntrusionPreventionExclusionItems struct{
	Items []IntrusionPreventionsExclusionItem `json:"items"`
	PagesByOffset `json:"pages,omitempty"`
}
func (e ExclusionDirections) ToPtr()*ExclusionDirections{
	return &e
}
type IntrusionPreventionsExclusionItem struct{
	// ID is exclusion ID - matches [a-f0-9]{64}
	ID *string `json:"id,omitempty"`
	// Type refers to exclusion type.
	// this value is always
	// intrusionPrevention
	Type *string `json:"type,omitempty"`
	// LocalPorts that are allowed
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

// IntrusionPreventionExclusionsList - Get all Intrusion Prevention exclusions.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/exclusions/intrusion-prevention
func (e *EndpointService) IntrusionPreventionExclusionsList(ctx context.Context, tenantID string,  tenantURL BaseURL, opts *ListByPageOffset) (*IntrusionPreventionExclusionItems, error) {
path := fmt.Sprintf("%ssettings/exclusions/intrusion-prevention", e.basePath)
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
ipe := new(IntrusionPreventionExclusionItems)
_, err = e.client.Do(ctx, req, ipe)
if err != nil {
return nil, err
}

return ipe, nil

}
