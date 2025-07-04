// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/terraform-schema/search"
	"sort"
	"strings"
)

func LoadSearch(path string, files map[string]*hcl.File) (*search.Meta, map[string]hcl.Diagnostics) {
	filenames := make([]string, 0)
	diags := make(map[string]hcl.Diagnostics, 0)

	mod := newDecodedSearch()
	for filename, f := range files {
		filenames = append(filenames, filename)
		if isSearchFile(filename) {
			diags[filename] = loadSearchFromFile(f, mod)
		}
	}

	sort.Strings(filenames)

	variables := make(map[string]search.Variable)
	for key, variable := range mod.Variables {
		variables[key] = *variable
	}

	lists := make(map[string]search.List)
	for key, list := range mod.List {
		lists[key] = *list
	}

	return &search.Meta{
		Path:      path,
		Filenames: filenames,
		Variables: variables,
		Lists:     lists,
	}, diags
}

func isSearchFile(name string) bool {
	return strings.HasSuffix(name, ".tfquery.hcl") ||
		strings.HasSuffix(name, ".tfquery.hcl.json")
}
