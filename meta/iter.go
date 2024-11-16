// SPDX-FileCopyrightText: 2019 The SourceHut API Contributors
//
// SPDX-License-Identifier: BSD-2-Clause

package meta

import (
	"git.sr.ht/~samwhited/sourcehut-go"
)

// SSHKeyIter is used for iterating over the account's authorized SSH keys.
type SSHKeyIter struct {
	*sourcehut.Iter
}

// Key returns the SSH key which the iterator is currently pointing to.
func (i SSHKeyIter) Key() SSHKey {
	return *(i.Current().(*SSHKey))
}

// PGPKeyIter is used for iterating over the account's PGP keys.
type PGPKeyIter struct {
	*sourcehut.Iter
}

// Key returns the PGP key which the iterator is currently pointing to.
func (i PGPKeyIter) Key() PGPKey {
	return *(i.Current().(*PGPKey))
}

// AuditLogIter is used for iterating over the account's PGP keys.
type AuditLogIter struct {
	*sourcehut.Iter
}

// Log returns the audit log entry which the iterator is currently pointing to.
func (i AuditLogIter) Log() AuditLog {
	return *(i.Current().(*AuditLog))
}
