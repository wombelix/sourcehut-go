// SPDX-FileCopyrightText: 2020 The SourceHut API Contributors
//
// SPDX-License-Identifier: BSD-2-Clause

package todo

import (
	"git.sr.ht/~samwhited/sourcehut-go"
)

// TrackerIter is used for iterating over a collection of issue trackers.
type TrackerIter struct {
	*sourcehut.Iter
}

// Tracker returns the issue tracker which the iterator is currently pointing
// to.
func (i TrackerIter) Tracker() *Tracker {
	return i.Current().(*Tracker)
}
