package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
)

func ModuleSchema(v *version.Version) *schema.BodySchema {
	return &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"data":      datasourceBlockSchema,
			"locals":    localsBlockSchema,
			"module":    moduleBlockSchema,
			"output":    outputBlockSchema,
			"provider":  providerBlockSchema,
			"resource":  resourceBlockSchema,
			"variable":  variableBlockSchema,
			"terraform": terraformBlockSchema,
		},
	}
}
