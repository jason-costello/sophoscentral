package sophoscentral

import (
	"context"
	"errors"
	"fmt"
)

type ExploitMitigationNewAppRequest struct{
	Paths []string `json:"paths"`
}

// ExploitMitigationApplicationAdd - Exclude a set of file paths from Exploit Mitigation.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/exploit-mitigation/applications
func (e *EndpointService) ExploitMitigationApplicationAdd(ctx context.Context, tenantID string,
	tenantURL BaseURL, item ExploitMitigationNewAppRequest) (*ExploitMitigationAppItem, error) {
	path := fmt.Sprintf("%ssettings/exploit-mitigation/applications", e.basePath)

 	if err := verifyExploitMitigationNewAppRequest(item); err != nil{
		return nil, err
	}


	req, err := e.client.NewRequest(ctx, "POST", path, &tenantURL, item)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")
	e.client.Token.SetAuthHeader(req)
	sei := new(ExploitMitigationAppItem)
	_, err = e.client.Do(ctx, req, sei)
	if err != nil {
		return nil, err
	}

	return sei, nil

}


func verifyExploitMitigationNewAppRequest(i ExploitMitigationNewAppRequest)error{


	if i.Paths == nil{
		return errors.New("nil paths")
	}

	for x := range i.Paths{
		if x > 100{
			return errors.New("only 100 paths allowed")
		}

		if len(i.Paths[x]) < 1 || len(i.Paths[x]) > 260{
			return errors.New(fmt.Sprintf("path index %d must be <=1 length <= 260", x))
		}

	}

	return nil
}
