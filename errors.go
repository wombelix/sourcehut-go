// Copyright 2019 The Sourcehut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package sourcehut

// Ensure that the build fails if Error and Errors don't implement error.
var _, _ error = (*Error)(nil), (*Errors)(nil)

// Error represents an individual error returned from a Sourcehut API call.
//
// API docs: https://man.sr.ht/api-conventions.md#error-responses
type Error struct {
	Field  string
	Reason string

	statusCode int
}

// StatusCode returns the HTTP status code of the request that unmarshaled this
// error.
// May not be set for Errors originating from code outside this package.
func (err Error) StatusCode() int {
	return err.statusCode
}

// Error satisfies the error interface for Error.
func (err Error) Error() string {
	return err.Reason
}

// Errors is a slice of Error's that itself implements error.
type Errors []Error

// Error satisfies the error interface for Errors.
func (err Errors) Error() string {
	return "Multiple API errors occured"
}

// StatusCode returns the HTTP status code of the request that unmarshaled this
// error.
// May not be set for Errors originating from code outside this package.
func (err Errors) StatusCode() int {
	if len(err) == 0 {
		return 0
	}
	return err[0].statusCode
}
