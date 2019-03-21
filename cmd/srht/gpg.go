// Copyright 2019 The SourceHut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strconv"

	"git.sr.ht/~samwhited/sourcehut-go"
	"git.sr.ht/~samwhited/sourcehut-go/meta"
	"mellium.im/cli"
)

func pgpCmd(srhtClient sourcehut.Client, env envVars) (*cli.Command, error) {
	client, err := meta.NewClient(
		meta.SrhtClient(srhtClient),
		meta.Base(env.meta),
	)
	if err != nil {
		return nil, err
	}

	return &cli.Command{
		Usage:       "pgp <command> [options]",
		Description: "Account PGP key commands.",
		Commands: []*cli.Command{
			deletePGPKeyCmd(client),
			getPGPKeyCmd(client),
			listPGPKeyCmd(client),
			newPGPKeyCmd(client),
		},
		Run: func(c *cli.Command, _ ...string) error {
			c.Help()
			return nil
		},
	}, nil
}

func getPGPKeyCmd(client *meta.Client) *cli.Command {
	return &cli.Command{
		Usage:       "get <id>",
		Description: `Show the PGP key with the given ID.`,
		Run: func(c *cli.Command, args ...string) error {
			if len(args) != 1 {
				c.Help()
				return errWrongArgs
			}
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			k, err := client.GetPGPKey(id)
			if err != nil {
				return err
			}
			// TODO: format?
			fmt.Printf("%+v\n", k)
			return nil
		},
	}
}

func deletePGPKeyCmd(client *meta.Client) *cli.Command {
	return &cli.Command{
		Usage:       "delete <id>",
		Description: `Delete the PGP key with the given ID.`,
		Run: func(c *cli.Command, args ...string) error {
			if len(args) != 1 {
				c.Help()
				return errWrongArgs
			}
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			return client.DeletePGPKey(id)
		},
	}
}

func newPGPKeyCmd(client *meta.Client) *cli.Command {
	return &cli.Command{
		Usage:       "new <key (authorized_keys format)>",
		Description: `Authorize a new PGP key.`,
		Run: func(c *cli.Command, args ...string) error {
			if len(args) != 1 {
				c.Help()
				return errWrongArgs
			}

			k, err := client.NewPGPKey(args[0])
			if err != nil {
				return err
			}
			// TODO: format?
			fmt.Printf("%+v\n", k)
			return nil
		},
	}
}

func listPGPKeyCmd(client *meta.Client) *cli.Command {
	return &cli.Command{
		Usage:       "list",
		Description: `List all authorized PGP keys.`,
		Run: func(c *cli.Command, args ...string) error {
			if len(args) != 0 {
				c.Help()
				return errWrongArgs
			}

			iter, err := client.ListPGPKeys()
			if err != nil {
				return err
			}
			for iter.Next() {
				// TODO: format?
				fmt.Printf("%+v\n", iter.Key())
			}
			return iter.Err()
		},
	}
}
