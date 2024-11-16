// SPDX-FileCopyrightText: 2019 The SourceHut API Contributors
//
// SPDX-License-Identifier: BSD-2-Clause

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
