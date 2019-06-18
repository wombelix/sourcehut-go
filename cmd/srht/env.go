// Copyright 2019 The Sourcehut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"git.sr.ht/~samwhited/sourcehut-go/git"
	"git.sr.ht/~samwhited/sourcehut-go/lists"
	"git.sr.ht/~samwhited/sourcehut-go/meta"
	"git.sr.ht/~samwhited/sourcehut-go/paste"
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
	meta  string
	lists string
	git   string
}

func (env envVars) String() string {
	redactedToken := "…"
	switch {
	case len(env.token) == 0:
		redactedToken = ""
	case len(env.token) > 8:
		redactedToken = env.token[:8] + redactedToken
	}
	return fmt.Sprintf(`SRHT_TOKEN	= %q
SRHT_META_BASE = %q
SRHT_PASTE_BASE	= %q
SRHT_LISTS_BASE = %q
SRHT_GIT_BASE = %q
`, redactedToken, env.meta, env.paste, env.lists, env.git)
}

func newEnv() envVars {
	return envVars{
		token: os.Getenv("SRHT_TOKEN"),
		paste: defEnv("SRHT_PASTE_BASE", paste.BaseURL),
		meta:  defEnv("SRHT_META_BASE", meta.BaseURL),
		lists: defEnv("SRHT_LISTS_BASE", lists.BaseURL),
		git:   defEnv("SRHT_GIT_BASE", git.BaseURL),
	}
}
