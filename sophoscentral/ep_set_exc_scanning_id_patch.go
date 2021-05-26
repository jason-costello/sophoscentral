package sophoscentral

import (
	"context"
	"errors"
	"fmt"
)
// ScanningExclusionUpdateItem - Patch object to update the exclusion.
type ScanningExclusionUpdateItem struct{
	// Value - Exclusion value to be updated. Behavioral and Detected Exploit exclusions do not support updating this value.
	Value *string `json:"value,omitempty"`
	// ScanMode Default value of scan mode is "onDemandAndOnAccess"
	// for exclusions of type path, posixPath and virtualPath,
	// "onAccess" for process, web, pua, amsi.
	// Behavioral and Detected Exploits (exploitMitigation)
	// type exclusions do not support a scan mode.
	ScanMode ScanMode `json:"scanMode,omitempty"`
	// Direction - allowed values are

	// Comment indicating why the exclusion was created.
	Comment *string `json:"comment,omitempty"`

}

// ScanningExclusionUpdate - Update a scanning exclusion by ID.
// https://api-{dataRegion}.central.sophos.com/endpoint/v1/settings/exclusions/scanning/{exclusionId}
func (e *EndpointService) ScanningExclusionUpdate(ctx context.Context, tenantID string,  tenantURL BaseURL, exclusionID string, updateItem ScanningExclusionUpdateItem) (*ScanningExclusionItem, error) {
	path := fmt.Sprintf("%ssettings/exclusions/scanning/%s", e.basePath, exclusionID)

	if len(*updateItem.Comment) > 100{
		return nil, errors.New("comment field is limited to 100 chars or less")
	}
	req, err := e.client.NewRequest(ctx, "PATCH", path, &tenantURL, updateItem)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")
	e.client.Token.SetAuthHeader(req)
	sei := new(ScanningExclusionItem)
	_, err = e.client.Do(ctx, req, sei)
	if err != nil {
		return nil, err
	}

	return sei, nil

}
