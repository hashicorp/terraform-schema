// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package backends

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func localBackend(v *version.Version) *schema.BodySchema {
	docsUrl := "https://www.terraform.io/docs/language/settings/backends/local.html"
	return &schema.BodySchema{
		Description: lang.Markdown("Local (filesystem) backend, locks state using system APIs"),
		HoverURL:    docsUrl,
		DocsLink: &schema.DocsLink{
			URL: docsUrl,
		},
		Attributes: map[string]*schema.AttributeSchema{
			"path": {
				Constraint:  schema.LiteralType{Type: cty.String},
				Description: lang.Markdown("The path to the tfstate file. This defaults to `terraform.tfstate` relative to the root module."),
				IsOptional:  true,
			},
			"workspace_dir": {
				Constraint:  schema.LiteralType{Type: cty.String},
				Description: lang.Markdown("The path to non-default workspaces."),
				IsOptional:  true,
			},
		},
	}
}
