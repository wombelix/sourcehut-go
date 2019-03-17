// Copyright 2019 The SourceHut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package sourcehut

// Ensure that the build fails if Error and Errors don't implement error.
var _, _ error = (*Error)(nil), (*Errors)(nil)

// Error represents an individual error returned from a SourceHut API call.
//
// API docs: https://man.sr.ht/api-conventions.md#error-responses
type Error struct {
	Field  string
	Reason string
}

// Error satisfies the error interface for Error.
func (err *Error) Error() string {
	return err.Reason
}

// Errors is a slice of Error's that itself implements error.
type Errors []Error

// Error satisfies the error interface for Errors.
func (err Errors) Error() string {
	// TODO: I don't love this errors collection thing. Is there a sane way to
	// implement error for it?
	// return "Multiple API errors occured"
	var s string
	for _, e := range err {
		s += e.Error() + " :: "
	}
	return s
}
