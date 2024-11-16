// SPDX-FileCopyrightText: 2020 The SourceHut API Contributors
//
// SPDX-License-Identifier: BSD-2-Clause

package todo

import (
	"time"

	"git.sr.ht/~wombelix/sourcehut-go"
)

// ShortTracker represents the unexpanded form of an issue tracker.
type ShortTracker struct {
	Name    string              `json:"name"`
	Owner   sourcehut.ShortUser `json:"owner"`
	Created time.Time           `json:"created"`
	Updated time.Time           `json:"updated"`
}

// Tracker represents the expanded form of an issue tracker.
type Tracker struct {
	ShortTracker

	Desc  string `json:"description"`
	Perms struct {
		Anonymous []string `json:"anonymous"`
		Submitter []string `json:"submitter"`
		User      []string `json:"user"`
	} `json:"default_permissions"`
}
