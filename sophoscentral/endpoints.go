package sophoscentral

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"

	"net/url"
)

// EndpointsService handles communication with the endpoint related methods of the sophos central API

type EndpointService service

// Endpoints contains the returned data from the EndpointsService.List
type Endpoints struct {
	// Item - slice of endpoints
	Items []Item `json:"items"`
	// Pages - Data returned to allow pagination
	// Endpoints utilizes Page-By-Key
	Pages *PagesByFromKey `json:"pages"`
}

// Item contains data for one endpoint
type Item struct {
	// ID is unique id for the endpoint
	ID *string `json:"id" faker:"uuid_hyphenated"`
	// Type - type of endpoint
	// The following values are allowed
	// computer, server, securityVM
	Type                    *string           `json:"type,omitempty"`
	Tenant                  *TenantEP         `json:"tenant,omitempty"`
	Hostname                *string           `json:"hostname,omitempty"`
	Health                  *Health           `json:"health,omitempty"`
	OS                      *OS               `json:"os,omitempty"`
	Ipv4Addresses           []string          `json:"ipv4Addresses,omitempty"`
	Ipv6Addresses           []string          `json:"ipv6Addresses,omitempty"`
	MACAddresses            []string          `json:"macAddresses,omitempty"`
	Group                   *Group            `json:"group,omitempty"`
	AssociatedPerson        *AssociatedPerson `json:"associatedPerson,omitempty"`
	TamperProtectionEnabled *bool             `json:"tamperProtectionEnabled,omitempty"`
	AssignedProducts        []AssignedProduct `json:"assignedProducts,omitempty"`
	LastSeenAt              *string           `json:"lastSeenAt,omitempty"`
	Encryption              *EncryptionEP     `json:"encryption,omitempty"`
	Lockdown                *Lockdown         `json:"lockdown,omitempty"`
}

type Pages struct{}

type IsolationStatus string

type EPCloudProvider string

type LockdownStatus string

type LockdownUpdateStatus string

type ENCStatus string

//func UnmarshalEndpoints(data []byte) (Endpoints, error) {
//	var r Endpoints
//	err := json.Unmarshal(data, &r)
//	if err != nil {
//		return Endpoints{}, fmt.Errorf("%s: %w", "unmarhal failed", err)
//	}
//	return r, err
//}

func (e *Endpoints) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

type EncryptionEP struct {
	Volumes []Volume `json:"volumes"`
}

type Volume struct {
	VolumeID string    `json:"volumeID" faker:"uuid_hyphenated"`
	Status   ENCStatus `json:"status"`
}

type AssignedProduct struct {
	Code    Code           `json:"code"`
	Version string         `json:"version"`
	Status  *string `json:"status"`
}

type AssociatedPerson struct {
	Name     *string `json:"name,omitempty"`
	ViaLogin string  `json:"viaLogin"`
	ID       *string `json:"id,omitempty"`
}

type Group struct {
	Name string `json:"name"`
}

type Health struct {
	Overall  *string  `json:"overall"`
	Threats  Threats  `json:"threats"`
	Services Services `json:"services"`
}

// Services Status of services on the endpoint.
type Services struct {
	// Status  Health status of an endpoint or a service running on an endpoint.
	// The following values are allowed:
	// good, suspicious, bad, unknown
	Status Overall `json:"status"`
	// ServiceDetails Details of services on the endpoint.
	ServiceDetails []ServiceDetail `json:"serviceDetails"`
}

// ServiceDetail Detail of a service on the endpoint.
type ServiceDetail struct {
	//Name service name
	Name *string `json:"name"`
	// Status of a service on an endpoint.
	Status *string `json:"status"`
}

type Threats struct {
	Status *string `json:"status"`
}

type Lockdown struct {
	Status       LockdownStatus       `json:"status"`
	UpdateStatus LockdownUpdateStatus `json:"updateStatus"`
}

type OS struct {
	IsServer     bool     `json:"isServer"`
	Platform     *string `json:"platform"`
	Name         string   `json:"name"`
	MajorVersion int64    `json:"majorVersion"`
	MinorVersion int64    `json:"minorVersion"`
	Build        *int64   `json:"build,omitempty"`
}

