// Copyright 2019 The Sourcehut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

// Package testlog is a log.Logger that proxies to the Log function on a
// testing.T.
//
// It is used to group log messages under the tests that generated them in test
// output and to only show those messages if the test that generated them
// failed.
package testlog

import (
	"log"
	"testing"
)

// New returns a new logger that logs to the provided testing.T.
func New(t testing.TB) *log.Logger {
	t.Helper()
	return log.New(testWriter{TB: t}, t.Name()+" ", log.LstdFlags|log.Lshortfile|log.LUTC)
}

type testWriter struct {
	testing.TB
}

func (tw testWriter) Write(p []byte) (int, error) {
	tw.Helper()
	tw.Logf("%s", p)
	return len(p), nil
}
