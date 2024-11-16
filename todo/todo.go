// SPDX-FileCopyrightText: 2020 The SourceHut API Contributors
//
// SPDX-License-Identifier: BSD-2-Clause

// Package todo provides easy API access to Sourcehut issue trackers.
package todo

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"git.sr.ht/~samwhited/sourcehut-go"
)

// BaseURL is the default public Sourcehut issue trackers API URL.
// It is exported for convenience.
const BaseURL = "https://todo.sr.ht/api/"

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

// Client handles communication with the issue tracker related methods of the
// Sourcehut API.
//
// API docs: https://man.sr.ht/todo.sr.ht/api.md
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

// GetUser returns information about the provided username, or the currently
// authenticated user if the username is empty.
func (c *Client) GetUser(username string) (sourcehut.User, error) {
	user := sourcehut.User{}
	_, err := c.do("GET", path.Join("user", username), "", nil, &user)
	return user, err
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

// NewTracker creates and returns a new repository from the provided template.
func (c *Client) NewTracker(name, description string) (*Tracker, error) {
	jsonTracker, err := json.Marshal(struct {
		Name string `json:"name"`
		Desc string `json:"description"`
	}{
		Name: name,
		Desc: description,
	})
	if err != nil {
		return nil, err
	}

	newTracker := &Tracker{}
	_, err = c.do("POST", "trackers", "application/json", bytes.NewReader(jsonTracker), newTracker)
	if err != nil {
		return nil, err
	}
	return newTracker, nil
}

// Tracker returns information about a specific issue tracker owned by the
// provided username.
// If an empty username is provided, the authenticated user is used.
func (c *Client) Tracker(username, tracker string) (*Tracker, error) {
	p := "trackers"
	if username != "" {
		p = "user/" + url.PathEscape(username) + "/trackers"
	}
	p = path.Join(p, url.PathEscape(tracker))

	newTracker := &Tracker{}
	_, err := c.do("GET", p, "", nil, newTracker)
	if err != nil {
		return nil, err
	}

	return newTracker, nil
}

// Trackers returns an iterator over all issue trackers owned by the provided
// username.
// If an empty username is provided, the authenticated user is used.
func (c *Client) Trackers(username string) (TrackerIter, error) {
	path := "trackers"
	if username != "" {
		path = "user/" + url.PathEscape(username) + "/trackers"
	}
	return c.trackers("GET", path, nil)
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

func (c *Client) trackers(method, u string, body io.Reader) (TrackerIter, error) {
	u = c.baseURL.String() + u
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return TrackerIter{}, err
	}
	iter := c.srhtClient.List(req, func() interface{} {
		return &Tracker{}
	})
	return TrackerIter{Iter: iter}, nil
}