type TenantEP struct {
	ID string `json:"id" faker:"uuid_hyphenated"`
}

type Code string

type Overall string

// ServiceDetailName Details of services on the endpoint.
//type ServiceDetailName string
//
//type ServiceDetailStatus string
//
//
//type InstalledState string
//
//type Platform string
//
//type TypeEP string

func (e Endpoints) String() string {
	return Stringify(e)
}

type EndpointListOptions struct {
	ListByFromKeyOptions
	// Sort defines how to sort the data
	// string should match (^[^:]+$)|(^[^:]+:(asc|desc)$)
	Sort string `url:"sort,omitempty"`

	// HealthStatus - find endpoints by status
	// The following values are allowed:
	// bad, good, suspicious, unknown
	HealthStatus string `url:"healthStatus,omitempty"`

	// Type - Find endpoints by type.
	// The following values are allowed:
	// computer, server, securityVm
	Type string `url:"type,omitempty"`

	// TamperProtectionEnabled Find endpoints by whether Tamper Protection is turned on.
	TamperProtectionEnabled bool `url:"tamperProtectionEnabled,omitempty"`

	// LockdownStatus - Find endpoints by lockdown status.
	// The following values are allowed:
	// creatingWhitelist, installing, locked, notInstalled, registering,
	// starting, stopping,unavailable, uninstalled, unlocked
	LockdownStatus string `url:"lockdownStatus,omitempty"`

	// LastSeenBefore - Find endpoints that were last seen before the given date
	// and time (UTC) or a duration relative to the current date and time (exclusive).
	// Examples:
	// 	  To get this:  3 days 4 hours 5 minutes and 0 seconds ago, value is case-sensitive
	//    Use string `-P3DT4H5M0S`
	//
	//		One Day from now
	//		P1D
	//
	//		2 hours ago
	//		-PT2H
	//
	//		20 minutes ago
	//		-PT20M
	//
	//		200 seconds from now
	//		PT200S

	LastSeenBefore string `url:"lastSeenBefore,omitempty"`

	// LastSeenAfter - Find endpoints that were last seen before the given date
	// and time (UTC) or a duration relative to the current date and time (exclusive).
	// Examples:
	// 	  To get this: 4 hours and 500 seconds from now, value is case-sensitive
	//    Use string `PT4H500S`
	//
	//		One Day from ago
	//		-P1D
	//
	//		2 hours from now
	//		PT2H
	//
	//		20 minutes from now
	//		PT20M
	//
	//		200 seconds ago
	//		-PT200S
	LastSeenAfter string `url:"lastSeenAfter,omitempty"`

	// IDs -Find endpoints with the specified IDs.
	IDs []string `url:"ids,omitempty"`

	// IsolationStatus - Find endpoints by isolation status.
	// The following values are allowed:
	//isolated, notIsolated
	IsolationStatus string `url:"isolationStatus,omitempty"`

	// HostnameContains - Find endpoints where the hostname contains the given string.
	hostnameContains string `url:"hostnameContains,omitempty"`

	// AssociatedPersonsContains - Find endpoints where the name of the person
	// associated with the endpoint contains the given string.
	AssociatedPersonContains string `url:"associatedPersonContains,omitempty"`

	// GroupNameContains - Find endpoints where the name of the group the endpoint is in contains the given string.
	GroupNameContains string `url:"groupNameContains,omitempty"`

	// Search Term to search for in the specified search fields.
	Search string `url:"search,omitempty"`

	// SearchFields - List of search fields for finding the given search term. Defaults to all applicable fields.
	// The following values are allowed:
	// hostname, groupName, associatedPersonName, ipAddresses
	// If not used default uses all values
	SearchFields []string `url:"searchFields,omitEmpty"`

	// IPAddresses - Find endpoints by IP addresses.
	IPAddresses []string `url:"ipAddresses,omitEmpty"`

	// Cloud - Find endpoints that are cloud instances. You must use URL encoding.
	// matches ^(aws|azure|gcp)|((aws|azure|gcp):([0-9a-zA-Z-_]{1,64}))|([0-9a-zA-Z-_]{1,64})$
	// Examples:
	// azure:42349c92,aws:i-3bc4829309
	// aws
	// aws, azure:4975692a
	// i-3bc4829309,42349c92
	Cloud []string `url:"cloud,omitempty"`

	// Fields - The fields to return in a partial response.
	Fields []string `url:"fields,omitempty"`

	// View - Type of view to be returned in response.
	// The following values are allowed:
	// basic, summary, full
	View []string `url:"view,omitempty"`
}

