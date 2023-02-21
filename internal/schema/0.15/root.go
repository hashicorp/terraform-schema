// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"

	v014_mod "github.com/hashicorp/terraform-schema/internal/schema/0.14"
)

func ModuleSchema(v *version.Version) *schema.BodySchema {
	bs := v014_mod.ModuleSchema(v)
	bs.Blocks["terraform"] = patchTerraformBlockSchema(bs.Blocks["terraform"])
	bs.Blocks["resource"].Body.Blocks["provisioner"].DependentBody = ProvisionerDependentBodies(v)
	bs.Blocks["resource"].Body.Blocks["connection"].DependentBody = ConnectionDependentBodies(v)

	return bs
}
