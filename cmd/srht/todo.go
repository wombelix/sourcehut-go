// Copyright 2020 The Sourcehut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"git.sr.ht/~samwhited/sourcehut-go"
	"git.sr.ht/~samwhited/sourcehut-go/todo"
	"mellium.im/cli"
)

func todoCmd(srhtClient sourcehut.Client, env envVars) (*cli.Command, error) {
	client, err := todo.NewClient(
		todo.SrhtClient(srhtClient),
		todo.Base(env.todo),
	)
	if err != nil {
		return nil, err
	}

	return &cli.Command{
		Usage:       "todo <command> [options]",
		Description: "Manipulate issue trackers.",
		Commands: []*cli.Command{
			getTODOUserCmd(client),
			listTrackersCmd(client),
			todoVersionCmd(client),
		},
		Run: func(c *cli.Command, _ ...string) error {
			c.Help()
			return nil
		},
	}, nil
}

func getTODOUserCmd(client *todo.Client) *cli.Command {
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

func todoVersionCmd(client *todo.Client) *cli.Command {
	return &cli.Command{
		Usage:       "version",
		Description: "Shows the version of the todo endpoint.",
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

func listTrackersCmd(client *todo.Client) *cli.Command {
	return &cli.Command{
		Usage:       "trackers [username]",
		Description: "List all issue trackers owned by the given username or the authenticated user.",
		Run: func(c *cli.Command, args ...string) error {
			var user string
			if len(args) > 0 {
				user = args[0]
			}

			trackers, err := client.Trackers(user)
			if err != nil {
				return err
			}
			for trackers.Next() {
				// TODO: formatting
				fmt.Printf("%+v\n", trackers.Tracker())
			}
			return trackers.Err()
		},
	}
}
