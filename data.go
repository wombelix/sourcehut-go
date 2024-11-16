// SPDX-FileCopyrightText: 2019 The SourceHut API Contributors
//
// SPDX-License-Identifier: BSD-2-Clause

package sourcehut

// ShortUser represents the unexpanded form of a user returned by most API
// endpoints.
type ShortUser struct {
	CanonicalName string `json:"canonical_name"`
	Name          string `json:"name"`
}

// User represents the expanded form of a user.
type User struct {
	ShortUser

	Email    string `json:"email"`
	URL      string `json:"url"`
	Location string `json:"location"`
	Bio      string `json:"bio"`
}
