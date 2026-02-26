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
			Type:       "data",
			LabelNames: []string{"data_type", "name"},
		},
		{
			Type:       "variable",
			LabelNames: []string{"name"},
		},
	},
}

var variableSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{
			Name: "description",
		},
		{
			Name: "type",
		},
		{
			Name: "default",
		},
		{
			Name: "sensitive",
		},
	},
}

var dataSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{
			Name: "attrs",
		},
	},
}
