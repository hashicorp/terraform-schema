// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package search

type Meta struct {
	Path      string
	Filenames []string

	Variables map[string]Variable
	Lists     map[string]List
}
