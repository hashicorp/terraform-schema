// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package backends

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func s3Backend(v *version.Version) *schema.BodySchema {
	// https://github.com/hashicorp/terraform/blob/v0.12.0/backend/remote-state/s3/backend.go
	docsUrl := "https://www.terraform.io/docs/language/settings/backends/s3.html"
	bodySchema := &schema.BodySchema{
		Description: lang.Markdown("Amazon S3 (with locking via DynamoDB)"),
		HoverURL:    docsUrl,
		DocsLink: &schema.DocsLink{
			URL: docsUrl,
		},
		Attributes: map[string]*schema.AttributeSchema{
			"bucket": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("The name of the S3 bucket"),
			},

			"key": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("The path to the state file inside the bucket"),
			},

			"region": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("AWS Region of the S3 Bucket and DynamoDB Table (if used)."),
			},

			"dynamodb_endpoint": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("A custom endpoint for the DynamoDB API"),
			},

			"endpoint": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("A custom endpoint for the S3 API"),
			},

			"iam_endpoint": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("A custom endpoint for the IAM API"),
			},

			"sts_endpoint": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("A custom endpoint for the STS API"),
			},

			"encrypt": {
				Constraint:  schema.LiteralType{Type: cty.Bool},
				IsOptional:  true,
				Description: lang.Markdown("Whether to enable server side encryption of the state file"),
			},

			"acl": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Canned ACL to be applied to the state file"),
			},

			"access_key": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("AWS access key"),
			},

			"secret_key": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("AWS secret key"),
			},

			"kms_key_id": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The ARN of a KMS Key to use for encrypting the state"),
			},

			"lock_table": {
				Constraint: schema.LiteralType{Type: cty.String},
				IsOptional: true,
				Description: lang.Markdown("DynamoDB table for state locking;\n\n" +
					"**DEPRECATED:** Please use the `dynamodb_table` attribute instead."),
				IsDeprecated: true,
			},

			"dynamodb_table": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("DynamoDB table for state locking and consistency"),
			},

			"profile": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("AWS profile name"),
			},

			"shared_credentials_file": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Path to a shared credentials file"),
			},

			"token": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("MFA token"),
			},

			"skip_credentials_validation": {
				Constraint:  schema.LiteralType{Type: cty.Bool},
				IsOptional:  true,
				Description: lang.Markdown("Skip the credentials validation via STS API."),
			},

			"skip_get_ec2_platforms": {
				Constraint:   schema.LiteralType{Type: cty.Bool},
				IsOptional:   true,
				Description:  lang.Markdown("Skip getting the supported EC2 platforms."),
				IsDeprecated: true,
			},

			"skip_region_validation": {
				Constraint:  schema.LiteralType{Type: cty.Bool},
				IsOptional:  true,
				Description: lang.Markdown("Skip static validation of region name."),
			},

			"skip_requesting_account_id": {
				Constraint:   schema.LiteralType{Type: cty.Bool},
				IsOptional:   true,
				Description:  lang.Markdown("Skip requesting the account ID."),
				IsDeprecated: true,
			},

			"skip_metadata_api_check": {
				Constraint:  schema.LiteralType{Type: cty.Bool},
				IsOptional:  true,
				Description: lang.Markdown("Skip the AWS Metadata API check."),
			},

			"role_arn": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The role to be assumed"),
			},

			"session_name": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The session name to use when assuming the role."),
			},

			"external_id": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The external ID to use when assuming the role"),
			},

			"assume_role_policy": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The permissions applied when assuming a role."),
			},

			"workspace_key_prefix": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The prefix applied to the non-default state path inside the bucket."),
			},

			"force_path_style": {
				Constraint:  schema.LiteralType{Type: cty.Bool},
				IsOptional:  true,
				Description: lang.Markdown("Force s3 to use path style api."),
			},

			"max_retries": {
				Constraint:  schema.LiteralType{Type: cty.Number},
				IsOptional:  true,
				Description: lang.Markdown("The maximum number of times an AWS API request is retried on retryable failure."),
			},
		},
	}

	if v.GreaterThanOrEqual(v0_12_8) {
		// https://github.com/hashicorp/terraform/commit/5e3c3baf
		bodySchema.Attributes["sse_customer_key"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("The base64-encoded encryption key to use for server-side encryption with customer-provided keys (SSE-C)."),
			IsSensitive: true,
		}
	}

	if v.GreaterThanOrEqual(v0_13_0) {
		// https://github.com/hashicorp/terraform/commit/ba081aa1

		delete(bodySchema.Attributes, "lock_table")
		delete(bodySchema.Attributes, "skip_get_ec2_platforms")
		delete(bodySchema.Attributes, "skip_requesting_account_id")

		bodySchema.Attributes["assume_role_duration_seconds"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.Number},
			IsOptional:  true,
			Description: lang.Markdown("Seconds to restrict the assume role session duration."),
		}
		bodySchema.Attributes["assume_role_policy_arns"] = &schema.AttributeSchema{
			Constraint:  schema.Set{Elem: schema.LiteralType{Type: cty.String}},
			IsOptional:  true,
			Description: lang.Markdown("Amazon Resource Names (ARNs) of IAM Policies describing further restricting permissions for the IAM Role being assumed."),
		}

		bodySchema.Attributes["assume_role_tags"] = &schema.AttributeSchema{
			Constraint: schema.Map{
				Elem: schema.LiteralType{Type: cty.String},
			},
			IsOptional:  true,
			Description: lang.Markdown("Assume role session tags."),
		}

		bodySchema.Attributes["assume_role_transitive_tag_keys"] = &schema.AttributeSchema{
			Constraint: schema.Set{
				Elem: schema.LiteralType{Type: cty.String},
			},
			IsOptional:  true,
			Description: lang.Markdown("Assume role session tag keys to pass to any subsequent sessions."),
		}
	}
	if v.GreaterThanOrEqual(v1_6_0) {
		bodySchema.Attributes["region"].IsRequired = false
		bodySchema.Attributes["region"].IsOptional = true

		// deprecations
		bodySchema.Attributes["assume_role_duration_seconds"].IsDeprecated = true
		bodySchema.Attributes["assume_role_policy_arns"].IsDeprecated = true
		bodySchema.Attributes["assume_role_policy"].IsDeprecated = true
		bodySchema.Attributes["assume_role_tags"].IsDeprecated = true
		bodySchema.Attributes["assume_role_transitive_tag_keys"].IsDeprecated = true
		bodySchema.Attributes["dynamodb_endpoint"].IsDeprecated = true
		bodySchema.Attributes["endpoint"].IsDeprecated = true
		bodySchema.Attributes["external_id"].IsDeprecated = true
		bodySchema.Attributes["force_path_style"].IsDeprecated = true
		bodySchema.Attributes["iam_endpoint"].IsDeprecated = true
		bodySchema.Attributes["role_arn"].IsDeprecated = true
		bodySchema.Attributes["session_name"].IsDeprecated = true
		bodySchema.Attributes["shared_credentials_file"].IsDeprecated = true
		bodySchema.Attributes["sts_endpoint"].IsDeprecated = true

		// new fields
		bodySchema.Attributes["allowed_account_ids"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.Set(cty.String)},
			IsOptional:  true,
			Description: lang.Markdown("List of allowed AWS account IDs."),
		}
		bodySchema.Attributes["ec2_metadata_service_endpoint"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("Address of the EC2 metadata service (IMDS) endpoint to use."),
		}
		bodySchema.Attributes["ec2_metadata_service_endpoint_mode"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("Mode to use in communicating with the metadata service."),
		}
		bodySchema.Attributes["endpoints"] = &schema.AttributeSchema{
			Constraint: schema.Object{
				Attributes: schema.ObjectAttributes{
					"dynamodb": &schema.AttributeSchema{
						Constraint:  schema.LiteralType{Type: cty.String},
						IsOptional:  true,
						Description: lang.Markdown("A custom endpoint for the DynamoDB API"),
					},

					"iam": &schema.AttributeSchema{
						Constraint:  schema.LiteralType{Type: cty.String},
						IsOptional:  true,
						Description: lang.Markdown("A custom endpoint for the IAM API"),
					},

					"s3": &schema.AttributeSchema{
						Constraint:  schema.LiteralType{Type: cty.String},
						IsOptional:  true,
						Description: lang.Markdown("A custom endpoint for the S3 API"),
					},

					"sts": &schema.AttributeSchema{
						Constraint:  schema.LiteralType{Type: cty.String},
						IsOptional:  true,
						Description: lang.Markdown("A custom endpoint for the STS API"),
					},
				},
			},
		}
		bodySchema.Attributes["forbidden_account_ids"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.Set(cty.String)},
			IsOptional:  true,
			Description: lang.Markdown("List of forbidden AWS account IDs."),
		}
		bodySchema.Attributes["sts_region"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("AWS region for STS."),
		}
		bodySchema.Attributes["retry_mode"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("Specifies how retries are attempted."),
		}
		bodySchema.Attributes["shared_config_files"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.Set(cty.String)},
			IsOptional:  true,
			Description: lang.Markdown("List of paths to shared config files"),
		}
		bodySchema.Attributes["shared_credentials_files"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.Set(cty.String)},
			IsOptional:  true,
			Description: lang.Markdown("List of paths to shared credentials files"),
		}
		bodySchema.Attributes["use_path_style"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.Bool},
			IsOptional:  true,
			Description: lang.Markdown("Enable path-style S3 URLs."),
		}
		bodySchema.Attributes["assume_role"] = &schema.AttributeSchema{
			Constraint: schema.Object{
				Attributes: schema.ObjectAttributes{
					"role_arn": &schema.AttributeSchema{
						Constraint:  schema.LiteralType{Type: cty.String},
						IsRequired:  true,
						Description: lang.Markdown("The role to be assumed."),
					},

					"duration": &schema.AttributeSchema{
						Constraint:  schema.LiteralType{Type: cty.String},
						IsOptional:  true,
						Description: lang.Markdown("The duration, between 15 minutes and 12 hours, of the role session. Valid time units are ns, us (or µs), ms, s, h, or m."),
					},

					"external_id": &schema.AttributeSchema{
						Constraint:  schema.LiteralType{Type: cty.String},
						IsOptional:  true,
						Description: lang.Markdown("The external ID to use when assuming the role"),
					},

					"policy": &schema.AttributeSchema{
						Constraint:  schema.LiteralType{Type: cty.String},
						IsOptional:  true,
						Description: lang.Markdown("IAM Policy JSON describing further restricting permissions for the IAM Role being assumed."),
					},

					"policy_arns": &schema.AttributeSchema{
						Constraint:  schema.LiteralType{Type: cty.Set(cty.String)},
						IsOptional:  true,
						Description: lang.Markdown("Amazon Resource Names (ARNs) of IAM Policies describing further restricting permissions for the IAM Role being assumed."),
					},

					"session_name": &schema.AttributeSchema{
						Constraint:  schema.LiteralType{Type: cty.String},
						IsOptional:  true,
						Description: lang.Markdown("The session name to use when assuming the role."),
					},

					"source_identity": &schema.AttributeSchema{
						Constraint:  schema.LiteralType{Type: cty.String},
						IsOptional:  true,
						Description: lang.Markdown("Source identity specified by the principal assuming the role."),
					},

					"tags": &schema.AttributeSchema{
						Constraint:  schema.LiteralType{Type: cty.Map(cty.String)},
						IsOptional:  true,
						Description: lang.Markdown("Assume role session tags."),
					},

					"transitive_tag_keys": &schema.AttributeSchema{
						Constraint:  schema.LiteralType{Type: cty.Set(cty.String)},
						IsOptional:  true,
						Description: lang.Markdown("Assume role session tag keys to pass to any subsequent sessions."),
					},
				},
			},
		}
		bodySchema.Attributes["assume_role_with_web_identity"] = &schema.AttributeSchema{
			Constraint: schema.Object{
				Attributes: schema.ObjectAttributes{
					"role_arn": &schema.AttributeSchema{
						Constraint:  schema.LiteralType{Type: cty.String},
						IsRequired:  true,
						Description: lang.Markdown("The role to be assumed."),
					},

					"duration": &schema.AttributeSchema{
						Constraint:  schema.LiteralType{Type: cty.String},
						IsOptional:  true,
						Description: lang.Markdown("The duration, between 15 minutes and 12 hours, of the role session. Valid time units are ns, us (or µs), ms, s, h, or m."),
					},

					"policy": &schema.AttributeSchema{
						Constraint:  schema.LiteralType{Type: cty.String},
						IsOptional:  true,
						Description: lang.Markdown("IAM Policy JSON describing further restricting permissions for the IAM Role being assumed."),
					},

					"policy_arns": &schema.AttributeSchema{
						Constraint:  schema.LiteralType{Type: cty.Set(cty.String)},
						IsOptional:  true,
						Description: lang.Markdown("Amazon Resource Names (ARNs) of IAM Policies describing further restricting permissions for the IAM Role being assumed."),
					},

					"session_name": &schema.AttributeSchema{
						Constraint:  schema.LiteralType{Type: cty.String},
						IsOptional:  true,
						Description: lang.Markdown("The session name to use when assuming the role."),
					},

					"web_identity_token": &schema.AttributeSchema{
						Constraint:  schema.LiteralType{Type: cty.String},
						IsOptional:  true,
						Description: lang.Markdown("Value of a web identity token from an OpenID Connect (OIDC) or OAuth provider."),
					},

					"web_identity_token_file": &schema.AttributeSchema{
						Constraint:  schema.LiteralType{Type: cty.String},
						IsOptional:  true,
						Description: lang.Markdown("File containing a web identity token from an OpenID Connect (OIDC) or OAuth provider."),
					},
				},
			},
		}
		bodySchema.Attributes["use_legacy_workflow"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.Bool},
			IsOptional:  true,
			Description: lang.Markdown("Use the legacy authentication workflow, preferring environment variables over backend configuration."),
		}
		bodySchema.Attributes["custom_ca_bundle"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("File containing custom root and intermediate certificates."),
		}
		bodySchema.Attributes["http_proxy"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("Address of an HTTP proxy to use when accessing the AWS API."),
		}
		bodySchema.Attributes["insecure"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.Bool},
			IsOptional:  true,
			Description: lang.Markdown("Whether to explicitly allow the backend to perform insecure SSL requests."),
		}
		bodySchema.Attributes["use_fips_endpoint"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.Bool},
			IsOptional:  true,
			Description: lang.Markdown("Force the backend to resolve endpoints with FIPS capability."),
		}
		bodySchema.Attributes["use_dualstack_endpoint"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.Bool},
			IsOptional:  true,
			Description: lang.Markdown("Force the backend to resolve endpoints with DualStack capability."),
		}
	}

	return bodySchema
}
