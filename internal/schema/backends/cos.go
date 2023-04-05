// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package backends

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func cosBackend(v *version.Version) *schema.BodySchema {
	// https://github.com/hashicorp/terraform/blob/v0.13.0/backend/remote-state/cos/backend.go
	docsUrl := "https://www.terraform.io/docs/language/settings/backends/cos.html"
	bodySchema := &schema.BodySchema{
		Description: lang.Markdown("Tencent Cloud Object Storage"),
		HoverURL:    docsUrl,
		DocsLink: &schema.DocsLink{
			URL: docsUrl,
		},
		Attributes: map[string]*schema.AttributeSchema{
			"secret_id": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("Secret id of Tencent Cloud"),
			},
			"secret_key": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				IsSensitive: true,
				Description: lang.Markdown("Secret key of Tencent Cloud"),
			},
			"region": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("The region of the COS bucket"),
			},
			"bucket": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("The name of the COS bucket"),
			},
			"prefix": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The directory for saving the state file in bucket"),
			},
			"key": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The path for saving the state file in bucket"),
			},
			"encrypt": {
				Constraint:  schema.LiteralType{Type: cty.Bool},
				IsOptional:  true,
				Description: lang.Markdown("Whether to enable server side encryption of the state file"),
			},
			"acl": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Object ACL to be applied to the state file"),
			},
		},
	}

	return bodySchema
}
