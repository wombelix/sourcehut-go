// Copyright 2019 The SourceHut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

// Package sourcehut provides access to the SourceHut HTTP API.
package sourcehut

import (
	"encoding/json"
	"net/http"
)

// Client is like http.Client except that it knows how to authenticate to the
// SourceHut API.
type Client struct {
	// User agent used when communicating with the SourceHut API.
	UserAgent string

	httpClient *http.Client
}

// NewClient returns a new SourceHut API client.
// If a nil httpClient is provided, http.DefaultClient will be used.
// To use API methods which require authentication, provide an http.Client that
// will perform the authentication for you.
func NewClient(c *http.Client, accessToken string) *Client {
	if c == nil {
		c = http.DefaultClient
	}

	return &Client{
		httpClient: c,
	}
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

	// TODO: decode common fields and check if this is an error.
	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return nil, err
	}

	return &Response{
		Response: resp,
	}, nil
}
