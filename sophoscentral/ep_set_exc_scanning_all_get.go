package sophoscentral

import (
	"context"
	"fmt"
)

type ScanningExclusionType string
const(
	Path ScanningExclusionType = "path"
	PosixPath ScanningExclusionType = "posixPath"
	VirtualPath ScanningExclusionType = "virtualPath"
	Process ScanningExclusionType = "process"
	Web ScanningExclusionType = "web"
	Pua ScanningExclusionType = "pua"
	ExploitMitigation ScanningExclusionType = "exploitMitigation"
	AMSI ScanningExclusionType = "amsi"
	Behavioral ScanningExclusionType = "behavioral"
)
type ScanMode string
const(
	OnDemand ScanMode = "onDemand"
	OnAccess ScanMode = "onAccess"
	OnDemandAndOnAccess ScanMode = "onDemandAndOnAccess"
)

type ScanningExclusionGetOptions struct {
	ListByPageOffset
	// Type - scanning exclusion type
	// the following values are allowed
	// path, posixPath, virtualPath,process,web
	// pua,exploitMitigation,amsi,behavioral
	Type ScanningExclusionType `url:"type,omitempty"`
	}

type ScanningExclusionItems struct{
	Item []ScanningExclusionItem `json:"items"`
	Pages PagesByOffset `json:"pages"`
}

type ScanningExclusionItem struct{
	// ID - Unique ID for the scanning exclusion setting.
	ID *string `json:"id,omitempty"`
	// LocalPorts that are allowed
	// Type - Scanning exclusion type.
	// The following values are allowed:
	// path, posixPath, virtualPath, process, web, pua, detectedExploit, amsi, behavioral
	Type *ScanningExclusionType `json:"type,omitempty"`
	// ScanMode Default value of scan mode is "onDemandAndOnAccess"
	// for exclusions of type path, posixPath and virtualPath,
	// "onAccess" for process, web, pua, amsi.
	// Behavioral and Detected Exploits (exploitMitigation)
	// type exclusions do not support a scan mode.
	ScanMode ScanMode `json:"scanMode,omitempty"`
	// Direction - allowed values are
	Direction *ExclusionDirections `json:"direction,omitempty"`
	// Description - Exclusion description added by the system.
	Description *string `json:"description,omitempty"`
	// Comment indicating why the exclusion was created.
	Comment *string `json:"comment,omitempty"`

}
// ScanningExclusionList - List all scanning exclusions.
// Pagination is using page by offset.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/exclusions/scanning
func (e *EndpointService) ScanningExclusionList(ctx context.Context, tenantID string,  tenantURL BaseURL, opts *ScanningExclusionGetOptions) (*ScanningExclusionItems, error) {
	path := fmt.Sprintf("%ssettings/exclusions/scanning", e.basePath)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil,  err
	}
	req, err := e.client.NewRequest(ctx, "GET", path, &tenantURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")
	e.client.Token.SetAuthHeader(req)
	sei := new(ScanningExclusionItems)
	_, err = e.client.Do(ctx, req, sei)
	if err != nil {
		return nil, err
	}

	return sei, nil

}
