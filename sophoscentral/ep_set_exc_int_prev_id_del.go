package sophoscentral


import (
	"context"
	"fmt"
)

// IntrusionPreventionExclusionDelete - delete intrusion prevention exclusion by id
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/exclusions/intrusion-prevention/{exclusionId}
func (e *EndpointService) IntrusionPreventionExclusionDelete(ctx context.Context, tenantID string, tenantURL BaseURL, exclusionID string) (*DeletedResponse, *Response, error) {
	path := fmt.Sprintf("%ssettings/exclusions/intrusion-prevention/%s", e.basePath, exclusionID)


	req, err := e.client.NewRequest(ctx, "DELETE", path, &tenantURL, nil)
	if err != nil {
		return nil,nil, err
	}

	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")

	ipe := new(DeletedResponse)
	resp, err := e.client.Do(ctx, req, ipe)
	if err != nil {
		return nil, resp,err
	}

	return ipe,resp,  nil


}