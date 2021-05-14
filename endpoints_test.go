package sophoscentral

import (
	"context"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestClient_GetEndpoints(t *testing.T) {
	type fields struct {
		ctx          context.Context
		logger       *logrus.Logger
		token        *oauth2.Token
		baseURL      *url.URL
		httpClient   *http.Client
		Partner      *PartnerService
		Organization *OrganizationService
		Tenant       *TenantService
	}
	type args struct {
		ctx         context.Context
		geoURL      string
		queryParams map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Endpoints
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				ctx:          tt.fields.ctx,
				logger:       tt.fields.logger,
				token:        tt.fields.token,
				baseURL:      tt.fields.baseURL,
				httpClient:   tt.fields.httpClient,
				Partner:      tt.fields.Partner,
				Organization: tt.fields.Organization,
				Tenant:       tt.fields.Tenant,
			}
			got, err := c.GetEndpoints(tt.args.ctx, tt.args.geoURL, tt.args.queryParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEndpoints() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEndpoints() got = %v, want %v", got, tt.want)
			}
		})
	}
}
