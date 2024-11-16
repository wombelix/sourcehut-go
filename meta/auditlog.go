// SPDX-FileCopyrightText: 2019 The SourceHut API Contributors
//
// SPDX-License-Identifier: BSD-2-Clause

package meta

import (
	"io"
	"net/http"
)

// ListAuditLog returns an iterator over all audit log entries available to the
// authenticated user.
func (c *Client) ListAuditLog() (AuditLogIter, error) {
	return c.auditLogs("GET", "user/audit-log", nil)
}

func (c *Client) auditLogs(method, u string, body io.Reader) (AuditLogIter, error) {
	u = c.baseURL.String() + u
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return AuditLogIter{}, err
	}
	iter := c.srhtClient.List(req, func() interface{} {
		return &AuditLog{}
	})
	return AuditLogIter{Iter: iter}, nil
}
