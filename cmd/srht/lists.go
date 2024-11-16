// SPDX-FileCopyrightText: 2019 The SourceHut API Contributors
//
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"fmt"
	"strings"

	"git.sr.ht/~samwhited/sourcehut-go"
	"git.sr.ht/~samwhited/sourcehut-go/lists"
	"mellium.im/cli"
)

func listsCmd(srhtClient sourcehut.Client, env envVars) (*cli.Command, error) {
	client, err := lists.NewClient(
		lists.SrhtClient(srhtClient),
		lists.Base(env.lists),
	)
	if err != nil {
		return nil, err
	}

	return &cli.Command{
		Usage:       "lists <command> [options]",
		Description: "Manipulate mailing lists.",
		Commands: []*cli.Command{
			getListsUserCmd(client),
			listEmailsCmd(client),
			listPostsCmd(client),
			listsVersionCmd(client),
		},
		Run: func(c *cli.Command, _ ...string) error {
			c.Help()
			return nil
		},
	}, nil
}

func getListsUserCmd(client *lists.Client) *cli.Command {
	return &cli.Command{
		Usage:       "user [username]",
		Description: `Show the named or authenticated user's profile.`,
		Run: func(c *cli.Command, args ...string) error {
			username := ""
			switch len(args) {
			case 0:
			case 1:
				username = args[0]
			default:
				c.Help()
				return nil
			}
			user, err := client.GetUser(username)
			if err != nil {
				return err
			}
			// TODO: format?
			fmt.Printf("%+v\n", user)
			return nil
		},
	}
}

func listsVersionCmd(client *lists.Client) *cli.Command {
	return &cli.Command{
		Usage:       "version",
		Description: "Shows the version of the lists endpoint.",
		Run: func(c *cli.Command, ids ...string) error {
			ver, err := client.Version()
			if err != nil {
				return err
			}
			fmt.Println(ver)
			return nil
		},
	}
}

func listEmailsCmd(client *lists.Client) *cli.Command {
	return &cli.Command{
		Usage:       "emails username",
		Description: "List all emails sent by the given user.",
		Run: func(c *cli.Command, args ...string) error {
			username := ""
			switch len(args) {
			case 1:
				username = args[0]
			default:
				c.Help()
				return nil
			}

			posts, err := client.ListEmails(username)
			if err != nil {
				return err
			}
			for posts.Next() {
				// TODO: formatting
				fmt.Printf("%+v\n", posts.Post())
			}
			return posts.Err()
		},
	}
}

func listPostsCmd(client *lists.Client) *cli.Command {
	return &cli.Command{
		Usage:       "posts username/listname",
		Description: "List all emails to the list owned by the given username.",
		Run: func(c *cli.Command, args ...string) error {
			parts := make([]string, 2)
			switch len(args) {
			case 1:
				parts = strings.SplitN(args[0], "/", 2)
			default:
				c.Help()
				return nil
			}

			posts, err := client.ListPosts(parts[0], parts[1])
			if err != nil {
				return err
			}
			for posts.Next() {
				// TODO: formatting
				fmt.Printf("%+v\n", posts.Post())
			}
			return posts.Err()
		},
	}
}
