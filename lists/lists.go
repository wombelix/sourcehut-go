// Package lists provides easy API access to SourceHut mailing lists.
package lists

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"git.sr.ht/~samwhited/sourcehut-go"
)

// BaseURL is the default public SourceHut mailing lists API URL.
// It is exported for convenience.
const BaseURL = "https://lists.sr.ht/api/"

// Client handles communication with the mailing lists related methods of the
// SourceHut API.
//
// API docs: https://man.sr.ht/lists.sr.ht/api.md
type Client struct {
	baseURL    *url.URL
	srhtClient *sourcehut.Client
}

// NewClient returns a new mailing list API client.
// If baseURL is empty, it defaults to the public SourceHut mailing lists API.
// If baseURL does not have a trailing slash, one is added automatically.
// If srhtClient is nil, a new client is created using http.DefaultClient.
func NewClient(baseURL string, srhtClient *sourcehut.Client) (*Client, error) {
	if baseURL == "" {
		baseURL = BaseURL
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}

	if srhtClient == nil {
		// TODO: access
		srhtClient = sourcehut.NewClient(nil, "")
	}

	return &Client{
		baseURL:    u,
		srhtClient: srhtClient,
	}, nil
}

// Version returns the version of the API.
//
// API docs: https://man.sr.ht/api-conventions.md#get-apiversion
func (c *Client) Version() (string, error) {
	var ver struct {
		Version string `json:"version"`
	}
	_, err := c.do("GET", "version", nil, &ver)
	return ver.Version, err
}

// do sends an API request and returns the API response.
// The response is unmarshaled into v if successful, or returned as an error
// value if an API error has occured.
func (c *Client) do(method, u string, body io.Reader, v interface{}) (*sourcehut.Response, error) {
	u = c.baseURL.String() + u
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}
	return c.srhtClient.Do(req, v)
}
