// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"github.com/hashicorp/hcl/v2"
)

var rootSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "list",
			LabelNames: []string{"type", "name"},
		},
		{
			Type:       "variable",
			LabelNames: []string{"name"},
		},
	},
}

var listSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{
			Name: "include_resource",
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
