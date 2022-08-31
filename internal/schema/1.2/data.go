package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
)

func datasourceLifecycleBlock() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.Markdown("Lifecycle customizations to set validity conditions of the datasource"),
		Body: &schema.BodySchema{
			Blocks: map[string]*schema.BlockSchema{
				"precondition": {
					Body: conditionBody(),
				},
				"postcondition": {
					Body: conditionBody(),
				},
			},
		},
	}
}
