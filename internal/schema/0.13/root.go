package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"

	v012_mod "github.com/hashicorp/terraform-schema/internal/schema/0.12"
)

func ModuleSchema(v *version.Version) *schema.BodySchema {
	bs := v012_mod.ModuleSchema(v)

	bs.Blocks["module"] = moduleBlockSchema
	bs.Blocks["provider"] = providerBlockSchema
	bs.Blocks["terraform"] = terraformBlockSchema

	return bs
}
