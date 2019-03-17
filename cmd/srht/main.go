// Copyright 2019 The SourceHut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"

	"git.sr.ht/~samwhited/sourcehut-go"
	"mellium.im/cli"
)

// Configuration that is manipulated at build time.
var (
	commit  string
	version string
)

const userAgent = "git.sr.ht/~samwhited/sourcehut-go/cmd/srht"

func main() {
	env := newEnv()
	srhtClient := sourcehut.NewClient(
		sourcehut.Token(env.token),
		sourcehut.UserAgent(userAgent),
	)

	paste, err := pasteCmd(srhtClient, env)
	if err != nil {
		log.Fatal("Paste URL could not be parsed.")
	}

	// Commands
	cmds := &cli.Command{
		Usage: os.Args[0],
	}
	cmds.Commands = []*cli.Command{
		aboutCmd(os.Stdout, version, commit, env),
		paste,
		cli.Help(cmds),
	}

	err = cmds.Exec(os.Args[1:]...)
	switch err {
	case cli.ErrInvalidCmd:
		helpCmd := cli.Help(cmds)
		helpCmd.Exec()
		os.Exit(1)
	case cli.ErrNoRun:
		helpCmd := cli.Help(cmds)
		helpCmd.Exec()
	case nil:
	default:
		log.Fatal(err)
	}
}
