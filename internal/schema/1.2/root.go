// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"

	v1_1_mod "github.com/hashicorp/terraform-schema/internal/schema/1.1"
)

func ModuleSchema(v *version.Version) *schema.BodySchema {
	bs := v1_1_mod.ModuleSchema(v)
	bs.Blocks["data"].Body.Blocks = map[string]*schema.BlockSchema{
		"lifecycle": datasourceLifecycleBlock(),
	}
	bs.Blocks["resource"].Body.Blocks["lifecycle"] = resourceLifecycleBlock()
	bs.Blocks["output"].Body.Blocks = map[string]*schema.BlockSchema{
		"lifecycle": outputLifecycleBlock(),
	}

	return bs
}

func conditionBody(enableSelfRefs bool) *schema.BodySchema {
	bs := &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{
			"condition": {
				Constraint: schema.AnyExpression{OfType: cty.Bool},
				IsRequired: true,
				Description: lang.Markdown("Condition, a boolean expression that should return `true` " +
					"if the intended assumption or guarantee is fulfilled or `false` if it is not."),
			},
			"error_message": {
				Constraint: schema.AnyExpression{OfType: cty.String},
				IsRequired: true,
				Description: lang.Markdown("Error message to return if the `condition` isn't met " +
					"(evaluates to `false`)."),
			},
		},
	}

	if enableSelfRefs {
		bs.Extensions = &schema.BodyExtensions{
			SelfRefs: true,
		}
	}

	return bs
}
