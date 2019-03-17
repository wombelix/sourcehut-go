// Copyright 2019 The SourceHut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package paste

import (
	"time"

	"git.sr.ht/~samwhited/sourcehut-go"
)

// Paste contains data about a set of files.
type Paste struct {
	ID      string              `json:"sha"`
	Created time.Time           `json:"created"`
	User    sourcehut.ShortUser `json:"user"`
	Files   []struct {
		ID   string `json:"blob_id"`
		Name string `json:"filename"`
	} `json:"files"`
}

// Blob contains data about an individual file in a paste.
type Blob struct {
	ID       string    `json:"sha"`
	Created  time.Time `json:"created"`
	Contents string    `json:"contents"`
}
