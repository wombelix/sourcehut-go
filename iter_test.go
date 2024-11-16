// SPDX-FileCopyrightText: 2019 The SourceHut API Contributors
//
// SPDX-License-Identifier: BSD-2-Clause

package sourcehut_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"git.sr.ht/~samwhited/sourcehut-go"
	"git.sr.ht/~samwhited/sourcehut-go/internal/testlog"
)

type iterTest struct {
	code int
	body []string
	vals []interface{}
	d    func() interface{}
	err  error
}

var iterTests = [...]iterTest{
	0: {code: 404, err: io.EOF},
	1: {
		code: 400,
		body: []string{`{"errors": [{"field": "f", "reason": "r"}]}`},
		err:  sourcehut.Error{Field: "f", Reason: "r"},
	},
	2: {
		code: 403,
		body: []string{`{"errors": [{"field": "f", "reason": "r"}, {"field": "f2", "reason": "r2"}]}`},
		err:  sourcehut.Errors{sourcehut.Error{Field: "f", Reason: "r"}, sourcehut.Error{Field: "f2", Reason: "r2"}},
	},
	3: {err: io.EOF},
	4: {body: []string{`{}`}, err: io.EOF},
	5: {
		body: []string{`{
  "results": [ {}, {} ],
  "results_per_page": 50,
  "total": 2
}`},
		d:    func() interface{} { return &struct{}{} },
		vals: []interface{}{&struct{}{}, &struct{}{}},
	},
	6: {
		body: []string{`{
  "next": "2",
  "results": [ {}, {} ],
  "results_per_page": 50,
  "total": 2
}`, `{
  "results": [ {}, {} ],
  "results_per_page": 50,
  "total": 2
}`},
		vals: []interface{}{make(map[string]interface{}), make(map[string]interface{}), make(map[string]interface{}), make(map[string]interface{})},
	},
	7: {
		body: []string{`{
  "results": [],
  "results_per_page": 50,
  "total": 2
}`},
	},
	8: {
		body: []string{`{
  "results": "wat",
  "results_per_page": 50,
  "total": 2
}`},
		err: sourcehut.ErrWantArray,
	},
}

func TestIter(t *testing.T) {
	for i, tc := range iterTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var served int
			server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if tc.code != 0 {
					w.WriteHeader(tc.code)
				}
				if len(tc.body) > 0 {
					_, err := w.Write([]byte(tc.body[served]))
					if err != nil {
						t.Fatalf("Error writing response body (this should never happen): %q", err)
					}
				} else {
					_, err := w.Write([]byte{})
					if err != nil {
						t.Fatalf("Error writing response body (this should never happen): %q", err)
					}
				}
				served++
			}))
			server.Config.ErrorLog = testlog.New(t)
			server.Start()
			client := server.Client()
			srhtClient := sourcehut.NewBaseClient(client)
			defer server.Close()

			doIterTest(t, server.URL, srhtClient, tc)
		})
	}
}

func errEq(statusCode int, e1, e2 sourcehut.Error) bool {
	if e1.Field != e2.Field {
		return false
	}
	if e1.Reason != e2.Reason {
		return false
	}
	return e1.StatusCode() == statusCode
}

func doIterTest(t *testing.T, u string, client sourcehut.Client, tc iterTest) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		t.Fatal(err)
	}
	iter := client.List(req, tc.d)
	var i int
	for iter.Next() {
		v := iter.Current()
		if len(tc.vals) <= i {
			t.Fatal("More values decoded than expected from iter")
		}
		if !reflect.DeepEqual(v, tc.vals[i]) {
			t.Fatalf("Unexpected value in response: want=`%+v', got=`%+v'", tc.vals[i], v)
		}
		i++
	}
	err = iter.Err()

	switch e := err.(type) {
	case sourcehut.Error:
		if !errEq(tc.code, e, tc.err.(sourcehut.Error)) {
			t.Fatalf("Unexpected err: want=%#v, got=%#v", tc.err, err)
		}
	case sourcehut.Errors:
		for i := 0; i < len(e); i++ {
			tcErr := tc.err.(sourcehut.Errors)
			if !errEq(tc.code, e[i], tcErr[i]) {
				t.Fatalf("Unexpected err: want=%#v, got=%#v", tcErr[i], e[i])
			}
		}
	}
	if err != nil && iter.Next() {
		t.Fatalf("Next unexpectedly returned true after error")
	}
}
