// +build !integration

package sophoscentral

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestEndpointService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	if client == nil {
		t.Fatal("client == nil")
		return
	}
	mux.HandleFunc("/endpoint/v1/endpoints", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
	"items": [{
		"id": "endpointid-guid",
		"type": "server",
		"tenant": {
			"id": "3d7a50a6-aee1-4193-a178-689bcb86f750"
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

	}],
	"pages": {
		"fromKey": "1",
		"size": 1,
		"total": 1,
		"items": 1,
		"maxSize": 500
	}
}`))

	})

	assert.NotNil(t, client.Endpoints)

	ctx := context.Background()
	epl := EndpointListOptions{
		ListByFromKeyOptions: ListByFromKeyOptions{
			PageFromKey: "1",
			PageSize:    500,
			PageTotal:   true,
		},
	}

	got, _, err := client.Endpoints.List(ctx, "3d7a50a6-aee1-4193-a178-689bcb86f750", "", &Endpoints{}, epl)
	if err != nil {
		assert.NoErrorf(t, err, "EndpointService.Get returned error: %v", err)
	}
	assert.NotNil(t, got)

	want := &Endpoints{
		Items: []Item{
			{ID: String("endpointid-guid"),
				Type:     String("server"),
				Tenant:   &TenantEP{ID: "3d7a50a6-aee1-4193-a178-689bcb86f750"},
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
		},
		Pages: &PagesByFromKey{
			FromKey: String("1"),
			Size:    Int(1),
			Total:   Int(1),
			Items:   Int(1),
			MaxSize: Int(500),
		},
	}
	assert.Equal(t, want, got)

	//const methodName = "Get"
	//
	//testNewRequestAndDoFailure(t, methodName, httpClient, func() (*Response, error) {
	//	got, resp, err := httpClient.Endpoints.Get(ctx, "tenantid-guid", "url", "endpointid-guid")
	//	assert.Nilf(t, got, "testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)

	//	return resp, err
	//})

}
func TestEndpointService_Get(t *testing.T) {

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
				w.Write([]byte(`{
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
				w.Write([]byte(`{
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
				w.Write([]byte(`{
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
			assert.NotNil(t, res)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantStatusCode, res.StatusCode)
		})
	}
}

