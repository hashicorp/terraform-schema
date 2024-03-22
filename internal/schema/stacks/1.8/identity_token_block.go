// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func identityTokenBlockSchema(v *version.Version) *schema.BlockSchema {
	/*
		Reference: https://github.com/hashicorp/tfc-agent/blob/main/core/components/stacks/tfdeploycfg/identity_token.go
		TODO:
			- Source better descriptions
			- Verify all attributes are added here
	*/
	return &schema.BlockSchema{
		Description: lang.PlainText("An identity token block is a definition of a JSON Web Token (JWT) that will be generated for a given deployment if referenced in the inputs for that deployment block. The block label defines the token name, which must be unique within the stack."),
		Address: &schema.BlockAddrSchema{
			FriendlyName: "deployment",
			ScopeId:      refscope.ProviderScope,
			AsReference:  true,
			Steps: []schema.AddrStep{
				schema.LabelStep{Index: 0},
				schema.AttrValueStep{Name: "alias", IsOptional: true},
			},
		},
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Provider},
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
