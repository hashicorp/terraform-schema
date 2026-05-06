// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func patchRunBlockSchema(bs *schema.BlockSchema) *schema.BlockSchema {
	bs.Body.Attributes["state_key"] = &schema.AttributeSchema{
		Description: lang.Markdown("An optional key to override the default state file used for this run block. Setting this value forces Terraform to use a specific state file identified by the given key, allowing state to be shared between run blocks. Read more on [module states](https://developer.hashicorp.com/terraform/language/tests#modules-state)"),
		Constraint:  schema.AnyExpression{OfType: cty.String},
		IsOptional:  true,
	}

	bs.Body.Blocks["override_resource"].Body.Attributes["override_during"] = overrideDuringAttributeSchema()
	bs.Body.Blocks["override_data"].Body.Attributes["override_during"] = overrideDuringAttributeSchema()
	bs.Body.Blocks["override_module"].Body.Attributes["override_during"] = overrideDuringAttributeSchema()

	return bs
}
