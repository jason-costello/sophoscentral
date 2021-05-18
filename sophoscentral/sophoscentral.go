//go:generate go run gen-accessors.go
//go:generate go run gen-stringify-test.go
package sophoscentral
// design influenced by https://github.com/google/go-github

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-querystring/query"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

const(
defaultBaseURL = "https://api.central.sophos.com/"
userAgent = "rax-sophoscentral"

)
var errNonNilContext = errors.New("context must be non-nil")

// Client manages communcation with the Sophos Central Api
type Client struct{
	Token *oauth2.Token

	client *http.Client

	BaseURL *url.URL
	UserAgent string
	common service
	// Services used to interact with different parts of Sophos Central api

	Common *CommonService
	Endpoints *EndpointService
	LiveDiscover *LiveDiscoverService
	Organization *OrganizationService
	Partner *PartnerService
	WhoAmI *WhoAmIService

}

type service struct{
	client *Client
}

// ListByPageOffset specifies the parameters to methods that support pagination by page offset value
type ListByPageOffset struct{
Page int `url:"page,omitempty"`
PageSize int `url:"pageSize,omitempty"`
PageTotal int `url:"pageTotal,omitempty"`
}

// ListByFromKeyOptions specifies the parameters to methods that support pagination by from-key value
type ListByFromKeyOptions struct{
	PageFromKey string `url:"pageFromKey,omitempty"`
	PageSize int `url:"pageSize,omitempty"`
	PageTotal bool `url:"pageTotal,omitempty"`
}


type PagesByOffset struct{
	Current *int `json:"current,omitempty"`
	Size *int `json:"size,omitempty"`
	Total *int `json:"total,omitempty"`
	Items *int `json:"items,omitempty"`
	MaxSize *int `json:"maxSize,omitempty"`

}
type PagesByFromKey struct {
		FromKey *string `json:"fromKey,omitempty"`
		NextKey *string `json:"nextKey,omitempty"`
		Size    *int    `json:"size,omitempty"`
		Total   *int    `json:"total,omitempty"`
		Items   *int    `json:"items,omitempty"`
		MaxSize *int    `json:"maxSize,omitempty"`
	}


func DecodePageKey(encodedKey string) string{


	decodeString, err := base64.StdEncoding.DecodeString(encodedKey) // to []byte
	if err!=nil || decodeString == nil{
		return ""
	}
	key := string(decodeString)
	if key[0] == '[' {
		key = key[1:]
	}
	if key[len(key)-1] == ']'{
		key = key[:len(key)-1]
	}
	return key



}
// addOptions adds the parameters in opts as URL query parameters to s. opts
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opts interface{}) (string, error) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opts)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}


// NewClient returns a new Sophos Central API client. If a nil httpClient is
// provided, a new http.Client will be used. To use API methods which require
// authentication, provide an http.Client that will perform the authentication
// for you (such as that provided by the golang.org/x/oauth2 library).
func NewClient(ctx context.Context, hc *http.Client, token *oauth2.Token) *Client {


	if ctx == nil{
		ctx = context.Background()
	}

	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{client: hc, BaseURL: baseURL, UserAgent: userAgent}
	c.common.client = c
	c.Common = (*CommonService)(&c.common)
	c.Endpoints = (*EndpointService)(&c.common)
	c.LiveDiscover = (*LiveDiscoverService)(&c.common)
	c.Organization = (*OrganizationService)(&c.common)
	c.Partner = (*PartnerService)(&c.common)
	c.WhoAmI = (*WhoAmIService)(&c.common)
	c.Token = token

	return c
}


// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}
// Response is a GitHub API response. This wraps the standard http.Response
// returned from GitHub and provides convenient access to things like
// pagination links.
type Response struct {
	*http.Response

	// These fields provide the page values for paginating through a set of
	// results. Any or all of these may be set to the zero value for
	// responses that are not part of a paginated set, or for which there
	// are no additional pages.
	//
	// These fields support what is called "offset pagination" and should
	// be used with the ListByPageOffset struct.
	NextPage  int
	PrevPage  int
	FirstPage int
	LastPage  int

	// Additionally, some APIs support "cursor pagination" instead of offset.
	// This means that a token points directly to the next record which
	// can lead to O(1) performance compared to O(n) performance provided
	// by offset pagination.
	//
	// For APIs that support cursor pagination (such as
	// TeamsService.ListIDPGroupsInOrganization), the following field
	// will be populated to point to the next page.
	//
	// To use this token, set ListByFromKeyOptions.Page to this value before
	// calling the endpoint again.
	NextPageToken string


}
// newResponse creates a new Response for the provided http.Response.
// r must not be nil.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	//response.populatePageValues()
	return response
}

