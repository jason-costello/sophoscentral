package sophoscentral

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestEndpointService_SetTamperProtection(t *testing.T) {

	tests := []struct {
		name           string
		ctx            context.Context
		path           string
		handlerFunc    func(http.ResponseWriter, *http.Request)
		client         *Client
		want           *TamperProtectionSettings
		wantStatusCode int
		wantErr        bool
	}{
		{
			name: "valid - no error returned",
			ctx:  context.Background(),
			path: "/endpoint/v1/endpoints/endpointid-guid/tamper-protection",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "POST")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(201)
				_, _ = w.Write([]byte(`{
  "enabled": true,
  "regeneratePassword": true
}`))
			},
			client:         nil,
			want:           &TamperProtectionSettings{Password: (*string)(nil), PreviousPasswords: []TPPreviousPasswords(nil), Enabled: true},
			wantStatusCode: 201,
			wantErr:        false,
		},

		{
			name: "invalid - 500 error returned",
			ctx:  context.Background(),
			path: "/endpoint/v1/endpoints/endpointid-guid/tamper-protection",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "POST")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(500)
			},
			client:         nil,
			want:           nil,
			wantStatusCode: 500,
			wantErr:        true,
		},
		{
			name: "invalid - 404 error returned",
			ctx:  context.Background(),
			path: "/endpoint/v1/endpoints/invalidguidid/tamper-protection",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "POST")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(404)
			},
			client:         nil,
			want:           nil,
			wantStatusCode: 404,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, mux, _, teardown := setup()
			defer teardown()
			tt.client = client
			mux.HandleFunc(tt.path, tt.handlerFunc)
			got, res, err := client.Endpoints.TamperProtectionToggle(tt.ctx, "tenantid-guid", "", "endpointid-guid")
			if (err != nil) != tt.wantErr {
				t.Errorf("TamperProtection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, res)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantStatusCode, res.StatusCode)
		})
	}
}
