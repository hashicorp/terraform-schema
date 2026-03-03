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

func deploymentBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.Markdown("Deployment"),
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.StaticStep{Name: "deployment"},
				schema.LabelStep{Index: 0},
			},
			FriendlyName: "deployment",
			ScopeId:      refscope.DeploymentScope,
			AsReference:  true,
		},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name, lang.TokenModifierDependent},
				Description:            lang.PlainText("Deployment Name"),
				IsDepKey:               true,
				Completable:            true,
			},
		},
		Body: &schema.BodySchema{
			Extensions: &schema.BodyExtensions{
				DynamicBlocks: true,
			},
			Attributes: map[string]*schema.AttributeSchema{
				"inputs": {
					Description: lang.Markdown("A mapping of stack variable names to values for this deployment. The keys of this map must correspond to the names of variables defined for the stack. The values must be valid HCL literals meeting the type constraint of those variables. Values are also expressions, currently with access to identity token references only"),
					IsOptional:  true,
					Constraint: schema.Map{
						Name: "map of variable references",
						Elem: schema.AnyExpression{OfType: cty.DynamicPseudoType},
					},
				},
				"deployment_group": {
					Description: lang.Markdown("A reference to the `deployment_group` block that this deployment belongs to"),
					IsOptional:  true,
					Constraint:  schema.Reference{OfScopeId: refscope.DeploymentGroupScope},
				},
				"destroy": {
					Description:  lang.Markdown("A boolean flag that indicates whether HCP Terraform should destroy this deployment"),
					IsOptional:   true,
					DefaultValue: schema.DefaultValue{Value: cty.False},
					Constraint:   schema.LiteralType{Type: cty.Bool},
				},
				"migrate": {
					Description:  lang.Markdown("A boolean flag set by `tf-migrate` CLI when migrating the state of an existing workspace to a specific deployment. Lets HCP Terraform know that other deployments can continue planning as normal while this deployment is migrating state"),
					IsOptional:   true,
					DefaultValue: schema.DefaultValue{Value: cty.False},
					Constraint:   schema.LiteralType{Type: cty.Bool},
				},
			},
		},
	}
}
