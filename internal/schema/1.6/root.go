// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"

	v1_5_mod "github.com/hashicorp/terraform-schema/internal/schema/1.5"
)

func ModuleSchema(v *version.Version) *schema.BodySchema {
	bs := v1_5_mod.ModuleSchema(v)
	bs.Blocks["import"] = importBlock()
	bs.Blocks["terraform"] = patchTerraformBlockSchema(bs.Blocks["terraform"])
	return bs
}
