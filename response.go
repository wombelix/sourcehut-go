// Copyright 2019 The SourceHut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package sourcehut

import (
	"net/http"
)

// Response is a SourceHut API response.
// This wraps the standard http.Response and provides convenient access to
// pagination links.
//
// API docs: https://man.sr.ht/api-conventions.md
type Response struct {
	*http.Response

	Next           int      `json:"next"`
	ResultsPerPage int      `json:"results_per_page"`
	Results        []string `json:"results"`
	Total          int      `json:"total"`
	Errors         []Error  `json:"errors"`
}

// Ensure that the build fails if Error
var _ error = (*Error)(nil)

// Error represents an individual error returned from a SourceHut API call.
//
// API docs: https://man.sr.ht/api-conventions.md#error-responses
type Error struct {
	Field  string
	Reason string
}

// Error satisfies the error interface for the Error type.
func (err *Error) Error() string {
	return err.Reason
}
