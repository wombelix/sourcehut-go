// SPDX-FileCopyrightText: 2019 The SourceHut API Contributors
//
// SPDX-License-Identifier: BSD-2-Clause

// Package paste provides easy API access to Sourcehut file pasting.
package paste

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"git.sr.ht/~samwhited/sourcehut-go"
)

// BaseURL is the default public Sourcehut paste API URL.
// It is exported for convenience.
const BaseURL = "https://paste.sr.ht/api/"

// Option is used to configure an API client.
type Option func(*Client) error

// SrhtClient returns an option that configures the client to use the provided
// sourcehut.Client for API requests.
// If unspecified, the default sourcehut.Client (with no options of its own) is
// used.
func SrhtClient(client sourcehut.Client) Option {
	return func(c *Client) error {
		c.srhtClient = client
		return nil
	}
}

// Base returns an option that configures the public Sourcehut API URL.
//
// If base does not have a trailing slash, one is added automatically.
// If unspecified, BaseURL is used.
func Base(base string) Option {
	return func(c *Client) error {
		if base == "" {
			base = BaseURL
		}
		u, err := url.Parse(base)
		if err != nil {
			return err
		}
		if !strings.HasSuffix(u.Path, "/") {
			u.Path += "/"
		}
		c.baseURL = u
		return nil
	}
}

// Client handles communication with the paste related methods of the Sourcehut
// API.
//
// API docs: https://man.sr.ht/paste.sr.ht/api.md
type Client struct {
	baseURL    *url.URL
	srhtClient sourcehut.Client
}

// NewClient returns a new paste API client.
func NewClient(opts ...Option) (*Client, error) {
	u, err := url.Parse(BaseURL)
	if err != nil {
		return nil, err
	}

	c := Client{
		baseURL: u,

		// TODO: with no access token, is this behavior useful?
		// Maybe this should be a required argument and not an option.
		srhtClient: sourcehut.NewClient(),
	}
	for _, opt := range opts {
		if err = opt(&c); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

// Version returns the version of the API.
//
// API docs: https://man.sr.ht/api-conventions.md#get-apiversion
func (c *Client) Version() (string, error) {
	var ver struct {
		Version string `json:"version"`
	}
	_, err := c.do("GET", "version", "", nil, &ver)
	return ver.Version, err
}

// List returns an iterator over all pastes owned by the authenticated user.
func (c *Client) List() (Iter, error) {
	return c.list("GET", "pastes", nil)
}

// Get returns information about a paste with the given ID.
func (c *Client) Get(id string) (Paste, error) {
	p := Paste{}
	_, err := c.do("GET", "pastes/"+url.PathEscape(id), "", nil, &p)
	return p, err
}

// New creates an new paste from the list of files.
func (c *Client) New(f Files) (Paste, error) {
	p := Paste{}
	buf := &bytes.Buffer{}
	e := json.NewEncoder(buf)
	err := e.Encode(struct {
		Files Files `json:"files"`
	}{
		Files: f,
	})
	if err != nil {
		return Paste{}, err
	}
	_, err = c.do("POST", "pastes", "application/json", buf, &p)
	return p, err
}

// GetBlob returns information about a particular file in a paste.
func (c *Client) GetBlob(id string) (Blob, error) {
	p := Blob{}
	_, err := c.do("GET", "blobs/"+url.PathEscape(id), "", nil, &p)
	return p, err
}

func (c *Client) do(method, u, contentType string, body io.Reader, v interface{}) (*http.Response, error) {
	u = c.baseURL.String() + u
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return c.srhtClient.Do(req, v)
}

func (c *Client) list(method, u string, body io.Reader) (Iter, error) {
	u = c.baseURL.String() + u
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return Iter{}, err
	}
	iter := c.srhtClient.List(req, func() interface{} {
		return &Paste{}
	})
	return Iter{Iter: iter}, nil
}
