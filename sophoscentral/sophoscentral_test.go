package sophoscentral

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"
	"time"
)
const (
	// baseURLPath is a non-empty Client.BaseURL path to use during tests,
	// to ensure relative URLs are used for all endpoints.
	baseURLPath = "/api-testing"
)
// setup sets up a test HTTP server along with a github.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative. It only makes a difference
	// when there's a non-empty base URL path. So, use that. See issue #752.
	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\tDid you accidentally use an absolute endpoint URL rather than relative?")
		fmt.Fprintln(os.Stderr, "\tSee https://github.com/google/go-github/issues/752 for information.")
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the GitHub client being tested and is
	// configured to use test server.
	client = NewClient(context.Background(), server.Client(), nil)
	url, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = url

	return client, mux, server.URL, server.Close
}

// openTestFile creates a new file with the given name and content for testing.
// In order to ensure the exact file name, this function will create a new temp
// directory, and create the file in that directory. It is the caller's
// responsibility to remove the directory and its contents when no longer needed.
func openTestFile(name, content string) (file *os.File, dir string, err error) {
	dir, err = ioutil.TempDir("", "sophoscentral")
	if err != nil {
		return nil, dir, err
	}

	file, err = os.OpenFile(path.Join(dir, name), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return nil, dir, err
	}

	fmt.Fprint(file, content)

	// close and re-open the file to keep file.Stat() happy
	file.Close()
	file, err = os.Open(file.Name())
	if err != nil {
		return nil, dir, err
	}

	return file, dir, err
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	t.Helper()
	want := url.Values{}
	for k, v := range values {
		want.Set(k, v)
	}

	r.ParseForm()
	if got := r.Form; !reflect.DeepEqual(got, want) {
		t.Errorf("Request parameters: %v, want %v", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	t.Helper()
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
	}
}

func testURLParseError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func testBody(t *testing.T, r *http.Request, want string) {
	t.Helper()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Error reading request body: %v", err)
	}
	if got := string(b); got != want {
		t.Errorf("request Body is %s, want %s", got, want)
	}
}

// Test whether the marshaling of v produces JSON that corresponds
// to the want string.
func testJSONMarshal(t *testing.T, v interface{}, want string) {
	t.Helper()
	// Unmarshal the wanted JSON, to verify its correctness, and marshal it back
	// to sort the keys.
	u := reflect.New(reflect.TypeOf(v)).Interface()
	if err := json.Unmarshal([]byte(want), &u); err != nil {
		t.Errorf("Unable to unmarshal JSON for %v: %v", want, err)
	}
	w, err := json.Marshal(u)
	if err != nil {
		t.Errorf("Unable to marshal JSON for %#v", u)
	}

	// Marshal the target value.
	j, err := json.Marshal(v)
	if err != nil {
		t.Errorf("Unable to marshal JSON for %#v", v)
	}

	if string(w) != string(j) {
		t.Errorf("json.Marshal(%q) returned %s, want %s", v, j, w)
	}
}

