// Copyright 2019 The SourceHut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

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
				if !strings.HasPrefix(username, "~") {
					username = "~" + username
				}
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
