package sophoscentral

import (
	"context"
	"fmt"
)

// IsolationExclusionUpdate - Updates an Isolation exclusion by ID.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/exclusions/isolation/{exclusionId}
func (e *EndpointService) IsolationExclusionUpdate(ctx context.Context, tenantID string,  tenantURL BaseURL, exclusionID string, exclusionItem IsolationExclusionItem) (*IsolationExclusionItem, *Response, error) {
	path := fmt.Sprintf("%ssettings/exclusions/isolation/%s", e.basePath, exclusionID)


	if err := verifyIsolationExclusionItem(exclusionItem);err!= nil{
		return nil, nil,err
	}


	req, err := e.client.NewRequest(ctx, "PATCH", path, &tenantURL,exclusionItem)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")

	iei := new(IsolationExclusionItem)
	resp, err := e.client.Do(ctx, req, iei)
	if err != nil {
		return nil,resp,  err
	}

	return iei, resp, nil

}
