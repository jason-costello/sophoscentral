package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
	"os"
	"reflect"
	"sophoscentral/sophoscentral"
	"strings"
)
var (
	client *github.Client

	// auth indicates whether tests are being run with an OAuth token.
	// Tests can use this flag to skip certain tests when run without auth.
	auth bool

	skipURLs = flag.Bool("skip_urls", false, "skip url fields")
)
func main() {
	flag.Parse()

	token := os.Getenv("SC_AUTH_TOKEN")
	if token == "" {
		print("!!! No OAuth token. Some tests won't run. !!!\n\n")
		client = github.NewClient(nil)
	} else {
		tc := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		))
		client = github.NewClient(tc)
		auth = true
	}

	for _, tt := range []struct {
		url string
		typ interface{}
	}{
		// https://api-{dataRegion}.central.sophos.com/common/v1/alerts

		//{"rate_limit", &github.RateLimits{}},
		{"endpoint/v1/endpoints", &sophoscentral.Endpoints{}},


	} {
		err := testType(tt.url, tt.typ)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	}
}
// testType fetches the JSON resource at urlStr and compares its keys to the
// struct fields of typ.
func testType(urlStr string, typ interface{}) error {
	slice := reflect.Indirect(reflect.ValueOf(typ)).Kind() == reflect.Slice

	req, err := client.NewRequest("GET", urlStr, nil)
	if err != nil {
		return err
	}

	// start with a json.RawMessage so we can decode multiple ways below
	raw := new(json.RawMessage)
	_, err = client.Do(context.Background(), req, raw)
	if err != nil {
		return err
	}

	// unmarshal directly to a map
	var m1 map[string]interface{}
	if slice {
		var s []map[string]interface{}
		err = json.Unmarshal(*raw, &s)
		if err != nil {
			return err
		}
		m1 = s[0]
	} else {
		err = json.Unmarshal(*raw, &m1)
		if err != nil {
			return err
		}
	}

	// unmarshal to typ first, then re-marshal and unmarshal to a map
	err = json.Unmarshal(*raw, typ)
	if err != nil {
		return err
	}

	var byt []byte
	if slice {
		// use first item in slice
		v := reflect.Indirect(reflect.ValueOf(typ))
		byt, err = json.Marshal(v.Index(0).Interface())
		if err != nil {
			return err
		}
	} else {
		byt, err = json.Marshal(typ)
		if err != nil {
			return err
		}
	}

	var m2 map[string]interface{}
	err = json.Unmarshal(byt, &m2)
	if err != nil {
		return err
	}

	// now compare the two maps
	for k, v := range m1 {
		if *skipURLs && strings.HasSuffix(k, "_url") {
			continue
		}
		if _, ok := m2[k]; !ok {
			fmt.Printf("%v missing field for key: %v (example value: %v)\n", reflect.TypeOf(typ), k, v)
		}
	}

	return nil
}
