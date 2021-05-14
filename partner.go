package sophoscentral

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"net/http"
)

type PartnerService struct {
	ID uuid.UUID
	BaseURL string
}


func (p *PartnerService) GetTenants(ctx context.Context, token *oauth2.Token, hc *http.Client) (TenantResponse, error){

url := "https://api.central.sophos.com/partner/v1/tenants"

req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
if err != nil{
	return TenantResponse{},err
}
token.SetAuthHeader(req)
req.Header.Set("X-Partner-ID", p.ID.String())
req.Header.Set("Content-Type", "application/json")

	b, err := MakeRequest(hc, req)
	if err != nil {
		return   TenantResponse{}, fmt.Errorf("%s: %w", ErrHttpDo, err)
	}

	return UnmarshalTenantResponse(b)

}

type TenantResponse struct {
	Items []TenantsResponseItem `json:"items"`
	Pages TenantsResponsePages `json:"pages"`
}

type TenantsResponseItem struct{
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	DataGeography DataGeography          `json:"dataGeography"`
	DataRegion    DataRegion                 `json:"dataRegion"`
	BillingType   BillingType             `json:"billingType"`
	Partner       TenantsResponsePartner `json:"partner"`
	ApiHost       string                 `json:"apiHost"`
	Status        string                 `json:"status"`
}

type TenantsResponsePartner struct{
	ID string `json:"id"`
}

type TenantsResponsePages struct{
	Current int `json:"current"`
	Size    int `json:"size"`
	Maxsize int `json:"maxSize"`
}
type DataGeography string
const (
	USGeo DataGeography = "US"
	IEGeo DataGeography = "IE"
	DEGeo DataGeography = "DE"
)
type DataRegion string
const(
	EU01 DataRegion = "eu01"
	EU02 DataRegion = "eu02"
	US01 DataRegion = "us01"
	US02 DataRegion = "us02"
	US03 DataRegion = "us03"
)
type BillingType string
const(
	Term BillingType = "term"
	Trial BillingType = "trial"
	Usage BillingType = "usage"
	)

func UnmarshalTenantResponse(data []byte) (TenantResponse, error) {
	var r TenantResponse
	err := json.Unmarshal(data, &r)
	if err != nil{
		return TenantResponse{}, fmt.Errorf("%s: %w", ErrUnmarshalFailed, err)
	}
	return r, err
}
