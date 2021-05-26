package sophoscentral

import (
	"context"
	"errors"
	"fmt"
)

// IntrusionPreventionExclusionAdd - Set Intrusion Prevention exclusion item.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/exclusions/intrusion-prevention
func (e *EndpointService) IntrusionPreventionExclusionAdd(ctx context.Context, tenantID string,  tenantURL BaseURL, i IntrusionPreventionsExclusionItem) (*IntrusionPreventionsExclusionItem, error) {
path := fmt.Sprintf("%ssettings/exclusions/intrusion-prevention", e.basePath)
/*
      https://developer.sophos.com/docs/endpoint-v1/1/routes/settings/exclusions/intrusion-prevention/post
   	constraints on each field noted at above link
*/

if err := verifyExclusionItem(i);err!= nil{
return nil, err
}



req, err := e.client.NewRequest(ctx, "POST", path, &tenantURL, i)
if err != nil {
return nil, err
}
req.Header.Set("X-Tenant-ID", tenantID)
req.Header.Set("Content-Type", "application/json")
e.client.Token.SetAuthHeader(req)
ipe := new(IntrusionPreventionsExclusionItem)
_, err = e.client.Do(ctx, req, ipe)
if err != nil {
return nil, err
}

return ipe, nil

}

func verifyExclusionItem(i IntrusionPreventionsExclusionItem) error{


	if len(i.LocalPorts) != 1{
		return  errors.New("must contain exactly 1 local port")

	}
	if len(i.RemotePorts) != 1 {
		return  errors.New("must contain exactly 1 remote port")
	}

	if *i.Direction != Inbound && *i.Direction != Outbound && *i.Direction != Both{
		return  errors.New("only inbound, outbound, both allowed as direction value")
	}

	if len(i.RemoteAddresses) != 1 {
		return  errors.New("must contain exactly 1 remote address")
	}
	return nil
}
