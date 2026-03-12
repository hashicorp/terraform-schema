// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"

	v1_9_test "github.com/hashicorp/terraform-schema/internal/schema/tests/1.9"
)

// TestSchema returns the static schema for a test
// configuration (*.tftest.hcl) file.
func TestSchema(v *version.Version) *schema.BodySchema {
	bs := v1_9_test.TestSchema(v)

	bs.Blocks["run"] = patchRunBlockSchema(bs.Blocks["run"])
	bs.Blocks["mock_provider"] = patchMockProviderBlockSchema(bs.Blocks["mock_provider"])
	bs.Blocks["mock_resource"] = patchMockResourceBlockSchema(bs.Blocks["mock_resource"])
	bs.Blocks["mock_data"] = patchMockDataBlockSchema(bs.Blocks["mock_data"])
	bs.Blocks["override_module"] = patchOverrideModuleBlockSchema(bs.Blocks["override_module"])
	bs.Blocks["override_resource"] = patchOverrideResourceBlockSchema(bs.Blocks["override_resource"])
	bs.Blocks["override_data"] = patchOverrideDataBlockSchema(bs.Blocks["override_data"])

	return bs
}

func overrideDuringAttributeSchema() *schema.AttributeSchema {
	return &schema.AttributeSchema{
		Description: lang.PlainText("Allows overriding the point in time where terraform generates data"),
		IsOptional:  true,
		Constraint: schema.OneOf{
			schema.Keyword{
				Keyword:     "apply",
				Description: lang.Markdown("Default behavior where data is generated during the apply operation and (known after apply) is returned during the plan"),
			},
			schema.Keyword{
				Keyword:     "plan",
				Description: lang.Markdown("Allows to generate data during the plan operation. The same data will be used during the apply"),
			},
		},
	}
}
