package sophoscentral

import (
	"context"
	"errors"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test_verifyUUIDs(t *testing.T) {
	type args struct {
		uids []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid uuids",
			args: args{uids: []string{uuid.New().String(), uuid.New().String(), uuid.New().String()}},
			want: true,
		},
		{
			name: "empty uuid slice",
			args: args{uids: []string{}},
			want: false,
		},
		{
			name: "all not uuids",
			args: args{uids: []string{"hello", "world"}},
			want: false,
		},
		{
			name: "one invalid uuid",
			args: args{uids: []string{uuid.New().String(), "i am not uuid", uuid.New().String()}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := areValidUUIDs(tt.args.uids); got != tt.want {
				t.Errorf("verifyUUIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_verifyAlertsQueryParams(t *testing.T) {
	type args struct {
		qp map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{qp: map[string]string{
				"groupKey":    "any string",
				"from":        time.Now().UTC().Format("2006-01-02T03:04:05.000MST"),
				"to":          time.Now().UTC().Format("2006-01-02T03:04:05.000MST"),
				"sort":        "unsure what text here asc",
				"product":     "endpoint",
				"category":    "azure",
				"severity":    "high",
				"ids":         strings.Join([]string{uuid.New().String(), uuid.New().String(), uuid.New().String()}, ","),
				"fields":      strings.Join([]string{"any string", "values", "right", "now"}, ","),
				"pageSize":    "50",
				"pageFromKey": "not sure what data type goes here",
				"pageTotal":   "true",
			},
			},
			wantErr: false,
		},
		{
			name: "all invalid values",
			args: args{qp: map[string]string{
				"groupKey":    "any string",
				"from":        time.Now().UTC().Format(time.Kitchen),
				"to":          time.Now().UTC().Format(time.ANSIC),
				"sort":        "",
				"product":     "ep",
				"category":    "a",
				"severity":    "tiny",
				"ids":         strings.Join([]string{"no guids here", "none here"}, ","),
				"fields":      strings.Join([]string{""}, ","),
				"pageSize":    "",
				"pageFromKey": "",
				"pageTotal":   "what, me worry??",
			},
			},
			wantErr: true,
		},
		{
			name: "valid",
			args: args{qp: map[string]string{
				"groupKey":    "any string",
				"from":        time.Now().UTC().Format("2006-01-02T03:04:05.000MST"),
				"to":          time.Now().UTC().Format("2006-01-02T03:04:05.000MST"),
				"sort":        "::%::",
				"product":     "endpoint",
				"category":    "azure",
				"severity":    "high",
				"ids":         strings.Join([]string{uuid.New().String(), uuid.New().String(), uuid.New().String()}, ","),
				"fields":      strings.Join([]string{"any string", "values", "right", "now"}, ","),
				"pageSize":    "101",
				"pageFromKey": "not sure what data type goes here",
				"pageTotal":   "neither",
			},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := verifyAlertsQueryParams(tt.args.qp); (err != nil) != tt.wantErr {
				t.Errorf("verifyAlertsQueryParams() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}


func TestClient_GetAlerts(t *testing.T) {
	 burl, _ := url.Parse("https://api-us02.central.sophos.com")


	type fields struct {
		ctx          context.Context
		logger       *logrus.Logger
		token        *oauth2.Token
		baseURL      *url.URL
		httpClient   *http.Client
		Partner      *PartnerService
		Organization *OrganizationService
		Tenant       *TenantService
	}
	type args struct {
		ctx         context.Context
		tenantID string
		queryParams map[string]string
	}
	tests := []struct {
		name    string
		s        *httptest.Server
		fields  fields
		args    args
		want    Alerts
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				ctx:         context.Background(),
				tenantID:    "5AC55058-622D-4929-8E5D-8FF554F312FE",
				queryParams: map[string]string{"pageSize": "1"},
			},
			fields: fields{
				ctx:    context.Background(),
				logger: nil,
				token: &oauth2.Token{
					AccessToken:  "i am a token value",
					TokenType:    "daft",
					RefreshToken: "i am a refresh token value",
					Expiry:       time.Now().Add(24 * 7 * 52 * 42 * time.Hour),
				},
				baseURL: burl,
				httpClient: httpClientWithRoundTripper(200, `{
    "id": "5AC55058-622D-4929-8E5D-8FF554F312FE",
    "idType": "tenant",
    "apiHosts": {
        "global": "https://api-us03.central.sophos.com"
    }
}`),
				Partner:      nil,
				Organization: nil,
				Tenant:       nil,
			},
			want:    Alerts{},
			wantErr: false,
		},
		{
			name: "nil context - failed request creation",
			args: args{
				ctx:         nil,
				tenantID:    "5AC55058-622D-4929-8E5D-8FF554F312FE",
				queryParams: map[string]string{"pageSize": "1"},
			},
			fields: fields{
				ctx:    context.Background(),
				logger: nil,
				token: &oauth2.Token{
					AccessToken:  "i am a token value",
					TokenType:    "daft",
					RefreshToken: "i am a refresh token value",
					Expiry:       time.Now().Add(24 * 7 * 52 * 42 * time.Hour),
				},
				baseURL: burl,
				httpClient: httpClientWithRoundTripper(200, `{
    "id": "5AC55058-622D-4929-8E5D-8FF554F312FE",
    "idType": "tenant",
    "apiHosts": {
        "global": "https://api-us03.central.sophos.com"
    }
}`),
				Partner:      nil,
				Organization: nil,
				Tenant:       nil,
			},
			want:    Alerts{},
			wantErr: true,
		},
		{
			name: "401_error",
			args: args{
				ctx:         nil,
				tenantID:    "5AC55058-622D-4929-8E5D-8FF554F312FE",
				queryParams: map[string]string{"pageSize": "1"},
			},
			fields: fields{
				ctx:    context.Background(),
				logger: nil,
				token: &oauth2.Token{
					AccessToken:  "i am a token value",
					TokenType:    "daft",
					RefreshToken: "i am a refresh token value",
					Expiry:       time.Now().Add(24 * 7 * 52 * 42 * time.Hour),
				},
				baseURL: burl,
				httpClient: httpClientWithRoundTripper(401, `{
    "id": "5AC55058-622D-4929-8E5D-8FF554F312FE",
    "idType": "tenant",
    "apiHosts": {
        "global": "https://api-us03.central.sophos.com"
    }
}`),
				Partner:      nil,
				Organization: nil,
				Tenant:       nil,
			},
			want:    Alerts{},
			wantErr: true,
		},
		{
			name: "err_500",
			args: args{
				ctx:         nil,
				tenantID:    "5AC55058-622D-4929-8E5D-8FF554F312FE",
				queryParams: map[string]string{"pageSize": "1"},
			},
			fields: fields{
				ctx:    context.Background(),
				logger: nil,
				token: &oauth2.Token{
					AccessToken:  "i am a token value",
					TokenType:    "daft",
					RefreshToken: "i am a refresh token value",
					Expiry:       time.Now().Add(24 * 7 * 52 * 42 * time.Hour),
				},
				baseURL: burl,
				httpClient: httpClientWithRoundTripper(500, `{
    "id": "5AC55058-622D-4929-8E5D-8FF554F312FE",
    "idType": "tenant",
    "apiHosts": {
        "global": "https://api-us03.central.sophos.com"
    }
}`),
				Partner:      nil,
				Organization: nil,
				Tenant:       nil,
			},
			want:    Alerts{},
			wantErr: true,
		},
		{
			name: "http.do returns error",
			args: args{
				ctx:         context.Background(),
				tenantID:    "5AC55058-622D-4929-8E5D-8FF554F312FE",
				queryParams: map[string]string{"pageSize": "1"},
			},
			fields: fields{
				ctx:    context.Background(),
				logger: nil,
				token: &oauth2.Token{
					AccessToken:  "i am a token value",
					TokenType:    "daft",
					RefreshToken: "i am a refresh token value",
					Expiry:       time.Now().Add(24 * 7 * 52 * 42 * time.Hour),
				},
				baseURL:      burl,
				httpClient:   httpClientWithErrorRoundTripper(),
				Partner:      nil,
				Organization: nil,
				Tenant:       nil,
			},
			want:    Alerts{},
			wantErr: true,
		},
		{
			name: "unmarshal error",
			args: args{
				ctx:         context.Background(),
				tenantID:    "5AC55058-622D-4929-8E5D-8FF554F312FE",
				queryParams: map[string]string{"pageSize": "1"},
			},
			fields: fields{
				ctx:    context.Background(),
				logger: nil,
				token: &oauth2.Token{
					AccessToken:  "i am a token value",
					TokenType:    "daft",
					RefreshToken: "i am a refresh token value",
					Expiry:       time.Now().Add(24 * 7 * 52 * 42 * time.Hour),
				},
				baseURL: burl,
				httpClient: httpClientWithRoundTripper(200, `{
    "id": "5AC55058-622D-4929-8E5D-8FF554F312FE",
    "idType": ,
    "apiHosts": {
        : "https://api-us03.central.sophos.com"
    }
}`),
				Partner:      nil,
				Organization: nil,
				Tenant:       nil,
			},
			want:    Alerts{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				ctx:          tt.fields.ctx,
				logger:       tt.fields.logger,
				token:        tt.fields.token,
				baseURL:      tt.fields.baseURL,
				httpClient:   tt.fields.httpClient,
				Partner:      tt.fields.Partner,
				Organization: tt.fields.Organization,
				Tenant:       tt.fields.Tenant,
			}
			if _, ok := tt.args.queryParams["pageSize"]; !ok{
				t.Errorf("map not initialized")
				return
			}
			got, err := c.GetAlerts(tt.args.ctx, TenantsResponseItem{
				ID:            "b9b62247-783c-4e59-93c8-8adaaa53c7b1",
				Name:          "TriCore Solutions",
				DataGeography: "US",
				DataRegion:    "us03",
				BillingType:   "usage",
				Partner:       TenantsResponsePartner{},
				ApiHost:       "https://api-us03.central.sophos.com",
				Status:        "active",
			},
			tt.args.queryParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAlerts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAlerts() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetAlert(t *testing.T) {
	burl, _ := url.Parse("https://api-us03.central.sophos.com")
	type fields struct {
		ctx          context.Context
		logger       *logrus.Logger
		token        *oauth2.Token
		baseURL      *url.URL
		httpClient   *http.Client
		Partner      *PartnerService
		Organization *OrganizationService
		Tenant       *TenantService
	}
	type args struct {
		ctx         context.Context
		tenant      TenantsResponseItem
		alertID     string
		queryParams map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    AlertItem
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				ctx:         context.Background(),
				tenant: TenantsResponseItem{
					ID:            "bccccccc-783c-4e59-93c8-8adaaa53c7b1",
					Name:          "TriCore Solutions",
					DataGeography: "US",
					DataRegion:    "us03",
					BillingType:   "usage",
					Partner:       TenantsResponsePartner{},
					ApiHost:       "https://api-us03.central.sophos.com",
					Status:        "active",
				},
				alertID: "fa17b218-6c9b-40bf-a72c-bdc4ed917fa5",
				queryParams: map[string]string{"pageSize": "1"},
			},
			fields: fields{
				ctx:    context.Background(),
				logger: nil,
				token: &oauth2.Token{
					AccessToken:  "i am a token value",
					TokenType:    "daft",
					RefreshToken: "i am a refresh token value",
					Expiry:       time.Now().Add(24 * 7 * 52 * 42 * time.Hour),
				},
				baseURL: burl ,
				httpClient: httpClientWithRoundTripper(200, `{
	"id": "fa17b218-6c9b-40bf-a72c-bdc4ed917fa5",
	"allowedActions": [
		"acknowledge"
	],
	"category": "general",
	"description": "Your API token SophosMigration2020 has expired",
	"groupKey": "zazo6UmVuZXdBcGlUb2tlbiwxLA",
	"managedAgent": {},
	"product": "other",
	"raisedAt": "2021-05-02T06:00:25.454Z",
	"severity": "high",
	"tenant": {
		"id": "bccccccc-783c-4e59-93c8-8adaaa53c7b1",
		"name": "TriCore Solutions"
	},
	"type": "Event::Task::RenewApiToken"
				
}`),
				Partner:      nil,
				Organization: nil,
				Tenant:       nil,
			},
			want:   AlertItem{
				ID:             "fa17b218-6c9b-40bf-a72c-bdc4ed917fa5",
				AllowedActions: []AllowedAction{Acknowledge},
				Category:       General,
				Description:    "Your API token SophosMigration2020 has expired",
				GroupKey:       "zazo6UmVuZXdBcGlUb2tlbiwxLA",
				ManagedAgent:   ManagedAgent{},
				Product:       Other,
				RaisedAt:       "2021-05-02T06:00:25.454Z",
				Severity:       High,
				Tenant:         Tenant{
					ID:   "bccccccc-783c-4e59-93c8-8adaaa53c7b1",
					Name: "TriCore Solutions",
				},
				Type:           "Event::Task::RenewApiToken",
			},
			wantErr:  false,
		},



		{
			name: "401 error",
			args: args{
				ctx:         context.Background(),
				tenant: TenantsResponseItem{
					ID:            "bccccccc-783c-4e59-93c8-8adaaa53c7b1",
					Name:          "TriCore Solutions",
					DataGeography: "US",
					DataRegion:    "us03",
					BillingType:   "usage",
					Partner:       TenantsResponsePartner{},
					ApiHost:       "https://api-us03.central.sophos.com",
					Status:        "active",
				},
				alertID: "fa17b218-6c9b-40bf-a72c-bdc4ed917fa5",
				queryParams: map[string]string{"pageSize": "1"},
			},
			fields: fields{
				ctx:    context.Background(),
				logger: nil,
				token: &oauth2.Token{
					AccessToken:  "i am a token value",
					TokenType:    "daft",
					RefreshToken: "i am a refresh token value",
					Expiry:       time.Now().Add(24 * 7 * 52 * 42 * time.Hour),
				},
				baseURL: burl ,
				httpClient: httpClientWithRoundTripper(401, ``),
				Partner:      nil,
				Organization: nil,
				Tenant:       nil,
			},
			want:   AlertItem{
				ID:             "fa17b218-6c9b-40bf-a72c-bdc4ed917fa5",
				AllowedActions: []AllowedAction{Acknowledge},
				Category:       General,
				Description:    "Your API token SophosMigration2020 has expired",
				GroupKey:       "zazo6UmVuZXdBcGlUb2tlbiwxLA",
				ManagedAgent:   ManagedAgent{},
				Product:       Other,
				RaisedAt:       "2021-05-02T06:00:25.454Z",
				Severity:       High,
				Tenant:         Tenant{
					ID:   "bccccccc-783c-4e59-93c8-8adaaa53c7b1",
					Name: "TriCore Solutions",
				},
				Type:           "Event::Task::RenewApiToken",
			},
			wantErr:  true,
		},



		{
			name: "500 error",
			args: args{
				ctx:         context.Background(),
				tenant: TenantsResponseItem{
					ID:            "bccccccc-783c-4e59-93c8-8adaaa53c7b1",
					Name:          "TriCore Solutions",
					DataGeography: "US",
					DataRegion:    "us03",
					BillingType:   "usage",
					Partner:       TenantsResponsePartner{},
					ApiHost:       "https://api-us03.central.sophos.com",
					Status:        "active",
				},
				alertID: "fa17b218-6c9b-40bf-a72c-bdc4ed917fa5",
				queryParams: map[string]string{"pageSize": "1"},
			},
			fields: fields{
				ctx:    context.Background(),
				logger: nil,
				token: &oauth2.Token{
					AccessToken:  "i am a token value",
					TokenType:    "daft",
					RefreshToken: "i am a refresh token value",
					Expiry:       time.Now().Add(24 * 7 * 52 * 42 * time.Hour),
				},
				baseURL: burl ,
				httpClient: httpClientWithRoundTripper(500, ``),
				Partner:      nil,
				Organization: nil,
				Tenant:       nil,
			},
			want:   AlertItem{
				ID:             "fa17b218-6c9b-40bf-a72c-bdc4ed917fa5",
				AllowedActions: []AllowedAction{Acknowledge},
				Category:       General,
				Description:    "Your API token SophosMigration2020 has expired",
				GroupKey:       "zazo6UmVuZXdBcGlUb2tlbiwxLA",
				ManagedAgent:   ManagedAgent{},
				Product:       Other,
				RaisedAt:       "2021-05-02T06:00:25.454Z",
				Severity:       High,
				Tenant:         Tenant{
					ID:   "bccccccc-783c-4e59-93c8-8adaaa53c7b1",
					Name: "TriCore Solutions",
				},
				Type:           "Event::Task::RenewApiToken",
			},
			wantErr:  true,
		},



		{
			name: "http.Do error",
			args: args{
				ctx:         context.Background(),
				tenant: TenantsResponseItem{
					ID:            "bccccccc-783c-4e59-93c8-8adaaa53c7b1",
					Name:          "TriCore Solutions",
					DataGeography: "US",
					DataRegion:    "us03",
					BillingType:   "usage",
					Partner:       TenantsResponsePartner{},
					ApiHost:       "https://api-us03.central.sophos.com",
					Status:        "active",
				},
				alertID: "fa17b218-6c9b-40bf-a72c-bdc4ed917fa5",
				queryParams: map[string]string{"pageSize": "1"},
			},
			fields: fields{
				ctx:    context.Background(),
				logger: nil,
				token: &oauth2.Token{
					AccessToken:  "i am a token value",
					TokenType:    "daft",
					RefreshToken: "i am a refresh token value",
					Expiry:       time.Now().Add(24 * 7 * 52 * 42 * time.Hour),
				},
				baseURL: burl ,
				httpClient: httpClientWithErrorRoundTripper(),
				Partner:      nil,
				Organization: nil,
				Tenant:       nil,
			},
			want:   AlertItem{
				ID:             "fa17b218-6c9b-40bf-a72c-bdc4ed917fa5",
				AllowedActions: []AllowedAction{Acknowledge},
				Category:       General,
				Description:    "Your API token SophosMigration2020 has expired",
				GroupKey:       "zazo6UmVuZXdBcGlUb2tlbiwxLA",
				ManagedAgent:   ManagedAgent{},
				Product:       Other,
				RaisedAt:       "2021-05-02T06:00:25.454Z",
				Severity:       High,
				Tenant:         Tenant{
					ID:   "bccccccc-783c-4e59-93c8-8adaaa53c7b1",
					Name: "TriCore Solutions",
				},
				Type:           "Event::Task::RenewApiToken",
			},
			wantErr:  true,
		},



		{
			name: "nil context error",
			args: args{
				ctx:         nil,
				tenant: TenantsResponseItem{
					ID:            "bccccccc-783c-4e59-93c8-8adaaa53c7b1",
					Name:          "TriCore Solutions",
					DataGeography: "US",
					DataRegion:    "us03",
					BillingType:   "usage",
					Partner:       TenantsResponsePartner{},
					ApiHost:       "https://api-us03.central.sophos.com",
					Status:        "active",
				},
				alertID: "fa17b218-6c9b-40bf-a72c-bdc4ed917fa5",
				queryParams: map[string]string{"pageSize": "1"},
			},
			fields: fields{
				ctx:    context.Background(),
				logger: nil,
				token: &oauth2.Token{
					AccessToken:  "i am a token value",
					TokenType:    "daft",
					RefreshToken: "i am a refresh token value",
					Expiry:       time.Now().Add(24 * 7 * 52 * 42 * time.Hour),
				},
				baseURL: burl ,
				httpClient: httpClientWithRoundTripper(200, `{
	"id": "fa17b218-6c9b-40bf-a72c-bdc4ed917fa5",
	"allowedActions": [
		"acknowledge"
	],
	"category": "general",
	"description": "Your API token SophosMigration2020 has expired",
	"groupKey": "zazo6UmVuZXdBcGlUb2tlbiwxLA",
	"managedAgent": {},
	"product": "other",
	"raisedAt": "2021-05-02T06:00:25.454Z",
	"severity": "high",
	"tenant": {
		"id": "bccccccc-783c-4e59-93c8-8adaaa53c7b1",
		"name": "TriCore Solutions"
	},
	"type": "Event::Task::RenewApiToken"
				
}`),
				Partner:      nil,
				Organization: nil,
				Tenant:       nil,
			},
			want:   AlertItem{
				ID:             "fa17b218-6c9b-40bf-a72c-bdc4ed917fa5",
				AllowedActions: []AllowedAction{Acknowledge},
				Category:       General,
				Description:    "Your API token SophosMigration2020 has expired",
				GroupKey:       "zazo6UmVuZXdBcGlUb2tlbiwxLA",
				ManagedAgent:   ManagedAgent{},
				Product:       Other,
				RaisedAt:       "2021-05-02T06:00:25.454Z",
				Severity:       High,
				Tenant:         Tenant{
					ID:   "bccccccc-783c-4e59-93c8-8adaaa53c7b1",
					Name: "TriCore Solutions",
				},
				Type:           "Event::Task::RenewApiToken",
			},
			wantErr:  false,
		},


		{
			name: "unmarshal error",
			args: args{
				ctx:         context.Background(),
				tenant: TenantsResponseItem{
					ID:            "bccccccc-783c-4e59-93c8-8adaaa53c7b1",
					Name:          "TriCore Solutions",
					DataGeography: "US",
					DataRegion:    "us03",
					BillingType:   "usage",
					Partner:       TenantsResponsePartner{},
					ApiHost:       "https://api-us03.central.sophos.com",
					Status:        "active",
				},
				alertID: "fa17b218-6c9b-40bf-a72c-bdc4ed917fa5",
				queryParams: map[string]string{"pageSize": "1"},
			},
			fields: fields{
				ctx:    context.Background(),
				logger: nil,
				token: &oauth2.Token{
					AccessToken:  "i am a token value",
					TokenType:    "daft",
					RefreshToken: "i am a refresh token value",
					Expiry:       time.Now().Add(24 * 7 * 52 * 42 * time.Hour),
				},
				baseURL: burl ,
				httpClient: httpClientWithRoundTripper(200, `{
	"id": "fa17b218-6c9b-40bf-a72c-bdc4ed917fa5",
	"allowedActions": [
		"acknowledge"
	],
	"category": "general",
	"description": "Your API token SophosMigration2020 has expired",
	"groupKey": "zazo6UmVuZXdBcGlUb2tlbiwxLA",
	"managedAgent": {},
	"her",
	"raisedAt": "2021-05-02T06:00:25.454Z",
	"severity": "high",
	" {
		"id": "bccccccc-783c-4e59-93c8-8adaaa53c7b1",
		"name": "TriCore Solutions"
	},
	"type": "Event::Task::RenewApiToken"
				
}`),
				Partner:      nil,
				Organization: nil,
				Tenant:       nil,
			},
			want:   AlertItem{
				ID:             "fa17b218-6c9b-40bf-a72c-bdc4ed917fa5",
				AllowedActions: []AllowedAction{Acknowledge},
				Category:       General,
				Description:    "Your API token SophosMigration2020 has expired",
				GroupKey:       "zazo6UmVuZXdBcGlUb2tlbiwxLA",
				ManagedAgent:   ManagedAgent{},
				Product:       Other,
				RaisedAt:       "2021-05-02T06:00:25.454Z",
				Severity:       High,
				Tenant:         Tenant{
					ID:   "bccccccc-783c-4e59-93c8-8adaaa53c7b1",
					Name: "TriCore Solutions",
				},
				Type:           "Event::Task::RenewApiToken",
			},
			wantErr:  true,
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				ctx:          tt.fields.ctx,
				logger:       tt.fields.logger,
				token:        tt.fields.token,
				baseURL:      tt.fields.baseURL,
				httpClient:   tt.fields.httpClient,
				Partner:      tt.fields.Partner,
				Organization: tt.fields.Organization,
				Tenant:       tt.fields.Tenant,
			}
			got, err := c.GetAlert(tt.args.ctx, tt.args.tenant, tt.args.alertID, tt.args.queryParams)
			if err != nil && tt.wantErr {
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAlert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAlert() got = %v, want %v", got, tt.want)
			}
		})
	}
}
func mustParseURL(s string) *url.URL{
	u, err := url.Parse(s)
	if err != nil{
		panic("failed must parse url")
	}
	return u
}
func mustParseTime(t string, layout string)time.Time {

	tt, err := time.Parse(layout, t)
	if err != nil{
		panic("failed must parse time")
	}
	return tt



}
func TestClient_RespondToAlert(t *testing.T) {
		a := assert.New(t)
	type fields struct {
		ctx          context.Context
		logger       *logrus.Logger
		token        *oauth2.Token
		baseURL      *url.URL
		httpClient   *http.Client
		Partner      *PartnerService
		Organization *OrganizationService
		Tenant       *TenantService
	}
	type args struct {
		ctx           context.Context
		tenant        TenantsResponseItem
		alertID       string
		action        AllowedAction
		actionMessage string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want AlertActionResponse
		wantErr bool
	}{
			{
				name: "one",
				fields: fields{
					ctx:          context.Background(),
					logger:       logrus.New(),
					token:        &oauth2.Token{
						AccessToken: faker.Jwt(),
						TokenType:    "token",
						RefreshToken: faker.Sentence(),
						Expiry:      time.Now().Add(24 * 7 * 52 * 42 * time.Hour),
					},
					baseURL:     mustParseURL(faker.URL()),
					httpClient:   httpClientWithRoundTripper(201, `{
										"id": "49310a33-4acc-409b-aafb-07b8bc06ef01",
					"alertID": "bc893b97-86a8-41aa-b65c-910e11505605",
					"action": "acknowledge",
					"status": "completed",
					"requestedAt": "2021-05-02T06:00:25.454Z",
					"completedAt": "2021-05-02T06:05:25.454Z",
					"startedAt": "2021-05-02T06:02:25.454Z",
					"result": "result string here"
					}`),

					Partner:      nil,
					Organization: nil,
					Tenant:       &TenantService{
						ID:   uuid.MustParse("49310a33-4acc-409b-aafb-07b8bc06ef01"),
					},

				},
				args: args{
					ctx:           context.Background(),
					tenant:        TenantsResponseItem{
						ID:            "49310a33-4acc-409b-aafb-07b8bc06ef01",
						Name:          "I am tenant name",
						DataGeography: USGeo,
						DataRegion:    US03,
						BillingType:   Usage,
						Partner:       TenantsResponsePartner{ID: "C37A4BC7-715A-48FD-AE03-D184A391B136"},
						ApiHost:       "https://api-us03.central.sophos.com",
						Status:        "",
					},
					alertID:       "bc893b97-86a8-41aa-b65c-910e11505605",
					action:        Acknowledge,
					actionMessage: "I am an action message",
				},
				want: AlertActionResponse{
					ID:          "49310a33-4acc-409b-aafb-07b8bc06ef01",
					AlertID:     "bc893b97-86a8-41aa-b65c-910e11505605",
					Action:      "acknowledge",
					Status:      Completed,
					RequestedAt: mustParseTime("2021-05-02T06:00:25.454Z", time.RFC3339),
					CompletedAt: mustParseTime("2021-05-02T06:05:25.454Z", time.RFC3339),
					StartedAt:   mustParseTime("2021-05-02T06:02:25.454Z", time.RFC3339),
					Result:      "result string here",
				},
				wantErr: false,
			},
		{
			name: "401 error",
			fields: fields{
				ctx:          context.Background(),
				logger:       logrus.New(),
				token:        &oauth2.Token{
					AccessToken: faker.Jwt(),
					TokenType:    "token",
					RefreshToken: faker.Sentence(),
					Expiry:      time.Now().Add(24 * 7 * 52 * 42 * time.Hour),
				},
				baseURL:     mustParseURL(faker.URL()),
				httpClient:   httpClientWithRoundTripper(401, ``),
				Partner:      nil,
				Organization: nil,
				Tenant:       &TenantService{
					ID:   uuid.MustParse("49310a33-4acc-409b-aafb-07b8bc06ef01"),
				},

			},
			args: args{
				ctx:           context.Background(),
				tenant:        TenantsResponseItem{
					ID:            "49310a33-4acc-409b-aafb-07b8bc06ef01",
					Name:          "I am tenant name",
					DataGeography: USGeo,
					DataRegion:    US03,
					BillingType:   Usage,
					Partner:       TenantsResponsePartner{ID: "C37A4BC7-715A-48FD-AE03-D184A391B136"},
					ApiHost:       "https://api-us03.central.sophos.com",
					Status:        "",
				},
				alertID:       "bc893b97-86a8-41aa-b65c-910e11505605",
				action:        Acknowledge,
				actionMessage: "I am an action message",
			},
			want: AlertActionResponse{},
			wantErr: true,
		},
		{
			name: "500 error",
			fields: fields{
				ctx:          context.Background(),
				logger:       logrus.New(),
				token:        &oauth2.Token{
					AccessToken: faker.Jwt(),
					TokenType:    "token",
					RefreshToken: faker.Sentence(),
					Expiry:      time.Now().Add(24 * 7 * 52 * 42 * time.Hour),
				},
				baseURL:     mustParseURL(faker.URL()),
				httpClient:   httpClientWithRoundTripper(500, ``),
				Partner:      nil,
				Organization: nil,
				Tenant:       &TenantService{
					ID:   uuid.MustParse("49310a33-4acc-409b-aafb-07b8bc06ef01"),
				},

			},
			args: args{
				ctx:           context.Background(),
				tenant:        TenantsResponseItem{
					ID:            "49310a33-4acc-409b-aafb-07b8bc06ef01",
					Name:          "I am tenant name",
					DataGeography: USGeo,
					DataRegion:    US03,
					BillingType:   Usage,
					Partner:       TenantsResponsePartner{ID: "C37A4BC7-715A-48FD-AE03-D184A391B136"},
					ApiHost:       "https://api-us03.central.sophos.com",
					Status:        "",
				},
				alertID:       "bc893b97-86a8-41aa-b65c-910e11505605",
				action:        Acknowledge,
				actionMessage: "I am an action message",
			},
			want: AlertActionResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				ctx:          tt.fields.ctx,
				logger:       tt.fields.logger,
				token:        tt.fields.token,
				baseURL:      tt.fields.baseURL,
				httpClient:   tt.fields.httpClient,
				Partner:      tt.fields.Partner,
				Organization: tt.fields.Organization,
				Tenant:       tt.fields.Tenant,
			}
			got, err := c.RespondToAlert(tt.args.ctx, tt.args.tenant, tt.args.alertID, tt.args.action, tt.args.actionMessage);
			if (err != nil) != tt.wantErr {
				ne := errors.Unwrap(err)
				t.Errorf("wrapped error: %v", err)
				t.Errorf("unwrapped error: %v", ne)
				return
			}

			a.Equal(tt.want, got)

		})
	}
}

func TestClient_AlertsSearch(t *testing.T) {
	a := assert.New(t)
	sp := Server
	serverPointer := &sp
	type fields struct {
		ctx          context.Context
		logger       *logrus.Logger
		token        *oauth2.Token
		baseURL      *url.URL
		httpClient   *http.Client
		Partner      *PartnerService
		Organization *OrganizationService
		Tenant       *TenantService
	}
	type args struct {
		ctx         context.Context
		tenantID    string
		geoURL      string
		asr         AlertSearchRequest
		queryParams map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Alerts
		wantErr bool
	}{
		{
			name: "one",
			fields: fields{
				ctx:          context.Background(),
				logger:       logrus.New(),
				token:        &oauth2.Token{
					AccessToken: faker.Jwt(),
					TokenType:    "token",
					RefreshToken: faker.Sentence(),
					Expiry:      time.Now().Add(24 * 7 * 52 * 42 * time.Hour),
				},
				baseURL:     mustParseURL(faker.URL()),
				httpClient:   httpClientWithRoundTripper(200, `{
	"items": [{
		"id": "bc893b97-86a8-41aa-b65c-910e11505605",
		"allowedActions": [
			"acknowledge"
		],
		"category": "policy",
		"description": "Real time protection disabled",
		"groupKey": "MixFdmVudDo6RW5kcG9pbnQ6OlNhdkRpc2FibGVkLDUxMyw",
		"managedAgent": {
			"id": "641d09de-f229-438f-bbcd-82c9bf6bfb58",
			"type": "server"
		},
		"product": "endpoint",
		"raisedAt": "2021-04-25T20:01:07.825Z",
		"severity": "high",
		"tenant": {
			"id": "49310a33-4acc-409b-aafb-07b8bc06ef01",
			"name": "OOP"
		},
		"type": "Event::Endpoint::SavDisabled"
	}]
}`),

				Partner:      nil,
				Organization: nil,
				Tenant:       &TenantService{
					ID:   uuid.MustParse("49310a33-4acc-409b-aafb-07b8bc06ef01"),
				},

			},
			args: args{
				ctx:           context.Background(),
				tenantID: "49310a33-4acc-409b-aafb-07b8bc06ef01",
				geoURL: "https://api-us03.central.sophos.com",
				asr:       AlertSearchRequest{
					Category:    []Category{Azure,AdSync},
					GroupKey:    "",
					Fields:      nil,
					From:        time.Time{},
					IDs:         []string{"bc893b97-86a8-41aa-b65c-910e11505605"},
					Product:     []Product{Endpoint},
					Severity:    []Severity{High},
					To:          time.Time{},
					PageFromKey: "",
					PageSize:    1,
					PageTotal:   false,
					Sort:        nil,
				},
				queryParams:       nil,
			},
			want: Alerts{
				Items: []AlertItem{
					AlertItem{
					ID:             "bc893b97-86a8-41aa-b65c-910e11505605",
					AllowedActions: []AllowedAction{Acknowledge},
					Category:       Policy,
					Description:    "Real time protection disabled",
					GroupKey:       "MixFdmVudDo6RW5kcG9pbnQ6OlNhdkRpc2FibGVkLDUxMyw",
					ManagedAgent:   ManagedAgent{
						ID:   getPointer("641d09de-f229-438f-bbcd-82c9bf6bfb58"),
						Type: serverPointer,
					},
					Product:        Endpoint,
					RaisedAt:       "2021-04-25T20:01:07.825Z",
					Severity:       High,
					Tenant:         Tenant{
						ID:   "49310a33-4acc-409b-aafb-07b8bc06ef01",
						Name: "OOP",
					},
					Type:           "Event::Endpoint::SavDisabled",
				}},
				Pages: Pages{},
			},
		},
	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				ctx:          tt.fields.ctx,
				logger:       tt.fields.logger,
				token:        tt.fields.token,
				baseURL:      tt.fields.baseURL,
				httpClient:   tt.fields.httpClient,
				Partner:      tt.fields.Partner,
				Organization: tt.fields.Organization,
				Tenant:       tt.fields.Tenant,
			}
			got, err := c.AlertsSearch(tt.args.ctx, tt.args.tenantID, tt.args.geoURL, tt.args.asr, tt.args.queryParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("AlertsSearch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			a.Equal(tt.want, got)

		})
	}
}

func getPointer(s string) *string{

	var ns string
	ns = s
	return &ns
}
