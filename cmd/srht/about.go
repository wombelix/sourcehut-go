// Copyright 2019 The Sourcehut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"mellium.im/cli"
)

func aboutCmd(w io.Writer, version, commit string, env envVars) *cli.Command {
	return &cli.Command{
		Usage:       "about",
		Description: "Show information about srht.",
		Run: func(c *cli.Command, _ ...string) error {
			fmt.Fprintf(w, `SourceHut (%s)

version:     %s
git hash:    %s
go version:  %s
go compiler: %s
platform:    %s/%s

Environment:

%v
`,
				os.Args[0],
				version, commit,
				runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH,
				env,
			)
			return nil
		},
	}
}
