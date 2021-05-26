package sophoscentral

import (
	"context"
	"errors"
	"fmt"
)

// ExploitMitigationApplicationDelete -Get detected exploits and the number of each detected exploit.
// page by offset pagination
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/exploit-mitigation/detected-exploits
func (e *EndpointService) ExploitMitigationApplicationDelete(ctx context.Context, tenantID, exploitMitigationAppId string,  tenantURL BaseURL) (*DeletedResponse, *Response, error) {

	path := fmt.Sprintf("%ssettings/exploit-mitigation/applications/%s", e.basePath, exploitMitigationAppId)

	if exploitMitigationAppId == ""{
		return nil, nil, errors.New("exploitMitigationAppId is empty")
	}

	req, err := e.client.NewRequest(ctx, "DELETE", path, &tenantURL, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("X-Tenant-ID", tenantID)

	sei := new(DeletedResponse)

	resp, err := e.client.Do(ctx, req, sei)
	if err != nil {
		return nil, resp,err
	}

	return sei, resp, nil

}
