// Copyright 2019 The SourceHut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package paste

import (
	"time"
)

import (
	"git.sr.ht/~samwhited/sourcehut-go"
)

// Iter is used for iterating over a collection of pastes.
type Iter struct {
	*sourcehut.Iter
}

// Paste returns the paste which the iterator is currently pointing to.
func (i Iter) Paste() *Paste {
	return i.Current().(*Paste)
}

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
