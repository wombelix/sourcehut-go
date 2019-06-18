// Copyright 2019 The Sourcehut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package paste

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
