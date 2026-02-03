// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"github.com/hashicorp/hcl/v2"
)

var rootSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "variable",
			LabelNames: []string{"name"},
		},
		{
			Type: "policy",
		},
		{
			Type:       "resource_policy",
			LabelNames: []string{"type", "name"},
		},
		{
			Type:       "provider_policy",
			LabelNames: []string{"type", "name"},
		},
		{
			Type:       "module_policy",
			LabelNames: []string{"source", "name"},
		},
	},
}

var policyBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{
			Name: "enforcement_level",
		},
	},
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type: "terraform_config",
		},
	},
}

var terraformConfigBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{
			Name: "required_version",
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
			Name: "sensitive",
		},
		{
			Name: "default",
		},
	},
}
