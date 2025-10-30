// Copyright (c) HashiCorp, Inc.
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

	bs.Blocks["run"].Body.Blocks["state_key"] = stateKeyBlockSchema()
	bs.Blocks["mock_provider"].Body.Blocks["override_during"] = overrideDuringBlockSchema()

	return bs
}
