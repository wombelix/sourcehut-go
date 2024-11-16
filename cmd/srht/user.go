// SPDX-FileCopyrightText: 2019 The SourceHut API Contributors
//
// SPDX-License-Identifier: BSD-2-Clause

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
		Description: "Get account information.",
		Commands: []*cli.Command{
			getUserCmd(client),
			listAuditLogsCmd(client),
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
		Description: `Show the authenticated users profile.`,
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
		Description: "Shows the version of the meta endpoint.",
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

func listAuditLogsCmd(client *meta.Client) *cli.Command {
	return &cli.Command{
		Usage:       "log",
		Description: `Lists audit logs.`,
		Run: func(c *cli.Command, args ...string) error {
			if len(args) != 0 {
				c.Help()
				return errWrongArgs
			}

			iter, err := client.ListAuditLog()
			if err != nil {
				return err
			}
			for iter.Next() {
				// TODO: format?
				fmt.Printf("%+v\n", iter.Log())
			}
			return iter.Err()
		},
	}
}
