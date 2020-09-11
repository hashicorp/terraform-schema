package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	universal "github.com/hashicorp/terraform-schema/internal/schema/universal"
)

// CoreModuleSchemaForVersion finds a module schema which is relevant
// for the given Terraform version.
// It will return error if such schema cannot be found.
func CoreModuleSchemaForVersion(version *version.Version) (*schema.BodySchema, error) {
	// TODO: Implement version matching once we actually have schemas
	return universal.Module, nil
}

// UniversalCoreModuleSchema returns a minimal universal module schema
// which is valid for any v0.12+ version of Terraform
func UniversalCoreModuleSchema() *schema.BodySchema {
	return universal.Module
}
