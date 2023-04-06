// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package backends

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func mantaBackend(v *version.Version) *schema.BodySchema {
	// https://github.com/hashicorp/terraform/blob/v0.12.0/backend/remote-state/manta/backend.go
	// https://github.com/hashicorp/terraform/blob/v1.0.0/internal/backend/remote-state/manta/backend.go
	// Docs:
	// https://github.com/hashicorp/terraform/blob/v1.0.0/website/docs/language/settings/backends/manta.html.md
	docsUrl := "https://www.terraform.io/docs/language/settings/backends/manta.html"
	bodySchema := &schema.BodySchema{
		Description: lang.Markdown("Manta (Triton Object Storage)"),
		HoverURL:    docsUrl,
		DocsLink: &schema.DocsLink{
			URL: docsUrl,
		},
		Attributes: map[string]*schema.AttributeSchema{
			"account": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("The name of the Manta account."),
			},

			"user": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The username of the Triton account used to authenticate with the Triton API."),
			},

			"url": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The Manta API Endpoint."),
			},

			"key_material": {
				Constraint: schema.LiteralType{Type: cty.String},
				IsOptional: true,
				Description: lang.Markdown("The private key of an SSH key associated with the Triton account to be used.\n\n" +
					"If this is not set, the private key corresponding to the fingerprint in `key_id` must be available via an SSH Agent."),
			},

			"key_id": {
				Constraint: schema.LiteralType{Type: cty.String},
				IsRequired: true,
				Description: lang.Markdown("The fingerprint of the public key matching the key specified in `key_path`.\n\n" +
					"It can be obtained via the command `ssh-keygen -l -E md5 -f /path/to/key`."),
			},

			"insecure_skip_tls_verify": {
				Constraint: schema.LiteralType{Type: cty.Bool},
				IsOptional: true,
				Description: lang.Markdown("This allows skipping TLS verification of the Triton endpoint. " +
					"It is useful when connecting to a temporary Triton installation such as Cloud-On-A-Laptop " +
					"which does not generally use a certificate signed by a trusted root CA."),
			},

			"path": {
				Constraint: schema.LiteralType{Type: cty.String},
				IsRequired: true,
				Description: lang.Markdown("The path relative to your private storage directory where the state file will be stored. " +
					"**Please Note:** If this path does not exist, then the backend will create this folder location as part of backend creation."),
			},

			"object_name": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The name of the state file"),
			},
		},
	}

	return bodySchema
}
