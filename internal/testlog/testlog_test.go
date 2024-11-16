// SPDX-FileCopyrightText: 2019 The SourceHut API Contributors
//
// SPDX-License-Identifier: BSD-2-Clause

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
