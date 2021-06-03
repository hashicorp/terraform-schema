package backends

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func etcdv3Backend(v *version.Version) *schema.BodySchema {
	// https://github.com/hashicorp/terraform/blob/v0.12.0/backend/remote-state/etcdv3/backend.go
	// https://github.com/hashicorp/terraform/blob/v1.0.0/internal/backend/remote-state/etcdv3/backend.go
	docsUrl := "https://www.terraform.io/docs/language/settings/backends/etcdv3.html"
	bodySchema := &schema.BodySchema{
		Description: lang.Markdown("etcd v3"),
		HoverURL:    docsUrl,
		DocsLink: &schema.DocsLink{
			URL: docsUrl,
		},
		Attributes: map[string]*schema.AttributeSchema{
			"endpoints": {
				Expr: schema.ExprConstraints{
					schema.ListExpr{
						Elem:     schema.LiteralTypeOnly(cty.String),
						MinItems: 1,
					},
				},
				IsRequired:  true,
				Description: lang.Markdown("Endpoints for the etcd cluster."),
			},

			"username": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("Username used to connect to the etcd cluster."),
			},

			"password": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("Password used to connect to the etcd cluster."),
			},

			"prefix": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("An optional prefix to be added to keys when to storing state in etcd."),
			},

			"lock": {
				Expr:        schema.LiteralTypeOnly(cty.Bool),
				IsOptional:  true,
				Description: lang.Markdown("Whether to lock state access."),
			},

			"cacert_path": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The path to a PEM-encoded CA bundle with which to verify certificates of TLS-enabled etcd servers."),
			},

			"cert_path": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The path to a PEM-encoded certificate to provide to etcd for secure client identification."),
			},

			"key_path": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The path to a PEM-encoded key to provide to etcd for secure client identification."),
			},
		},
	}

	return bodySchema
}
