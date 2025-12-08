// Copyright IBM Corp. 2020, 2025
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"

	v012_mod "github.com/hashicorp/terraform-schema/internal/schema/0.12"
	v1_3_mod "github.com/hashicorp/terraform-schema/internal/schema/1.3"
	v1_4_mod "github.com/hashicorp/terraform-schema/internal/schema/1.4"
	v1_8_mod "github.com/hashicorp/terraform-schema/internal/schema/1.8"
)

func ModuleSchema(v *version.Version) *schema.BodySchema {
	bs := v1_8_mod.ModuleSchema(v)

	bs.Blocks["removed"].Body.Blocks["provisioner"] = v012_mod.ProvisionerBlock(v)
	bs.Blocks["removed"].Body.Blocks["provisioner"].DependentBody = v1_4_mod.ProvisionerDependentBodies(v)
	bs.Blocks["removed"].Body.Blocks["provisioner"].Body.Blocks["connection"].DependentBody = v1_3_mod.ConnectionDependentBodies(v)
	bs.Blocks["removed"].Body.Blocks["provisioner"].Body.Attributes["when"] = &schema.AttributeSchema{
		Constraint: schema.OneOf{
			schema.Keyword{
				Keyword:     "destroy",
				Description: lang.Markdown("Run the provisioner when the resource is destroyed"),
			},
		},
		IsOptional:  true,
		Description: lang.Markdown("When to run the provisioner - `removed` resources can only be destroyed."),
	}

	bs.Blocks["removed"].Body.Blocks["connection"] = v012_mod.ConnectionBlock(v)
	bs.Blocks["removed"].Body.Blocks["connection"].DependentBody = v1_3_mod.ConnectionDependentBodies(v)

	return bs
}
