package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
)

func provisionerBlock(v *version.Version) *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.Markdown("Provisioner to model specific actions on the local machine or on a remote machine " +
			"in order to prepare servers or other infrastructure objects for service"),
		Labels: []*schema.LabelSchema{
			{
				Name:        "type",
				Description: lang.PlainText("Type of provisioner to use, e.g. `remote-exec` or `file`"),
				IsDepKey:    true,
				Completable: true,
			},
		},
		Body: &schema.BodySchema{
			HoverURL: "https://www.terraform.io/docs/language/resources/provisioners/syntax.html",
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
				"connection": connectionBlock(v),
			},
		},
		DependentBody: ProvisionerDependentBodies(v),
	}
}
