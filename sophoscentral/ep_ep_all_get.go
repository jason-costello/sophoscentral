package sophoscentral

import (
	"context"
	"encoding/json"
	"errors"
)

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
	Code    Code    `json:"code"`
	Version string  `json:"version"`
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
	Status string `json:"status"`
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
	IsServer     bool    `json:"isServer"`
	Platform     *string `json:"platform"`
	Name         string  `json:"name"`
	MajorVersion int64   `json:"majorVersion"`
	MinorVersion int64   `json:"minorVersion"`
	Build        *int64  `json:"build,omitempty"`
}

type TenantEP struct {
	ID string `json:"id" faker:"uuid_hyphenated"`
}

type Code string

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
	HostnameContains string `url:"hostnameContains,omitempty"`

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

// List gathers page of endpoints for a tenant ID
// https://api-{region}.central.sophos.com/endpoint/v1/endpoints
func (e *EndpointService) List(ctx context.Context, tenantID string, tenantURL BaseURL, opts EndpointListOptions) (*Endpoints, []*Response, error) {

	// url path to call
	path := e.basePath + "endpoints"
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}
	if tenantID == "" {
		return nil, nil, errors.New("empty tenantID")
	}

	var tenantURLptr *BaseURL
	if tenantURL == "" {
		tenantURLptr = nil
	} else {
		tenantURLptr = &tenantURL
	}

	var responses []*Response
	req, err := e.client.NewRequest(ctx, "GET", path, tenantURLptr, nil)
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


	return &eps, responses, nil

}
