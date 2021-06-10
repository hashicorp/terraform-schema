package backends

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func remoteBackend(v *version.Version) *schema.BodySchema {
	docsUrl := "https://www.terraform.io/docs/language/settings/backends/remote.html"
	return &schema.BodySchema{
		Description: lang.Markdown("Remote backend to store state and run operations in Terraform Cloud."),
		HoverURL:    docsUrl,
		DocsLink: &schema.DocsLink{
			URL: docsUrl,
		},
		Attributes: map[string]*schema.AttributeSchema{
			"hostname": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The remote backend hostname to connect to (defaults to `app.terraform.io`)."),
			},
			"organization": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsRequired:  true,
				Description: lang.Markdown("The name of the organization containing the targeted workspace(s)."),
			},
			"token": {
				Expr:       schema.LiteralTypeOnly(cty.String),
				IsOptional: true,
				Description: lang.Markdown("The token used to authenticate with the remote backend. If credentials for the " +
					"host are configured in the CLI Config File, then those will be used instead."),
			},
		},
		Blocks: map[string]*schema.BlockSchema{
			"workspaces": {
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"name": {
							Expr:       schema.LiteralTypeOnly(cty.String),
							IsOptional: true,
							Description: lang.Markdown("A workspace name used to map the default workspace to a named remote workspace. " +
								"When configured only the default workspace can be used. This option conflicts " +
								"with `prefix`"),
						},
						"prefix": {
							Expr:       schema.LiteralTypeOnly(cty.String),
							IsOptional: true,
							Description: lang.Markdown("A prefix used to filter workspaces using a single configuration. New workspaces " +
								"will automatically be prefixed with this prefix. If omitted only the default " +
								"workspace can be used. This option conflicts with `name`"),
						},
					},
				},
				Type: schema.BlockTypeObject,
			},
		},
	}
}
