package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func resourceBlockSchema(v *version.Version) *schema.BlockSchema {
	bs := &schema.BlockSchema{
		Labels: []*schema.LabelSchema{
			{
				Name:        "type",
				Description: lang.PlainText("Resource Type"),
				IsDepKey:    true,
			},
			{
				Name:        "name",
				Description: lang.PlainText("Reference Name"),
			},
		},
		Description: lang.PlainText("A resource block declares a resource of a given type with a given local name. The name is " +
			"used to refer to this resource from elsewhere in the same Terraform module, but has no significance " +
			"outside of the scope of a module."),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"provider": {
					Expr:        schema.ExprConstraints{},
					IsOptional:  true,
					Description: lang.Markdown("Reference to a `provider` configuration block, e.g. `mycloud.west` or `mycloud`"),
					IsDepKey:    true,
				},
				"count": {
					Expr:        schema.LiteralTypeOnly(cty.Number),
					IsOptional:  true,
					Description: lang.Markdown("Number of instances of this resource, e.g. `3`"),
				},
				"depends_on": {
					Expr: schema.ExprConstraints{
						schema.TupleConsExpr{
							Name: "set of references",
						},
					},
					IsOptional:  true,
					Description: lang.Markdown("Set of references to hidden dependencies, e.g. other resources or data sources"),
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"lifecycle":  lifecycleBlock,
				"connection": connectionBlock,
			},
		},
	}

	if v.GreaterThanOrEqual(v0_12_6) {
		bs.Body.Attributes["for_each"] = &schema.AttributeSchema{
			Expr: schema.ExprConstraints{
				schema.LiteralTypeExpr{Type: cty.Set(cty.DynamicPseudoType)},
				schema.LiteralTypeExpr{Type: cty.Map(cty.DynamicPseudoType)},
			},
			IsOptional:  true,
			Description: lang.Markdown("A set or a map where each item represents an instance of this resource"),
		}
	}

	return bs
}

var lifecycleBlock = &schema.BlockSchema{
	Description: lang.Markdown("Lifecycle customizations to change default resource behaviours during apply"),
	Body: &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{
			"create_before_destroy": {
				Expr:       schema.LiteralTypeOnly(cty.Bool),
				IsOptional: true,
				Description: lang.Markdown("Whether to reverse the default order of operations (destroy -> create) during apply " +
					"when the resource requires replacement (cannot be updated in-place)"),
			},
			"prevent_destroy": {
				Expr:       schema.LiteralTypeOnly(cty.Bool),
				IsOptional: true,
				Description: lang.Markdown("Whether to prevent accidental destruction of the resource and cause Terraform " +
					"to reject with an error any plan that would destroy the resource"),
			},
			"ignore_changes": {
				Expr: schema.ExprConstraints{
					schema.TupleConsExpr{},
				},
				IsOptional:  true,
				Description: lang.Markdown("A set of fields (references) of which to ignore changes to, e.g. `tags`"),
			},
		},
	},
}

var provisionerBlock = &schema.BlockSchema{
	Description: lang.Markdown("Provisioner to model specific actions on the local machine or on a remote machine " +
		"in order to prepare servers or other infrastructure objects for service"),
	Labels: []*schema.LabelSchema{
		{
			Name:        "type",
			Description: lang.PlainText("Type of provisioner to use, e.g. `remote-exec` or `file`"),
			IsDepKey:    true,
		},
	},
	Body: &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{
			"when": {
				Expr: schema.ExprConstraints{
					schema.KeywordExpr{
						Keyword:     "create",
						Description: lang.Markdown("Run the provisioner when the resource is created"),
					},
					schema.KeywordExpr{
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
				Expr: schema.ExprConstraints{
					schema.KeywordExpr{
						Keyword:     "fail",
						Description: lang.Markdown("Raise an error and stop applying (the default behavior). If this is a creation provisioner, taint the resource."),
					},
					schema.KeywordExpr{
						Keyword:     "continue",
						Description: lang.Markdown("Ignore the error and continue with creation or destruction"),
					},
				},
				Description: lang.Markdown("What to do when the provisioner run fails to finish - `fail` (default), " +
					"or `continue` (ignore the error)"),
			},
		},
		Blocks: map[string]*schema.BlockSchema{
			"connection": connectionBlock,
		},
	},
}

var connectionBlock = &schema.BlockSchema{
	Description: lang.Markdown("Connection block describing how the provisioner connects to the given instance"),
	MaxItems:    1,
	Body: &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{
			"type": {
				Expr: schema.ExprConstraints{
					schema.LiteralValue{
						Val:         cty.StringVal("ssh"),
						Description: lang.Markdown("Use SSH to connect and provision the instance"),
					},
					schema.LiteralValue{
						Val:         cty.StringVal("winrm"),
						Description: lang.Markdown("Use WinRM to connect and provision the instance"),
					},
				},
				IsOptional:  true,
				Description: lang.Markdown("Connection type to use - `ssh` (default) or `winrm`"),
			},
		},
	},
}
