// SPDX-FileCopyrightText: 2019 The SourceHut API Contributors
//
// SPDX-License-Identifier: BSD-2-Clause

// Package git provides easy API access to Sourcehut Git repositories.
package git

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"git.sr.ht/~wombelix/sourcehut-go"
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
	_, err := c.do("GET", "version", "", nil, &ver)
	return ver.Version, err
}

// Repo returns information about a specific repository owned by the provided
// username.
// If an empty username is provided, the authenticated user is used.
func (c *Client) Repo(username, repo string) (*Repo, error) {
	p := "repos"
	if username != "" {
		p = url.PathEscape(username) + "/repos"
	}
	p = path.Join(p, url.PathEscape(repo))

	newRepo := &Repo{}
	_, err := c.do("GET", p, "", nil, newRepo)
	if err != nil {
		return nil, err
	}

	return newRepo, nil
}

// DeleteRepo removes a repository.
func (c *Client) DeleteRepo(repo string) error {
	_, err := c.do("DELETE", path.Join("repos", url.PathEscape(repo)), "", nil, nil)
	return err
}

// NewRepo creates and returns a new repository from the provided template.
func (c *Client) NewRepo(name, description string, visibility RepoVisibility) (*Repo, error) {
	jsonRepo, err := json.Marshal(struct {
		Name string `json:"name"`
		Desc string `json:"description"`
		Visi string `json:"visibility"`
	}{
		Name: name,
		Desc: description,
		Visi: string(visibility),
	})
	if err != nil {
		return nil, err
	}

	newRepo := &Repo{}
	_, err = c.do("POST", "repos", "application/json", bytes.NewReader(jsonRepo), newRepo)
	if err != nil {
		return nil, err
	}
	return newRepo, nil
}

// UpdateRepo updates an existing repository.
// Only the name, description, and visibility will be modified.
//
// If repo.Name differs from oldName, a redirect from the old name to the new
// name.
func (c *Client) UpdateRepo(oldName string, repo *Repo) error {
	updateData := make(map[string]interface{})

	// Only include name if it's different from oldName (for renaming)
	if repo.Name != "" && repo.Name != oldName {
		updateData["name"] = repo.Name
	}

	// Always include description, allows empty as well
	updateData["description"] = repo.Description

	if repo.Visibility != "" {
		// Validate visibility value
		if repo.Visibility != VisibilityPublic && repo.Visibility != VisibilityUnlisted && repo.Visibility != VisibilityPrivate {
			return fmt.Errorf("invalid visibility: %s (must be public, unlisted, or private)", repo.Visibility)
		}
		updateData["visibility"] = string(repo.Visibility)
	}

	jsonRepo, err := json.Marshal(updateData)
	if err != nil {
		return err
	}

	p := path.Join("repos", url.PathEscape(oldName))
	_, err = c.do("PUT", p, "application/json", bytes.NewReader(jsonRepo), nil)
	return err
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
	_, err := c.do("GET", path.Join("user", username), "", nil, &user)
	return user, err
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
