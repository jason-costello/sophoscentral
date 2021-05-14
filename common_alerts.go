package sophoscentral

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/*

Implementation for sophos central COMMON API
https://developer.sophos.com/docs/common-v1/1/overview
Resource Types
URIs are relative to https://api-{dataRegion}.central.sophos.com/common/v1"," unless otherwise noted.

GET		/alerts
GET		/alerts/{alertId}
POST	/alerts/{alertId}/actions
POST	/alerts/search

GET		/directory/user-groups
POST	/directory/user-groups
GET		/directory/user-groups/{groupId}
PATCH	/directory/user-groups/{groupId}
DELETE	/directory/user-groups/{groupId}
GET		/directory/user-groups/{groupId}/users
POST	/directory/user-groups/{groupId}/users
DELETE	/directory/user-groups/{groupId}/users
DELETE	/directory/user-groups/{groupId}/users/{userId}

GET		/users/{userId}
PATCH	/users/{userId}
DELETE	/users/{userId}
GET		/users/{userId}/{groups}
POST	/users/{userId}/{groups}
DELETE	/users/{userId}/{groups}
GET		/users/{userId}/groups/{groupID}
*/

/* Wrappers for Alerts.
Allowed query params: groupKey"," from"," to"," sort"," product"," category"," severity"," ids"," fields"," pageSize"," pageFromkey"," pageTotal
*/


