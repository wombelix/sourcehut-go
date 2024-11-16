// SPDX-FileCopyrightText: 2019 The SourceHut API Contributors
//
// SPDX-License-Identifier: BSD-2-Clause

// Package lists provides easy API access to Sourcehut mailing lists.
package lists

import (
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"git.sr.ht/~wombelix/sourcehut-go"
)

// BaseURL is the default public Sourcehut mailing lists API URL.
// It is exported for convenience.
const BaseURL = "https://lists.sr.ht/api/"

// Option is used to configure an API client.
type Option func(*Client) error

// SrhtClient returns an option that configures the client to use the provided
// sourcehut.Client for API requests.
// If unspecified, a default sourcehut.Client (with no options of its own) is
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

// Client handles communication with the mailing lists related methods of the
// Sourcehut API.
//
// API docs: https://man.sr.ht/lists.sr.ht/api.md
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

// List returns an iterator over all mailing lists owned by the provided
// username.
// If an empty username is provided, the authenticated user is used.
func (c *Client) List(username string) (ListIter, error) {
	path := "lists"
	if username != "" {
		path = "user/" + url.PathEscape(username) + "/lists"
	}
	return c.lists("GET", path, nil)
}

// ListPosts returns the posts in a mailing list owned by the given username.
func (c *Client) ListPosts(username, listname string) (PostIter, error) {
	p := path.Join("user", username, "lists", listname, "posts")
	return c.posts("GET", p, nil)
}

// GetUser returns information about the provided username, or the currently
// authenticated user if the username is empty.
func (c *Client) GetUser(username string) (sourcehut.User, error) {
	user := sourcehut.User{}
	_, err := c.do("GET", path.Join("user", username), nil, &user)
	return user, err
}

// ListEmails returns all emails sent by the provided user.
func (c *Client) ListEmails(username string) (PostIter, error) {
	return c.posts("GET", path.Join("user", username, "emails"), nil)
}

func (c *Client) do(method, u string, body io.Reader, v interface{}) (*http.Response, error) {
	u = c.baseURL.String() + u
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}
	return c.srhtClient.Do(req, v)
}

func (c *Client) lists(method, u string, body io.Reader) (ListIter, error) {
	u = c.baseURL.String() + u
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return ListIter{}, err
	}
	iter := c.srhtClient.List(req, func() interface{} {
		return &List{}
	})
	return ListIter{Iter: iter}, nil
}

func (c *Client) posts(method, u string, body io.Reader) (PostIter, error) {
	u = c.baseURL.String() + u
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return PostIter{}, err
	}
	iter := c.srhtClient.List(req, func() interface{} {
		return &Post{}
	})
	return PostIter{Iter: iter}, nil
}
