// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"github.com/hashicorp/hcl/v2"
)

// decodedTest is the type representing a decoded Terraform test.
type decodedTest struct {
}

func newDecodedTest() *decodedTest {
	return &decodedTest{}
}

// loadTestFromFile reads given file, interprets it and stores in given test
// This is useful for any caller which does tokenization/parsing on its own
// e.g. because it will reuse these parsed files later for more detailed
// interpretation.
func loadTestFromFile(file *hcl.File, _ *decodedTest) hcl.Diagnostics {
	var diags hcl.Diagnostics

	content, _, contentDiags := file.Body.PartialContent(rootSchema)
	diags = append(diags, contentDiags...)
	for _, block := range content.Blocks {
		switch block.Type {
		// TODO? decode mock_provider blocks (they have a source attribute)
		// TODO? decode run -> module blocks (they have a source attribute)
		}
	}

	return diags
}
