package sophoscentral

//import (
//	"context"
//	"fmt"
//	"github.com/stretchr/testify/assert"
//	"golang.org/x/oauth2"
//	"golang.org/x/oauth2/clientcredentials"
//	"net/http"
//	"net/url"
//	"reflect"
//	"testing"
//)
//func GetAuthToken(ctx context.Context, clientID, clientSecret, tokenURL string) (*oauth2.Token, error) {
//	config := &clientcredentials.Config{
//		ClientID:       clientID,
//		ClientSecret:   clientSecret,
//		Scopes:         []string{"token"},
//		TokenURL:       tokenURL,
//		EndpointParams: url.Values{},
//	}
//
//
//	return config.TokenSource(ctx).Token()
//}
//
//func TestEndpointService_List(t *testing.T) {
//	a := assert.New(t)
//	ctx := context.Background()

//	authToken, err := GetAuthToken(ctx, cidResult, csResult, "https://id.sophos.com/api/v2/oauth2/token" )
//	a.NoError(err)
//	a.NotEmpty(authToken)
//
//	oauthConfig := oauth2.Config{

//		Endpoint:     oauth2.Endpoint{
//			TokenURL:  "https://id.sophos.com/api/v2/oauth2/token",
//			AuthStyle:	oauth2.AuthStyleInParams,
//		},
//		Scopes:       []string{"token"},
//	}
//
//	httpClient := oauthConfig.Client(ctx, authToken)
//
//	client := NewClient(ctx, httpClient, authToken)
//
//	type args struct {
//		ctx       context.Context
//		tenantID  string
//		tenantURL string
//		opts      EndpointListOptions
//	}
//	tests := []struct {
//		name    string
//		e       *EndpointService
//		args    args
//		want    *Endpoints
//		want1   []*Response
//		wantErr bool
//	}{
//		{
//			name: "one",
//			e: 	client.Endpoints,
//			args: args{
//				ctx:       context.Background(),
//				tenantID:  "b9b62247-783c-4e59-93c8-8adaaa53c7b1",
//				tenantURL: "https://api-us03.central.sophos.com/endpoint/v1/endpoints",
//				opts:      EndpointListOptions{
//					HealthStatus:      "",
//					Type:              "",
//					ListByFromKeyOptions: ListByFromKeyOptions{
//						PageFromKey: "",
//						PageSize:    500,
//						PageTotal: true  ,
//					},
//				},
//			},
//			want: &Endpoints{},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, _, err := tt.e.List(tt.args.ctx, tt.args.tenantID, tt.args.tenantURL, &Endpoints{}, tt.args.opts)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			a.NotNil(got.Item)
//			a.NotEmpty(got.Item)
//
//		})
//	}
//}
//
//func TestEndpointService_Get(t *testing.T) {
//	client, mux, _, teardown := setup()
//	defer teardown()
//
//	mux.HandleFunc("/endpoints", func(w http.ResponseWriter, r *http.Request) {
//		testMethod(t, r, "GET")
//		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
//	})
//
//	ctx := context.Background()
//	got, _, err := client.Endpoints.List(ctx, "xxxx", "tenanturl", &Endpoints{}, EndpointListOptions{})
//	if err != nil {
//		t.Errorf("Repositories.List returned error: %v", err)
//	}
//
//	want := []*Item{{ID: String("1")}, {ID: String("2")}}
//	if !reflect.DeepEqual(got, want) {
//		t.Errorf("Repositories.List returned %+v, want %+v", got, want)
//	}
//
//	const methodName = "List"
//	testBadOptions(t, methodName, func() (err error) {
//		_, _, err = client.Endpoints.List(ctx, "\n", "\n", nil, EndpointListOptions{})
//		return err
//	})
//
//	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
//		got, resp, err := client.Endpoints.List(ctx, "\n", "\n", nil, EndpointListOptions{})
//		if got != nil {
//			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
//		}
//		return resp[0], err
//	})
//}
