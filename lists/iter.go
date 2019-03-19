// Copyright 2019 The SourceHut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package lists

import (
	"git.sr.ht/~samwhited/sourcehut-go"
)

// ListIter is used for iterating over a collection of mailing lists.
type ListIter struct {
	*sourcehut.Iter
}

// List returns the mailing list which the iterator is currently pointing to.
func (i ListIter) List() *List {
	return i.Current().(*List)
}

// PostIter is used for iterating over a collection of mailing list posts.
type PostIter struct {
	*sourcehut.Iter
}

// Post returns the post which the iterator is currently pointing to.
func (i PostIter) Post() *Post {
	return i.Current().(*Post)
}
