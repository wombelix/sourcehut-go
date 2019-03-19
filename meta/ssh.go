// Copyright 2019 The SourceHut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package meta

import (
	"bytes"
	"encoding/json"
	"strconv"
)

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
