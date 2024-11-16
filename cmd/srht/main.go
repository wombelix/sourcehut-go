// SPDX-FileCopyrightText: 2019 The SourceHut API Contributors
//
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"log"
	"os"

	"git.sr.ht/~wombelix/sourcehut-go"
	"mellium.im/cli"
)

// Configuration that is manipulated at build time.
var (
	commit  string
	version string
)

const userAgent = "git.sr.ht/~wombelix/sourcehut-go/cmd/srht"

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	env := newEnv()
	srhtClient := sourcehut.NewClient(
		sourcehut.Token(env.token),
		sourcehut.UserAgent(userAgent),
	)

	user, err := userCmd(srhtClient, env)
	if err != nil {
		logger.Fatal("Meta URL could not be parsed.")
	}
	key, err := keyCmd(srhtClient, env)
	if err != nil {
		logger.Fatal("Meta URL could not be parsed.")
	}
	pgp, err := pgpCmd(srhtClient, env)
	if err != nil {
		logger.Fatal("Meta URL could not be parsed.")
	}
	paste, err := pasteCmd(srhtClient, env)
	if err != nil {
		logger.Fatal("Paste URL could not be parsed.")
	}
	lists, err := listsCmd(srhtClient, env)
	if err != nil {
		logger.Fatal("Lists URL could not be parsed.")
	}
	git, err := gitCmd(srhtClient, env)
	if err != nil {
		logger.Fatal("Git URL could not be parsed.")
	}
	todo, err := todoCmd(srhtClient, env)
	if err != nil {
		logger.Fatal("TODO URL could not be parsed.")
	}

	// Commands
	cmds := &cli.Command{
		Usage: os.Args[0],
	}
	cmds.Commands = []*cli.Command{
		aboutCmd(os.Stdout, version, commit, env),
		git,
		key,
		lists,
		paste,
		pgp,
		user,
		todo,
		cli.Help(cmds),
	}

	err = cmds.Exec(os.Args[1:]...)
	switch err {
	case cli.ErrInvalidCmd:
		helpCmd := cli.Help(cmds)
		if err = helpCmd.Exec(); err != nil {
			logger.Fatalf("Error showing help: %q", err)
		}
		os.Exit(1)
	case cli.ErrNoRun:
		helpCmd := cli.Help(cmds)
		if err = helpCmd.Exec(); err != nil {
			logger.Fatalf("Error showing help: %q", err)
		}
	case nil:
	default:
		logger.Fatal(err)
	}
}
