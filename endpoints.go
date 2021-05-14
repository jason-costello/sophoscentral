package sophoscentral

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetEndpoints(ctx context.Context, tenantID string, geoURL string, queryParams map[string]string) (Endpoints, error){
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/endpoints

reqURL := fmt.Sprintf("%s/endpoint/v1/endpoints", geoURL)
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil{
		fmt.Println("err: ", err.Error())
		return Endpoints{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID",tenantID)
	c.token.SetAuthHeader(req)

	if queryParams != nil {

		q := req.URL.Query()
		for k, v := range queryParams {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
		// defer body.Close()
	}
	b, err := MakeRequest(c.httpClient, req)
	if err != nil {
		return   Endpoints{}, fmt.Errorf("%s: %w", ErrHttpDo, err)
	}

	endpoints, err  := UnmarshalEndpoints(b)
	if err != nil {
		return Endpoints{}, err
	}

	if endpoints.Pages.Total > endpoints.Pages.Current{
		pf, err := NewPaginationFields(ctx, c.httpClient, req, endpoints.Pages)
		if err == nil {
			allPagesBytes, err := GetRemainingPages(pf)
			if err == nil {
				endpoints = collectEndpoints(allPagesBytes,endpoints)
			} else {
				fmt.Println("error collecting additional pages")
			}
		}
	}

	return endpoints, nil

}

func collectEndpoints(pb [][]byte, allEP Endpoints) Endpoints{

	for _, p := range pb{

		endpoints, err  := UnmarshalEndpoints(p)
		if err == nil {
			allEP.Item = append(allEP.Item, endpoints.Item...)
		}
	}

	return allEP

}

type IsolationStatus string
const(
	Isolated IsolationStatus = "isolated"
	NotIsolated = "notIsolated"
)

type EPCloudProvider string
const (
	CPAWS EPCloudProvider= "aws"
	CPAzure = "azure"
)

type LockdownStatus string
const(
	CreatingWhiteList LockdownStatus = "creatingWhiteList"
	Installing = "installing"
	Locked = "locked"
	LDSNotInstalled = "notInstalled"
	Registering = "registering"
	Starting = "starting"
	Stopping = "stopping"
	Unavailable = "unavailable"
	Uninstalled = "uninstalled"
	Unlocked = "unlocked"
)
type LockdownUpdateStatus string
const(
	UpToDate LockdownUpdateStatus = "upToDate"
	LUSUpdating LockdownUpdateStatus= "updating"
	RebootRequired LockdownUpdateStatus= "rebootRequired"
	notInstalled LockdownUpdateStatus= "notInstalled"


)

type ENCStatus string
const(
	NotEncrypted ENCStatus = "notEncrypted"
	Encrypted = "encrypted"
	Encrypting = "encrypting"
	NotSupported = "notSupported"
	Suspended = "suspended"
	EncUnknown = "unknown"
)



func UnmarshalEndpoints(data []byte) (Endpoints, error) {
	var r Endpoints
	err := json.Unmarshal(data, &r)
	if err != nil{
		return Endpoints{}, fmt.Errorf("%s: %w", ErrUnmarshalFailed, err)
	}
	return r, err
}



func (r *Endpoints) Marshal() ([]byte, error) {
	return json.Marshal(r)
}


type Endpoints struct{
	Item []EndpointItem `json:"items"`
	Pages Pages `json:"pages"`
}
type EndpointItem struct {
	ID                      string            `json:"id"`
	Type                    TypeEP              `json:"type"`
	Tenant                  Tenant            `json:"tenant"`
	Hostname                string            `json:"hostname"`
	Health                  *Health           `json:"health,omitempty"`
	OS                      OS                `json:"os"`
	Ipv4Addresses           []string          `json:"ipv4Addresses"`
	Ipv6Addresses           []string          `json:"ipv6Addresses,omitempty"`
	MACAddresses            []string          `json:"macAddresses,omitempty"`
	Group                   Group             `json:"group"`
	AssociatedPerson        *AssociatedPerson `json:"associatedPerson,omitempty"`
	TamperProtectionEnabled *bool             `json:"tamperProtectionEnabled,omitempty"`
	AssignedProducts        []AssignedProduct `json:"assignedProducts,omitempty"`
	LastSeenAt              string            `json:"lastSeenAt"`
	Encryption				*EncryptionEP		`json:"encryption"`
	Lockdown                *Lockdown         `json:"lockdown,omitempty"`
}
type EncryptionEP struct{

	Volumes []Volume `json:"volumes"`

}

type Volume struct{
	VolumeID string `json:"volumeID"`
	Status ENCStatus `json:"status"`
}

type AssignedProduct struct {
	Code    Code                  `json:"code"`
	Version string                `json:"version"`
	Status  InstalledState `json:"status"`
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
	Overall  Overall  `json:"overall"`
	Threats  Threats  `json:"threats"`
	Services Services `json:"services"`
}

type Services struct {
	Status         Overall         `json:"status"`
	ServiceDetails []ServiceDetail `json:"serviceDetails"`
}

type ServiceDetail struct {
	Name   ServiceDetailName   `json:"name"`
	Status ServiceDetailStatus `json:"status"`
}

type Threats struct {
	Status Overall `json:"status"`
}

type Lockdown struct {
	Status       LockdownStatus `json:"status"`
	UpdateStatus LockdownUpdateStatus `json:"updateStatus"`
}

type OS struct {
	IsServer     bool     `json:"isServer"`
	Platform     Platform `json:"platform"`
	Name         string   `json:"name"`
	MajorVersion int64    `json:"majorVersion"`
	MinorVersion int64    `json:"minorVersion"`
	Build        *int64   `json:"build,omitempty"`
}

type TenantEP struct {
	ID string `json:"id"`
}


type Code string
const (
	CoreAgent Code = "coreAgent"
	InterceptX Code= "interceptX"
	EndpointProtection Code= "endpointProtection"
	DeviceEncryption Code= "deviceEncryption"
	MTR Code= "mtr"
)


type Overall string
const (

	Good Overall = "good"
	Suspicious Overall = "suspicious"
	Bad Overall = "bad"
	Unknown Overall = "unknown"
)

type ServiceDetailName string
const (
	FileDetection ServiceDetailName = "File Detection"
	HitmanProAlertService ServiceDetailName = "HitmanPro.Alert service"
	SophosAntiVirus ServiceDetailName = "Sophos Anti-Virus"
	SophosAntiVirusStatusReporter ServiceDetailName = "Sophos Anti-Virus Status Reporter"
	SophosAutoUpdateService ServiceDetailName = "Sophos AutoUpdate Service"
	SophosCleanService ServiceDetailName = "Sophos Clean Service"
	SophosDeviceControlService ServiceDetailName = "Sophos Device Control Service"
	SophosEndpointDefense ServiceDetailName = "Sophos Endpoint Defense"
	SophosEndpointDefenseService ServiceDetailName = "Sophos Endpoint Defense Service"
	SophosFileIntegrityMonitoring ServiceDetailName = "Sophos File Integrity Monitoring"
	SophosFileScanner ServiceDetailName = "Sophos File Scanner"
	SophosFileScannerService ServiceDetailName = "Sophos File Scanner Service"
	SophosMCSAgent ServiceDetailName = "Sophos MCS Agent"
	SophosMCSClient ServiceDetailName = "Sophos MCS Client"
	SophosNetworkThreatProtection ServiceDetailName = "Sophos Network Threat Protection"
	SophosSafestoreService ServiceDetailName = "Sophos Safestore Service"
	SophosSystemProtectionService ServiceDetailName = "Sophos System Protection Service"
	SophosWebControlService ServiceDetailName = "Sophos Web Control Service"
	SophosWebIntelligenceFilterService ServiceDetailName = "Sophos Web Intelligence Filter Service"
	SophosWebIntelligenceService ServiceDetailName = "Sophos Web Intelligence Service"
)

type ServiceDetailStatus string
const (
	Running ServiceDetailStatus = "running"
	Stopped ServiceDetailStatus = "stopped"
	Missing ServiceDetailStatus = "missing"
)

type InstalledState string
const (
	NotInstalled InstalledState = "notInstalled"
	Installed InstalledState = "installed"
)


type Platform string
const (
	Windows Platform = "windows"
	Linux Platform = "linux"
	MacOS Platform = "macos"
)

type TypeEP string
const (
	ServerEP TypeEP = "server"
	ComputerEP TypeEP = "computer"
	SecurityVMEP TypeEP = "securityVm"
)