// Test how bad options are handled. Method f under test should
// return an error.
func testBadOptions(t *testing.T, methodName string, f func() error) {
	t.Helper()
	if methodName == "" {
		t.Error("testBadOptions: must supply method methodName")
	}
	if err := f(); err == nil {
		t.Errorf("bad options %v err = nil, want error", methodName)
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient(nil, http.DefaultClient, nil)

	if got, want := c.BaseURL.String(), defaultBaseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
	if got, want := c.UserAgent, userAgent; got != want {
		t.Errorf("NewClient UserAgent is %v, want %v", got, want)
	}

	nhc := new(http.Client)
	c2 := NewClient(nil, nhc, nil)
	if c.client == c2.client {
		t.Error("NewClient returned same http.Clients, but they should differ")
	}
}

func TestNewRequest_invalidJSON(t *testing.T) {
	c := NewClient(nil, nil, nil)

	type T struct {
		A map[interface{}]interface{}
	}
	_, err := c.NewRequest("GET", ".", &T{})

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if err, ok := err.(*json.UnsupportedTypeError); !ok {
		t.Errorf("Expected a JSON error; got %#v.", err)
	}
}

func TestNewRequest_badURL(t *testing.T) {
	c := NewClient(nil, nil, nil)
	_, err := c.NewRequest("GET", ":", nil)
	testURLParseError(t, err)
}

func TestNewRequest_badMethod(t *testing.T) {
	c := NewClient(nil, nil, nil)
	if _, err := c.NewRequest("BOGUS\nMETHOD", ".", nil); err == nil {
		t.Fatal("NewRequest returned nil; expected error")
	}
}

// ensure that no User-Agent header is set if the client's UserAgent is empty.
// This caused a problem with Google's internal http client.
func TestNewRequest_emptyUserAgent(t *testing.T) {
	c := NewClient(nil, nil, nil)
	c.UserAgent = ""
	req, err := c.NewRequest("GET", ".", nil)
	if err != nil {
		t.Fatalf("NewRequest returned unexpected error: %v", err)
	}
	if _, ok := req.Header["User-Agent"]; ok {
		t.Fatal("constructed request contains unexpected User-Agent header")
	}
}

// If a nil body is passed to sophoscentral.NewRequest, make sure that nil is also
// passed to http.NewRequest. In most cases, passing an io.Reader that returns
// no content is fine, since there is no difference between an HTTP request
// body that is an empty string versus one that is not set at all. However in
// certain cases, intermediate systems may treat these differently resulting in
// subtle errors.
func TestNewRequest_emptyBody(t *testing.T) {
	c := NewClient(nil, nil, nil)
	req, err := c.NewRequest("GET", ".", nil)
	if err != nil {
		t.Fatalf("NewRequest returned unexpected error: %v", err)
	}
	if req.Body != nil {
		t.Fatalf("constructed request contains a non-nil Body")
	}
}

func TestNewRequest_errorForNoTrailingSlash(t *testing.T) {
	tests := []struct {
		rawurl    string
		wantError bool
	}{
		{rawurl: "https://example.com/api/v3", wantError: true},
		{rawurl: "https://example.com/api/v3/", wantError: false},
	}
	c := NewClient(nil, nil, nil)
	for _, test := range tests {
		u, err := url.Parse(test.rawurl)
		if err != nil {
			t.Fatalf("url.Parse returned unexpected error: %v.", err)
		}
		c.BaseURL = u
		if _, err := c.NewRequest(http.MethodGet, "test", nil); test.wantError && err == nil {
			t.Fatalf("Expected error to be returned.")
		} else if !test.wantError && err != nil {
			t.Fatalf("NewRequest returned unexpected error: %v.", err)
		}
	}
}

//
//func TestDo(t *testing.T) {
//	client, mux, _, teardown := setup()
//	defer teardown()
//
//	type foo struct {
//		A string
//	}
//
//	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//		testMethod(t, r, "GET")
//		fmt.Fprint(w, `{"A":"a"}`)
//	})
//
//	req, _ := client.NewRequest("GET", ".", nil)
//	body := new(foo)
//	ctx := context.Background()
//	resp, err := client.Do(ctx, req, body)
//	if err != nil{
//		t.Error(err)
//	}
//	bb := resp.Body
//
//	want := &foo{"a"}
//	if !reflect.DeepEqual(bb, want) {
//		t.Errorf("Response body = %v, want %v", body, want)
//	}
//}

func TestDo_nilContext(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	req, _ := client.NewRequest("GET", ".", nil)
	_, err := client.Do(nil, req, nil)

	if !errors.Is(err, errNonNilContext) {
		t.Errorf("Expected context must be non-nil error")
	}
}

//func TestDo_httpError(t *testing.T) {
//	client, mux, _, teardown := setup()
//	defer teardown()
//
//	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//		http.Error(w, "Bad Request", 400)
//	})
//
//	req, _ := client.NewRequest("GET", ".", nil)
//	ctx := context.Background()
//	resp, err := client.Do(ctx, req, nil)
//
//	if err == nil {
//		t.Fatal("Expected HTTP 400 error, got no error.")
//	}
//	if resp.StatusCode != 400 {
//		t.Errorf("Expected HTTP 400 error, got %d status code.", resp.StatusCode)
//	}
//}

//func TestDo_noContent(t *testing.T) {
//	client, mux, _, teardown := setup()
//	defer teardown()
//
//	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(http.StatusNoContent)
//	})
//
//	var body json.RawMessage
//
//	req, _ := client.NewRequest("GET", ".", nil)
//	ctx := context.Background()
//	_, err := client.Do(ctx, req, &body)
//	if err != nil {
//		t.Fatalf("Do returned unexpected error: %v", err)
//	}
//}

func TestSanitizeURL(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"/?a=b", "/?a=b"},
		{"/?a=b&client_secret=secret", "/?a=b&client_secret=REDACTED"},
		{"/?a=b&client_id=id&client_secret=secret", "/?a=b&client_id=id&client_secret=REDACTED"},
	}

	for _, tt := range tests {
		inURL, _ := url.Parse(tt.in)
		want, _ := url.Parse(tt.want)

		if got := sanitizeURL(inURL); !reflect.DeepEqual(got, want) {
			t.Errorf("sanitizeURL(%v) returned %v, want %v", tt.in, got, want)
		}
	}
}

