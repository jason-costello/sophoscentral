package sophoscentral

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Error         string
	Message       string
	CorrelationID string
	Code          string
	CreatedAt     time.Time
	RequestID     string
	DocURL        string
}


type Pages struct{
	Current int `json:"current,omitempty"`
	FromKey string `json:"fromKey,omitempty"`
	NextKey string `json:"nextKey,omitempty"`
	Size int `json:"size,omitempty"`
	Total int `json:"total,omitempty"`
	Items int `json:"items,omitempty"`
	MaxSize int `json:"maxSize,omitempty"`
}


var ErrInvalidTenantID = errors.New("invalid tenant id")
var ErrInvalidOrganizationID = errors.New("invalid organization id")
var ErrInvalidPartnerID = errors.New("invalid partner id")
var ErrAlertID = errors.New("invalid alert id")
var ErrUnmarshalFailed = errors.New("failed to unmarshal")
var ErrMarshalFailed = errors.New("failed to marshal")
var ErrFailedToCreateRequest = errors.New("failed to create new request")
var ErrHttpDo = errors.New("failed to action request")
var ErrReadBody = errors.New("failed to read body")
var Err500Returned = errors.New("500 type status code returned")
var Err400Returned = errors.New("400 type status code returned")
var ErrInvalidQueryParams = errors.New("query params failed to verify")

func MakeRequest(hc *http.Client, req *http.Request)([]byte, error){

	resp, err := hc.Do(req)
	if err != nil {
		return  nil, fmt.Errorf("%s: %w", ErrHttpDo, err)
	}
	if resp.StatusCode >= 400 &&  resp.StatusCode < 500 {
		return nil, fmt.Errorf("%s: %w", Err400Returned, errors.New(resp.Status))
	}

	if resp.StatusCode >= 500 &&  resp.StatusCode < 600 {
		return nil, fmt.Errorf("%s: %w", Err500Returned, errors.New(resp.Status))
	}

	// defer body.Close()
	defer io.CopyN(io.Discard, resp.Body, 1000)
	b, err := io.ReadAll(resp.Body)
	if err != nil{
		return nil, err
	}


	return b, nil


}
func MustMakeRequest(hc *http.Client, req *http.Request) []byte {
	fmt.Println("making request: " , req.URL.String())
	resp, err := hc.Do(req)
	if err != nil {
		return  nil
	}
	if resp.StatusCode >= 400 &&  resp.StatusCode < 500 {
		return nil
	}

	if resp.StatusCode >= 500 &&  resp.StatusCode < 600 {
		return nil
	}

	// defer body.Close()
	defer io.CopyN(io.Discard, resp.Body, 1000)
	b, err := io.ReadAll(resp.Body)
	if err != nil{
		return nil
	}
	return b


}
