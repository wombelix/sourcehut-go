// SPDX-FileCopyrightText: 2024 Dominik Wombacher <dominik@wombacher.cc>
//
// SPDX-License-Identifier: BSD-2-Clause

package logger

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	logLevelEnVar = "SRHT_LOGLEVEL"
)

var (
	Log *log.Logger
)

func init() {
	Log = log.New()

	Log.SetOutput(os.Stdout)

	switch strings.ToUpper(os.Getenv(logLevelEnVar)) {
	default:
		Log.SetLevel(log.InfoLevel)
	case "INFO":
		Log.SetLevel(log.InfoLevel)
	case "ERROR":
		Log.SetLevel(log.ErrorLevel)
	case "DEBUG":
		Log.SetLevel(log.DebugLevel)
		Log.SetReportCaller(true)
	}
}