// BareDo  populatePageValues parses the HTTP Link response headers and populates the
//// various pagination link values in the Response.
//func (r *Response) populatePageValues() {
//
//	/* Page by key example
//
//	  "pages": {
//	    "fromKey": "WyI4ODc4NjA5ZS1iZmE2LTRhOGQtOGM3Zi02YjYwNTAzZjA0NWQiXQ==",
//	    "size": 500,
//	    "total": 2,
//	    "items": 929,
//	    "maxSize": 500
//	  }
//
//	 */
//
//	if links, ok := r.Response.Header["Link"]; ok && len(links) > 0 {
//		for _, link := range strings.Split(links[0], ",") {
//			segments := strings.Split(strings.TrimSpace(link), ";")
//
//			// link must at least have href and rel
//			if len(segments) < 2 {
//				continue
//			}
//
//			// ensure href is properly formatted
//			if !strings.HasPrefix(segments[0], "<") || !strings.HasSuffix(segments[0], ">") {
//				continue
//			}
//
//			// try to pull out page parameter
//			url, err := url.Parse(segments[0][1 : len(segments[0])-1])
//			if err != nil {
//				continue
//			}
//			page := url.Query().Get("page")
//			if page == "" {
//				continue
//			}
//
//			for _, segment := range segments[1:] {
//				switch strings.TrimSpace(segment) {
//				case `rel="next"`:
//					if r.NextPage, err = strconv.Atoi(page); err != nil {
//						r.NextPageToken = page
//					}
//				case `rel="prev"`:
//					r.PrevPage, _ = strconv.Atoi(page)
//				case `rel="first"`:
//					r.FirstPage, _ = strconv.Atoi(page)
//				case `rel="last"`:
//					r.LastPage, _ = strconv.Atoi(page)
//				}
//
//			}
//		}
//	}
//}
// The provided ctx must be non-nil, if it is nil an error is returned. If it is
// canceled or times out, ctx.Err() will be returned.
func (c *Client) BareDo(ctx context.Context, req *http.Request) (*Response, error) {
	if ctx == nil {
		return nil, errNonNilContext
	}
	req = withContext(ctx, req)


	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// If the error type is *url.Error, sanitize its URL before returning.
		if e, ok := err.(*url.Error); ok {
			if tURL, err := url.Parse(e.URL); err == nil {
				e.URL = sanitizeURL(tURL).String()
				return nil, e
			}
		}

		return nil, err
	}

	response := newResponse(resp)


	err = CheckResponse(resp)
	if err != nil {
		defer resp.Body.Close()
		// Special case for AcceptedErrors. If an AcceptedError
		// has been encountered, the response's payload will be
		// added to the AcceptedError and returned.

		aerr, ok := err.(*AcceptedError)
		if ok {
			b, readErr := ioutil.ReadAll(resp.Body)
			if readErr != nil {
				return response, readErr
			}

			aerr.Raw = b
			err = aerr
		}
	}
	return response, err
}


// Do - The provided ctx must be non-nil, if it is nil an error is returned. If it
// is canceled or times out, ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.BareDo(ctx, req)
	if err != nil {

		if err.Error() == "Unauthorized"{
		} else {
			return resp, err
		}
	}
	defer resp.Body.Close()

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, resp.Body)
	default:
		decErr := json.NewDecoder(resp.Body).Decode(v)
		if decErr == io.EOF {
			decErr = nil // ignore EOF errors caused by empty response body
		}
		if decErr != nil {
			err = decErr
		}
	}



	return resp, err
}

// compareHttpResponse returns whether two http.Response objects are equal or not.
// Currently, only StatusCode is checked. This function is used when implementing the
// Is(error) bool interface for the custom error types in this package.
func compareHttpResponse(r1, r2 *http.Response) bool {
	if r1 == nil && r2 == nil {
		return true
	}

	if r1 != nil && r2 != nil {
		return r1.StatusCode == r2.StatusCode
	}
	return false
}



