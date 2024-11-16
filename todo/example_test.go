// SPDX-FileCopyrightText: 2020 The SourceHut API Contributors
//
// SPDX-License-Identifier: BSD-2-Clause

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
