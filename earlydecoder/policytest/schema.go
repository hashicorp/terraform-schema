// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"github.com/hashicorp/hcl/v2"
)

var rootSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type: "policytest",
		},
		{
			Type: "inputs",
		},
		{
			Type: "locals",
		},
		{
			Type:       "resource",
			LabelNames: []string{"resource_type", "test_case_name"},
		},
		{
			Type:       "module",
			LabelNames: []string{"module_source", "test_case_name"},
		},
		{
			Type:       "provider",
			LabelNames: []string{"provider_type", "test_case_name"},
		},
		{
			Type:       "data",
			LabelNames: []string{"data_type", "name"},
		},
	},
}
