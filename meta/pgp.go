// Copyright 2019 The Sourcehut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package meta

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

// GetPGPKey returns the PGP key with the provided ID.
func (c *Client) GetPGPKey(id int64) (PGPKey, error) {
	key := PGPKey{}
	_, err := c.do("GET", "user/pgp-keys/"+strconv.FormatInt(id, 10), "", nil, &key)
	return key, err
}

// DeletePGPKey deletes the PGP key with the provided ID.
func (c *Client) DeletePGPKey(id int64) error {
	_, err := c.do("DELETE", "user/pgp-keys/"+strconv.FormatInt(id, 10), "", nil, nil)
	return err
}

// NewPGPKey creates a new PGP key.
// The key should be in authorized_keys format.
func (c *Client) NewPGPKey(k string) (PGPKey, error) {
	key := PGPKey{}
	jsonKey, err := json.Marshal(struct {
		Key string `json:"pgp-key"`
	}{
		Key: k,
	})
	if err != nil {
		return key, err
	}

	_, err = c.do("POST", "user/pgp-keys", "application/json", bytes.NewReader(jsonKey), &key)
	return key, err
}

// ListPGPKeys returns an iterator over all PGP keys authorized on the users
// account.
func (c *Client) ListPGPKeys() (PGPKeyIter, error) {
	return c.pgpKeys("GET", "user/pgp-keys", nil)
}

func (c *Client) pgpKeys(method, u string, body io.Reader) (PGPKeyIter, error) {
	u = c.baseURL.String() + u
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return PGPKeyIter{}, err
	}
	iter := c.srhtClient.List(req, func() interface{} {
		return &PGPKey{}
	})
	return PGPKeyIter{Iter: iter}, nil
}
