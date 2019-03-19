// Copyright 2019 The SourceHut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

// Package meta provides easy API access to SourceHut account info.
package meta

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"git.sr.ht/~samwhited/sourcehut-go"
)

// BaseURL is the default public SourceHut user API URL.
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

// Base returns an option that configures the public SourceHut API URL.
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

// Client handles communication with the user related methods of the SourceHut
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

// GetSSH returns the SSH key with the provided ID.
func (c *Client) GetSSHKey(id int64) (SSHKey, error) {
	key := SSHKey{}
	_, err := c.do("GET", "user/ssh-keys/"+strconv.FormatInt(id, 10), "", nil, &key)
	return key, err
}

// DeleteSSH deletes the SSH key with the provided ID.
func (c *Client) DeleteSSHKey(id int64) error {
	_, err := c.do("DELETE", "user/ssh-keys/"+strconv.FormatInt(id, 10), "", nil, nil)
	return err
}

// NewSSHKey creates a new SSH key.
// The key should be in authorized_keys format.
func (c *Client) NewSSHKey(k string) (SSHKey, error) {
	key := SSHKey{}
	jsonKey, err := json.Marshal(struct {
		Key string `json:"ssh-key"`
	}{
		Key: k,
	})
	if err != nil {
		return key, err
	}

	_, err = c.do("POST", "user/ssh-keys", "application/json", bytes.NewReader(jsonKey), &key)
	return key, err
}

// ListSSHKeys returns an iterator over all SSH keys authorized on the users
// account.
func (c *Client) ListSSHKeys() (SSHKeyIter, error) {
	return c.sshKeys("GET", "user/ssh-keys", nil)
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

func (c *Client) sshKeys(method, u string, body io.Reader) (SSHKeyIter, error) {
	u = c.baseURL.String() + u
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return SSHKeyIter{}, err
	}
	iter := c.srhtClient.List(req, func() interface{} {
		return &SSHKey{}
	})
	return SSHKeyIter{Iter: iter}, nil
}
