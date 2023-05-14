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
		Blocks: make(map[string]*schema.BlockSchema, 0),
	}

	if v.GreaterThanOrEqual(v1_3_0) {
		// See https://github.com/hashicorp/terraform/pull/31425/files
		bodySchema.Attributes["accelerate"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.Bool},
			IsOptional:  true,
			Description: lang.Markdown("Whether to enable global Acceleration"),
		}
	}

	if v.GreaterThanOrEqual(v1_4_0) {
		// See https://github.com/hashicorp/terraform/pull/32631/files
		bodySchema.Attributes["secret_id"].IsRequired = false
		bodySchema.Attributes["secret_id"].IsOptional = true
		bodySchema.Attributes["secret_key"].IsRequired = false
		bodySchema.Attributes["secret_key"].IsOptional = true

		bodySchema.Attributes["security_token"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("TencentCloud Security Token of temporary access credentials. It can be sourced from the `TENCENTCLOUD_SECURITY_TOKEN` environment variable. Notice: for supported products, please refer to: [temporary key supported products](https://intl.cloud.tencent.com/document/product/598/10588)."),
			IsSensitive: true,
		}

		bodySchema.Blocks["assume_role"] = &schema.BlockSchema{
			Type:        schema.BlockTypeSet,
			MaxItems:    1,
			Description: lang.Markdown("The `assume_role` block. If provided, terraform will attempt to assume this role using the supplied credentials."),
			Body: &schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{
					"role_arn": {
						Constraint:  schema.LiteralType{Type: cty.String},
						IsRequired:  true,
						Description: lang.Markdown("The ARN of the role to assume. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_ARN`."),
					},
					"session_name": {
						Constraint:  schema.LiteralType{Type: cty.String},
						IsRequired:  true,
						Description: lang.Markdown("The session name to use when making the AssumeRole call. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SESSION_NAME`."),
					},
					"session_duration": {
						Constraint:  schema.LiteralType{Type: cty.Number},
						IsRequired:  true,
						Description: lang.Markdown("The duration of the session when making the AssumeRole call. Its value ranges from 0 to 43200(seconds), and default is 7200 seconds. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SESSION_DURATION`."),
					},
					"policy": {
						Constraint:  schema.LiteralType{Type: cty.String},
						IsOptional:  true,
						Description: lang.Markdown("A more restrictive policy when making the AssumeRole call. Its content must not contains `principal` elements. Notice: more syntax references, please refer to: [policies syntax logic](https://intl.cloud.tencent.com/document/product/598/10603)."),
					},
				},
			},
		}
	}

	return bodySchema
}
