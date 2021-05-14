package sophoscentral

import (
	"context"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestPartnerService_GetTenants(t *testing.T) {
	type fields struct {
		ID      uuid.UUID
		BaseURL string
	}
	type args struct {
		ctx   context.Context
		token *oauth2.Token
		hc    *http.Client
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     TenantResponse
		respCode string
		wantErr  bool
	}{
		{
			name: "one",
			fields: fields{
				ID:      uuid.UUID{},
				BaseURL: "",
			},
			args: args{
				ctx: context.Background(),
				token: &oauth2.Token{
					AccessToken:  "access token",
					TokenType:    "token",
					RefreshToken: "refresher",
					Expiry:       time.Now().Add(24 * 7 * 52 * 42 * time.Hour),
				},
				hc: httpClientWithRoundTripper(200, `{
  "items": [
    {
      "id": "03b43abe-4f41-4734-b6d6-70b2fbdc2504",
      "name": "Father Oz",
      "dataGeography": "US",
      "dataRegion": "us03",
      "billingType": "usage",
      "partner": {
		"id": "d2ba043d-7fcd-4158-a861-1ec2c01f3d14"		
      },
      "apiHost": "https://api-us03.central.sophos.com",
      "status": "active"
    }
  ],
  "pages": {
    "current": 1,
    "size": 50,
    "maxSize": 100
  }
}`,
)},
			want: TenantResponse{
				Items: []TenantsResponseItem{
					{
						ID:            "03b43abe-4f41-4734-b6d6-70b2fbdc2504",
						Name:          "Father Oz",
						DataGeography: "US",
						DataRegion:    "us03",
						BillingType:   "usage",
						Partner:       TenantsResponsePartner{ID:"d2ba043d-7fcd-4158-a861-1ec2c01f3d14" },
						ApiHost:       "https://api-us03.central.sophos.com",
						Status:        "active",
					},
				},
				Pages: TenantsResponsePages{
					Current: 1,
					Size:    50,
					Maxsize: 100,
				},
			},
			respCode: "200",
			wantErr: false,
		},
		{
			name: "500 error",
			fields: fields{
				ID:      uuid.UUID{},
				BaseURL: "",
			},
			args: args{
				ctx: context.Background(),
				token: &oauth2.Token{
					AccessToken:  "access token",
					TokenType:    "token",
					RefreshToken: "refresher",
					Expiry:       time.Now().Add(24 * 7 * 52 * 42 * time.Hour),
				},
				hc: httpClientWithRoundTripper(500, `{
  "error": "i am your error",
  "message": "i am your error message",
	"correlation_id": "i should be a guid",
	"code": "500",
	"created_at": "2021-05-04T12:00:00.000UTC",
	"request_id": "I should be a guid too",
	"doc_url": "rtfm"
}`,
				)},
				want: TenantResponse{},
respCode:             "500",
			wantErr:  true,
		},

	}

		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PartnerService{
				ID:      tt.fields.ID,
				BaseURL: tt.fields.BaseURL,
			}

			got, err := p.GetTenants(tt.args.ctx, tt.args.token, tt.args.hc)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTenants() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTenants() got = %v, want %v", got, tt.want)
			}


		})
	}
}
