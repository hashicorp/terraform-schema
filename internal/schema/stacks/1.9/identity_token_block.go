// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func identityTokenBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.PlainText("An identity token block is a definition of a JSON Web Token (JWT) that will be generated for a given deployment if referenced in the inputs for that deployment block. The block label defines the token name, which must be unique within the stack."),
		Labels: []*schema.LabelSchema{
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name, lang.TokenModifierDependent},
				Description:            lang.PlainText("Identity name"),
				IsDepKey:               true,
				Completable:            true,
			},
		},
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"audience": {
					Description: lang.Markdown("The audience(s) that tokens generated with this configuration block will be generated with. Audience(s) are the resource(s)/server(s) that the token is intended for. With an audience claim, the cloud service authorizing the workload can be confident that the token is being presented intentionally to that service"),
					IsOptional:  true,
					Constraint: schema.List{
						// TODO: Is a list correct for this attribute?
						Elem: schema.AnyExpression{OfType: cty.String},
					},
				},
			},
		},
	}
}
