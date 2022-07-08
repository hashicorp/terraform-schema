package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
)

var outputLifecycleBlock = &schema.BlockSchema{
	Description: lang.Markdown("Lifecycle customizations, to set a validity condition of the output"),
	Body: &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"precondition": {
				Body: conditionBody,
			},
		},
	},
}
