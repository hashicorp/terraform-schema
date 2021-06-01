package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"

	v014_mod "github.com/hashicorp/terraform-schema/internal/schema/0.14"
)

// See https://github.com/hashicorp/terraform/tree/v0.15.0/builtin/provisioners

func ProvisionerDependentBodies(v *version.Version) map[schema.SchemaKey]*schema.BodySchema {
	return map[schema.SchemaKey]*schema.BodySchema{
		labelKey("file"):        v014_mod.FileProvisioner,
		labelKey("local-exec"):  v014_mod.LocalExecProvisioner,
		labelKey("remote-exec"): v014_mod.RemoteExecProvisioner,
	}
}

func labelKey(value string) schema.SchemaKey {
	return schema.NewSchemaKey(schema.DependencyKeys{
		Labels: []schema.LabelDependent{{Index: 0, Value: value}},
	})
}
