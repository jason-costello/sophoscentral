package sophoscentral

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestNewClient(t *testing.T) {

	a := assert.New(t)
	token := &oauth2.Token{
		AccessToken:  "i am a token value",
		TokenType:    "daft",
		RefreshToken: "i am a refresh token value",
		Expiry:       time.Now().Add(24 * 7 * 52 * 42 * time.Hour),
	}
	// Test all Nil except token.  all nil values will be replaced
	// with defaults from NewClient
	client, err := NewClient(nil, nil, token, nil)
	a.IsType(&Client{}, client)
	a.NoError(err)

	// Test all non-nil values
	client, err = NewClient(context.Background(), &http.Client{}, token, logrus.New())
	a.IsType(&Client{}, client)
	a.NoError(err)

	// creating an httpClient to pass in as option
	httpClient := func(c *Client) {
		h := http.DefaultClient
		c.httpClient = h
		c.httpClient.Timeout = 60 * time.Second
	}

	// Test all Nil except token.  passing httpClient as option
	client, err = NewClient(nil, nil, token, nil, httpClient)
	a.IsType(&Client{}, client)
	a.NoError(err)
	// verify http timeout is same as set above
	var ht time.Duration
	if client != nil {
		ht = client.httpClient.Timeout
	}
	a.Equal(ht, 60*time.Second)

	// set context value to make sure client context is set to
	// what is passed in
	ctxVal := "ah, yep"
	cb := context.Background()
	client, err = NewClient(context.WithValue(cb, "tester", ctxVal), http.DefaultClient, token, nil)
	a.IsType(&Client{}, client)
	a.NoError(err)
	rv := client.ctx.Value("tester").(string)
	a.Equal("ah, yep", rv)

	// all nil.  token will return an error because it is nil
	client, err = NewClient(nil, nil, nil, nil, httpClient)
	a.Error(err)
	a.Equal(ErrNilToken, err)

	// all nil.  token will return an error because it is empty
	client, err = NewClient(nil, nil, &oauth2.Token{}, nil, httpClient)
	a.Error(err)
	a.Equal(ErrEmptyToken, err)

}

// ExampleNewClient shows how to create a new sophos central client
func ExampleNewClient() {
	// don't use an empty token - problems will be had
	// just doing this for example
	token := &oauth2.Token{
		AccessToken:  "i am a token value",
		TokenType:    "daft",
		RefreshToken: "i am a refresh token value",
		Expiry:       time.Now().Add(24 * 7 * 52 * 42 * time.Hour),
	}
	client, _ := NewClient(context.Background(), &http.Client{}, token, logrus.New())
	fmt.Println(client.token.TokenType)
	// Output: daft

}

