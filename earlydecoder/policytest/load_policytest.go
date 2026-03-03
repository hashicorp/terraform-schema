// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"github.com/hashicorp/hcl/v2"
)

// decodedPolicyTest is the type representing a decoded Terraform policytest.
type decodedPolicyTest struct {
	//Variables   map[string]*policytest.Variable
}

func newDecodedPolicyTest() *decodedPolicyTest {
	return &decodedPolicyTest{
		//Variables:   make(map[string]*policytest.Variable),
	}
}

// loadPolicyTestFromFile reads given file, interprets it and stores in given PolicyTest
// This is useful for any caller which does tokenization/parsing on its own
// e.g. because it will reuse these parsed files later for more detailed
// interpretation.
func loadPolicyTestFromFile(file *hcl.File, decodedPolicyTest *decodedPolicyTest) hcl.Diagnostics {
	var diags hcl.Diagnostics
	_, _, contentDiags := file.Body.PartialContent(rootSchema)
	diags = append(diags, contentDiags...)
	return diags
}
