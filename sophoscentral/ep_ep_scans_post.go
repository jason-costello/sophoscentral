package sophoscentral

import (
	"context"
	"fmt"
)


type EndpointScanResponseStatus string
const Requested EndpointScanResponseStatus = "requested"

type EndpointScanResponse struct{
	ID string `json:"id"`
	Status EndpointScanResponseStatus `json:"status"`
	RequestedAt string `json:"requestedAt"`
}
// EndpointScan - Sends a request to the specified endpoint to perform or configure a scan.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/endpoints/{endpointId}/scans
// post empty body to endpointID
func (e *EndpointService) EndpointScan(ctx context.Context, endpointID, tenantID string, tenantURL BaseURL) (*EndpointScanResponse, *Response, error) {
	path := fmt.Sprintf("%sendpoints/%s/scans", e.basePath, endpointID)

	req, err := e.client.NewRequest(ctx, "POST", path, &tenantURL, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("X-Tenant-ID", tenantID)

	tps := new(EndpointScanResponse)
	resp, err := e.client.Do(ctx, req, tps)
	if err != nil {
		return nil, resp, err
	}

	return tps, resp, nil

}
