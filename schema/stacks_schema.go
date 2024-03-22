package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	stack_1_8 "github.com/hashicorp/terraform-schema/internal/schema/stacks/1.8"
)

func CoreStackSchema(v *version.Version) *schema.BodySchema {
	return stack_1_8.StackSchema(v)
}

func CoreDeploySchema(v *version.Version) *schema.BodySchema {
	return stack_1_8.DeploymentSchema(v)
}
