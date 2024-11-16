// SPDX-FileCopyrightText: 2020 The SourceHut API Contributors
//
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"fmt"
	"strings"

	"git.sr.ht/~wombelix/sourcehut-go"
	"git.sr.ht/~wombelix/sourcehut-go/todo"
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
		Usage:       "trackers [username] [tracker]",
		Description: "List issue trackers owned by the given username or the authenticated user.",
		Run: func(c *cli.Command, args ...string) error {
			var (
				user        string
				trackerName string
			)
			// TODO: can we name a tracker with a "~"? If sourcehut doesn't prevent
			// this it could lead to a bug here. Maybe make username an option on all
			// commands that take it? eg. -u?
			if len(args) > 0 {
				if strings.HasPrefix(args[0], "~") {
					user = args[0]
					if len(args) > 1 {
						trackerName = args[1]
					}
				} else {
					trackerName = args[0]
				}
			}

			if trackerName == "" {
				trackers, err := client.Trackers(user)
				if err != nil {
					return err
				}
				for trackers.Next() {
					// TODO: formatting
					fmt.Printf("%+v\n", trackers.Tracker())
				}
				return trackers.Err()
			}

			tracker, err := client.Tracker(user, trackerName)
			if err != nil {
				return err
			}
			// TODO: formatting
			fmt.Printf("%+v\n", tracker)
			return nil
		},
	}
}
