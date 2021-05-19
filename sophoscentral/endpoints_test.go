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
		ListByFromKeyOptions:     ListByFromKeyOptions{
			PageFromKey: "1",
			PageSize:    500,
			PageTotal:   true,
		},
	}
	got, _, err := client.Endpoints.List(ctx, "3d7a50a6-aee1-4193-a178-689bcb86f750", "https://tenurl.com", &Endpoints{}, epl)
	if err != nil {
		assert.NoErrorf(t, err,"EndpointService.Get returned error: %v", err)
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
						Status: Overall("good"),
						ServiceDetails: []ServiceDetail{
							{
								Name:   "HitmanPro.Alert service",
								Status: "running",
							},
							{
								Name:   "Sophos Anti-Virus",
								Status: "running",
							},
						},
					},
				},
				OS: &OS{
					IsServer:     true,
					Platform:     "windows",
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
						Status:  "installed",
					},
					{
						Code:    "endpointProtection",
						Version: "10.8.9.2",
						Status:  "installed",
					},
					{
						Code:    "interceptX",
						Version: "2.10.18",
						Status:  "installed",
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
	//testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
	//	got, resp, err := client.Endpoints.Get(ctx, "tenantid-guid", "url", "endpointid-guid")
	//	assert.Nilf(t, got, "testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)

	//	return resp, err
	//})

}

func TestEndpointService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/endpoint/v1/endpoints/endpointid-guid", func(w http.ResponseWriter, r *http.Request) {
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

	})

	assert.NotNil(t, client.Endpoints)

	ctx := context.Background()

	got, _, err := client.Endpoints.Get(ctx, "tenantid-guid", "", "endpointid-guid")
	if err != nil {
		assert.NoErrorf(t, err,"EndpointService.Get returned error: %v", err)
	}
	assert.NotNil(t, got)



	want := &Item{
		ID:                      String("endpointid-guid"),
		Type:                    String("server"),
		Tenant:                  &TenantEP{ID: "tenantid-guid"},
		Hostname:                String("myserverhostname"),
		Health:                  &Health{
			Overall:  String("suspicious"),
			Threats:  Threats{Status: String("good")},
			Services: Services{
				Status:         Overall("good"),
				ServiceDetails: []ServiceDetail{
					{
						Name:   "HitmanPro.Alert service",
						Status: "running",
					},
					{
						Name:   "Sophos Anti-Virus",
						Status: "running",
					},
				},
			},
		},
		OS:                      &OS{
			IsServer:     true,
			Platform:    "windows",
			Name:         "Windows Server 2012 R2 Standard",
			MajorVersion: 6,
			MinorVersion: 3,
			Build:        Int64(9600),
		},
		Ipv4Addresses:           []string{"10.127.181.8", "110.205.25.8"},
		Ipv6Addresses:           []string{"fc80::25cc:2ae5:d417:54a0"},
		MACAddresses:            []string{       "00:5a:52:83:54:AB", "00:5a:53:8C:28:9D"},
		Group:                   &Group{Name: "Bolton's Skinning Service"},
		AssociatedPerson:        &AssociatedPerson{
			Name:     nil,
			ViaLogin: `MyCompy\\tywinthelion`,
			ID:       nil,
		},
		TamperProtectionEnabled: Bool(true),
		AssignedProducts:        []AssignedProduct{
			{
				Code:    "coreAgent",
				Version: "2.10.8",
				Status:  "installed",
			},
			{
				Code:    "endpointProtection",
				Version: "10.8.9.2",
				Status:  "installed",
			},
			{
				Code:    "interceptX",
				Version: "2.10.18",
				Status:  "installed",
			},
		},

		LastSeenAt:              String("2020-12-28T10:46:57.393Z"),
		Encryption:              nil,
		Lockdown:                &Lockdown{
			Status:       "notInstalled",
			UpdateStatus: "notInstalled",
		},
	}
	assert.Equal(t, want, got)

	//const methodName = "Get"
	//
	//testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
	//	got, resp, err := client.Endpoints.Get(ctx, "tenantid-guid", "url", "endpointid-guid")
	//	assert.Nilf(t, got, "testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)

	//	return resp, err
	//})

}
