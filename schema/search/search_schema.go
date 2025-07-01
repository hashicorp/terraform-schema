// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	search_1_14 "github.com/hashicorp/terraform-schema/internal/schema/search/1.14"
)

// CoreSearchSchemaForVersion finds a schema for search configuration files
// that is relevant for the given Terraform version.
// It will return an error if such schema cannot be found.
func CoreSearchSchemaForVersion(v *version.Version) (*schema.BodySchema, error) {
	return search_1_14.SearchSchema(v), nil
}
