package backends

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func httpBackend(v *version.Version) *schema.BodySchema {
	// https://github.com/hashicorp/terraform/blob/v0.12.0/backend/remote-state/http/backend.go
	docsUrl := "https://www.terraform.io/docs/language/settings/backends/http.html"
	bodySchema := &schema.BodySchema{
		Description: lang.Markdown("HTTP (REST)"),
		HoverURL:    docsUrl,
		DocsLink: &schema.DocsLink{
			URL: docsUrl,
		},
		Attributes: map[string]*schema.AttributeSchema{
			"address": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsRequired:  true,
				Description: lang.Markdown("The address of the REST endpoint"),
			},
			"update_method": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("HTTP method to use when updating state"),
			},
			"lock_address": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The address of the lock REST endpoint"),
			},
			"unlock_address": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The address of the unlock REST endpoint"),
			},
			"lock_method": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The HTTP method to use when locking"),
			},
			"unlock_method": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The HTTP method to use when unlocking"),
			},
			"username": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The username for HTTP basic authentication"),
			},
			"password": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The password for HTTP basic authentication"),
			},
			"skip_cert_verification": {
				Expr:        schema.LiteralTypeOnly(cty.Bool),
				IsOptional:  true,
				Description: lang.Markdown("Whether to skip TLS verification."),
			},
		},
	}

	if v.GreaterThanOrEqual(v0_12_2) {
		// https://github.com/hashicorp/terraform/commit/5b6b1663
		bodySchema.Attributes["retry_max"] = &schema.AttributeSchema{
			Expr:        schema.LiteralTypeOnly(cty.Number),
			IsOptional:  true,
			Description: lang.Markdown("The number of HTTP request retries."),
		}
		bodySchema.Attributes["retry_wait_min"] = &schema.AttributeSchema{
			Expr:        schema.LiteralTypeOnly(cty.Number),
			IsOptional:  true,
			Description: lang.Markdown("The minimum time in seconds to wait between HTTP request attempts."),
		}
		bodySchema.Attributes["retry_wait_max"] = &schema.AttributeSchema{
			Expr:        schema.LiteralTypeOnly(cty.Number),
			IsOptional:  true,
			Description: lang.Markdown("The maximum time in seconds to wait between HTTP request attempts."),
		}
	}

	return bodySchema
}
