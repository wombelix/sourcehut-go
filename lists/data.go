// Copyright 2019 The Sourcehut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package lists

import (
	"time"

	"git.sr.ht/~samwhited/sourcehut-go"
)

// ShortList represents the unexpanded form of a mailing list.
type ShortList struct {
	Name  string              `json:"name"`
	Owner sourcehut.ShortUser `json:"owner"`
}

// List represents the expanded form of a mailing list.
type List struct {
	ShortList

	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Desc    string    `json:"description"`
	Perms   struct {
		// TODO: "browse", "reply", and "post" are valid permissions.
		// Make a type for this.
		NonSubscriber []string `json:"nonsubscriber"`
		Subscriber    []string `json:"subscriber"`
		Account       []string `json:"account"`
	} `json:"permissions"`
}

// ShortPost represents the unexpanded form of an email message.
type ShortPost struct {
	ID        int64                `json:"id"`
	Created   time.Time            `json:"created"`
	List      ShortList            `json:"list"`
	MessageID string               `json:"message_id"`
	ParentID  int64                `json:"parent_id"`
	Sender    *sourcehut.ShortUser `json:"sender"`
	Subject   string               `json:"subject"`
	ThreadID  int64                `json:"thread_id"`
}

// Post represents the expanded form of an email message.
type Post struct {
	ShortPost

	Patch        bool   `json:"is_patch"`
	PullRequest  bool   `json:"is_request_pull"`
	Replies      int64  `json:"replies"`
	Participants int64  `json:"participants"`
	Envelope     string `json:"envelope"`
}
