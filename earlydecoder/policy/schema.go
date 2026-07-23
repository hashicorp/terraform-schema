// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"github.com/hashicorp/hcl/v2"
)

var rootSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
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
		{
			Type:       "input",
			LabelNames: []string{"name"},
		},
	},
}

var policyBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{},
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type: "terraform_config",
		},
		{
			Type: "plugins",
		},
		{
			Type: "required_providers",
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

var resourcePolicyBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "filter"},
		{Name: "enforcement_level"},
		{Name: "operations"},
	},
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type: "locals",
		},
		{
			Type: "enforce",
		},
	},
}

var providerPolicyBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "filter"},
		{Name: "enforcement_level"},
	},
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type: "locals",
		},
		{
			Type: "enforce",
		},
	},
}

var modulePolicyBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "filter"},
		{Name: "enforcement_level"},
	},
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type: "locals",
		},
		{
			Type: "enforce",
		},
	},
}

var inputSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "description"},
		{Name: "type"},
		{Name: "default"},
		{Name: "sensitive"},
		{Name: "nullable"},
	},
}
