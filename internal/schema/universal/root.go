package schema

import (
	"github.com/hashicorp/hcl-lang/schema"
)

var Module = &schema.BodySchema{
	Blocks: map[string]*schema.BlockSchema{
		"provider": providerBlockSchema,
		"resource": resourceBlockSchema,
		"data":     datasourceBlockSchema,
		"locals":   localsBlockSchema,
		"output":   outputBlockSchema,
		"variable": variableBlockSchema,
		"module":   moduleBlockSchema,
	},
}
