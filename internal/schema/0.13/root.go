// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"

	v012_mod "github.com/hashicorp/terraform-schema/internal/schema/0.12"
)

var v0_13_4 = version.Must(version.NewVersion("0.13.4"))

func ModuleSchema(v *version.Version) *schema.BodySchema {
	bs := v012_mod.ModuleSchema(v)

	bs.Blocks["module"] = moduleBlockSchema()
	bs.Blocks["provider"] = providerBlockSchema()
	bs.Blocks["terraform"] = terraformBlockSchema(v)
	bs.Blocks["resource"].Body.Blocks["provisioner"].DependentBody = ProvisionerDependentBodies(v)

	return bs
}
