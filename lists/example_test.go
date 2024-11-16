// SPDX-FileCopyrightText: 2019 The SourceHut API Contributors
//
// SPDX-License-Identifier: BSD-2-Clause

package lists_test

import (
	"log"

	"git.sr.ht/~wombelix/sourcehut-go"
	"git.sr.ht/~wombelix/sourcehut-go/lists"
)

func ExamplePostIter() {
	srhtClient := sourcehut.NewClient(sourcehut.Token("<personal access token>"))
	listClient, _ := lists.NewClient(lists.SrhtClient(srhtClient))

	iter, _ := listClient.ListPosts("~sircmpwn", "sr.ht-dev")
	for iter.Next() {
		p := iter.Post()
		log.Printf("Post %d: %q\n", p.ID, p.Subject)
	}
	if err := iter.Err(); err != nil {
		log.Fatalf("Error fetching posts: %q", err)
	}
}
