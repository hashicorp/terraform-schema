// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"

	v1_15_mod "github.com/hashicorp/terraform-schema/internal/schema/1.14"
)

func ModuleSchema(v *version.Version) *schema.BodySchema {
	bs := v1_15_mod.ModuleSchema(v)

	bs.Blocks["terraform"] = patchTerraformBlockSchema(bs.Blocks["terraform"])

	return bs
}
