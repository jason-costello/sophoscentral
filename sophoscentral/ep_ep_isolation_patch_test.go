package sophoscentral

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)


func TestEndpointService_UpdateIsolationSetting(t *testing.T) {

	tests := []struct {
		name           string
		ctx            context.Context
		path           string
		handlerFunc    func(http.ResponseWriter, *http.Request)
		client         *Client
		payload 	   IsolationSettingsUpdate
		want           *IsolationSettings
		wantStatusCode int
		wantErr        bool
	}{
		{
			name: "valid delete",
			ctx: context.Background(),
			path: "/endpoint/v1/endpoints/id/isolation",
			handlerFunc: func(w http.ResponseWriter, r *http.Request){
				testMethod(t,r,"PATCH")
				w.WriteHeader(200)
				_, _ = w.Write([]byte(`{"deleted": true}`))
			},
			client: nil,
			payload: IsolationSettingsUpdate{
				Enabled: Bool(true),
				Comment: String("some string"),
			},
			want: &IsolationSettings{},

			wantStatusCode: 200,
			wantErr: false,
		},
		{
			name: "valid delete",
			ctx: context.Background(),

			path: "/endpoint/v1/endpoints/id/isolation",
			handlerFunc: func(w http.ResponseWriter, r *http.Request){
				testMethod(t,r,"PATCH")
				w.WriteHeader(500)
				_, _ = w.Write([]byte(``))
			},
			client: nil,
			payload: IsolationSettingsUpdate{},
			want:nil,
			wantStatusCode: 500,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, mux, _, teardown := setup()
			defer teardown()
			tt.client = client
			mux.HandleFunc(tt.path, tt.handlerFunc)
			got,  res, err := client.Endpoints.UpdateIsolationSetting(tt.ctx, "tenantid-guid", "", "id", tt.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if res == nil {
				assert.NotNil(t, res)
				return
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantStatusCode, res.StatusCode)
		})
	}
}
