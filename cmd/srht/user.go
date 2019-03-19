// Copyright 2019 The SourceHut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"git.sr.ht/~samwhited/sourcehut-go"
	"git.sr.ht/~samwhited/sourcehut-go/meta"
	"mellium.im/cli"
)

func userCmd(srhtClient sourcehut.Client, env envVars) (*cli.Command, error) {
	client, err := meta.NewClient(
		meta.SrhtClient(srhtClient),
		meta.Base(env.meta),
	)
	if err != nil {
		return nil, err
	}

	return &cli.Command{
		Usage:       "user <command> [options]",
		Description: "Get account information",
		Commands: []*cli.Command{
			getUserCmd(client),
			metaVersionCmd(client),
		},
		Run: func(c *cli.Command, _ ...string) error {
			c.Help()
			return nil
		},
	}, nil
}

func getUserCmd(client *meta.Client) *cli.Command {
	return &cli.Command{
		Usage:       "get",
		Description: `Show the authenticated users profile`,
		Run: func(c *cli.Command, _ ...string) error {
			user, err := client.GetUser()
			if err != nil {
				return err
			}
			// TODO: format?
			fmt.Printf("%+v\n", user)
			return nil
		},
	}
}

func metaVersionCmd(client *meta.Client) *cli.Command {
	return &cli.Command{
		Usage:       "version",
		Description: "Shows the version of the meta endpoint",
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