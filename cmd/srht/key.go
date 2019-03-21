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

var (
	errWrongArgs = fmt.Errorf("Wrong number of arguments")
)

func keyCmd(srhtClient sourcehut.Client, env envVars) (*cli.Command, error) {
	client, err := meta.NewClient(
		meta.SrhtClient(srhtClient),
		meta.Base(env.meta),
	)
	if err != nil {
		return nil, err
	}

	return &cli.Command{
		Usage:       "key <command> [options]",
		Description: "Account SSH key commands.",
		Commands: []*cli.Command{
			deleteSSHKeyCmd(client),
			getSSHKeyCmd(client),
			listSSHKeyCmd(client),
			newSSHKeyCmd(client),
		},
		Run: func(c *cli.Command, _ ...string) error {
			c.Help()
			return nil
		},
	}, nil
}

func getSSHKeyCmd(client *meta.Client) *cli.Command {
	return &cli.Command{
		Usage:       "get <id>",
		Description: `Show the SSH key with the given ID.`,
		Run: func(c *cli.Command, args ...string) error {
			if len(args) != 1 {
				c.Help()
				return errWrongArgs
			}
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			k, err := client.GetSSHKey(id)
			if err != nil {
				return err
			}
			// TODO: format?
			fmt.Printf("%+v\n", k)
			return nil
		},
	}
}

func deleteSSHKeyCmd(client *meta.Client) *cli.Command {
	return &cli.Command{
		Usage:       "delete <id>",
		Description: `Delete the SSH key with the given ID.`,
		Run: func(c *cli.Command, args ...string) error {
			if len(args) != 1 {
				c.Help()
				return errWrongArgs
			}
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			return client.DeleteSSHKey(id)
		},
	}
}

func newSSHKeyCmd(client *meta.Client) *cli.Command {
	return &cli.Command{
		Usage:       "new <key (authorized_keys format)>",
		Description: `Authorize a new SSH key.`,
		Run: func(c *cli.Command, args ...string) error {
			if len(args) != 1 {
				c.Help()
				return errWrongArgs
			}

			k, err := client.NewSSHKey(args[0])
			if err != nil {
				return err
			}
			// TODO: format?
			fmt.Printf("%+v\n", k)
			return nil
		},
	}
}

func listSSHKeyCmd(client *meta.Client) *cli.Command {
	return &cli.Command{
		Usage:       "list",
		Description: `List all authorized SSH keys.`,
		Run: func(c *cli.Command, args ...string) error {
			if len(args) != 0 {
				c.Help()
				return errWrongArgs
			}

			iter, err := client.ListSSHKeys()
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