func TestCheckResponse(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(strings.NewReader(`{"message":"m"}`)),
	}
	err := CheckResponse(res).(*ErrorResponse)

	if err == nil {
		t.Errorf("Expected error response.")
	}

	want := &ErrorResponse{
		Response: res,
		Message:  "m"}

	if !errors.Is(err, want) {
		t.Errorf("Error = %#v, want %#v", err, want)
	}
}

func TestCompareHttpResponse(t *testing.T) {
	testcases := map[string]struct {
		h1       *http.Response
		h2       *http.Response
		expected bool
	}{
		"both are nil": {
			expected: true,
		},
		"both are non nil - same StatusCode": {
			expected: true,
			h1:       &http.Response{StatusCode: 200},
			h2:       &http.Response{StatusCode: 200},
		},
		"both are non nil - different StatusCode": {
			expected: false,
			h1:       &http.Response{StatusCode: 200},
			h2:       &http.Response{StatusCode: 404},
		},
		"one is nil, other is not": {
			expected: false,
			h2:       &http.Response{},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			v := compareHttpResponse(tc.h1, tc.h2)
			if tc.expected != v {
				t.Errorf("Expected %t, got %t for (%#v, %#v)", tc.expected, v, tc.h1, tc.h2)
			}
		})
	}
}

// ensure that we properly handle API errors that do not contain a response body
func TestCheckResponse_noBody(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}
	err := CheckResponse(res).(*ErrorResponse)

	if err == nil {
		t.Errorf("Expected error response.")
	}

	want := &ErrorResponse{
		Response: res,
	}
	if !errors.Is(err, want) {
		t.Errorf("Error = %#v, want %#v", err, want)
	}
}

func TestCheckResponse_unexpectedErrorStructure(t *testing.T) {
	httpBody := `{"message":"m"}`
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(strings.NewReader(httpBody)),
	}
	err := CheckResponse(res).(*ErrorResponse)

	if err == nil {
		t.Errorf("Expected error response.")
	}

	want := &ErrorResponse{
		Response: res,
		Message:  "m",
		Errors:   "",
	}
	if !errors.Is(err, want) {
		t.Errorf("Error = %#v, want %#v", err, want)
	}
	data, err2 := ioutil.ReadAll(err.Response.Body)
	if err2 != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	if got := string(data); got != httpBody {
		t.Errorf("ErrorResponse.Response.Body = %q, want %q", got, httpBody)
	}
}

func TestErrorResponse_Error(t *testing.T) {
	res := &http.Response{Request: &http.Request{}}
	err := ErrorResponse{Message: "m", Response: res}
	if err.Error() == "" {
		t.Errorf("Expected non-empty ErrorResponse.Error()")
	}
}

func TestError_Error(t *testing.T) {
	err := Error{}
	if err.Error() == "" {
		t.Errorf("Expected non-empty Error.Error()")
	}
}

func TestFormatRateReset(t *testing.T) {
	d := 120*time.Minute + 12*time.Second
	got := formatRateReset(d)
	want := "[rate reset in 120m12s]"
	if got != want {
		t.Errorf("Format is wrong. got: %v, want: %v", got, want)
	}

	d = 14*time.Minute + 2*time.Second
	got = formatRateReset(d)
	want = "[rate reset in 14m02s]"
	if got != want {
		t.Errorf("Format is wrong. got: %v, want: %v", got, want)
	}

	d = 2*time.Minute + 2*time.Second
	got = formatRateReset(d)
	want = "[rate reset in 2m02s]"
	if got != want {
		t.Errorf("Format is wrong. got: %v, want: %v", got, want)
	}

	d = 12 * time.Second
	got = formatRateReset(d)
	want = "[rate reset in 12s]"
	if got != want {
		t.Errorf("Format is wrong. got: %v, want: %v", got, want)
	}

	d = -1 * (2*time.Hour + 2*time.Second)
	got = formatRateReset(d)
	want = "[rate limit was reset 120m02s ago]"
	if got != want {
		t.Errorf("Format is wrong. got: %v, want: %v", got, want)
	}
}

