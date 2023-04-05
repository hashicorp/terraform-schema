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

func resourceBlockSchema(v *version.Version) *schema.BlockSchema {
	bs := &schema.BlockSchema{
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.LabelStep{Index: 0},
				schema.LabelStep{Index: 1},
			},
			FriendlyName:         "resource",
			ScopeId:              refscope.ResourceScope,
			AsReference:          true,
			DependentBodyAsData:  true,
			InferDependentBody:   true,
			DependentBodySelfRef: true,
		},
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Resource},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "type",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
				Description:            lang.PlainText("Resource Type"),
				IsDepKey:               true,
				Completable:            true,
			},
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("Reference Name"),
			},
		},
		Description: lang.PlainText("A resource block declares a resource of a given type with a given local name. The name is " +
			"used to refer to this resource from elsewhere in the same Terraform module, but has no significance " +
			"outside of the scope of a module."),
		Body: &schema.BodySchema{
			Extensions: &schema.BodyExtensions{
				Count:         true,
				ForEach:       true, // for_each was introduced in 0.12.6, but for simplicity we report it for all 0.12+
				DynamicBlocks: true,
			},
			Attributes: map[string]*schema.AttributeSchema{
				"provider": {
					Constraint:             schema.Reference{OfScopeId: refscope.ProviderScope},
					IsOptional:             true,
					Description:            lang.Markdown("Reference to a `provider` configuration block, e.g. `mycloud.west` or `mycloud`"),
					IsDepKey:               true,
					SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent},
				},
				"depends_on": {
					Constraint: schema.Set{
						Elem: schema.OneOf{
							schema.Reference{OfScopeId: refscope.DataScope},
							schema.Reference{OfScopeId: refscope.ModuleScope},
							schema.Reference{OfScopeId: refscope.ResourceScope},
						},
					},
					IsOptional:  true,
					Description: lang.Markdown("Set of references to hidden dependencies, e.g. other resources or data sources"),
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"lifecycle":   lifecycleBlock(),
				"connection":  connectionBlock(v),
				"provisioner": provisionerBlock(v),
			},
		},
	}

	return bs
}

func lifecycleBlock() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.Markdown("Lifecycle customizations to change default resource behaviours during apply"),
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
			},
		},
	}
}