// List gathers all endpoints for a tenant ID
// https://developer.sophos.com/docs/endpoint-v1/1/routes/endpoints/get
func (e *EndpointService) List(ctx context.Context, tenantID, tenantURL string, endpoints *Endpoints, opts EndpointListOptions) (*Endpoints, []*Response, error) {
	if tenantID == "" {
		return nil, nil, errors.New("empty tenantID")
	}

	if _, err := uuid.Parse(tenantID); err != nil {
		return nil, nil, errors.New("invalid tenant id")
	}

	if tenantURL == "" {
		return nil, nil, errors.New("empty tenant url")
	}

	if _, err := url.Parse(tenantURL); err != nil {
		return nil, nil, errors.New("invalid tenant url")
	}

	var responses []*Response

	req, err := e.client.NewRequest("GET", "endpoint/v1/endpoints", nil)
	if err != nil {
		return nil, nil, err
	}

	var eps Endpoints

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", tenantID)
	resp, err := e.client.Do(ctx, req, &eps)
	if err != nil {
		return nil, nil, err
	}
	responses = append(responses, resp)

	endpoints.Items = append(endpoints.Items, eps.Items...)
	endpoints.Pages = eps.Pages

	if eps.Pages.GetNextKey() != "" {
		opts = EndpointListOptions{
			HealthStatus: "",
			Type:         "",
			ListByFromKeyOptions: ListByFromKeyOptions{
				PageFromKey: eps.Pages.GetNextKey(),
				PageSize:    500,
				PageTotal:   true,
			},
		}

		rep, res, err := e.List(ctx, tenantID, tenantURL, endpoints, opts)
		if resp.StatusCode == 429 {
			time.Sleep(1 * time.Second)
			return e.List(ctx, tenantID, tenantURL, endpoints, opts)
		}
		return rep, res, err

	}

	return &eps, responses, nil

}

// Get fetches an endpoint
// https://developer.sophos.com/docs/endpoint-v1/1/routes/endpoints/%7BendpointId%7D/get
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/endpoints/{endpointId}
func (e *EndpointService) Get(ctx context.Context, tenantID, tenantURL, endpointID string) (*Item, *Response, error) {

	u := fmt.Sprintf("endpoint/v1/endpoints/%s", endpointID)

	if e == nil {
		return nil, nil, errors.New("nil ep client")
	}
	if e.client == nil {

		return nil, nil, errors.New("nil  client")

	}

	req, err := e.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("X-Tenant-ID", tenantID)

	ep := new(Item)
	resp, err := e.client.Do(ctx, req, ep)
	if err != nil {
		return nil, resp, err
	}
	return ep, resp, nil

}


type TamperProtectionSettings struct {
	Password          *string               `json:"password,omitempty"`
	PreviousPasswords []TPPreviousPasswords `json:"previousPasswords,omitempty"`
	Enabled           bool                  `json:"enabled,omitempty"`
}
type TPPreviousPasswords struct {
	Password      *string `json:"password,omitempty"`
	InvalidatedAt *string `json:"invalidatedAt,omitempty"`
}

// TamperProtection fetches the TamperProtection settings for a specific endpoint
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/endpoints/{endpointId}/tamper-protection
func (e *EndpointService) TamperProtection(ctx context.Context, tenantID, tenantURL, endpointID string) (*TamperProtectionSettings, *Response, error) {

	u := fmt.Sprintf("endpoint/v1/endpoints/%s/tamper-protection", endpointID)

	if e == nil {
		return nil, nil, errors.New("nil ep client")
	}
	if e.client == nil {
		return nil, nil, errors.New("nil  client")
	}

	req, err := e.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("X-Tenant-ID", tenantID)

	tps := new(TamperProtectionSettings)
	resp, err := e.client.Do(ctx, req, tps)
	if err != nil {
		return nil, resp, err
	}
	return tps, resp, nil

}
