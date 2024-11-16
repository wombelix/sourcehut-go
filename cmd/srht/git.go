// SPDX-FileCopyrightText: 2019 The SourceHut API Contributors
//
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"fmt"

	"git.sr.ht/~wombelix/sourcehut-go"
	"git.sr.ht/~wombelix/sourcehut-go/git"
	"mellium.im/cli"
)

func gitCmd(srhtClient sourcehut.Client, env envVars) (*cli.Command, error) {
	client, err := git.NewClient(
		git.SrhtClient(srhtClient),
		git.Base(env.git),
	)
	if err != nil {
		return nil, err
	}

	return &cli.Command{
		Usage:       "git <command> [options]",
		Description: "Manipulate Git repos.",
		Commands: []*cli.Command{
			gitReposCmd(client),
			gitVersionCmd(client),
		},
		Run: func(c *cli.Command, _ ...string) error {
			c.Help()
			return nil
		},
	}, nil
}

func gitVersionCmd(client *git.Client) *cli.Command {
	return &cli.Command{
		Usage:       "version",
		Description: "Shows the version of the Git endpoint.",
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

func gitReposCmd(client *git.Client) *cli.Command {
	return &cli.Command{
		Usage:       "repos username",
		Description: "List all of the users repos.",
		Run: func(c *cli.Command, args ...string) error {
			username := ""
			switch len(args) {
			case 1:
				username = args[0]
			default:
				c.Help()
				return nil
			}

			repos, err := client.Repos(username)
			if err != nil {
				return err
			}
			for repos.Next() {
				// TODO: formatting
				fmt.Printf("%+v\n", repos.Repo())
			}
			return repos.Err()
		},
	}
}
