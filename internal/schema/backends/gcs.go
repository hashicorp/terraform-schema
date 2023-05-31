// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package backends

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func gcsBackend(v *version.Version) *schema.BodySchema {
	// https://github.com/hashicorp/terraform/blob/v0.12.0/backend/remote-state/gcs/backend.go
	docsUrl := "https://www.terraform.io/docs/language/settings/backends/gcs.html"
	bodySchema := &schema.BodySchema{
		Description: lang.Markdown("Google Cloud Storage"),
		HoverURL:    docsUrl,
		DocsLink: &schema.DocsLink{
			URL: docsUrl,
		},
		Attributes: map[string]*schema.AttributeSchema{
			"bucket": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("The name of the Google Cloud Storage bucket"),
			},

			"path": {
				Constraint:   schema.LiteralType{Type: cty.String},
				IsOptional:   true,
				Description:  lang.Markdown("Path of the default state file;\nDEPRECATED: Use the `prefix` option instead"),
				IsDeprecated: true,
			},

			"prefix": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The directory where state files will be saved inside the bucket"),
			},

			"credentials": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Google Cloud JSON Account Key"),
			},

			"encryption_key": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("A 32 byte base64 encoded 'customer supplied encryption key' used to encrypt all state."),
			},
		},
	}

	if v.GreaterThanOrEqual(v0_12_10) {
		// https://github.com/hashicorp/terraform/commit/f6c90c1d
		bodySchema.Attributes["access_token"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("An OAuth2 token used for GCP authentication"),
		}
	}

	if v.GreaterThanOrEqual(v0_14_0) {
		// https://github.com/hashicorp/terraform/commit/c43731a0
		bodySchema.Attributes["impersonate_service_account"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("The service account to impersonate for all Google API Calls"),
		}

		bodySchema.Attributes["impersonate_service_account_delegates"] = &schema.AttributeSchema{
			Constraint: schema.List{
				Elem: schema.LiteralType{Type: cty.String},
			},
			IsOptional:  true,
			Description: lang.Markdown("The delegation chain for the impersonated service account"),
		}
	}

	if v.GreaterThanOrEqual(v0_15_0) {
		// https://github.com/hashicorp/terraform/commit/3b9c5e5b
		delete(bodySchema.Attributes, "path")
	}

	if v.GreaterThanOrEqual(v1_4_0) {
		// https://github.com/hashicorp/terraform/commit/89ef27d3
		bodySchema.Attributes["storage_custom_endpoint"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("A URL containing three parts: the protocol, the DNS name pointing to a Private Service Connect endpoint, and the path for the Cloud Storage API (`/storage/v1/b`, [see here](https://cloud.google.com/storage/docs/json_api/v1/buckets/get#http-request))."),
		}
		// https://github.com/hashicorp/terraform/commit/d43ec0f30
		bodySchema.Attributes["kms_encryption_key"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("A Cloud KMS key ('customer managed encryption key') used when reading and writing state files in the bucket. Format should be 'projects/{{project}}/locations/{{location}}/keyRings/{{keyRing}}/cryptoKeys/{{name}}'."),
		}
	}

	return bodySchema
}
