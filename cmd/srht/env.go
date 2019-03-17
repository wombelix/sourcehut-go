// Copyright 2019 The SourceHut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
)

// BUG(ssw): Tool does not load a config file or source .env files if present.

func defEnv(key, def string) string {
	env := os.Getenv(key)
	if env == "" {
		return def
	}
	return env
}

type envVars struct {
	token string
	paste string
}

func (env envVars) String() string {
	redactedToken := "…"
	if len(env.token) > 8 {
		redactedToken = env.token[:8] + redactedToken
	}
	return fmt.Sprintf(`SRHT_TOKEN	= %q
SRHT_PASTE_BASE	= %q
`, redactedToken, env.paste)
}

func newEnv() envVars {
	return envVars{
		token: os.Getenv("SRHT_TOKEN"),
		paste: defEnv("SRHT_PASTE_BASE", "https://paste.sr.ht/api/"),
	}
}
