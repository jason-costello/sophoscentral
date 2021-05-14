package sophoscentral

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

var ErrNilToken = errors.New("token is nil")
var ErrEmptyToken = errors.New("token is empty")

type Client struct {
	ctx          context.Context
	logger       *logrus.Logger
	token        *oauth2.Token
	baseURL      *url.URL
	httpClient   *http.Client
	Partner      *PartnerService
	Organization *OrganizationService
	Tenant       *TenantService
}

// NewClient returns a SophosCentral client that can be used to access all functionality
// of the documented Api.
func NewClient(ctx context.Context, httpClient *http.Client, token *oauth2.Token, logger *logrus.Logger, options ...func(*Client)) (*Client, error) {

	c := Client{}

	for _, option := range options {
		option(&c)
	}

	if c.ctx == nil {
		if ctx == nil {
			ctx = context.Background()
		}
		c.ctx = ctx
	}

	if c.httpClient == nil {
		if httpClient == nil {
			httpClient = http.DefaultClient
		}
		c.httpClient = httpClient
	}

	if c.logger == nil {
		if logger == nil {
			logger = logrus.New()
		}
		c.logger = logger
	}

	if c.token == nil {
		if token == nil {
			return nil, ErrNilToken
		}
		if token.AccessToken == "" {
			return nil, ErrEmptyToken
		}
		c.token = token
	}

	if c.Partner == nil {
		c.Partner = &PartnerService{}
	}

	if c.Organization == nil {
		c.Organization = &OrganizationService{}
	}

	if c.Tenant == nil {
		c.Tenant = &TenantService{}
	}

	var err error
	c.baseURL, err = url.Parse("https://api.central.sophos.com")
	if err != nil {
		return nil, err
	}

	return &c, nil
}

type EntityResponse struct {
	ID       string   `json:"id"`
	IDType   string   `json:"idType"`
	ApiHosts ApiHosts `json:"apiHosts"`
}
type ApiHosts struct {
	Global     string `json:"global,omitempty"`
	DataRegion string `json:"dataRegion,omitempty"`
}

func unmarshalEntityUUID(b []byte) (EntityResponse, error) {
	var er EntityResponse

	err := json.Unmarshal(b, &er)
	if err != nil {
		return EntityResponse{}, err
	}
	return er, nil

}

// WhoAmI maps directly to the whoami/v1 api call
// and returns a guid that identifies the type of entity that has
// authenticated
func (c *Client) WhoAmI() error {

	u := c.baseURL.String() + "/whoami/v1"
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return fmt.Errorf("failed to create whoami request: %w", err)
	}
	c.token.SetAuthHeader(req)
	b, err := MakeRequest(c.httpClient, req)
	if err != nil {
		return   fmt.Errorf("%s: %w", ErrHttpDo, err)
	}
	er, err := unmarshalEntityUUID(b)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if err := c.setTypeID(er); err != nil {
		return err
	}

	return nil
}

func (c *Client) setTypeID(er EntityResponse) error {

	var err error
	switch er.IDType {
	case "partner":
		c.Partner.ID, err = uuid.Parse(er.ID)
		if err != nil {
			return fmt.Errorf("invalid partner id - failed to parase uuid: %w", err)
		}
	case "organization":
		c.Organization.ID, err = uuid.Parse(er.ID)
		if err != nil {
			return fmt.Errorf("invalid organization id - failed to parase uuid: %w", err)
		}
	case "tenant":
		c.Tenant.ID, err = uuid.Parse(er.ID)
		if err != nil {
			return fmt.Errorf("invalid tenant id - failed to parase uuid: %w", err)
		}

	}

	return nil

}
