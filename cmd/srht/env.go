// SPDX-FileCopyrightText: 2019 The SourceHut API Contributors
//
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"fmt"
	"os"

	"git.sr.ht/~samwhited/sourcehut-go/git"
	"git.sr.ht/~samwhited/sourcehut-go/lists"
	"git.sr.ht/~samwhited/sourcehut-go/meta"
	"git.sr.ht/~samwhited/sourcehut-go/paste"
	"git.sr.ht/~samwhited/sourcehut-go/todo"
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
	todo  string
}

func (env envVars) String() string {
	redactedToken := "…"
	switch {
	case len(env.token) == 0:
		redactedToken = ""
	case len(env.token) > 8:
		redactedToken = env.token[:8] + redactedToken
	}
	return fmt.Sprintf(`SRHT_TOKEN      = %q
SRHT_META_BASE  = %q
SRHT_PASTE_BASE = %q
SRHT_LISTS_BASE = %q
SRHT_GIT_BASE   = %q
SRHT_TODO_BASE  = %q
`, redactedToken, env.meta, env.paste, env.lists, env.git)
}

func newEnv() envVars {
	return envVars{
		token: os.Getenv("SRHT_TOKEN"),
		paste: defEnv("SRHT_PASTE_BASE", paste.BaseURL),
		meta:  defEnv("SRHT_META_BASE", meta.BaseURL),
		lists: defEnv("SRHT_LISTS_BASE", lists.BaseURL),
		git:   defEnv("SRHT_GIT_BASE", git.BaseURL),
		todo:  defEnv("SRHT_TODO_BASE", todo.BaseURL),
	}
}
