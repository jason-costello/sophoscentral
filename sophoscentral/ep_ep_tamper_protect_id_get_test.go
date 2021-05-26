package sophoscentral

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestEndpointService_TamperProtection(t *testing.T) {

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
				testMethod(t, r, "GET")
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{
					  "password": "password",
					  "previousPasswords": [
						{
						  "password": "password",
						  "invalidatedAt": "invalidatedAt"
						},
						{
						  "password": "password",
						  "invalidatedAt": "invalidatedAt"
						}
					  ],
					  "enabled": true
					}`))
			},
			client: nil,
			want: &TamperProtectionSettings{
				Password: String("password"),
				PreviousPasswords: []TPPreviousPasswords{
					{
						Password:      String("password"),
						InvalidatedAt: String("invalidatedAt"),
					},
					{
						Password:      String("password"),
						InvalidatedAt: String("invalidatedAt"),
					},
				},
				Enabled: true,
			},
			wantStatusCode: 200,
			wantErr:        false,
		},
		{
			name: "nil client - error returned",
			ctx:  context.Background(),
			path: "/endpoint/v1/endpoints/endpointid-guid/tamper-protection",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(500)
				_, _ = w.Write([]byte(`{
				"error": "string",
  "message": "string",
  "correlationId": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "code": "string",
  "createdAt": "string",
  "requestId": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "docUrl": "string"
}`))
			},
			client:         nil,
			want:           nil,
			wantStatusCode: 500,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, mux, _, teardown := setup()
			defer teardown()
			tt.client = client
			mux.HandleFunc(tt.path, tt.handlerFunc)
			got, res, err := client.Endpoints.TamperProtectionGet(tt.ctx, "tenantid-guid", "", "endpointid-guid")
			if (err != nil) != tt.wantErr {
				t.Errorf("TamperProtection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if res == nil{
				assert.NotNil(t, res)
				return
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantStatusCode, res.StatusCode)
		})
	}
}
