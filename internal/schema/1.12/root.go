// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"

	v1_10_mod "github.com/hashicorp/terraform-schema/internal/schema/1.10"
)

func ModuleSchema(v *version.Version) *schema.BodySchema {
	bs := v1_10_mod.ModuleSchema(v)

	// Override the import block with the new schema that supports identity attribute
	bs.Blocks["import"] = importBlock()

	return bs
}
