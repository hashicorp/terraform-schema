// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package backends

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func swiftBackend(v *version.Version) *schema.BodySchema {
	// https://github.com/hashicorp/terraform/blob/v0.12.0/backend/remote-state/swift/backend.go
	docsUrl := "https://www.terraform.io/docs/language/settings/backends/swift.html"
	bodySchema := &schema.BodySchema{
		Description: lang.Markdown("Swift (OpenStack object storage)"),
		HoverURL:    docsUrl,
		DocsLink: &schema.DocsLink{
			URL: docsUrl,
		},
		Attributes: map[string]*schema.AttributeSchema{
			"auth_url": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("The Identity authentication URL."),
			},

			"user_id": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("User ID to login with."),
			},

			"user_name": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Username to login with."),
			},

			"tenant_id": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The ID of the Tenant (Identity v2) or Project (Identity v3) to login with."),
			},

			"tenant_name": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The name of the Tenant (Identity v2) or Project (Identity v3) to login with."),
			},

			"password": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				IsSensitive: true,
				Description: lang.Markdown("Password to login with."),
			},

			"token": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Authentication token to use as an alternative to username/password."),
			},

			"domain_id": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The ID of the Domain to scope to (Identity v3)."),
			},

			"domain_name": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The name of the Domain to scope to (Identity v3)."),
			},

			"region_name": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("The name of the Region to use."),
			},

			"insecure": {
				Constraint:  schema.LiteralType{Type: cty.Bool},
				IsOptional:  true,
				Description: lang.Markdown("Trust self-signed certificates."),
			},

			"endpoint_type": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The catalog endpoint type to use."),
			},

			"cacert_file": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("A Custom CA certificate."),
			},

			"cert": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("A client certificate to authenticate with."),
			},

			"key": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("A client private key to authenticate with."),
			},

			"path": {
				Constraint:   schema.LiteralType{Type: cty.String},
				IsOptional:   true,
				Description:  lang.Markdown("Swift container path to use; **DEPRECATED:** Use `container` instead."),
				IsDeprecated: true,
			},

			"container": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Swift container to create"),
			},

			"archive_path": {
				Constraint:   schema.LiteralType{Type: cty.String},
				IsOptional:   true,
				Description:  lang.Markdown("Swift container path to archive state to; **DEPRECATED:** Use `archive_container` instead"),
				IsDeprecated: true,
			},

			"archive_container": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Swift container to archive state to."),
			},

			"expire_after": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Archive object expiry duration."),
			},

			"lock": {
				Constraint:  schema.LiteralType{Type: cty.Bool},
				IsOptional:  true,
				Description: lang.Markdown("Lock state access"),
			},
		},
	}

	if v.GreaterThanOrEqual(v0_12_2) {
		// https://github.com/hashicorp/terraform/commit/d8343aa9
		bodySchema.Attributes["user_domain_name"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("The name of the domain where the user resides (Identity v3)."),
		}
		bodySchema.Attributes["user_domain_id"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("The ID of the domain where the user resides (Identity v3)."),
		}
		bodySchema.Attributes["project_domain_name"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("The name of the domain where the project resides (Identity v3)."),
		}
		bodySchema.Attributes["project_domain_id"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("The ID of the domain where the project resides (Identity v3)."),
		}
		bodySchema.Attributes["default_domain"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("The name of the Domain ID to scope to if no other domain is specified. Defaults to `default` (Identity v3)."),
		}
		bodySchema.Attributes["cloud"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("An entry in a `clouds.yaml` file to use."),
		}

		// https://github.com/hashicorp/terraform/commit/d06609dd
		bodySchema.Attributes["application_credential_id"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("Application Credential ID to login with."),
		}

		bodySchema.Attributes["application_credential_name"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("Application Credential name to login with."),
		}

		bodySchema.Attributes["application_credential_secret"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("Application Credential secret to login with."),
		}
	}

	if v.GreaterThanOrEqual(v0_12_4) {
		// https://github.com/hashicorp/terraform/commit/cd7bfba1
		bodySchema.Attributes["state_name"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("Name of state object in container"),
		}
	}

	if v.GreaterThanOrEqual(v0_13_0) {
		// https://github.com/hashicorp/terraform/commit/bd344f9d
		bodySchema.Attributes["auth_url"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("The Identity authentication URL."),
		}
		bodySchema.Attributes["region_name"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("The name of the Region to use."),
		}
		bodySchema.Attributes["swauth"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.Bool},
			IsOptional:  true,
			Description: lang.Markdown("Use Swift's authentication system instead of Keystone."),
		}
		bodySchema.Attributes["allow_reauth"] = &schema.AttributeSchema{
			Constraint: schema.LiteralType{Type: cty.Bool},
			IsOptional: true,
			Description: lang.Markdown("If set to `true`, OpenStack authorization will be perfomed\n" +
				"automatically, if the initial auth token get expired. This is useful,\n" +
				"when the token TTL is low or the overall Terraform provider execution\n" +
				"time expected to be greater than the initial token TTL."),
		}
		bodySchema.Attributes["max_retries"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.Number},
			IsOptional:  true,
			Description: lang.Markdown("How many times HTTP connection should be retried until giving up."),
		}
		bodySchema.Attributes["disable_no_cache_header"] = &schema.AttributeSchema{
			Constraint: schema.LiteralType{Type: cty.Bool},
			IsOptional: true,
			Description: lang.Markdown("If set to `true`, the HTTP `Cache-Control: no-cache` header will " +
				"not be added by default to all API requests."),
		}
	}

	return bodySchema
}