//func TestEndpointService_Get(t *testing.T) {
//	client, mux, _, teardown := setup()
//	defer teardown()
//
//	mux.HandleFunc("/endpoint/v1/endpoints/endpointid-guid", func(w http.ResponseWriter, r *http.Request) {
//		testMethod(t, r, "GET")
//		w.Header().Set("Content-Type", "application/json")
//		w.Write([]byte(`{
//	"id": "endpointid-guid",
//	"type": "server",
//	"tenant": {
//		"id": "tenantid-guid"
//	},
//	"hostname": "myserverhostname",
//	"health": {
//		"overall": "suspicious",
//		"threats": {
//			"status": "good"
//		},
//		"services": {
//			"status": "good",
//			"serviceDetails": [{
//					"name": "HitmanPro.Alert service",
//					"status": "running"
//				},
//				{
//					"name": "Sophos Anti-Virus",
//					"status": "running"
//				}
//			]
//		}
//	},
//	"os": {
//		"isServer": true,
//		"platform": "windows",
//		"name": "Windows Server 2012 R2 Standard",
//		"majorVersion": 6,
//		"minorVersion": 3,
//		"build": 9600
//	},
//	"ipv4Addresses": [
//		"10.127.181.8",
//		"110.205.25.8"
//	],
//	"ipv6Addresses": [
//		"fc80::25cc:2ae5:d417:54a0"
//	],
//	"macAddresses": [
//		"00:5a:52:83:54:AB",
//		"00:5a:53:8C:28:9D"
//	],
//	"group": {
//		"name": "Bolton's Skinning Service"
//	},
//	"associatedPerson": {
//		"viaLogin": "MyCompy\\\\tywinthelion"
//	},
//	"tamperProtectionEnabled": true,
//	"assignedProducts": [{
//			"code": "coreAgent",
//			"version": "2.10.8",
//			"status": "installed"
//		},
//		{
//			"code": "endpointProtection",
//			"version": "10.8.9.2",
//			"status": "installed"
//		},
//		{
//			"code": "interceptX",
//			"version": "2.10.18",
//			"status": "installed"
//		}
//	],
//	"lastSeenAt": "2020-12-28T10:46:57.393Z",
//	"lockdown": {
//		"status": "notInstalled",
//		"updateStatus": "notInstalled"
//	}
//}`))
//
//	})
//
//	assert.NotNil(t, client.Endpoints)
//
//	ctx := context.Background()
//
//	got, _, err := client.Endpoints.Get(ctx, "tenantid-guid", "", "endpointid-guid")
//	if err != nil {
//		assert.NoErrorf(t, err, "EndpointService.Get returned error: %v", err)
//	}
//	assert.NotNil(t, got)
//
//	want := &Item{
//		ID:       String("endpointid-guid"),
//		Type:     String("server"),
//		Tenant:   &TenantEP{ID: "tenantid-guid"},
//		Hostname: String("myserverhostname"),
//		Health: &Health{
//			Overall: String("suspicious"),
//			Threats: Threats{Status: String("good")},
//			Services: Services{
//				Status: "good",
//				ServiceDetails: []ServiceDetail{
//					{
//						Name:   String("HitmanPro.Alert service"),
//						Status: String("running"),
//					},
//					{
//						Name:   String("Sophos Anti-Virus"),
//						Status: String("running"),
//					},
//				},
//			},
//		},
//		OS: &OS{
//			IsServer:     true,
//			Platform:     String("windows"),
//			Name:         "Windows Server 2012 R2 Standard",
//			MajorVersion: 6,
//			MinorVersion: 3,
//			Build:        Int64(9600),
//		},
//		Ipv4Addresses: []string{"10.127.181.8", "110.205.25.8"},
//		Ipv6Addresses: []string{"fc80::25cc:2ae5:d417:54a0"},
//		MACAddresses:  []string{"00:5a:52:83:54:AB", "00:5a:53:8C:28:9D"},
//		Group:         &Group{Name: "Bolton's Skinning Service"},
//		AssociatedPerson: &AssociatedPerson{
//			Name:     nil,
//			ViaLogin: `MyCompy\\tywinthelion`,
//			ID:       nil,
//		},
//		TamperProtectionEnabled: Bool(true),
//		AssignedProducts: []AssignedProduct{
//			{
//				Code:    "coreAgent",
//				Version: "2.10.8",
//				Status:  String("installed"),
//			},
//			{
//				Code:    "endpointProtection",
//				Version: "10.8.9.2",
//				Status:  String("installed"),
//			},
//			{
//				Code:    "interceptX",
//				Version: "2.10.18",
//				Status:  String("installed"),
//			},
//		},
//
//		LastSeenAt: String("2020-12-28T10:46:57.393Z"),
//		Encryption: nil,
//		Lockdown: &Lockdown{
//			Status:       "notInstalled",
//			UpdateStatus: "notInstalled",
//		},
//	}
//	assert.Equal(t, want, got)
//
//	//const methodName = "Get"
//	//
//	//testNewRequestAndDoFailure(t, methodName, httpClient, func() (*Response, error) {
//	//	got, resp, err := httpClient.Endpoints.Get(ctx, "tenantid-guid", "url", "endpointid-guid")
//	//	assert.Nilf(t, got, "testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
//
//	//	return resp, err
//	//})
//
//}

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
				w.Write([]byte(`{
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
				w.Write([]byte(`{
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
			got, res, err := client.Endpoints.TamperProtection(tt.ctx, "tenantid-guid", "", "endpointid-guid")
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
				w.Write([]byte(`{
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
			got, res, err := client.Endpoints.SetTamperProtection(tt.ctx, "tenantid-guid", "", "endpointid-guid")
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
				w.Write([]byte(`{"enabled": true}`))
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
				w.Write([]byte(`{"enabled": false}`))
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
				w.Write([]byte(`{
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
