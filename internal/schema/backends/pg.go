// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package backends

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func pgBackend(v *version.Version) *schema.BodySchema {
	// https://github.com/hashicorp/terraform/blob/v0.12.0/backend/remote-state/pg/backend.go
	docsUrl := "https://www.terraform.io/docs/language/settings/backends/pg.html"
	bodySchema := &schema.BodySchema{
		Description: lang.Markdown("PostgreSQL (v10+)"),
		HoverURL:    docsUrl,
		DocsLink: &schema.DocsLink{
			URL: docsUrl,
		},
		Attributes: map[string]*schema.AttributeSchema{
			"conn_str": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("Postgres connection string; a `postgres://` URL"),
			},

			"schema_name": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Name of the automatically managed Postgres schema to store state"),
			},
		},
	}

	if v.GreaterThanOrEqual(v0_12_8) {
		// https://github.com/hashicorp/terraform/commit/be5280e4
		bodySchema.Attributes["skip_schema_creation"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.Bool},
			IsOptional:  true,
			Description: lang.Markdown("If set to `true`, Terraform won't try to create the Postgres schema"),
		}
	}

	if v.GreaterThanOrEqual(v0_14_0) {
		// https://github.com/hashicorp/terraform/commit/12a0a21c
		bodySchema.Attributes["skip_table_creation"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.Bool},
			IsOptional:  true,
			Description: lang.Markdown("If set to `true`, Terraform won't try to create the Postgres table"),
		}

		bodySchema.Attributes["skip_index_creation"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.Bool},
			IsOptional:  true,
			Description: lang.Markdown("If set to `true`, Terraform won't try to create the Postgres index"),
		}
	}

	return bodySchema
}
