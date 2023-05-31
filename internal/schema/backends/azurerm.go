// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("The name of the storage account."),
			},

			"container_name": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("The container name."),
			},

			"key": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("The blob key."),
			},

			"environment": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The Azure cloud environment."),
			},

			"access_key": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The access key."),
			},

			"sas_token": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("A SAS Token used to interact with the Blob Storage Account."),
			},

			"resource_group_name": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The resource group name."),
			},

			"client_id": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The Client ID."),
			},

			"client_secret": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The Client Secret."),
			},

			"subscription_id": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The Subscription ID."),
			},

			"tenant_id": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The Tenant ID."),
			},

			"use_msi": {
				Constraint:  schema.LiteralType{Type: cty.Bool},
				IsOptional:  true,
				Description: lang.Markdown("Should Managed Service Identity be used?"),
			},

			"msi_endpoint": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The Managed Service Identity Endpoint."),
			},

			"endpoint": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("A custom Endpoint used to access the Azure Resource Manager API's."),
			},

			// Deprecated fields
			"arm_client_id": {
				Constraint:   schema.LiteralType{Type: cty.String},
				IsOptional:   true,
				Description:  lang.Markdown("Replaced by `client_id`"),
				IsDeprecated: true,
			},

			"arm_client_secret": {
				Constraint:   schema.LiteralType{Type: cty.String},
				IsOptional:   true,
				Description:  lang.Markdown("Replaced by `client_secret`"),
				IsDeprecated: true,
			},

			"arm_subscription_id": {
				Constraint:   schema.LiteralType{Type: cty.String},
				IsOptional:   true,
				Description:  lang.Markdown("Replaced by `subscription_id`"),
				IsDeprecated: true,
			},

			"arm_tenant_id": {
				Constraint:   schema.LiteralType{Type: cty.String},
				IsOptional:   true,
				Description:  lang.Markdown("Replaced by `tenant_id`"),
				IsDeprecated: true,
			},
		},
	}

	if v.GreaterThanOrEqual(v0_13_0) {
		// https://github.com/hashicorp/terraform/commit/0f85b283
		bodySchema.Attributes["snapshot"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.Bool},
			IsOptional:  true,
			Description: lang.Markdown("Enable/Disable automatic blob snapshotting"),
		}
	}

	if v.GreaterThanOrEqual(v0_13_1) {
		// https://github.com/hashicorp/terraform/commit/0d34e5d9
		bodySchema.Attributes["client_certificate_password"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("The password associated with the Client Certificate specified in `client_certificate_path`"),
		}
		bodySchema.Attributes["client_certificate_path"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("The path to the PFX file used as the Client Certificate when authenticating as a Service Principal"),
		}
		// https://github.com/hashicorp/terraform/commit/23b4c2db
		bodySchema.Attributes["metadata_host"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
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
			Constraint:  schema.LiteralType{Type: cty.Bool},
			IsOptional:  true,
			Description: lang.Markdown("Should Terraform use AzureAD Authentication to access the Blob?"),
		}
	}

	// This attribute was both introduced and deprecated in 1.2 🙈
	// https://github.com/hashicorp/terraform/commit/9f710558 (introduction)
	// https://github.com/hashicorp/terraform/commit/2eb9118c (deprecation)
	if v.GreaterThanOrEqual(v1_2_0) {
		bodySchema.Attributes["use_microsoft_graph"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.Bool},
			IsOptional:  true,
			Description: lang.Markdown("This field now defaults to `true` and will be removed in v1.3 of Terraform Core due to the deprecation of ADAL by Microsoft."),
		}
		if v.GreaterThanOrEqual(v1_3_0) {
			// See https://github.com/hashicorp/terraform/commit/05528e8c (removal)
			delete(bodySchema.Attributes, "use_microsoft_graph")
		}
	}

	return bodySchema
}
