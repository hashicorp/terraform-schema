// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"github.com/hashicorp/hcl/v2"
)

var rootSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "component",
			LabelNames: []string{"name"},
		},
		{
			Type: "required_providers",
		},
		{
			Type:       "variable",
			LabelNames: []string{"name"},
		},
		{
			Type:       "output",
			LabelNames: []string{"name"},
		},
	},
}

var componentSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{
			Name: "source",
		},
		{
			Name: "version",
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

var outputSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{
			Name: "description",
		},
		{
			Name: "value",
		},
		{
			Name: "sensitive",
		},
	},
}
