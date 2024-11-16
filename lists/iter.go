// SPDX-FileCopyrightText: 2019 The SourceHut API Contributors
//
// SPDX-License-Identifier: BSD-2-Clause

package lists

import (
	"git.sr.ht/~wombelix/sourcehut-go"
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
