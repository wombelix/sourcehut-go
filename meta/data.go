// Copyright 2019 The Sourcehut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package meta

import (
	"time"

	"git.sr.ht/~samwhited/sourcehut-go"
)

// User expands on the standard user struct.
type User struct {
	sourcehut.User

	UsePGPKey string `json:"use_pgp_key"`
}

// SSHKey contains information about an SSH key.
type SSHKey struct {
	ID          int64               `json:"id"`
	Authorized  time.Time           `json:"authorized"`
	Comment     string              `json:"comment"`
	Fingerprint string              `json:"fingerprint"`
	Key         string              `json:"key"`
	Owner       sourcehut.ShortUser `json:"owner"`
	LastUsed    time.Time           `json:"last_used"`
}

// PGPKey contains information about an PGP key.
type PGPKey struct {
	ID         int64               `json:"id"`
	Authorized time.Time           `json:"authorized"`
	Email      string              `json:"email"`
	KeyID      string              `json:"key_id"`
	Key        string              `json:"key"`
	LastUsed   time.Time           `json:"last_used"`
	Owner      sourcehut.ShortUser `json:"owner"`
}

// AuditLog represents a single entry in the audit log.
type AuditLog struct {
	ID      int64     `json:"id"`
	IP      string    `json:"ip"`
	Action  string    `json:"action"`
	Details string    `json:"details"`
	Created time.Time `json:"created"`
}
