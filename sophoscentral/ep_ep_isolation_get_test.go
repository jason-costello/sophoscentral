package sophoscentral

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

	func TestEndpointService_GetIsolationSettings(t *testing.T) {

		tests := []struct {
			name           string
			ctx            context.Context
			path           string
			handlerFunc    func(http.ResponseWriter, *http.Request)
			client         *Client
			want           *Item
			wantStatusCode int
			wantErr        bool
		}{
			{
				name: "valid - no error returned",
				ctx:  context.Background(),
				path: "/endpoint/v1/endpoints/endpointid-guid",
				handlerFunc: func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{
	"id": "endpointid-guid",
	"type": "server",
	"tenant": {
		"id": "tenantid-guid"
	},
	"hostname": "myserverhostname",
	"health": {
		"overall": "suspicious",
		"threats": {
			"status": "good"
		},
		"services": {
			"status": "good",
			"serviceDetails": [{
					"name": "HitmanPro.Alert service",
					"status": "running"
				},
				{
					"name": "Sophos Anti-Virus",
					"status": "running"
				}
			]
		}
	},
	"os": {
		"isServer": true,
		"platform": "windows",
		"name": "Windows Server 2012 R2 Standard",
		"majorVersion": 6,
		"minorVersion": 3,
		"build": 9600
	},
	"ipv4Addresses": [
		"10.127.181.8",
		"110.205.25.8"
	],
	"ipv6Addresses": [
		"fc80::25cc:2ae5:d417:54a0"
	],
	"macAddresses": [
		"00:5a:52:83:54:AB",
		"00:5a:53:8C:28:9D"
	],
	"group": {
		"name": "Bolton's Skinning Service"
	},
	"associatedPerson": {
		"viaLogin": "MyCompy\\\\tywinthelion"
	},
	"tamperProtectionEnabled": true,
	"assignedProducts": [{
			"code": "coreAgent",
			"version": "2.10.8",
			"status": "installed"
		},
		{
			"code": "endpointProtection",
			"version": "10.8.9.2",
			"status": "installed"
		},
		{
			"code": "interceptX",
			"version": "2.10.18",
			"status": "installed"
		}
	],
	"lastSeenAt": "2020-12-28T10:46:57.393Z",
	"lockdown": {
		"status": "notInstalled",
		"updateStatus": "notInstalled"
	}
}`))
				},
				client: nil,
				want: &Item{
					ID:       String("endpointid-guid"),
					Type:     String("server"),
					Tenant:   &TenantEP{ID: "tenantid-guid"},
					Hostname: String("myserverhostname"),
					Health: &Health{
						Overall: String("suspicious"),
						Threats: Threats{Status: String("good")},
						Services: Services{
							Status: "good",
							ServiceDetails: []ServiceDetail{
								{
									Name:   String("HitmanPro.Alert service"),
									Status: String("running"),
								},
								{
									Name:   String("Sophos Anti-Virus"),
									Status: String("running"),
								},
							},
						},
					},
					OS: &OS{
						IsServer:     true,
						Platform:     String("windows"),
						Name:         "Windows Server 2012 R2 Standard",
						MajorVersion: 6,
						MinorVersion: 3,
						Build:        Int64(9600),
					},
					Ipv4Addresses: []string{"10.127.181.8", "110.205.25.8"},
					Ipv6Addresses: []string{"fc80::25cc:2ae5:d417:54a0"},
					MACAddresses:  []string{"00:5a:52:83:54:AB", "00:5a:53:8C:28:9D"},
					Group:         &Group{Name: "Bolton's Skinning Service"},
					AssociatedPerson: &AssociatedPerson{
						Name:     nil,
						ViaLogin: `MyCompy\\tywinthelion`,
						ID:       nil,
					},
					TamperProtectionEnabled: Bool(true),
					AssignedProducts: []AssignedProduct{
						{
							Code:    "coreAgent",
							Version: "2.10.8",
							Status:  String("installed"),
						},
						{
							Code:    "endpointProtection",
							Version: "10.8.9.2",
							Status:  String("installed"),
						},
						{
							Code:    "interceptX",
							Version: "2.10.18",
							Status:  String("installed"),
						},
					},

					LastSeenAt: String("2020-12-28T10:46:57.393Z"),
					Encryption: nil,
					Lockdown: &Lockdown{
						Status:       "notInstalled",
						UpdateStatus: "notInstalled",
					},
				},

				wantStatusCode: 200,
				wantErr:        false,
			},
			{
				name: "500 error returned",
				ctx:  context.Background(),
				path: "/endpoint/v1/endpoints/endpointid-guid",
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
			{
				name: "429 error returned",
				ctx:  context.Background(),
				path: "/endpoint/v1/endpoints/endpointid-guid",
				handlerFunc: func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(429)
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
				wantStatusCode: 429,
				wantErr:        true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				client, mux, _, teardown := setup()
				defer teardown()
				tt.client = client
				mux.HandleFunc(tt.path, tt.handlerFunc)
				got, res, err := client.Endpoints.Get(tt.ctx, "tenantid-guid", "", "endpointid-guid")
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
