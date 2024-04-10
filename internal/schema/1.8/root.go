// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"

	v1_7_mod "github.com/hashicorp/terraform-schema/internal/schema/1.7"
)

func ModuleSchema(v *version.Version) *schema.BodySchema {
	bs := v1_7_mod.ModuleSchema(v)
	bs.Blocks["terraform"] = patchTerraformBlockSchema(bs.Blocks["terraform"])
	return bs
}
