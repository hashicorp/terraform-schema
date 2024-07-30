// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"

	v1_7_test "github.com/hashicorp/terraform-schema/internal/schema/tests/1.7"
)

// TestSchema returns the static schema for a test
// configuration (*.tftest.hcl) file.
func TestSchema(v *version.Version) *schema.BodySchema {
	bs := v1_7_test.TestSchema(v)

	// Removes the version attribute
	bs.Blocks["provider"].Body.Attributes = map[string]*schema.AttributeSchema{
		"alias": {
			Constraint:  schema.LiteralType{Type: cty.String},
			IsOptional:  true,
			Description: lang.Markdown("Alias for using the same provider with different configurations for different resources, e.g. `eu-west`"),
		},
	}

	return bs
}
