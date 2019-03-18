// Copyright 2019 The SourceHut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package meta

import (
	"git.sr.ht/~samwhited/sourcehut-go"
)

// User expands on the standard user struct.
type User struct {
	sourcehut.User

	UsePGPKey string `json:"use_pgp_key"`
}
