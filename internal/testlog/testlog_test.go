// Copyright 2019 The SourceHut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package testlog_test

import (
	"git.sr.ht/~samwhited/sourcehut-go/internal/testlog"
)

import (
	"testing"
)

func TestLog(t *testing.T) {
	logger := testlog.New(t)
	logger.Println("Logging should not cause a test failure")
}
