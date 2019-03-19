// Copyright 2019 The SourceHut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"git.sr.ht/~samwhited/sourcehut-go"
	"git.sr.ht/~samwhited/sourcehut-go/paste"
	"mellium.im/cli"
)

func pasteCmd(srhtClient sourcehut.Client, env envVars) (*cli.Command, error) {
	client, err := paste.NewClient(
		paste.SrhtClient(srhtClient),
		paste.Base(env.paste),
	)
	if err != nil {
		return nil, err
	}

	return &cli.Command{
		Usage:       "paste <command> [options]",
		Description: "Create or download pastes",
		Commands: []*cli.Command{
			getBlob(client),
			getPasteCmd(client),
			listPasteCmd(client),
			pasteVersionCmd(client),
		},
		Run: func(c *cli.Command, _ ...string) error {
			c.Help()
			return nil
		},
	}, nil
}

func getBlob(client *paste.Client) *cli.Command {
	var (
		saveBlobs bool
		zipName   string
	)
	flags := flag.NewFlagSet("blob", flag.ContinueOnError)
	flags.BoolVar(&saveBlobs, "O", false, "Write blob contents to the current working directory")
	flags.StringVar(&zipName, "o", "", "Write blob contents to the named zip file")

	return &cli.Command{
		Usage: "blob [options] id [id2 id3…]",
		Flags: flags,
		Description: `Print or download a file from a paste

By default file contents are written to stdout.
If blobs are written to a zip file or to the current working directory and a
file with the same name already exists, it will be truncated.
`,
		Run: func(c *cli.Command, args ...string) error {
			err := flags.Parse(args)
			if err != nil {
				return err
			}

			ids := flags.Args()
			if len(ids) == 0 {
				c.Help()
				return nil
			}

			// Create a zip file if -o was specified.
			var zipWriter *zip.Writer
			if zipName != "" {
				zipFile, err := os.Create(zipName)
				if err != nil {
					return fmt.Errorf("Error creating output file %q: %q", zipName, err)
				}
				defer zipFile.Close()
				zipWriter = zip.NewWriter(zipFile)
				defer zipWriter.Close()
			}

			for _, id := range ids {
				blob, err := client.GetBlob(id)
				if err != nil {
					// TODO: should this exit immediately, finish but return a non-zero
					// status, etc?
					fmt.Fprintf(os.Stderr, "Error fetching blob %s: %q\n", id, err)
					continue
				}

				if zipWriter != nil {
					w, err := zipWriter.CreateHeader(&zip.FileHeader{
						Name:     blob.ID,
						Modified: blob.Created,
					})
					if err != nil {
						return err
					}
					_, err = io.WriteString(w, blob.Contents)
					if err != nil {
						return fmt.Errorf("Error writing blob %s to %q: %q", blob.ID, zipName, err)
					}
					if saveBlobs {
						err = ioutil.WriteFile(blob.ID, []byte(blob.Contents), 0644)
						if err != nil {
							return fmt.Errorf("Error writing blob %s to disk: %q", blob.ID, err)
						}
					}
					continue
				}

				// TODO: how should blobs be formatted?
				fmt.Printf("%+v\n", blob)
			}
			return nil
		},
	}
}

func pasteVersionCmd(client *paste.Client) *cli.Command {
	return &cli.Command{
		Usage:       "version",
		Description: "Shows the version of the paste endpoint",
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

func listPasteCmd(client *paste.Client) *cli.Command {
	return &cli.Command{
		Usage:       "list",
		Description: "List pastes owned by the authenticated user",
		Run: func(c *cli.Command, ids ...string) error {
			iter, err := client.List()
			if err != nil {
				return err
			}

			for iter.Next() {
				paste := iter.Paste()
				// TODO: how should pastes be formatted?
				fmt.Printf("%+v\n", paste)
			}
			return iter.Err()
		},
	}
}

func getPasteCmd(client *paste.Client) *cli.Command {
	return &cli.Command{
		Usage:       "get id [id2 id3…]",
		Description: "Show one or more pastes",
		Run: func(c *cli.Command, ids ...string) error {
			if len(ids) == 0 {
				c.Help()
				return nil
			}

			for _, id := range ids {
				paste, err := client.Get(id)
				if err != nil {
					// TODO: should this exit immediately, finish but return a non-zero
					// status, etc?
					fmt.Fprintf(os.Stderr, "Error fetching paste %s: %q\n", id, err)
					continue
				}

				// TODO: how should pastes be formatted?
				fmt.Printf("%+v\n", paste)
			}
			return nil
		},
	}
}
