package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
)

var localsBlockSchema = &schema.BlockSchema{
	Description: lang.Markdown("Local values assigning names to expressions, so you can use these multiple times without repetition\n" +
		"e.g. `service_name = \"forum\"`"),
}
