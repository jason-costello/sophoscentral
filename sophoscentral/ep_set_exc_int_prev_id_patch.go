package sophoscentral

import (
	"context"
	"fmt"
)

// IntrusionPreventionExclusionUpdate - Get  Intrusion Prevention exclusion by id.
//https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/exclusions/intrusion-prevention/{exclusionId}
func (e *EndpointService) IntrusionPreventionExclusionUpdate(ctx context.Context, tenantID string,  tenantURL BaseURL, exclusionID string, exclusionItem IntrusionPreventionsExclusionItem) (*IntrusionPreventionsExclusionItem, error) {
	path := fmt.Sprintf("%ssettings/exclusions/intrusion-prevention/%s", e.basePath, exclusionID)


	if err := verifyExclusionItem(exclusionItem);err!= nil{
		return nil, err
	}


	req, err := e.client.NewRequest(ctx, "PATCH", path, &tenantURL,exclusionItem)
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
