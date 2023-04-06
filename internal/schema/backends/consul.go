// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package backends

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func consulBackend(v *version.Version) *schema.BodySchema {
	// https://github.com/hashicorp/terraform/blob/v0.12.0/backend/remote-state/consul/backend.go
	// https://github.com/hashicorp/terraform/blob/v1.0.0/internal/backend/remote-state/consul/backend.go
	docsUrl := "https://www.terraform.io/docs/language/settings/backends/consul.html"
	bodySchema := &schema.BodySchema{
		Description: lang.Markdown("Consul KV store"),
		HoverURL:    docsUrl,
		DocsLink: &schema.DocsLink{
			URL: docsUrl,
		},
		Attributes: map[string]*schema.AttributeSchema{
			"path": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("Path to store state in Consul"),
			},

			"access_token": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Access token for a Consul ACL"),
			},

			"address": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Address to the Consul Cluster"),
			},

			"scheme": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Scheme to communicate to Consul with"),
			},

			"datacenter": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Datacenter to communicate with"),
			},

			"http_auth": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("HTTP Auth in the format of 'username:password'"),
			},

			"gzip": {
				Constraint:  schema.LiteralType{Type: cty.Bool},
				IsOptional:  true,
				Description: lang.Markdown("Compress the state data using gzip"),
			},

			"lock": {
				Constraint:  schema.LiteralType{Type: cty.Bool},
				IsOptional:  true,
				Description: lang.Markdown("Lock state access"),
			},

			"ca_file": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("A path to a PEM-encoded certificate authority used to verify the remote agent's certificate."),
			},

			"cert_file": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("A path to a PEM-encoded certificate provided to the remote agent; requires use of key_file."),
			},

			"key_file": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("A path to a PEM-encoded private key, required if cert_file is specified."),
			},
		},
	}

	return bodySchema
}
