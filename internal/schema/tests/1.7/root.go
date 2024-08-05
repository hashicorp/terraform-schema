// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"

	v1_6_test "github.com/hashicorp/terraform-schema/internal/schema/tests/1.6"
)

// TestSchema returns the static schema for a test
// configuration (*.tftest.hcl) file.
func TestSchema(v *version.Version) *schema.BodySchema {
	bs := v1_6_test.TestSchema(v)

	bs.Blocks["mock_provider"] = mockProviderBlockSchema()
	bs.Blocks["override_resource"] = overrideResourceBlockSchema()
	bs.Blocks["override_data"] = overrideDataBlockSchema()
	bs.Blocks["override_module"] = overrideModuleBlockSchema()

	bs.Blocks["run"].Body.Blocks["override_resource"] = overrideResourceBlockSchema()
	bs.Blocks["run"].Body.Blocks["override_data"] = overrideDataBlockSchema()
	bs.Blocks["run"].Body.Blocks["override_module"] = overrideModuleBlockSchema()

	return bs
}

// MockSchema returns the static schema for a mock
// configuration (*.tfmock.hcl) file.
func MockSchema(_ *version.Version) *schema.BodySchema {
	return &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"mock_resource":     mockResourceBlockSchema(),
			"mock_data":         mockDataBlockSchema(),
			"override_resource": overrideResourceBlockSchema(),
			"override_data":     overrideDataBlockSchema(),
		},
	}
}
