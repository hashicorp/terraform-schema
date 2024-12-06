// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"github.com/hashicorp/hcl/v2"
	tftest "github.com/hashicorp/terraform-schema/test"
)

func LoadTest(path string, filename string, file *hcl.File) (*tftest.Meta, hcl.Diagnostics) {
	mod := newDecodedTest()
	diags := loadTestFromFile(file, mod)

	// TODO: a lot of mapping to do here from decoded test to tftest.Meta

	return &tftest.Meta{
		Path: path,
	}, diags
}
