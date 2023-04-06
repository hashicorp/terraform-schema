// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package backends

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func artifactoryBackend(v *version.Version) *schema.BodySchema {
	// https://github.com/hashicorp/terraform/blob/v0.12.0/backend/remote-state/artifactory/backend.go
	// https://github.com/hashicorp/terraform/blob/v1.0.0/internal/backend/remote-state/artifactory/backend.go
	// Docs:
	// https://github.com/hashicorp/terraform/blob/v1.0.0/website/docs/language/settings/backends/artifactory.html.md
	docsUrl := "https://www.terraform.io/docs/language/settings/backends/artifactory.html"
	return &schema.BodySchema{
		Description: lang.Markdown("Artifactory"),
		HoverURL:    docsUrl,
		DocsLink: &schema.DocsLink{
			URL: docsUrl,
		},
		Attributes: map[string]*schema.AttributeSchema{
			"username": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("Username"),
			},
			"password": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("Password"),
			},
			"url": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("Artfactory base URL (i.e. URL without repo and subpath)"),
			},
			"repo": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("The repository name"),
			},
			"subpath": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("Path within the repository"),
			},
		},
	}
}
