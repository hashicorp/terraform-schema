// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package backends

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func etcdv2Backend(v *version.Version) *schema.BodySchema {
	// https://github.com/hashicorp/terraform/blob/v0.12.0/backend/remote-state/etcdv2/backend.go
	// https://github.com/hashicorp/terraform/blob/v1.0.0/internal/backend/remote-state/etcdv2/backend.go
	docsUrl := "https://www.terraform.io/docs/language/settings/backends/etcd.html"
	bodySchema := &schema.BodySchema{
		Description: lang.Markdown("etcd v2.x"),
		HoverURL:    docsUrl,
		DocsLink: &schema.DocsLink{
			URL: docsUrl,
		},
		Attributes: map[string]*schema.AttributeSchema{
			"path": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("The path where to store the state"),
			},
			"endpoints": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("A space-separated list of the etcd endpoints"),
			},
			"username": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Username"),
			},
			"password": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Password"),
			},
		},
	}

	return bodySchema
}
