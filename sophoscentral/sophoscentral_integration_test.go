// +build integration

package sophoscentral_test

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sophoscentral/sophoscentral"
	"strings"
	"testing"
)



var skipURLs = flag.Bool("skip_urls", false, "skip url fields")

func TestReturnedTypes(t *testing.T){

	// auth indicates whether tests are being run with an OAuth token.
	// Tests can use this flag to skip certain tests when run without auth.

	a := assert.New(t)
	ctx := context.Background()
	cid := os.Getenv("SC_CLIENT_ID")
	cs := os.Getenv("SC_CLIENT_SEC")
	tid := os.Getenv("TYPE_TESTING_TENANT_ID_PASS")
	eid := os.Getenv("TYPE_TESTING_EP_ID_PASS")
	config := &clientcredentials.Config{
		ClientID:       cid,
		ClientSecret:   cs,
		Scopes:         []string{"token"},
		TokenURL:       "https://id.sophos.com/api/v2/oauth2/token",
		EndpointParams: url.Values{},
	}

	token, err := config.Token(ctx)
	a.NoError(err)

	oauthConfig := oauth2.Config{
		Endpoint: oauth2.Endpoint{
			TokenURL:  "https://id.sophos.com/api/v2/oauth2/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
		Scopes: []string{"token"},
	}


	test := []struct {
		url         string
		typ         interface{}
		hClient     *http.Client
		token       *oauth2.Token
		serviceType interface{}
		headers     map[string]string
	}{
		{
			url:         "https://api-us03.central.sophos.com/endpoint/v1/endpoints",
			typ:         &sophoscentral.Endpoints{},
			hClient:     oauthConfig.Client(ctx, token),
			token:       token,
			serviceType: sophoscentral.EndpointService{},
			headers: map[string]string{
				"Content-Type": "application/json",
				"X-Tenant-ID":  tid,
			},
		},
		{
			url:         fmt.Sprintf("https://api-us03.central.sophos.com/endpoint/v1/endpoints/%s", eid),
			typ:         &sophoscentral.Item{},
			hClient:     oauthConfig.Client(ctx, token),
			token:       token,
			serviceType: &sophoscentral.EndpointService{},
			headers: map[string]string{
				"Content-Type": "application/json",
				"X-Tenant-ID":  tid,
			},
		},
	}


	for _, tt := range test {

		err := testType(a, tt.url, tt.typ, tt.hClient, tt.token,  tt.headers)
		a.NoError(err)

	}

}



// testType fetches the JSON resource at urlStr and compares its keys to the
// struct fields of typ.
func testType(a *assert.Assertions, urlStr string, typ interface{}, hClient *http.Client, token *oauth2.Token, headers map[string]string) error {
	ctx := context.Background()
	slice := reflect.Indirect(reflect.ValueOf(typ)).Kind() == reflect.Slice
	client := sophoscentral.NewClient(ctx, hClient, token)

	req, err := client.NewRequest("GET", urlStr, nil)
	a.NoError(err)

	for k,v := range headers{
		req.Header.Set(k,v)
	}
	//req.Header.Set("Authorization", fmt.Sprintf("Bearer %s",token.AccessToken))
	// start with a json.RawMessage so we can decode multiple ways below
	raw := new(json.RawMessage)
	_, err = client.Do(context.Background(), req, raw)
	a.NoError(err)

	// unmarshal directly to a map
	var m1 map[string]interface{}
	if slice {
		var s []map[string]interface{}
		err = json.Unmarshal(*raw, &s)
		a.NoError(err)
		m1 = s[0]
	} else {
		err = json.Unmarshal(*raw, &m1)
		a.NoError(err)

	}

	// unmarshal to typ first, then re-marshal and unmarshal to a map
	err = json.Unmarshal(*raw, typ)
	a.NoError(err)


	var byt []byte
	if slice {
		// use first item in slice
		v := reflect.Indirect(reflect.ValueOf(typ))
		byt, err = json.Marshal(v.Index(0).Interface())
		a.NoError(err)

	} else {
		byt, err = json.Marshal(typ)
		a.NoError(err)

	}

	var m2 map[string]interface{}
	err = json.Unmarshal(byt, &m2)
	a.NoError(err)


	// now compare the two maps
	for k, v := range m1 {
		if *skipURLs && strings.HasSuffix(k, "_url") {
			continue
		}

		_, ok := m2[k]
		a.Truef(ok, "%v missing field for key: %v (example value: %v)\n", reflect.TypeOf(typ), k, sophoscentral.PrettyPrint(v))
		//if !ok {
		//	fmt.Printf("%v missing field for key: %v (example value: %v)\n", reflect.TypeOf(typ), k, PrettyPrint(v))
		//}
	}

	return nil
}