/*
An ErrorResponse reports one or more errors caused by an API request.
Sophos Central API docs: https://developer.sophos.com/intro (search page for 'Error response object')
*/
type ErrorResponse struct {
	Response *http.Response
	Errors         string    `json:"error,omitempty"`
	Message       string    `json:"message,omitempty"`
	CorrelationId string    `json:"correlationId,omitempty"`
	Code          string    `json:"code,omitempty"`
	CreatedAt     time.Time `json:"createdAt,omitempty"`
	RequestId     string    `json:"requestId,omitempty"`
	DocumentationURL        string    `json:"docUrl,omitempty"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %+v",
		r.Response.Request.Method, sanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode, r.Message, r.Errors)
}


// Is returns whether the provided error equals this error.
func (r *ErrorResponse) Is(target error) bool {
	v, ok := target.(*ErrorResponse)
	if !ok {
		return false
	}

	if r.Message != v.Message || (r.DocumentationURL != v.DocumentationURL) ||
		!compareHttpResponse(r.Response, v.Response) {
		return false
	}

	// Compare Errors.
	if len(r.Errors) != len(v.Errors) {
		return false
	}
	for idx := range r.Errors {
		if r.Errors[idx] != v.Errors[idx] {
			return false
		}
	}



	return true
}


// TwoFactorAuthError occurs when using HTTP Basic Authentication for a user
// that has two-factor authentication enabled. The request can be reattempted
// by providing a one-time password in the request.
type TwoFactorAuthError ErrorResponse

func (r *TwoFactorAuthError) Error() string { return (*ErrorResponse)(r).Error() }

// RateLimitError occurs when GitHub returns 403 Forbidden response with a rate limit
// remaining value of 0.
type RateLimitError struct {
	Rate     Rate           // Rate specifies last known rate limit for the client
	Response *http.Response // HTTP response that caused this error
	Message  string         `json:"message"` // error message
}

func (r *RateLimitError) Error() string {
	return fmt.Sprintf("%v %v: %d %v %v",
		r.Response.Request.Method, sanitizeURL(r.Response.Request.URL),
	r.Response.StatusCode, r.Message, formatRateReset(time.Until(r.Rate.Reset.Time)))
}

// Is returns whether the provided error equals this error.
func (r *RateLimitError) Is(target error) bool {
	v, ok := target.(*RateLimitError)
	if !ok {
		return false
	}

	return r.Rate == v.Rate &&
		r.Message == v.Message &&
		compareHttpResponse(r.Response, v.Response)
}

// AcceptedError occurs when GitHub returns 202 Accepted response with an
// empty body, which means a job was scheduled on the GitHub side to process
// the information needed and cache it.
// Technically, 202 Accepted is not a real error, it's just used to
// indicate that results are not ready yet, but should be available soon.
// The request can be repeated after some time.
type AcceptedError struct {
	// Raw contains the response body.
	Raw []byte
}
func withContext(ctx context.Context, req *http.Request) *http.Request {
	return req.WithContext(ctx)
}


func (*AcceptedError) Error() string {
	return "job scheduled on GitHub side; try again later"
}

// Is returns whether the provided error equals this error.
func (ae *AcceptedError) Is(target error) bool {
	v, ok := target.(*AcceptedError)
	if !ok {
		return false
	}
	return bytes.Compare(ae.Raw, v.Raw) == 0
}

// AbuseRateLimitError occurs when Sophos Central returns 429 Forbidden response"
type AbuseRateLimitError struct {
	Response *http.Response // HTTP response that caused this error
	Message  string         `json:"message"` // error message

	// RetryAfter is provided with some abuse rate limit errors. If present,
	// it is the amount of time that the client should wait before retrying.
	// Otherwise, the client should try again later (after an unspecified amount of time).
	RetryAfter *time.Duration
}

func (r *AbuseRateLimitError) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, sanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode, r.Message)
}

// Is returns whether the provided error equals this error.
func (r *AbuseRateLimitError) Is(target error) bool {
	v, ok := target.(*AbuseRateLimitError)
	if !ok {
		return false
	}

	return r.Message == v.Message &&
		r.RetryAfter == v.RetryAfter &&
		compareHttpResponse(r.Response, v.Response)
}

// sanitizeURL redacts the client_secret parameter from the URL which may be
// exposed to the user.
func sanitizeURL(uri *url.URL) *url.URL {
	if uri == nil {
		return nil
	}
	params := uri.Query()
	if len(params.Get("client_secret")) > 0 {
		params.Set("client_secret", "REDACTED")
		uri.RawQuery = params.Encode()
	}
	return uri
}

/*
Error reports more details on an individual error in an ErrorResponse.

*/
type Error struct {
	Resource string `json:"resource"` // resource on which the error occurred
	Field    string `json:"field"`    // field on which the error occurred
	Code     string `json:"code"`     // validation error code
	Message  string `json:"message"`  // Message describing the error. Errors with Code == "custom" will always have this set.
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v error caused by %v field on %v resource",
		e.Code, e.Field, e.Resource)
}

func (e *Error) UnmarshalJSON(data []byte) error {
	type aliasError Error // avoid infinite recursion by using type alias.
	if err := json.Unmarshal(data, (*aliasError)(e)); err != nil {
		return json.Unmarshal(data, &e.Message) // data can be json string.
	}
	return nil
}

// CheckResponse checks the API response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range or equal to 202 Accepted.
// API error responses are expected to have response
// body, and a JSON response body that maps to ErrorResponse.
//
// The error type will be *RateLimitError for rate limit exceeded errors,
// *AcceptedError for 202 Accepted status codes,
// and *TwoFactorAuthError for two-factor authentication errors.
func CheckResponse(r *http.Response) error {
	if r.StatusCode == http.StatusAccepted {
		return &AcceptedError{}
	}
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	// Re-populate error response body because Sophos Central error responses could be undocumented
	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	switch {
	case r.StatusCode == http.StatusUnauthorized:
		return (*TwoFactorAuthError)(errorResponse)

	case r.StatusCode == http.StatusTooManyRequests:
		abuseRateLimitError := &AbuseRateLimitError{
			Response: errorResponse.Response,
			Message:  errorResponse.Message,
		}

			retryAfterSeconds :=  getRandIntBetween(100, 750) // Error handling is noop.
			retryAfter := time.Duration(retryAfterSeconds) * time.Millisecond
			abuseRateLimitError.RetryAfter = &retryAfter

		return abuseRateLimitError

	default:
		return errorResponse
	}
}
// setRandomTimeout generates a random value between 100 and 750 milliseconds to wait until available for retry
// if error is generated for some reason when generating random value a default of 500 ms is used.
func getRandIntBetween(min, max int) int {
var delay = 0
	for delay < min {
		delay = rand.Intn(max)
	}

	return delay
}
// parseBoolResponse determines the boolean result from a GitHub API response.
// Several GitHub API methods return boolean responses indicated by the HTTP
// status code in the response (true indicated by a 204, false indicated by a
// 404). This helper function will determine that result and hide the 404
// error if present. Any other error will be returned through as-is.
//func parseBoolResponse(err error) (bool, error) {
//	if err == nil {
//		return true, nil
//	}
//
//	if err, ok := err.(*ErrorResponse); ok && err.Response.StatusCode == http.StatusNotFound {
//		// Simply false. In this one case, we do not pass the error through.
//		return false, nil
//	}
//
//	// some other real error occurred
//	return false, err
//}

// Rate represents the rate limit for the current client.
type Rate struct {
	// The number of requests per hour the client is currently limited to.
	Limit int `json:"limit"`

	// The number of remaining requests the client can make this hour.
	Remaining int `json:"remaining"`

	// The time at which the current rate limit will reset.
	Reset Timestamp `json:"reset"`
}

func (r Rate) String() string {
	return Stringify(r)
}
//
// RateLimits represents the rate limits for the current client.
type RateLimits struct {
	Core *Rate `json:"core"`

}

func (r RateLimits) String() string {
	return Stringify(r)
}

//type rateLimitCategory uint8

//const (
//	coreCategory rateLimitCategory = iota
//	searchCategory
//
//	categories // An array of this length will be able to contain all rate limit categories.
//)

// category returns the rate limit category of the endpoint, determined by Request.URL.Path.
//func category(path string) rateLimitCategory {
//	switch {
//	default:
//		return coreCategory
//	case strings.HasPrefix(path, "/search/"):
//		return searchCategory
//	}
//}

// RateLimits returns the rate limits for the current client.
//func (c *Client) RateLimits(ctx context.Context) (*RateLimits, *Response, error) {
//	req, err := c.NewRequest("GET", "rate_limit", nil)
//	if err != nil {
//		return nil, nil, err
//	}
//
//	response := new(struct {
//		Resources *RateLimits `json:"resources"`
//	})
//	resp, err := c.Do(ctx, req, response)
//	if err != nil {
//		return nil, resp, err
//	}
//
//	if response.Resources != nil {
//		c.rateMu.Lock()
//		if response.Resources.Core != nil {
//			c.rateLimits[coreCategory] = *response.Resources.Core
//		}
//		if response.Resources.Search != nil {
//			c.rateLimits[searchCategory] = *response.Resources.Search
//		}
//		c.rateMu.Unlock()
//	}
//
//	return response.Resources, resp, nil
//}




// formatRateReset formats d to look like "[rate reset in 2s]" or
// "[rate reset in 87m02s]" for the positive durations. And like "[rate limit was reset 87m02s ago]"
// for the negative cases.
func formatRateReset(d time.Duration) string {
	isNegative := d < 0
	if isNegative {
		d *= -1
	}
	secondsTotal := int(0.5 + d.Seconds())
	minutes := secondsTotal / 60
	seconds := secondsTotal - minutes*60

	var timeString string
	if minutes > 0 {
		timeString = fmt.Sprintf("%dm%02ds", minutes, seconds)
	} else {
		timeString = fmt.Sprintf("%ds", seconds)
	}

	if isNegative {
		return fmt.Sprintf("[rate limit was reset %v ago]", timeString)
	}
	return fmt.Sprintf("[rate reset in %v]", timeString)
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool { return &v }

// Int is a helper routine that allocates a new int value
// to store v and returns a pointer to it.
func Int(v int) *int { return &v }

// Int64 is a helper routine that allocates a new int64 value
// to store v and returns a pointer to it.
func Int64(v int64) *int64 { return &v }

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string { return &v }

func PrettyPrint(i interface{}) string {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}
