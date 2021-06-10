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
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsRequired:  true,
				Description: lang.Markdown("The name of the S3 bucket"),
			},

			"key": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsRequired:  true,
				Description: lang.Markdown("The path to the state file inside the bucket"),
			},

			"region": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsRequired:  true,
				Description: lang.Markdown("AWS Region of the S3 Bucket and DynamoDB Table (if used)."),
			},

			"dynamodb_endpoint": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("A custom endpoint for the DynamoDB API"),
			},

			"endpoint": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("A custom endpoint for the S3 API"),
			},

			"iam_endpoint": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("A custom endpoint for the IAM API"),
			},

			"sts_endpoint": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("A custom endpoint for the STS API"),
			},

			"encrypt": {
				Expr:        schema.LiteralTypeOnly(cty.Bool),
				IsOptional:  true,
				Description: lang.Markdown("Whether to enable server side encryption of the state file"),
			},

			"acl": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("Canned ACL to be applied to the state file"),
			},

			"access_key": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("AWS access key"),
			},

			"secret_key": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("AWS secret key"),
			},

			"kms_key_id": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The ARN of a KMS Key to use for encrypting the state"),
			},

			"lock_table": {
				Expr:       schema.LiteralTypeOnly(cty.String),
				IsOptional: true,
				Description: lang.Markdown("DynamoDB table for state locking;\n\n" +
					"**DEPRECATED:** Please use the `dynamodb_table` attribute instead."),
				IsDeprecated: true,
			},

			"dynamodb_table": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("DynamoDB table for state locking and consistency"),
			},

			"profile": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("AWS profile name"),
			},

			"shared_credentials_file": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("Path to a shared credentials file"),
			},

			"token": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("MFA token"),
			},

			"skip_credentials_validation": {
				Expr:        schema.LiteralTypeOnly(cty.Bool),
				IsOptional:  true,
				Description: lang.Markdown("Skip the credentials validation via STS API."),
			},

			"skip_get_ec2_platforms": {
				Expr:         schema.LiteralTypeOnly(cty.Bool),
				IsOptional:   true,
				Description:  lang.Markdown("Skip getting the supported EC2 platforms."),
				IsDeprecated: true,
			},

			"skip_region_validation": {
				Expr:        schema.LiteralTypeOnly(cty.Bool),
				IsOptional:  true,
				Description: lang.Markdown("Skip static validation of region name."),
			},

			"skip_requesting_account_id": {
				Expr:         schema.LiteralTypeOnly(cty.Bool),
				IsOptional:   true,
				Description:  lang.Markdown("Skip requesting the account ID."),
				IsDeprecated: true,
			},

			"skip_metadata_api_check": {
				Expr:        schema.LiteralTypeOnly(cty.Bool),
				IsOptional:  true,
				Description: lang.Markdown("Skip the AWS Metadata API check."),
			},

			"role_arn": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The role to be assumed"),
			},

			"session_name": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The session name to use when assuming the role."),
			},

			"external_id": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The external ID to use when assuming the role"),
			},

			"assume_role_policy": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The permissions applied when assuming a role."),
			},

			"workspace_key_prefix": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The prefix applied to the non-default state path inside the bucket."),
			},

			"force_path_style": {
				Expr:        schema.LiteralTypeOnly(cty.Bool),
				IsOptional:  true,
				Description: lang.Markdown("Force s3 to use path style api."),
			},

			"max_retries": {
				Expr:        schema.LiteralTypeOnly(cty.Number),
				IsOptional:  true,
				Description: lang.Markdown("The maximum number of times an AWS API request is retried on retryable failure."),
			},
		},
	}

	if v.GreaterThanOrEqual(v0_12_8) {
		// https://github.com/hashicorp/terraform/commit/5e3c3baf
		bodySchema.Attributes["sse_customer_key"] = &schema.AttributeSchema{
			Expr:        schema.LiteralTypeOnly(cty.String),
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
			Expr:        schema.LiteralTypeOnly(cty.Number),
			IsOptional:  true,
			Description: lang.Markdown("Seconds to restrict the assume role session duration."),
		}
		bodySchema.Attributes["assume_role_policy_arns"] = &schema.AttributeSchema{
			Expr: schema.ExprConstraints{
				schema.SetExpr{Elem: schema.LiteralTypeOnly(cty.String)},
			},
			IsOptional:  true,
			Description: lang.Markdown("Amazon Resource Names (ARNs) of IAM Policies describing further restricting permissions for the IAM Role being assumed."),
		}

		bodySchema.Attributes["assume_role_tags"] = &schema.AttributeSchema{
			Expr: schema.ExprConstraints{
				schema.MapExpr{
					Elem: schema.LiteralTypeOnly(cty.String),
				},
			},
			IsOptional:  true,
			Description: lang.Markdown("Assume role session tags."),
		}

		bodySchema.Attributes["assume_role_transitive_tag_keys"] = &schema.AttributeSchema{
			Expr: schema.ExprConstraints{
				schema.SetExpr{Elem: schema.LiteralTypeOnly(cty.String)},
			},
			IsOptional:  true,
			Description: lang.Markdown("Assume role session tag keys to pass to any subsequent sessions."),
		}
	}

	return bodySchema
}
