// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
)

func provisionerBlock(v *version.Version) *schema.BlockSchema {
	return &schema.BlockSchema{
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Provisioner},
		Description: lang.Markdown("Provisioner to model specific actions on the local machine or on a remote machine " +
			"in order to prepare servers or other infrastructure objects for service"),
		Labels: []*schema.LabelSchema{
			{
				Name:                   "type",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
				Description:            lang.PlainText("Type of provisioner to use, e.g. `remote-exec` or `file`"),
				IsDepKey:               true,
				Completable:            true,
			},
		},
		Body: &schema.BodySchema{
			Extensions: &schema.BodyExtensions{
				DynamicBlocks: true,
				SelfRefs:      true,
			},
			HoverURL: "https://www.terraform.io/docs/language/resources/provisioners/syntax.html",
			Attributes: map[string]*schema.AttributeSchema{
				"when": {
					Constraint: schema.OneOf{
						schema.Keyword{
							Keyword:     "create",
							Description: lang.Markdown("Run the provisioner when the resource is created"),
						},
						schema.Keyword{
							Keyword:     "destroy",
							Description: lang.Markdown("Run the provisioner when the resource is destroyed"),
						},
					},
					IsOptional: true,
					Description: lang.Markdown("When to run the provisioner - `create` or `destroy`, defaults to `create` " +
						"(i.e. after creation of the resource)"),
				},
				"on_failure": {
					IsOptional: true,
					Constraint: schema.OneOf{
						schema.Keyword{
							Keyword:     "fail",
							Description: lang.Markdown("Raise an error and stop applying (the default behavior). If this is a creation provisioner, taint the resource."),
						},
						schema.Keyword{
							Keyword:     "continue",
							Description: lang.Markdown("Ignore the error and continue with creation or destruction"),
						},
					},
					Description: lang.Markdown("What to do when the provisioner run fails to finish - `fail` (default), " +
						"or `continue` (ignore the error)"),
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"connection": connectionBlock(v),
			},
		},
		DependentBody: ProvisionerDependentBodies(v),
	}
}
