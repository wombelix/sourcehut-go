// Copyright 2019 The Sourcehut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package git

import (
	"os"
	"time"
)

// RepoVisibility is the visibility level of a repo.
type RepoVisibility string

// Supported visibility levels.
const (
	VisibilityPublic   RepoVisibility = "public"
	VisibilityUnlisted RepoVisibility = "unlisted"
	VisibilityPrivate  RepoVisibility = "private"
)

// Repo represents a repository.
type Repo struct {
	ID          int64          `json:"id"`
	Created     time.Time      `json:"created"`
	Subject     string         `json:"subject"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Visibility  RepoVisibility `json:"visibility"`
}

// Commit is a single commit in a repo.
type Commit struct {
	ID        string    `json:"id"`
	ShortID   string    `json:"short_id"`
	Author    Author    `json:"author"`
	Committer Author    `json:"committer"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Tree      string    `json:"tree"`
	Signature *struct {
		Signature string `json:"signature"`
		Data      string `json:"data"`
	} `json:"signature"`
	Parents []string `json:"Parents"`
}

// Author is information about the author or committer of a commit.
type Author struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// TreeType is a supported type of tree (tree or blob).
type TreeType string

// Valid and supported tree type's.
const (
	TypeTree TreeType = "tree"
	TypeBlob TreeType = "blob"
)

// Tree is a tree within a commit.
type Tree struct {
	ID      string `json:"id"`
	ShortID string `json:"short_id"`
	Entries []struct {
		ID   string      `json:"id"`
		Name string      `json:"name"`
		Type TreeType    `json:"type"`
		Mode os.FileMode `json:"mode"`
	} `json:"entries"`
}

// Ref is a reference to an object in a repo.
type Ref struct {
	Name   string `json:"name"`
	Target string `json:"target"`
}
