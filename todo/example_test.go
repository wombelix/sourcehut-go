// Copyright 2020 The Sourcehut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package todo_test

import (
	"log"

	"git.sr.ht/~samwhited/sourcehut-go"
	"git.sr.ht/~samwhited/sourcehut-go/todo"
)

func ExampleTrackerIter() {
	srhtClient := sourcehut.NewClient(sourcehut.Token("<personal access token>"))
	todoClient, _ := todo.NewClient(todo.SrhtClient(srhtClient))

	iter, _ := todoClient.Trackers("~sircmpwn")
	for iter.Next() {
		p := iter.Tracker()
		log.Printf("Tracker %s: %s\n", p.Name, p.Desc)
	}
	if err := iter.Err(); err != nil {
		log.Fatalf("Error fetching posts: %q", err)
	}
}
