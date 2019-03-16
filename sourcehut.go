// Copyright 2019 The SourceHut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

// Package sourcehut provides access to the SourceHut HTTP API.
package sourcehut

import (
	"encoding/json"
	"net/http"
)

// Option is used to configure a SourceHut API client.
type Option func(*Client)

// UserAgent returns an option that configures the client to use the provided
// user agent when making API requests.
func UserAgent(ua string) Option {
	return func(c *Client) {
		c.userAgent = ua
	}
}

// Token returns an option that configures the client to use the provided access
// token when making API requests.
// If no token is provided, the client can only make requests that do not
// require authentication.
func Token(t string) Option {
	return func(c *Client) {
		c.accessToken = t
	}
}

// HTTPClient returns an option that configures the client to use the provided
// http.Client for HTTP requests.
// By default, http.DefaultClient is used.
func HTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

// Client is like http.Client except that it knows how to authenticate to the
// SourceHut API.
type Client struct {
	userAgent   string
	accessToken string
	httpClient  *http.Client
}

// NewClient returns a new SourceHut API client configured with the provided
// options.
func NewClient(opts ...Option) *Client {
	c := Client{
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

// Do sends an API request and returns the API response.
// The response is unmarshaled into v if successful, or returned as an error
// value if an API error has occured.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if c.accessToken != "" {
		// TODO: do we need to sanitize this to prevent header injection in case the
		// user takes this value from somewhere they shouldn't?
		req.Header.Set("Authorization", c.accessToken)
	}
	if c.userAgent != "" {
		// TODO: do we need to sanitize this to prevent header injection in case the
		// user takes this value from somewhere they shouldn't?
		req.Header.Set("User-Agent", c.userAgent)
	}

	// TODO: decode common fields and check if this is an error.
	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return nil, err
	}

	return &Response{
		Response: resp,
	}, nil
}
