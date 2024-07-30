// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"sort"

	"github.com/hashicorp/hcl/v2"
	tftest "github.com/hashicorp/terraform-schema/test"
)

func LoadTest(path string, files map[string]*hcl.File) (*tftest.Meta, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	filenames := make([]string, 0)

	mod := newDecodedTest()
	for filename, f := range files {
		filenames = append(filenames, filename)
		fDiags := loadTestFromFile(f, mod)
		diags = append(diags, fDiags...)
	}

	sort.Strings(filenames)

	return &tftest.Meta{
		Path:      path,
		Filenames: filenames,
	}, diags
}
