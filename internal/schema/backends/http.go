// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("The address of the REST endpoint"),
			},
			"update_method": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("HTTP method to use when updating state"),
			},
			"lock_address": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The address of the lock REST endpoint"),
			},
			"unlock_address": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The address of the unlock REST endpoint"),
			},
			"lock_method": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The HTTP method to use when locking"),
			},
			"unlock_method": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The HTTP method to use when unlocking"),
			},
			"username": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The username for HTTP basic authentication"),
			},
			"password": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The password for HTTP basic authentication"),
			},
			"skip_cert_verification": {
				Constraint:  schema.LiteralType{Type: cty.Bool},
				IsOptional:  true,
				Description: lang.Markdown("Whether to skip TLS verification."),
			},
		},
	}

	if v.GreaterThanOrEqual(v0_12_2) {
		// https://github.com/hashicorp/terraform/commit/5b6b1663
		bodySchema.Attributes["retry_max"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.Number},
			IsOptional:  true,
			Description: lang.Markdown("The number of HTTP request retries."),
		}
		bodySchema.Attributes["retry_wait_min"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.Number},
			IsOptional:  true,
			Description: lang.Markdown("The minimum time in seconds to wait between HTTP request attempts."),
		}
		bodySchema.Attributes["retry_wait_max"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.Number},
			IsOptional:  true,
			Description: lang.Markdown("The maximum time in seconds to wait between HTTP request attempts."),
		}
	}

	if v.GreaterThanOrEqual(v1_4_0) {
		// https://github.com/hashicorp/terraform/commit/75e5ae27
		bodySchema.Attributes["client_ca_certificate_pem"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("A PEM-encoded CA certificate chain used by the client to verify server certificates during TLS authentication."),
		}
		bodySchema.Attributes["client_certificate_pem"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("A PEM-encoded certificate used by the server to verify the client during mutual TLS (mTLS) authentication."),
		}
		bodySchema.Attributes["client_private_key_pem"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("A PEM-encoded private key, required if client_certificate_pem is specified."),
		}
	}

	return bodySchema
}