func Test_unmarshalEntityUUID(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    EntityResponse
		wantErr bool
	}{
		{
			name: "partner",
			args: args{b: []byte(`{
    "id": "5AC55058-622D-4929-8E5D-8FF554F312FE",
    "idType": "partner",
    "apiHosts": {
        "global": "https://api.central.sophos.com"
    }
}`)},
			want: EntityResponse{
				ID:     "5AC55058-622D-4929-8E5D-8FF554F312FE",
				IDType: "partner",
				ApiHosts: ApiHosts{
					Global:     "https://api.central.sophos.com",
					DataRegion: "",
				},
			},
			wantErr: false,
		},

		{
			name: "organization",
			args: args{b: []byte(`{
    "id": "C37A4BC7-715A-48FD-AE03-D184A391B136",
    "idType": "organization",
    "apiHosts": {
        "global": "https://api.central.sophos.com"
    }
}`)},
			want: EntityResponse{
				ID:     "C37A4BC7-715A-48FD-AE03-D184A391B136",
				IDType: "organization",
				ApiHosts: ApiHosts{
					Global:     "https://api.central.sophos.com",
					DataRegion: "",
				},
			},
			wantErr: false,
		},

		{
			name: "tenant",
			args: args{b: []byte(`{
    "id": "7E29B4C2-E68A-4617-BFE1-844666B5300F",
    "idType": "tenant",
    "apiHosts": {
        "global": "https://api.central.sophos.com",
        "dataRegion": "https://api-us02.central.sophos.com"
    }
}`)},
			want: EntityResponse{
				ID:     "7E29B4C2-E68A-4617-BFE1-844666B5300F",
				IDType: "tenant",
				ApiHosts: ApiHosts{
					Global:     "https://api.central.sophos.com",
					DataRegion: "https://api-us02.central.sophos.com",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unmarshalEntityUUID(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshalEntityUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unmarshalEntityUUID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_WhoAmI(t *testing.T) {
	// TODO - need to add tests for 400/500 responses
	//server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
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
	tests := []struct {
		name     string
		s        *httptest.Server
		fields   fields
		wantID   uuid.UUID
		wantType string
		wantErr  bool
	}{
		{
			name: "partner",
			fields: fields{
				ctx:    nil,
				logger: nil,
				token: &oauth2.Token{
					AccessToken:  "i am a token value",
					TokenType:    "daft",
					RefreshToken: "i am a refresh token value",
					Expiry:       time.Now().Add(24 * 7 * 52 * 42 * time.Hour),
				},
				baseURL: nil,
				httpClient: httpClientWithRoundTripper(200, `{
    "id": "5AC55058-622D-4929-8E5D-8FF554F312FE",
    "idType": "partner",
    "apiHosts": {
        "global": "https://api.central.sophos.com"
    }
}`),
				Partner:      nil,
				Organization: nil,
				Tenant:       nil,
			},
			wantID:   uuid.MustParse("5AC55058-622D-4929-8E5D-8FF554F312FE"),
			wantType: "partner",
			wantErr:  false,
		},

		{
			name: "organization",
			fields: fields{
				ctx:    nil,
				logger: nil,
				token: &oauth2.Token{
					AccessToken:  "i am a token value",
					TokenType:    "daft",
					RefreshToken: "i am a refresh token value",
					Expiry:       time.Now().Add(24 * 7 * 52 * 42 * time.Hour),
				},
				baseURL: nil,
				httpClient: httpClientWithRoundTripper(200, `{
    "id": "C37A4BC7-715A-48FD-AE03-D184A391B136",
    "idType": "organization",
    "apiHosts": {
        "global": "https://api.central.sophos.com"
    }
}`),
				Partner:      nil,
				Organization: nil,
				Tenant:       nil,
			},
			wantID:   uuid.MustParse("C37A4BC7-715A-48FD-AE03-D184A391B136"),
			wantType: "organization",
			wantErr:  false,
		},

		{
			name: "tenant",
			fields: fields{
				ctx:    nil,
				logger: nil,
				token: &oauth2.Token{
					AccessToken:  "i am a token value",
					TokenType:    "daft",
					RefreshToken: "i am a refresh token value",
					Expiry:       time.Now().Add(24 * 7 * 52 * 42 * time.Hour),
				},
				baseURL: nil,
				httpClient: httpClientWithRoundTripper(200, `{
    "id": "7E29B4C2-E68A-4617-BFE1-844666B5300F",
    "idType": "tenant",
    "apiHosts": {
        "global": "https://api.central.sophos.com",
        "dataRegion": "https://api-us02.central.sophos.com"
    }
}`),
				Partner:      nil,
				Organization: nil,
				Tenant:       nil,
			},
			wantID:   uuid.MustParse("7E29B4C2-E68A-4617-BFE1-844666B5300F"),
			wantType: "tenant",
			wantErr:  false,
		},
		{
			name: "unmarshal error",
			fields: fields{
				ctx:    nil,
				logger: nil,
				token: &oauth2.Token{
					AccessToken:  "i am a token value",
					TokenType:    "daft",
					RefreshToken: "i am a refresh token value",
					Expiry:       time.Now().Add(24 * 7 * 52 * 42 * time.Hour),
				},
				baseURL:      nil,
				httpClient:   httpClientWithRoundTripper(200, `{ ""}`),
				Partner:      nil,
				Organization: nil,
				Tenant:       nil,
			},
			wantID:   uuid.MustParse("5AC55058-622D-4929-8E5D-8FF554F312FE"),
			wantType: "tenant",
			wantErr:  true,
		},
		{
			name: "500 returned",
			fields: fields{
				ctx:    nil,
				logger: nil,
				token: &oauth2.Token{
					AccessToken:  "i am a token value",
					TokenType:    "daft",
					RefreshToken: "i am a refresh token value",
					Expiry:       time.Now().Add(24 * 7 * 52 * 42 * time.Hour),
				},
				baseURL:      nil,
				httpClient:   httpClientWithRoundTripper(500, ""),
				Partner:      nil,
				Organization: nil,
				Tenant:       nil,
			},
			wantID:   uuid.UUID{},
			wantType: "",
			wantErr:  true,
		},

	{
	name: "http client error",
		fields: fields{
		ctx:    nil,
		logger: nil,
		token: &oauth2.Token{
			AccessToken:  "i am a token value",
			TokenType:    "daft",
			RefreshToken: "i am a refresh token value",
			Expiry:       time.Now().Add(24 * 7 * 52 * 42 * time.Hour),
		},
		baseURL:      nil,
		httpClient:   httpClientWithErrorRoundTripper(),
		Partner:      nil,
		Organization: nil,
		Tenant:       nil,
	},
		wantID:   uuid.UUID{},
		wantType: "",
		wantErr:  true,
	},
}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			c, err := NewClient(tt.fields.ctx, tt.fields.httpClient, tt.fields.token, nil)
			a.NoError(err)

			err = c.WhoAmI()

			if (err != nil) != tt.wantErr {
				t.Errorf("WhoAmI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				return
			}

			if tt.wantType == "partner" {
				a.Equal(tt.wantID.String(), c.Partner.ID.String())
			}
			if tt.wantType == "organization" {
				a.Equal(tt.wantID.String(), c.Organization.ID.String())
			}
			if tt.wantType == "tenant" {
				a.Equal(tt.wantID.String(), c.Tenant.ID.String())
			}

		})
	}
}

type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}


func httpClientWithRoundTripper(statusCode int, response string) *http.Client {
	return &http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: statusCode,
				Body:       ioutil.NopCloser(bytes.NewBufferString(response)),
			}
		}),
	}
}


type roundTripWithErrorFunc func(req *http.Request) error

func (f roundTripWithErrorFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, f(req)}


func httpClientWithErrorRoundTripper() *http.Client {
	return &http.Client{
		Transport: roundTripWithErrorFunc(func(req *http.Request)error {
			return errors.New("error")
		}),
	}
}
