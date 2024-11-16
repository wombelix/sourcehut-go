// SPDX-FileCopyrightText: 2019 The SourceHut API Contributors
//
// SPDX-License-Identifier: BSD-2-Clause

// Package meta provides easy API access to Sourcehut account info.
package meta

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"git.sr.ht/~samwhited/sourcehut-go"
)

// BaseURL is the default public Sourcehut user API URL.
// It is exported for convenience.
const BaseURL = "https://meta.sr.ht/api/"

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

// Client handles communication with the user related methods of the Sourcehut
// API.
//
// API docs: https://man.sr.ht/meta.sr.ht/user-api.md
type Client struct {
	baseURL    *url.URL
	srhtClient sourcehut.Client
}

// NewClient returns a new API client.
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

// GetUser returns information about the currently authenticated user.
func (c *Client) GetUser() (User, error) {
	user := User{}
	_, err := c.do("GET", "user/profile", "", nil, &user)
	return user, err
}

// ProfileParams is like sourcehut.User except that it omits the username fields
// and allows nil values for some fields that should not be updated.
type ProfileParams struct {
	Email    *string `json:"email,omitempty"`
	URL      *string `json:"url,omitempty"`
	Location *string `json:"location,omitempty"`
	Bio      *string `json:"bio,omitempty"`
}

// UpdateUser sets information about the user.
// Nil values indicate that the field should not be updated.
// If the email field is updated it will trigger a confirmation email.
func (c *Client) UpdateUser(user ProfileParams) (User, error) {
	newUser := User{}
	j, err := json.Marshal(user)
	if err != nil {
		return newUser, err
	}
	_, err = c.do("PUT", "user/profile", "application/json", bytes.NewReader(j), &newUser)
	return newUser, err
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
