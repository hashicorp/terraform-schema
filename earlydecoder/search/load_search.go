// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/terraform-schema/search"
)

type decodedSearch struct {
	List      map[string]*search.List
	Variables map[string]*search.Variable
}

func newDecodedSearch() *decodedSearch {
	return &decodedSearch{
		List:      make(map[string]*search.List),
		Variables: make(map[string]*search.Variable),
	}
}

func loadSearchFromFile(file *hcl.File, _ *decodedSearch) hcl.Diagnostics {
	var diags hcl.Diagnostics
	// TODO Implementation
	return diags
}
