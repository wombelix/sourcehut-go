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
type Option func(*Transport)

// UserAgent returns an option that configures the client to use the provided
// user agent when making API requests.
func UserAgent(ua string) Option {
	return func(rt *Transport) {
		rt.userAgent = ua
	}
}

// Token returns an option that configures the client to use the provided access
// token when making API requests.
// If no token is provided, the client can only make requests that do not
// require authentication.
func Token(t string) Option {
	return func(rt *Transport) {
		rt.accessToken = t
	}
}

// RoundTripper returns an option that configures the client to use the provided
// http.RoundTripper for HTTP requests.
// If unspecified, http.DefaultTransport is used.
func RoundTripper(rt http.RoundTripper) Option {
	return func(t *Transport) {
		t.baseRT = rt
	}
}

// Transport is an http.RoundTripper wrapping a base RoundTripper and adding a
// SourceHut API authorization header or user agent.
//
// Transport is a low-level mechanism.
// Most code will use the NewClient method instead.
type Transport struct {
	userAgent   string
	accessToken string
	baseRT      http.RoundTripper
}

// NewTransport returns an http.RoundTripper that is configured with the
// provided options.
func NewTransport(opts ...Option) *Transport {
	rt := Transport{}
	for _, opt := range opts {
		opt(&rt)
	}
	return &rt
}

// RoundTrip authorizes and authenticates the request with an
// access token from Transport's Source.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqBodyClosed := false
	if req.Body != nil {
		defer func() {
			if !reqBodyClosed {
				req.Body.Close()
			}
		}()
	}

	// TODO: clone req.

	if t.accessToken != "" {
		// TODO: do we need to sanitize this to prevent header injection in case the
		// user takes this value from somewhere they shouldn't?
		req.Header.Set("Authorization", "token "+t.accessToken)
	}
	if t.userAgent != "" {
		// TODO: do we need to sanitize this to prevent header injection in case the
		// user takes this value from somewhere they shouldn't?
		req.Header.Set("User-Agent", t.userAgent)
	}
	return t.base().RoundTrip(req)
}

// CancelRequest cancels an in-flight request by closing its connection.
//
// If the underlying http.RoundTripper does not support cancelation,
// CancelRequest is a noop.
func (t *Transport) CancelRequest(req *http.Request) {
	if c, ok := t.base().(interface {
		CancelRequest(*http.Request)
	}); ok {
		c.CancelRequest(req)
	}
}

func (t *Transport) base() http.RoundTripper {
	if t.baseRT != nil {
		return t.baseRT
	}
	return http.DefaultTransport
}

// Client is like http.Client except that it knows how to authenticate to the
// SourceHut API.
type Client struct {
	httpClient *http.Client
}

// NewBaseClient returns a new SourceHut API client configured to use the
// provided http.Client to perform HTTP requests.
//
// To add authentication use NewClient or provide a base client that is
// authenticated with the SourceHut API.
func NewBaseClient(base *http.Client) Client {
	return Client{httpClient: base}
}

// NewClient returns a new SourceHut API client configured with the provided
// options.
func NewClient(opts ...Option) Client {
	return Client{
		httpClient: &http.Client{
			Transport: NewTransport(opts...),
		},
	}
}

// Do sends an API request and returns the API response.
// The response is unmarshaled into v if successful, or returned as an error
// value if an API error has occured.
func (c Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.do(req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// List returns an iterator that can transparently make API requests to a
// paginated endpoint.
// Each item will be decoded into the value returned from a call to d.
// If d is nil, a map[string]interface{} is created for each item.
//
// No HTTP request will be issued until iteration is started by a call to Next.
func (c Client) List(req *http.Request, d func() interface{}) *Iter {
	return &Iter{req: req, c: c, into: d}
}

func (c Client) do(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		e := struct {
			Errors Errors `json:"errors"`
		}{}
		err = json.NewDecoder(resp.Body).Decode(&e)
		if err != nil {
			return resp, err
		}
		switch len(e.Errors) {
		case 0:
			return resp, nil
		case 1:
			return resp, e.Errors[0]
		}
		return resp, e.Errors
	}

	return resp, nil
}
