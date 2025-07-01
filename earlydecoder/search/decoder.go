// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"github.com/hashicorp/hcl/v2"
	tftest "github.com/hashicorp/terraform-schema/test"
	"sort"
)

func LoadTest(path string, files map[string]*hcl.File) (*tftest.Meta, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	filenames := make([]string, 0)

	mod := newDecodedSearch()
	for filename, f := range files {
		filenames = append(filenames, filename)
		fDiags := loadSearchFromFile(f, mod)
		diags = append(diags, fDiags...)
	}

	sort.Strings(filenames)

	return &tftest.Meta{
		Path:      path,
		Filenames: filenames,
	}, diags
}
