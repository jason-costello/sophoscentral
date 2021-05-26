package sophoscentral

import (
	"context"
	"fmt"
)

// ExploitMitigationAppItems - Page of applications protected by Exploit Mitigation.
type ExploitMitigationAppItems struct{
	// ExploitMitigationAppItem - list of applications protected by Exploit Mitigation.
	Item []ExploitMitigationAppItem `json:"items"`
	// Pages - pagination data
	Pages PagesByOffset `json:"pages"`
}
type ExploitMitigationAppItem struct{

	// ID of an Exploit Mitigation application.
	ID *string `json:"id,omitempty"`
	// Name of this Exploit Mitigation application.
	// length of name must fall within 1 ≤ length ≤ 1000
	Name *string `json:"name,omitempty"`
	// Paths included in this Exploit Mitigation application.
	// must contain at most 100 items
	// length of each path must fall within 1 <= length <= 260
	Paths []string `json:"paths,omitempty"`
	// Exploit Mitigation category ID.
	// The following values are allowed:
	// browsers, exclude, java, media, office, plugins, test, other
	Category ExploitMitigationCategoryID `json:"category,omitempty"`
	// Whether the application was detected by the system or added by the user.
	// The following values are allowed:
	// detected, custom
	Type ExploitMitigationAppType `json:"type,omitempty"`
	// Modifications made to the detected Exploit Mitigation Application.
	// This object does not apply to when type is custom
	Modifications ExploitMitigationAppMods `json:"modifications,omitempty"`
}


// ExploitMitigationAppMods potential modifications for a detected Exploit Mitigation Application.
// This object does not apply to when type is custom.
type ExploitMitigationAppMods struct {
	Protected bool `json:"protected"`
	Settings  struct {
		ASLR            bool `json:"ASLR"`
		BannedAPI       bool `json:"BannedAPI"`
		BottomUpASLR    bool `json:"BottomUpASLR"`
		Caller          bool `json:"Caller"`
		DEP             bool `json:"DEP"`
		HeapSpray       bool `json:"HeapSpray"`
		IAF             bool `json:"IAF"`
		Intruder        bool `json:"Intruder"`
		KbdGuard        bool `json:"KbdGuard"`
		LoadLib         bool `json:"LoadLib"`
		LockdownAutorun bool `json:"LockdownAutorun"`
		LockdownNewFile bool `json:"LockdownNewFile"`
		NullPage        bool `json:"NullPage"`
		SEHOP           bool `json:"SEHOP"`
		StackExec       bool `json:"StackExec"`
		StackPivot      bool `json:"StackPivot"`
	} `json:"settings"`
}
// ExploitMitigationCategoryID - Exploit Mitigation category ID.
type ExploitMitigationCategoryID string
const(
	Browsers ExploitMitigationCategoryID ="browsers"
	Excluded ExploitMitigationCategoryID ="excluded"
	Java ExploitMitigationCategoryID ="java"
	Media ExploitMitigationCategoryID ="media"
	Office ExploitMitigationCategoryID ="office"
	Plugins ExploitMitigationCategoryID ="plugins"
	Test ExploitMitigationCategoryID ="test"
	Other ExploitMitigationCategoryID ="other"

)

type ExploitMitigationAppType string
const(
	Detected ExploitMitigationAppType = "detected"
	Custom   ExploitMitigationAppType = "custom"
)


type ExploitMitigationApplicationGetOptions struct {
	ListByPageOffset
	// Type - Exploit Mitigation Application type.
	// the following values are allowed:
	// 'detected', 'custom'
	Type ExploitMitigationAppType `url:"type,omitempty"`

	// Whether or not Exploit Mitigation Application has been customized.
	Modified bool `url:"modified,omitempty"`
}
// ExploitMitigationApplicationList - List all Exploit Mitigation settings for all protected applications.
// Pagination is using page by offset.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/exploit-mitigation/applications
func (e *EndpointService) ExploitMitigationApplicationList(ctx context.Context, tenantID string,  tenantURL BaseURL, opts *ExploitMitigationApplicationGetOptions) (*ScanningExclusionItems, error) {
	path := fmt.Sprintf("%ssettings/exploit-mitigation/applications", e.basePath)
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
