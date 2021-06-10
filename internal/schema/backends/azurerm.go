package backends

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func azureRmBackend(v *version.Version) *schema.BodySchema {
	// https://github.com/hashicorp/terraform/blob/v0.12.0/backend/remote-state/azure/backend.go
	docsUrl := "https://www.terraform.io/docs/language/settings/backends/azurerm.html"
	bodySchema := &schema.BodySchema{
		Description: lang.Markdown("Azure Blob Storage"),
		HoverURL:    docsUrl,
		DocsLink: &schema.DocsLink{
			URL: docsUrl,
		},
		Attributes: map[string]*schema.AttributeSchema{
			"storage_account_name": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsRequired:  true,
				Description: lang.Markdown("The name of the storage account."),
			},

			"container_name": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsRequired:  true,
				Description: lang.Markdown("The container name."),
			},

			"key": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsRequired:  true,
				Description: lang.Markdown("The blob key."),
			},

			"environment": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The Azure cloud environment."),
			},

			"access_key": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The access key."),
			},

			"sas_token": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("A SAS Token used to interact with the Blob Storage Account."),
			},

			"resource_group_name": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The resource group name."),
			},

			"client_id": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The Client ID."),
			},

			"client_secret": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The Client Secret."),
			},

			"subscription_id": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The Subscription ID."),
			},

			"tenant_id": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The Tenant ID."),
			},

			"use_msi": {
				Expr:        schema.LiteralTypeOnly(cty.Bool),
				IsOptional:  true,
				Description: lang.Markdown("Should Managed Service Identity be used?"),
			},

			"msi_endpoint": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("The Managed Service Identity Endpoint."),
			},

			"endpoint": {
				Expr:        schema.LiteralTypeOnly(cty.String),
				IsOptional:  true,
				Description: lang.Markdown("A custom Endpoint used to access the Azure Resource Manager API's."),
			},

			// Deprecated fields
			"arm_client_id": {
				Expr:         schema.LiteralTypeOnly(cty.String),
				IsOptional:   true,
				Description:  lang.Markdown("Replaced by `client_id`"),
				IsDeprecated: true,
			},

			"arm_client_secret": {
				Expr:         schema.LiteralTypeOnly(cty.String),
				IsOptional:   true,
				Description:  lang.Markdown("Replaced by `client_secret`"),
				IsDeprecated: true,
			},

			"arm_subscription_id": {
				Expr:         schema.LiteralTypeOnly(cty.String),
				IsOptional:   true,
				Description:  lang.Markdown("Replaced by `subscription_id`"),
				IsDeprecated: true,
			},

			"arm_tenant_id": {
				Expr:         schema.LiteralTypeOnly(cty.String),
				IsOptional:   true,
				Description:  lang.Markdown("Replaced by `tenant_id`"),
				IsDeprecated: true,
			},
		},
	}

	if v.GreaterThanOrEqual(v0_13_0) {
		// https://github.com/hashicorp/terraform/commit/0f85b283
		bodySchema.Attributes["snapshot"] = &schema.AttributeSchema{
			Expr:        schema.LiteralTypeOnly(cty.Bool),
			IsOptional:  true,
			Description: lang.Markdown("Enable/Disable automatic blob snapshotting"),
		}
	}

	if v.GreaterThanOrEqual(v0_13_1) {
		// https://github.com/hashicorp/terraform/commit/0d34e5d9
		bodySchema.Attributes["client_certificate_password"] = &schema.AttributeSchema{
			Expr:        schema.LiteralTypeOnly(cty.String),
			IsOptional:  true,
			Description: lang.Markdown("The password associated with the Client Certificate specified in `client_certificate_path`"),
		}
		bodySchema.Attributes["client_certificate_path"] = &schema.AttributeSchema{
			Expr:        schema.LiteralTypeOnly(cty.String),
			IsOptional:  true,
			Description: lang.Markdown("The path to the PFX file used as the Client Certificate when authenticating as a Service Principal"),
		}
		// https://github.com/hashicorp/terraform/commit/23b4c2db
		bodySchema.Attributes["metadata_host"] = &schema.AttributeSchema{
			Expr:        schema.LiteralTypeOnly(cty.String),
			IsRequired:  true,
			Description: lang.Markdown("The Metadata URL which will be used to obtain the Cloud Environment."),
		}
	}

	if v.GreaterThanOrEqual(v0_15_0) {
		// https://github.com/hashicorp/terraform/commit/b263e688
		delete(bodySchema.Attributes, "arm_client_id")
		delete(bodySchema.Attributes, "arm_client_secret")
		delete(bodySchema.Attributes, "arm_subscription_id")
		delete(bodySchema.Attributes, "arm_tenant_id")
		bodySchema.Attributes["use_azuread_auth"] = &schema.AttributeSchema{
			Expr:        schema.LiteralTypeOnly(cty.Bool),
			IsOptional:  true,
			Description: lang.Markdown("Should Terraform use AzureAD Authentication to access the Blob?"),
		}
	}

	return bodySchema
}
