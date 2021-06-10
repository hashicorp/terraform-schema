package backends

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func gcsBackend(v *version.Version) *schema.BodySchema {
	// https://github.com/hashicorp/terraform/blob/v0.12.0/backend/remote-state/gcs/backend.go
	docsUrl := "https://www.terraform.io/docs/language/settings/backends/gcs.html"
	bodySchema := &schema.BodySchema{
		Description: lang.Markdown("Google Cloud Storage"),
		HoverURL:    docsUrl,
		DocsLink: &schema.DocsLink{
			URL: docsUrl,
		},
		Attributes: map[string]*schema.AttributeSchema{
			"bucket": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsRequired:  true,
				Description: lang.Markdown("The name of the Google Cloud Storage bucket"),
			},

			"path": {
				Expr:         schema.LiteralTypeOnly(cty.String),
				IsOptional:   true,
				Description:  lang.Markdown("Path of the default state file;\nDEPRECATED: Use the `prefix` option instead"),
				IsDeprecated: true,
			},

			"prefix": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The directory where state files will be saved inside the bucket"),
			},

			"credentials": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("Google Cloud JSON Account Key"),
			},

			"encryption_key": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("A 32 byte base64 encoded 'customer supplied encryption key' used to encrypt all state."),
			},
		},
	}

	if v.GreaterThanOrEqual(v0_12_10) {
		// https://github.com/hashicorp/terraform/commit/f6c90c1d
		bodySchema.Attributes["access_token"] = &schema.AttributeSchema{
			Expr:        schema.LiteralTypeOnly(cty.String),
			IsOptional:  true,
			Description: lang.Markdown("An OAuth2 token used for GCP authentication"),
		}
	}

	if v.GreaterThanOrEqual(v0_14_0) {
		// https://github.com/hashicorp/terraform/commit/c43731a0
		bodySchema.Attributes["impersonate_service_account"] = &schema.AttributeSchema{
			Expr:        schema.LiteralTypeOnly(cty.String),
			IsOptional:  true,
			Description: lang.Markdown("The service account to impersonate for all Google API Calls"),
		}

		bodySchema.Attributes["impersonate_service_account_delegates"] = &schema.AttributeSchema{
			Expr: schema.ExprConstraints{
				schema.ListExpr{Elem: schema.LiteralTypeOnly(cty.String)},
			},
			IsOptional:  true,
			Description: lang.Markdown("The delegation chain for the impersonated service account"),
		}
	}

	if v.GreaterThanOrEqual(v0_15_0) {
		// https://github.com/hashicorp/terraform/commit/3b9c5e5b
		delete(bodySchema.Attributes, "path")
	}

	return bodySchema
}
