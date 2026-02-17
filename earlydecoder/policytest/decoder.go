// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"sort"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/terraform-schema/policytest"
)

func LoadPolicyTest(path string, files map[string]*hcl.File) (*policytest.Meta, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	filenames := make([]string, 0)

	mod := newDecodedPolicyTest()
	for filename, f := range files {
		filenames = append(filenames, filename)
		fDiags := loadPolicyTestFromFile(f, mod)
		diags = append(diags, fDiags...)
	}

	sort.Strings(filenames)

	return &policytest.Meta{
		Path:      path,
		Filenames: filenames,
	}, diags
}
