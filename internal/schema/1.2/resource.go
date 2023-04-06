// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/zclconf/go-cty/cty"
)

func resourceLifecycleBlock() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.Markdown("Lifecycle customizations to change default resource behaviours during plan or apply"),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"create_before_destroy": {
					Constraint: schema.LiteralType{Type: cty.Bool},
					IsOptional: true,
					Description: lang.Markdown("Whether to reverse the default order of operations (destroy -> create) during apply " +
						"when the resource requires replacement (cannot be updated in-place)"),
				},
				"prevent_destroy": {
					Constraint: schema.LiteralType{Type: cty.Bool},
					IsOptional: true,
					Description: lang.Markdown("Whether to prevent accidental destruction of the resource and cause Terraform " +
						"to reject with an error any plan that would destroy the resource"),
				},
				"ignore_changes": {
					Constraint: schema.OneOf{
						schema.Set{
							// TODO: expose reference targets via attribute-only address
						},
						schema.Keyword{
							Keyword: "all",
							Description: lang.Markdown("Ignore all attributes, which means that Terraform can create" +
								" and destroy the remote object but will never propose updates to it"),
						},
					},
					IsOptional:  true,
					Description: lang.Markdown("A set of fields (references) of which to ignore changes to, e.g. `tags`"),
				},
				"replace_triggered_by": {
					Constraint: schema.Set{
						Elem: schema.Reference{
							OfScopeId: refscope.ResourceScope,
						},
					},
					IsOptional: true,
					Description: lang.Markdown("Set of references to any other resources which when changed cause " +
						"this resource to be proposed for replacement"),
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"precondition": {
					Body: conditionBody(false),
				},
				"postcondition": {
					Body: conditionBody(true),
				},
			},
		},
	}
}
