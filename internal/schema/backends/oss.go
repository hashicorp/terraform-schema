// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package backends

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func ossBackend(v *version.Version) *schema.BodySchema {
	// https://github.com/hashicorp/terraform/blob/v0.12.2/backend/remote-state/oss/backend.go
	docsUrl := "https://www.terraform.io/docs/language/settings/backends/oss.html"
	bodySchema := &schema.BodySchema{
		Description: lang.Markdown("Alibaba Cloud Object Storage Service"),
		HoverURL:    docsUrl,
		DocsLink: &schema.DocsLink{
			URL: docsUrl,
		},
		Attributes: map[string]*schema.AttributeSchema{
			"access_key": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Alibaba Cloud Access Key ID"),
			},

			"secret_key": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Alibaba Cloud Access Secret Key"),
			},

			"security_token": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Alibaba Cloud Security Token"),
			},

			"region": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The region of the OSS bucket."),
			},
			"tablestore_endpoint": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("A custom endpoint for the TableStore API"),
			},
			"endpoint": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("A custom endpoint for the OSS API"),
			},

			"bucket": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("The name of the OSS bucket"),
			},

			"prefix": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The directory where state files will be saved inside the bucket"),
			},

			"key": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The path of the state file inside the bucket"),
			},

			"tablestore_table": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("TableStore table for state locking and consistency"),
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

	if v.GreaterThanOrEqual(v0_12_6) {
		// https://github.com/hashicorp/terraform/commit/a490dfa4
		bodySchema.Attributes["assume_role"] = &schema.AttributeSchema{
			IsOptional: true,
			Constraint: schema.Object{
				Attributes: schema.ObjectAttributes{
					"role_arn": {
						Constraint:  schema.LiteralType{Type: cty.String},
						IsRequired:  true,
						Description: lang.Markdown("The ARN of a RAM role to assume prior to making API calls."),
					},
					"session_name": {
						Constraint:  schema.LiteralType{Type: cty.String},
						IsOptional:  true,
						Description: lang.Markdown("The session name to use when assuming the role."),
					},
					"policy": {
						Constraint:  schema.LiteralType{Type: cty.String},
						IsOptional:  true,
						Description: lang.Markdown("The permissions applied when assuming a role. You cannot use this policy to grant permissions which exceed those of the role that is being assumed."),
					},
					"session_expiration": {
						Constraint:  schema.LiteralType{Type: cty.Number},
						IsOptional:  true,
						Description: lang.Markdown("The time after which the established session for assuming role expires."),
					},
				},
			},
		}
	}

	if v.GreaterThanOrEqual(v0_12_8) {
		// https://github.com/hashicorp/terraform/commit/b69c0b41
		bodySchema.Attributes["shared_credentials_file"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("This is the path to the shared credentials file. If this is not set and a profile is specified, `~/.aliyun/config.json` will be used."),
		}
		bodySchema.Attributes["profile"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("This is the Alibaba Cloud profile name as set in the shared credentials file. It can also be sourced from the `ALICLOUD_PROFILE` environment variable."),
		}
	}

	if v.GreaterThanOrEqual(v0_12_14) {
		// https://github.com/hashicorp/terraform/commit/bfae6271
		bodySchema.Attributes["ecs_role_name"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("The RAM Role Name attached on a ECS instance for API operations. You can retrieve this from the 'Access Control' section of the Alibaba Cloud console."),
		}
	}

	return bodySchema
}
