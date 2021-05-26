package sophoscentral

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestEndpointService_IsolationExclusionUpdate(t *testing.T) {

		tests := []struct {
			name           string
			ctx            context.Context
			path           string
			handlerFunc    func(http.ResponseWriter, *http.Request)
			client         *Client
			payload 	   IsolationExclusionItem
			want           *IsolationExclusionItem
			wantStatusCode int
			wantErr        bool
		}{
			{
				name: "valid delete",
				ctx: context.Background(),
				path: "/endpoint/v1/settings/exclusions/isolation/id",
				handlerFunc: func(w http.ResponseWriter, r *http.Request){
					testMethod(t,r,"PATCH")
					w.WriteHeader(200)
					_, _ = w.Write([]byte(`{
  "remotePorts": [
    39501,
    39501
  ],
  "localPorts": [
    5249,
    5249
  ],
  "comment": "comment",
  "id": "id",
  "type": "isolation",
  "direction": "inbound",
  "remoteAddresses": [
    "remoteAddresses",
    "remoteAddresses"
  ]
}`))
				},
				client: nil,
				payload: IsolationExclusionItem{

					LocalPorts:      []uint16{65535},
					RemotePorts:     []uint16{65535},
					Direction:       Inbound.ToPtr(),
					RemoteAddresses: []string{"string"},
					Comment:         String("string"),
				},
				want: &IsolationExclusionItem{
					Type:            String("isolation"),
					LocalPorts:      []uint16{5249,5249},
					RemotePorts:     []uint16{39501,39501},
					Direction:       Inbound.ToPtr(),
					RemoteAddresses: []string{"remoteAddresses", "remoteAddresses"},
					Comment:         String("comment"),
				},
				wantStatusCode: 200,
				wantErr: false,
			},
			{
				name: "valid delete",
				ctx: context.Background(),

				path: "/endpoint/v1/settings/exclusions/isolation/id",
				handlerFunc: func(w http.ResponseWriter, r *http.Request){
					testMethod(t,r,"PATCH")
					w.WriteHeader(500)
					_, _ = w.Write([]byte(``))
				},
				client: nil,
				payload: IsolationExclusionItem{

					LocalPorts:      []uint16{65535},
					RemotePorts:     []uint16{65535},
					Direction:       Inbound.ToPtr(),
					RemoteAddresses: []string{"string"},
					Comment:         String("string"),
				},				want: nil,
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
				got,  resp, err := client.Endpoints.IsolationExclusionUpdate(tt.ctx, "tenantid-guid", "", "id", tt.payload)
				if (err != nil) != tt.wantErr {
					t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if resp == nil {
					assert.NotNil(t, resp)
					return
				}
				assert.Equal(t, tt.want, got)
				assert.Equal(t, tt.wantStatusCode, resp.StatusCode)
			})
		}
	}
