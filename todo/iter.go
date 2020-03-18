// Copyright 2020 The Sourcehut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

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
