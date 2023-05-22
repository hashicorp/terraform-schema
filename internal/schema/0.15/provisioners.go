// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"

	v014_mod "github.com/hashicorp/terraform-schema/internal/schema/0.14"
)

var (
	FileProvisioner       = v014_mod.FileProvisioner
	LocalExecProvisioner  = v014_mod.LocalExecProvisioner
	RemoteExecProvisioner = v014_mod.RemoteExecProvisioner
)

// See https://github.com/hashicorp/terraform/tree/v0.15.0/builtin/provisioners

func ProvisionerDependentBodies(v *version.Version) map[schema.SchemaKey]*schema.BodySchema {
	return map[schema.SchemaKey]*schema.BodySchema{
		labelKey("file"):        FileProvisioner,
		labelKey("local-exec"):  LocalExecProvisioner,
		labelKey("remote-exec"): RemoteExecProvisioner,
	}
}

func labelKey(value string) schema.SchemaKey {
	return schema.NewSchemaKey(schema.DependencyKeys{
		Labels: []schema.LabelDependent{{Index: 0, Value: value}},
	})
}
