// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"

	v1_9_test "github.com/hashicorp/terraform-schema/internal/schema/tests/1.9"
)

// TestSchema returns the static schema for a test
// configuration (*.tftest.hcl) file.
func TestSchema(v *version.Version) *schema.BodySchema {
	bs := v1_9_test.TestSchema(v)

	bs.Blocks["run"] = patchRunBlockSchema(bs.Blocks["run"])
	bs.Blocks["mock_provider"] = patchMockProviderBlockSchema(bs.Blocks["mock_provider"])
	bs.Blocks["override_module"] = patchOverrideModuleBlockSchema(bs.Blocks["override_module"])
	bs.Blocks["override_resource"] = patchOverrideResourceBlockSchema(bs.Blocks["override_resource"])
	bs.Blocks["override_data"] = patchOverrideDataBlockSchema(bs.Blocks["override_data"])

	return bs
}
