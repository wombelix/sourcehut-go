// SPDX-FileCopyrightText: 2019 The SourceHut API Contributors
//
// SPDX-License-Identifier: BSD-2-Clause

package git

import (
	"git.sr.ht/~wombelix/sourcehut-go"
)

// RepoIter is used for iterating over a collection of repos.
type RepoIter struct {
	*sourcehut.Iter
}

// Repo returns the repo which the iterator is currently pointing to.
func (i RepoIter) Repo() *Repo {
	return i.Current().(*Repo)
}
