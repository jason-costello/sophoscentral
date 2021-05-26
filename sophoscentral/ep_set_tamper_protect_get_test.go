package sophoscentral

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestEndpointService_GlobalTamperProtectionStatus(t *testing.T) {

	tests := []struct {
		name           string
		ctx            context.Context
		path           string
		handlerFunc    func(http.ResponseWriter, *http.Request)
		client         *Client
		want           *GlobalTamperProtectionEnabled
		wantStatusCode int
		wantErr        bool
	}{
		{
			name: "valid - true returned",
			ctx:  context.Background(),
			path: "/endpoint/v1/endpoints/settings/tamper-protection",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(200)
				_, _ = w.Write([]byte(`{"enabled": true}`))
			},
			client:         nil,
			want:           &GlobalTamperProtectionEnabled{Enabled: true},
			wantStatusCode: 200,
			wantErr:        false,
		},

		{
			name: "valid - false returned",
			ctx:  context.Background(),
			path: "/endpoint/v1/endpoints/settings/tamper-protection",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(200)
				_, _ = w.Write([]byte(`{"enabled": false}`))
			},
			client:         nil,
			want:           &GlobalTamperProtectionEnabled{Enabled: false},
			wantStatusCode: 200,
			wantErr:        false,
		},
		{
			name: "invalid - 500 returned",
			ctx:  context.Background(),
			path: "/endpoint/v1/endpoints/settings/tamper-protection",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(500)
				_, _ = w.Write([]byte(`{
  "error": "can't get global setting'",
  "message": "I'm tired and can't be bothered",
  "correlationId": "59763C8E-B687-47D0-8F7B-88113425CE3B",
  "code": "US4c5",
  "createdAt": "2019-08-15T11:25:45.987Z",
  "requestId": "6DB1D8AC-1BFA-448B-8439-5486E6D25A74",
  "docUrl": "http://docs.sophos.com/central/api-docs/en-us/errors.html#lazylazylazy"
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
			got, res, err := client.Endpoints.GlobalTamperProtectionStatus(tt.ctx, "tenantid-guid", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("TamperProtection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, res)
			assert.Equal(t, tt.want, got)
			if res != nil {
				assert.Equal(t, tt.wantStatusCode, res.StatusCode)
			}
		})
	}
}
