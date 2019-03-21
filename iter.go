// Copyright 2019 The SourceHut API Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package sourcehut

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

var (
	errWantArray = errors.New("expected json array in response")
)

// Response is a SourceHut API response.
// This wraps the standard http.Response and provides convenient access to
// pagination links.
//
// API docs: https://man.sr.ht/api-conventions.md
type Response struct {
	*http.Response `json:"-"`

	Next           string          `json:"next"`
	Results        json.RawMessage `json:"results"`
	ResultsPerPage int64           `json:"results_per_page"`
	Total          int64           `json:"total"`
}

// Iter provides a convenient API for iterating over the elements returned from
// paginated list API calls.
// Successive calls to the Next method step through each item in the list,
// fetching pages as needed.
type Iter struct {
	resp *Response
	v    interface{}
	err  error
	d    *json.Decoder
	into func() interface{}
	req  *http.Request
	c    Client
}

// Current returns the most recent item visited by the iterator.
func (i *Iter) Current() interface{} {
	return i.v
}

// Err returns the last error encountered by the iterator.
// It will only return a non-nil value if the previous call to Next returned
// false (but Next returning false does not guarantee that Err will return a
// non-nil value).
func (i *Iter) Err() error {
	return i.err
}

// Next advances the iterator to the next item in the list and makes it
// available through the Current method.
// When the end of the list is reached it returns False.
func (i *Iter) Next() bool {
	if i.err != nil {
		return false
	}

	if i.into == nil {
		i.v = make(map[string]interface{})
	} else {
		i.v = i.into()
	}

	if i.d == nil || !i.d.More() {
		// We're out of JSON to decode, fetch the next page if there is one…
		if i.resp != nil {
			if i.resp.Next == "" {
				return false
			}

			// TODO: clone req
			q := i.req.URL.Query()
			q.Set("start", i.resp.Next)
			i.req.URL.RawQuery = q.Encode()
		}

		resp, err := i.c.do(i.req)
		if err != nil {
			i.err = err
			return false
		}
		i.resp = &Response{Response: resp}
		defer i.resp.Body.Close()

		i.err = json.NewDecoder(i.resp.Body).Decode(i.resp)
		if i.err != nil {
			return false
		}
		i.d = json.NewDecoder(strings.NewReader(string(i.resp.Results)))
		// Advance past the first token so that we can treat the array as a stream.
		tok, err := i.d.Token()
		if err != nil {
			i.err = err
			return false
		}
		if delim, ok := tok.(json.Delim); !ok || delim != '[' {
			i.err = errWantArray
			return false
		}
	}

	// If we just fetched the list, but there's nothing in it ('[]'), don't try to
	// decode an element.
	if !i.d.More() {
		return false
	}

	i.err = i.d.Decode(&i.v)
	return i.err == nil
}