func TestRateLimitError(t *testing.T) {
	u, err := url.Parse("https://example.com")
	if err != nil {
		t.Fatal(err)
	}

	r := &RateLimitError{
		Response: &http.Response{
			Request:    &http.Request{Method: "PUT", URL: u},
			StatusCode: http.StatusTooManyRequests,
		},
		Message: "<msg>",
	}
	if got, want := r.Error(), "PUT https://example.com: 429 <msg> [rate limit was reset"; !strings.Contains(got, want) {
		t.Errorf("RateLimitError = %q, want %q", got, want)
	}
}

func TestAcceptedError(t *testing.T) {
	a := &AcceptedError{}
	if got, want := a.Error(), "try again later"; !strings.Contains(got, want) {
		t.Errorf("AcceptedError = %q, want %q", got, want)
	}
}

func TestAbuseRateLimitError(t *testing.T) {
	u, err := url.Parse("https://example.com")
	if err != nil {
		t.Fatal(err)
	}

	r := &AbuseRateLimitError{
		Response: &http.Response{
			Request:    &http.Request{Method: "PUT", URL: u},
			StatusCode: http.StatusTooManyRequests,
		},
		Message: "<msg>",
	}
	if got, want := r.Error(), "PUT https://example.com: 429 <msg>"; got != want {
		t.Errorf("AbuseRateLimitError = %q, want %q", got, want)
	}
}

func TestAddOptions_QueryValues(t *testing.T) {
	if _, err := addOptions("yo", ""); err == nil {
		t.Error("addOptions err = nil, want error")
	}
}

//func TestBareDo_returnsOpenBody(t *testing.T) {
//
//	client, mux, _, teardown := setup()
//	defer teardown()
//
//	expectedBody := "Hello from the other side !"
//
//	mux.HandleFunc("/test-url", func(w http.ResponseWriter, r *http.Request) {
//		testMethod(t, r, "GET")
//		fmt.Fprint(w, expectedBody)
//	})
//
//	ctx := context.Background()
//	req, err := client.NewRequest("GET", "test-url", nil)
//	if err != nil {
//		t.Fatalf("client.NewRequest returned error: %v", err)
//	}
//
//	resp, err := client.BareDo(ctx, req)
//	if err != nil {
//		t.Fatalf("client.BareDo returned error: %v", err)
//	}
//
//	got, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		t.Fatalf("ioutil.ReadAll returned error: %v", err)
//	}
//	if string(got) != expectedBody {
//		t.Fatalf("Expected %q, got %q", expectedBody, string(got))
//	}
//	if err := resp.Body.Close(); err != nil {
//		t.Fatalf("resp.Body.Close() returned error: %v", err)
//	}
//}
// Test function under NewRequest failure and then s.client.Do failure.
// Method f should be a regular call that would normally succeed, but
// should return an error when NewRequest or s.client.Do fails.
func testNewRequestAndDoFailure(t *testing.T, methodName string, client *Client, f func() (*Response, error)) {
	t.Helper()
	if methodName == "" {
		t.Error("testNewRequestAndDoFailure: must supply method methodName")
	}

	client.BaseURL.Path = ""
	resp, err := f()
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' %v resp = %#v, want nil", methodName, resp)
	}
	if err == nil {
		t.Errorf("client.BaseURL.Path='' %v err = nil, want error", methodName)
	}

	//client.BaseURL.Path = "/api-v3/"
	//client.rateLimits[0].Reset.Time = time.Now().Add(10 * time.Minute)
	resp, err = f()
	if want := http.StatusForbidden; resp == nil || resp.Response.StatusCode != want {
		if resp != nil {
			t.Errorf("rate.Reset.Time > now %v resp = %#v, want StatusCode=%v", methodName, resp.Response, want)
		} else {
			t.Errorf("rate.Reset.Time > now %v resp = nil, want StatusCode=%v", methodName, want)
		}
	}
	if err == nil {
		t.Errorf("rate.Reset.Time > now %v err = nil, want error", methodName)
	}
}
