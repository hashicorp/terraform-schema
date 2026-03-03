// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func stackBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.Markdown("Stack block sources an entire component configuration from the HCP Terraform private registry, including all of the Stack's components and providers"),
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.StaticStep{Name: "stack"},
				schema.LabelStep{Index: 0},
			},
			FriendlyName: "stack",
			ScopeId:      refscope.StackScope,
			AsReference:  true,
		},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Stack Name"),
			},
		},
		Body: &schema.BodySchema{
			Extensions: &schema.BodyExtensions{
				ForEach: true,
			},
			Attributes: map[string]*schema.AttributeSchema{
				"source": {
					Description: lang.Markdown("The component configuration to source from the HCP Terraform private registry. The source address format is `<NAMESPACE>/<NAME>`, where `<NAMESPACE>` is your organization's name and `<NAME>` is the component configuration's name"),
					IsRequired:  true,
					Constraint:  schema.LiteralType{Type: cty.String},
				},
				"version": {
					Description: lang.Markdown("The component configuration version to use from the HCP Terraform private registry. Accepts a version constraint string"),
					IsRequired:  true,
					Constraint:  schema.LiteralType{Type: cty.String},
				},
				"inputs": {
					Description: lang.Markdown("A mapping of input variable names to values for the component configuration. The keys in this map must correspond to the variable names defined in the component configuration"),
					IsRequired:  true,
					Constraint: schema.Map{
						Name: "map of input references",
						Elem: schema.AnyExpression{OfType: cty.DynamicPseudoType},
					},
				},
				"depends_on": {
					Description: lang.Markdown("Optionally specify explicit dependencies for stacks in a stack configuration, which must also be used when determining an order of operations for stacks"),
					IsOptional:  true,
					Constraint: schema.Set{
						Elem: schema.OneOf{
							schema.Reference{OfScopeId: refscope.StackScope},
							schema.Reference{OfScopeId: refscope.ComponentScope},
						},
					},
				},
			},
		},
	}
}
