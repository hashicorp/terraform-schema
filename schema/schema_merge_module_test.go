// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/hashicorp/terraform-schema/module"
	"github.com/hashicorp/terraform-schema/registry"
	"github.com/zclconf/go-cty-debug/ctydebug"
)

type uninstalledRemoteModuleReader struct{}

func (r *uninstalledRemoteModuleReader) ProviderSchema(modPath string, addr tfaddr.Provider, vc version.Constraints) (*ProviderSchema, error) {
	return nil, nil
}

func (r *uninstalledRemoteModuleReader) RegistryModuleMeta(addr tfaddr.Module, cons version.Constraints) (*registry.ModuleData, error) {
	return nil, fmt.Errorf("module not available")
}

func (r *uninstalledRemoteModuleReader) DeclaredModuleCalls(modPath string) (map[string]module.DeclaredModuleCall, error) {
	rawSource := "git::https://github.com/example/module.git"
	return map[string]module.DeclaredModuleCall{
		"app": {
			LocalName:     "app",
			RawSourceAddr: rawSource,
			SourceAddr:    module.RemoteSourceAddr("git::https://github.com/example/module.git"),
		},
	}, nil
}

func (r *uninstalledRemoteModuleReader) InstalledModulePath(rootPath string, normalizedSource string) (string, bool) {
	return "", false
}

func (r *uninstalledRemoteModuleReader) LocalModuleMeta(modPath string) (*module.Meta, error) {
	return nil, fmt.Errorf("module not available")
}

func TestSchemaMerger_resolveModuleDependentBody_uninstalledRemoteModule(t *testing.T) {
	sm := NewSchemaMerger(testCoreSchema())
	sm.SetStateReader(&uninstalledRemoteModuleReader{})
	sm.SetTerraformVersion(v0_15_0)

	meta := &module.Meta{
		Path:      "/test",
		Filenames: []string{"main.tf"},
	}
	mergedSchema, err := sm.SchemaForModule(meta)
	if err != nil {
		t.Fatal(err)
	}

	moduleSchema := mergedSchema.Blocks["module"]
	if len(moduleSchema.DependentBody) != 1 {
		t.Fatalf("expected 1 dependent body schema, got %d", len(moduleSchema.DependentBody))
	}

	var depSchema *schema.BodySchema
	for _, bodySchema := range moduleSchema.DependentBody {
		depSchema = bodySchema
		break
	}

	expected := schemaForUninstalledModuleBlock(module.DeclaredModuleCall{
		LocalName: "app",
	})
	if diff := cmp.Diff(expected, depSchema, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("schema mismatch: %s", diff)
	}
}