// GetAlerts accepts allowed query params and will return all alerts.  Default page size is 50
// and max page size is 100.
func (c *Client) GetAlerts(ctx context.Context, tenant TenantsResponseItem, queryParams map[string]string) (Alerts, error) {
	// https://api-{dataRegion}.central.sophos.com/common/v1/alerts
	//reqURL := tenant.ApiHost + "/common/v1/alerts"
	reqURL := tenant.ApiHost + "/common/v1/alerts"
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil{
		fmt.Println("err: ", err.Error())
		return Alerts{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID",tenant.ID)
	c.token.SetAuthHeader(req)

	q := req.URL.Query()
	for k,v := range queryParams {
		q.Add( k,v)
	}
	req.URL.RawQuery = q.Encode()
	// defer body.Close()


	b, err := MakeRequest(c.httpClient, req)
	if err != nil {
		return   Alerts{}, fmt.Errorf("%s: %w", ErrHttpDo, err)
	}

	alerts, err  := UnmarshalAlerts(b)
	if err != nil {
		return Alerts{}, err
	}
	return alerts, nil
}

// GetAlert accepts allowed query params and will return one alert by id.
func (c *Client) GetAlert(ctx context.Context, tenant TenantsResponseItem, alertID string, queryParams map[string]string) (AlertItem, error) {
	// https://api-{dataRegion}.central.sophos.com/common/v1/alerts/{alertID}
	//reqURL := tenant.ApiHost + "/common/v1/alerts/{alertID}"

	if ctx == nil{
		ctx = context.Background()
	}
	if _, err := uuid.Parse(tenant.ID); err != nil{
		return AlertItem{}, fmt.Errorf("%s: %w", ErrInvalidTenantID, err)
	}

	if _, err := uuid.Parse(alertID); err != nil{
		return AlertItem{}, fmt.Errorf("%s: %w", ErrAlertID, err)
	}

	reqURL := fmt.Sprintf("%s/%s/%s", tenant.ApiHost,"common/v1/alerts", alertID)

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil{
		fmt.Println("err: ", err.Error())
		return AlertItem{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID",tenant.ID)
	c.token.SetAuthHeader(req)

	q := req.URL.Query()
	for k,v := range queryParams {
		q.Add( k,v)
	}
	req.URL.RawQuery = q.Encode()




	b, err := MakeRequest(c.httpClient, req)
	if err != nil {
		return   AlertItem{}, fmt.Errorf("%s: %w", ErrHttpDo, err)
	}


	alert, err  := UnmarshalAlertItem(b)
	if err != nil {
		return AlertItem{}, err
	}
	return alert, nil
}


// RespondToAlert accepts allowed  tenant, alertID, and action to take.
// The request is posted and the status returned
func (c *Client) RespondToAlert(ctx context.Context, tenant TenantsResponseItem, alertID string, action AllowedAction, actionMessage string)  (AlertActionResponse, error) {
// https://api-{dataRegion}.central.sophos.com/common/v1/alerts/{alertId}/actions
	//reqURL := tenant.ApiHost + "/common/v1/alerts/{alertID}/actions"

	if _, err := uuid.Parse(tenant.ID); err != nil{
		return  AlertActionResponse{}, fmt.Errorf("%s: %w", ErrInvalidTenantID, err)
	}

	if _, err := uuid.Parse(alertID); err != nil{
		return   AlertActionResponse{},  fmt.Errorf("%s: %w", ErrAlertID, err)
	}


	rta := RespondToAlertAction{
		Action: action,
		Message: actionMessage,
	}

	payload, err := rta.Marshal()
	if err != nil {
		return  AlertActionResponse{},  fmt.Errorf("%s: %s", ErrMarshalFailed, err)
	}

	reqURL := fmt.Sprintf("%s/common/v1/alerts/%s/actions", tenant.ApiHost,alertID)

	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewBuffer(payload))
	if err != nil{
		return    AlertActionResponse{}, fmt.Errorf("%s: %w", ErrFailedToCreateRequest, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID",tenant.ID)
	c.token.SetAuthHeader(req)

	b, err := MakeRequest(c.httpClient, req)
	if err != nil {
		return   AlertActionResponse{}, fmt.Errorf("%s: %w", ErrHttpDo, err)
	}

	aar, err := UnmarshalAlertActionResponse(b)
	return aar, nil
}

// AlertsSearch accepts allowed  tenant, alertID, and action to take.
// The request is posted and the status returned
func (c *Client) AlertsSearch(ctx context.Context, tenantID, geoURL string, asr AlertSearchRequest, queryParams map[string]string)  (Alerts, error) {
	// https://api-{dataRegion}.central.sophos.com/common/v1/alerts/search
	//reqURL := tenant.ApiHost + "/common/v1/alerts/search"

	if _, err := uuid.Parse(tenantID); err != nil{
		return  Alerts{}, fmt.Errorf("%s: %w", ErrInvalidTenantID, err)
	}

	if geoURL == ""{
		return Alerts{}, errors.New("invalid geoURL")
	}


	reqURL := fmt.Sprintf("%s/common/v1/alerts/search", geoURL)


	payload, err := json.Marshal(asr)
	if err != nil{
		 return Alerts{}, fmt.Errorf("%s: %w", ErrMarshalFailed, err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewBuffer(payload))
	if err != nil{
		return    Alerts{}, fmt.Errorf("%s: %w", ErrFailedToCreateRequest, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID",tenantID)

	if err := verifyAlertsQueryParams(queryParams); err == nil {

		q := req.URL.Query()
		for k, v := range queryParams {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	c.token.SetAuthHeader(req)

	b, err := MakeRequest(c.httpClient, req)
	if err != nil {
		return   Alerts{}, fmt.Errorf("%s: %w", ErrHttpDo, err)
	}

	alerts, err := UnmarshalAlerts(b)
	return alerts, nil
}

func verifyAlertsQueryParams(qp map[string]string) error {

	var errMsgs []string

	if qp == nil{
		return errors.New("query params are nil")
	}
	for k, v := range qp {
		switch k {
		case "groupKey":
			if !isValidGroupKey(v){
				errMsgs = append(errMsgs, fmt.Sprintln("groupKey invalid"))
			}
		case "from":
			if !isValidTime(v) {
				errMsgs = append(errMsgs, fmt.Sprintln("from time invalid"))
			}
		case "to":
			if !isValidTime(v) {
				errMsgs = append(errMsgs, fmt.Sprintln("to time invalid"))
			}
		case "sort":
			if !isValidSort(v) {
				errMsgs = append(errMsgs, fmt.Sprintln("sort value is invalid"))
			}
		case "product":
			if !isValidProduct(v){
				errMsgs = append(errMsgs, fmt.Sprintln("product is invalid"))
			}
		case "category":
			if !isValidCategory(v){
				errMsgs = append(errMsgs, fmt.Sprintln("category is invalid"))
			}

		case "severity":
			if !isValidSeverity(v){
				errMsgs = append(errMsgs, fmt.Sprintln("severity is invalid"))
			}
		case "ids":
			if !areValidUUIDs(strings.Split(v, ",")) {
				errMsgs = append(errMsgs, fmt.Sprintln("ids are invalid"))
			}
		case "fields":
			if !isValidFields(v){
				errMsgs = append(errMsgs, fmt.Sprintln("fields is invalid"))
			}
		case "pageSize":
			if !isValidPageSize(v){
				errMsgs = append(errMsgs, fmt.Sprintln("pageSize is invalid"))
			}
		case "pageFromKey":
			if !isValidPageFromKey(v){
				errMsgs = append(errMsgs, fmt.Sprintln("pageFromKey is invalid"))
			}
		case "pageTotal":
			if !isValidPageTotal(v){
				errMsgs = append(errMsgs, fmt.Sprintln("pageTotal is invalid"))
			}

		}

	}

	if len(errMsgs) < 1 {
		return nil
	}
	return fmt.Errorf("%s: %w", ErrInvalidQueryParams, errors.New(strings.Join(errMsgs, "\n")))

}
func isValidGroupKey(s string) bool{

	if s == ""{
		return true
	}
	return true

}
func isValidTime(ts string) bool{
	var timeFormat = "2006-01-02T03:04:05.000MST"

	_, err := time.Parse(timeFormat, ts)
	if err != nil {
	return false
	}
	return true


}
func isValidSort(s string) bool{
	matched, err := regexp.Match(`(^[^:]+$)|(^[^:]+:(asc|desc)$)`, []byte(s))
	if err != nil {
		return false
	}

	return matched

}
func isValidProduct(s string) bool{
	ls := strings.ToLower(s)
	nls := Product(ls)
	if nls != ""{
		return true
	}
	return false
}
func isValidCategory(s string) bool{
	ls := strings.ToLower(s)
	nls := Category(ls)
	if nls != ""{
		return true
	}
	return false
}
func isValidSeverity(s string) bool{
ls := strings.ToLower(s)
 nls := Severity(ls)
 if nls != ""{
 	return true
 }
 return false
}
func isValidFields(s string) bool{
	if s == ""{
		return true
	}
	return true

}
func isValidPageSize(s string) bool{

	ps, err := strconv.Atoi(s)
	if err != nil{
		return false
	}

	if ps < 1 || ps > 100{
		return false
	}

	return true

}
func isValidPageFromKey(s string) bool{

	if s == ""{
		return true
	}
	return true

}
func isValidPageTotal(s string) bool{
	if strings.ToLower(s) == "false" || strings.ToLower(s) == "true"{
		return true
	}
 	return false


}
func areValidUUIDs(uuids []string) bool {
	if len(uuids) < 1 {
		return false
	}
	for _, x := range uuids {

		_, err := uuid.Parse(x)
		if err != nil {
			return false
		}
	}
	return true
}
func (r *Alerts) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
func (r *RespondToAlertAction) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalAlerts(data []byte) (Alerts, error) {
	var r Alerts
	err := json.Unmarshal(data, &r)
	if err != nil{
		return Alerts{}, fmt.Errorf("%s: %w", ErrUnmarshalFailed, err)
	}
	return r, err
}
func UnmarshalAlertItem(data []byte) (AlertItem, error){
	var r AlertItem
	err := json.Unmarshal(data, &r)
	if err != nil{
		return AlertItem{}, fmt.Errorf("%s: %w", ErrUnmarshalFailed, err)
	}
	return r, nil
}
func UnmarshalAlertActionResponse(data []byte) (AlertActionResponse, error) {
	var r AlertActionResponse
	err := json.Unmarshal(data, &r)
	if err != nil{
		return AlertActionResponse{}, fmt.Errorf("%s: %w", ErrUnmarshalFailed, err)
	}
	return r, err
}
type Alerts struct {
	Items []AlertItem `json:"items"`
	Pages Pages  `json:"pages"`
}
type AlertItem struct {
	ID             string          `json:"id"`
	AllowedActions []AllowedAction `json:"allowedActions"`
	Category       Category          `json:"category"`
	Description    string          `json:"description"`
	GroupKey       string          `json:"groupKey"`
	ManagedAgent   ManagedAgent    `json:"managedAgent"`
	Product        Product         `json:"product"`
	RaisedAt       string          `json:"raisedAt"`
	Severity       Severity        `json:"severity"`
	Tenant         Tenant          `json:"tenant"`
	Type           AlertType          `json:"type"`
}
type ManagedAgent struct {
	ID   *string  `json:"id,omitempty"`
	Type *Product `json:"type,omitempty"`
}
type Tenant struct {
	ID   string `json:"id"`
	Name Name   `json:"name"`
}
type AllowedAction string
const (
	Acknowledge AllowedAction = "acknowledge"
	CleanPua AllowedAction = "cleanPua"
	CleanVirus AllowedAction = "cleanVirus"
	AuthPua AllowedAction = "authPua"
	ClearThreat AllowedAction = "clearThreat"
	ClearHmpa AllowedAction = "clearHmpa"
	SendMsgPua AllowedAction = "sendMsgPua"
	SendMsgThreat AllowedAction = "sendMsgThreat"

)
type Product string
const (
	Other Product = "other"
	Server Product = "server"
	Endpoint Product = "endpoint"
	Mobile Product = "mobile"
	Encryption Product = "encryption"
	EmailGateway Product = "emailGateway"
	WebGateway Product = "webGateway"
	PhishThreat Product = "phishThreat"
	Wireless Product = "wireless"
	IAAS Product = "iaas"
	Firewall Product = "firewall"

)
type Severity string
const (
	High Severity = "high"
	Medium Severity = "medium"
	Low Severity = "low"
)
type Category string
const(
	Azure Category = "azure"
	AdSync Category = "adsync"
	ApplicationControl Category = "applicationControl"
	AppReputation Category = "appreputation"
	BlockListed Category = "blocklisted"
	Connectivity Category = "connectivity"
	CWG Category = "cwg"
	DENC Category = "denc"
	DownloadReputation Category = "downloadreputation"
	EndpointFirewall Category = "endpointfirewall"
	Fenc Category = "fenc"
	ForensicSnapshot Category = "forensicsnapshot"
	General Category = "general"
	Iaas Category = "iaas"
	IaasAzure Category = "iaasazure"
	Isolation Category = "isolation"
	Malware Category = "malware"
	Mtr Category = "mtr"
	Mobiles Category = "mobiles"
	Policy Category = "policy"
	Protection Category = "protection"
	Pua Category = "pua"
	RuntimeDetections Category = "runtimedetections"
	Security Category = "security"
	Smc Category = "smc"
	SystemHealth Category = "systemhealth"
	Uav Category = "uav"
	Uncategorized Category = "uncategorized"
	Updating Category = "updating"
	Utm Category = "utm"
	Virt Category = "virt"
	WirelessCategory Category = "wirelesscategory"
	XGEmail Category = "xgemail"
)
type Name string
type AlertType string
const(
	AlertMobile AlertType = "mobile"
	Computer AlertType = "computer"
	ServerAlert AlertType = "server"
	SecurityVM AlertType = "securityVM"
	UTM AlertType = "utm"
	AccessPoint AlertType = "accessPoint"
	WirelessNetwork AlertType = "wirelessNetwork"
	Mailbox AlertType = "mailbox"
	Slec AlertType = "slec"
	XGFirewall AlertType = "xgFirewall"



)
type RespondToAlertAction struct{
	Action AllowedAction `json:"action"`
	Message string `json:"message"`
}
type AlertActionResponse struct{
	ID string `json:"id"`
	AlertID string `json:"alertID"`
	Action AllowedAction `json:"action"`
	Status AlertActionStatus `json:"status"`
	RequestedAt time.Time `json:"requestedAt"`
	CompletedAt time.Time `json:"completedAt,omitempty"`
	StartedAt time.Time `json:"startedAt,omitempty"`
	Result string `json:"result,omitempty"`
}
type AlertActionStatus string
const(
	Requested AlertActionStatus = "requested"
	Started AlertActionStatus= "started"
	Completed AlertActionStatus = "completed"
)

type AlertSearchRequest struct {
	Category []Category `json:"category"`
	GroupKey string `json:"groupKey"`
	Fields []string `json:"fields"`
	From time.Time `json:"from"`
	IDs []string `json:"ids"`
	Product []Product `json:"product"`
	Severity []Severity `json:"severity"`
	To time.Time `json:"to"`
	PageFromKey string `json:"pageFromKey"`
	PageSize int `json:"pageSize"`
	PageTotal bool `json:"pageTotal"`
	Sort	[]string `json:"sort"`
}
