// Copyright 2019 The Sourcehut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

// Package git provides easy API access to Sourcehut Git repositories.
package git

import (
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"git.sr.ht/~samwhited/sourcehut-go"
)

// BaseURL is the default public Sourcehut API URL.
// It is exported for convenience.
const BaseURL = "https://git.sr.ht/api/"

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

// Client handles communication with the Sourcehut API.
//
// API docs: https://man.sr.ht/git.sr.ht/api.md
type Client struct {
	baseURL    *url.URL
	srhtClient sourcehut.Client
}

// NewClient returns a new mailing list API client.
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
	_, err := c.do("GET", "version", nil, &ver)
	return ver.Version, err
}

// Repos returns an iterator over all repos owned by the provided username.
// If an empty username is provided, the authenticated user is used.
func (c *Client) Repos(username string) (RepoIter, error) {
	path := "repos"
	if username != "" {
		path = url.PathEscape(username) + "/repos"
	}
	return c.repos("GET", path, nil)
}

// GetUser returns information about the provided username, or the currently
// authenticated user if the username is empty.
func (c *Client) GetUser(username string) (sourcehut.User, error) {
	user := sourcehut.User{}
	_, err := c.do("GET", path.Join("user", username), nil, &user)
	return user, err
}

func (c *Client) do(method, u string, body io.Reader, v interface{}) (*http.Response, error) {
	u = c.baseURL.String() + u
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}
	return c.srhtClient.Do(req, v)
}

func (c *Client) repos(method, u string, body io.Reader) (RepoIter, error) {
	u = c.baseURL.String() + u
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return RepoIter{}, err
	}
	iter := c.srhtClient.List(req, func() interface{} {
		return &Repo{}
	})
	return RepoIter{Iter: iter}, nil
}
