// Copyright 2019 The Sourcehut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package git

import (
	"git.sr.ht/~samwhited/sourcehut-go"
)

// RepoIter is used for iterating over a collection of repos.
type RepoIter struct {
	*sourcehut.Iter
}

// Repo returns the repo which the iterator is currently pointing to.
func (i RepoIter) Repo() *Repo {
	return i.Current().(*Repo)
}
