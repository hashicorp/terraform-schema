// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	tfjson "github.com/hashicorp/terraform-json"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/hashicorp/terraform-schema/earlydecoder"
	"github.com/hashicorp/terraform-schema/internal/addr"
	"github.com/hashicorp/terraform-schema/internal/schema/tokmod"
	"github.com/hashicorp/terraform-schema/module"
	"github.com/hashicorp/terraform-schema/registry"
	"github.com/zclconf/go-cty-debug/ctydebug"
	"github.com/zclconf/go-cty/cty"
)

var (
	v0_12_0 = version.Must(version.NewVersion("0.12.0"))
	v0_13_0 = version.Must(version.NewVersion("0.13.0"))
	v0_15_0 = version.Must(version.NewVersion("0.15.0"))
)

func TestSchemaMerger_SchemaForModule_noCoreSchema(t *testing.T) {
	sm := NewSchemaMerger(nil)

	_, err := sm.SchemaForModule(nil)
	if err == nil {
		t.Fatal("expected error for nil core schema")
	}

	if !errors.Is(err, coreSchemaRequiredErr{}) {
		t.Fatalf("unexpected error: %#v", err)
	}
}

func TestSchemaMerger_SchemaForModule_noProviderSchema(t *testing.T) {
	testCoreSchema := &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"provider": {
				Labels: []*schema.LabelSchema{
					{Name: "name"},
				},
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"alias": {Constraint: schema.LiteralType{Type: cty.String}, IsOptional: true},
					},
				},
			},
			"resource": {
				Labels: []*schema.LabelSchema{
					{Name: "type"},
					{Name: "name"},
				},
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"count": {Constraint: schema.LiteralType{Type: cty.Number}, IsOptional: true},
					},
				},
			},
			"data": {
				Labels: []*schema.LabelSchema{
					{Name: "type"},
					{Name: "name"},
				},
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"count": {Constraint: schema.LiteralType{Type: cty.Number}, IsOptional: true},
					},
				},
			},
			"module": {
				Labels: []*schema.LabelSchema{
					{Name: "name"},
				},
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"source": {
							Constraint: schema.LiteralType{Type: cty.String},
							IsRequired: true,
							IsDepKey:   true,
						},
						"version": {
							Constraint: schema.LiteralType{Type: cty.String},
							IsOptional: true,
						},
					},
				},
			},
		},
	}
	sm := NewSchemaMerger(testCoreSchema)

	_, err := sm.SchemaForModule(&module.Meta{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSchemaMerger_SchemaForModule_providerNameMatch(t *testing.T) {
	testCoreSchema := &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"provider": {
				Labels: []*schema.LabelSchema{
					{Name: "name"},
				},
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"alias": {Constraint: schema.LiteralType{Type: cty.String}, IsOptional: true},
					},
				},
			},
			"resource": {
				Labels: []*schema.LabelSchema{
					{Name: "type"},
					{Name: "name"},
				},
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"count": {Constraint: schema.LiteralType{Type: cty.Number}, IsOptional: true},
					},
				},
			},
			"data": {
				Labels: []*schema.LabelSchema{
					{Name: "type"},
					{Name: "name"},
				},
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"count": {Constraint: schema.LiteralType{Type: cty.Number}, IsOptional: true},
					},
				},
			},
			"module": {
				Labels: []*schema.LabelSchema{
					{Name: "name"},
				},
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"source": {
							Constraint: schema.LiteralType{Type: cty.String},
							IsRequired: true,
							IsDepKey:   true,
						},
						"version": {
							Constraint: schema.LiteralType{Type: cty.String},
							IsOptional: true,
						},
					},
				},
			},
		},
	}
	sm := NewSchemaMerger(testCoreSchema)
	sm.SetSchemaReader(&testJsonSchemaReader{
		ps: &tfjson.ProviderSchemas{
			FormatVersion: "1.0",
			Schemas: map[string]*tfjson.ProviderSchema{
				"registry.terraform.io/hashicorp/data": {
					ConfigSchema: &tfjson.Schema{},
					DataSourceSchemas: map[string]*tfjson.Schema{
						"data": {
							Block: &tfjson.SchemaBlock{
								Attributes: map[string]*tfjson.SchemaAttribute{
									"foobar": {
										AttributeType: cty.Bool,
										Optional:      true,
									},
								},
							},
						},
					},
				},
			},
		},
	})

	givenBodySchema, err := sm.SchemaForModule(&module.Meta{
		ProviderReferences: map[module.ProviderRef]tfaddr.Provider{
			{LocalName: "data"}: addr.NewDefaultProvider("data"),
		},
		ProviderRequirements: module.ProviderRequirements{
			addr.NewDefaultProvider("data"): version.MustConstraints(version.NewConstraint("1.0")),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	expectedBodySchema := &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"provider": {
				Labels: []*schema.LabelSchema{
					{Name: "name"},
				},
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"alias": {Constraint: schema.LiteralType{Type: cty.String}, IsOptional: true},
					},
				},
				DependentBody: map[schema.SchemaKey]*schema.BodySchema{
					`{"labels":[{"index":0,"value":"data"}]}`: {
						Blocks:     map[string]*schema.BlockSchema{},
						Attributes: map[string]*schema.AttributeSchema{},
						Detail:     "hashicorp/data",
						DocsLink: &schema.DocsLink{
							URL:     "https://registry.terraform.io/providers/hashicorp/data/latest/docs",
							Tooltip: "hashicorp/data Documentation",
						},
						HoverURL: "https://registry.terraform.io/providers/hashicorp/data/latest/docs",
					},
				},
			},
			"resource": {
				Labels: []*schema.LabelSchema{
					{Name: "type"},
					{Name: "name"},
				},
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"count": {Constraint: schema.LiteralType{Type: cty.Number}, IsOptional: true},
					},
				},
				DependentBody: map[schema.SchemaKey]*schema.BodySchema{},
			},
			"data": {
				Labels: []*schema.LabelSchema{
					{Name: "type"},
					{Name: "name"},
				},
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"count": {Constraint: schema.LiteralType{Type: cty.Number}, IsOptional: true},
					},
				},
				DependentBody: map[schema.SchemaKey]*schema.BodySchema{
					`{"labels":[{"index":0,"value":"data"}],"attrs":[{"name":"provider","expr":{"addr":"data"}}]}`: {
						Blocks: map[string]*schema.BlockSchema{},
						Attributes: map[string]*schema.AttributeSchema{
							"foobar": {
								IsOptional: true,
								Constraint: schema.AnyExpression{OfType: cty.Bool},
							},
						},
						Detail: "hashicorp/data",
					},
					`{"labels":[{"index":0,"value":"data"}]}`: {
						Blocks: map[string]*schema.BlockSchema{},
						Attributes: map[string]*schema.AttributeSchema{
							"foobar": {
								IsOptional: true,
								Constraint: schema.AnyExpression{OfType: cty.Bool},
							},
						},
						Detail: "hashicorp/data",
					},
				},
			},
			"module": {
				Labels: []*schema.LabelSchema{
					{Name: "name"},
				},
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"source": {
							Constraint: schema.LiteralType{Type: cty.String},
							IsRequired: true,
							IsDepKey:   true,
						},
						"version": {
							Constraint: schema.LiteralType{Type: cty.String},
							IsOptional: true,
						},
					},
				},
				DependentBody: map[schema.SchemaKey]*schema.BodySchema{},
			},
		},
	}

	if diff := cmp.Diff(expectedBodySchema, givenBodySchema, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("schema mismatch: %s", diff)
	}
}

func TestSchemaMerger_SchemaForModule_twiceMerged(t *testing.T) {
	testCoreSchema := &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"provider": {
				Labels: []*schema.LabelSchema{
					{
						Name: "name",
						SemanticTokenModifiers: lang.SemanticTokenModifiers{
							tokmod.Name,
							lang.TokenModifierDependent,
						},
					},
				},
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Provider},
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"alias": {
							Constraint: schema.LiteralType{Type: cty.String},
							IsOptional: true,
						},
					},
				},
			},
			"resource": {
				Labels: []*schema.LabelSchema{
					{
						Name: "type",
						SemanticTokenModifiers: lang.SemanticTokenModifiers{
							tokmod.Type,
							lang.TokenModifierDependent,
						},
					},
					{
						Name:                   "name",
						SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
					},
				},
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Resource},
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"count": {
							Constraint: schema.AnyExpression{OfType: cty.Number},
							IsOptional: true,
						},
					},
				},
			},
			"data": {
				Labels: []*schema.LabelSchema{
					{
						Name: "type",
						SemanticTokenModifiers: lang.SemanticTokenModifiers{
							tokmod.Type,
							lang.TokenModifierDependent,
						},
					},
					{
						Name:                   "name",
						SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
					},
				},
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Data},
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"count": {
							Constraint: schema.AnyExpression{OfType: cty.Number},
							IsOptional: true,
						},
					},
				},
			},
			"module": {
				Labels: []*schema.LabelSchema{
					{Name: "name", SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name}},
				},
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Module},
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"source": {
							Constraint:             schema.LiteralType{Type: cty.String},
							IsRequired:             true,
							IsDepKey:               true,
							SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent},
						},
						"version": {
							Constraint: schema.LiteralType{Type: cty.String},
							IsOptional: true,
						},
					},
				},
			},
		},
	}
	sm := NewSchemaMerger(testCoreSchema)
	sr := testSchemaReader(t, filepath.Join("testdata", "provider-schemas-0.15.json"), false)
	sm.SetSchemaReader(sr)

	vc := version.MustConstraints(version.NewConstraint("0.0.0"))

	mergedSchema, err := sm.SchemaForModule(&module.Meta{
		Path: "testdata",
		ProviderReferences: map[module.ProviderRef]tfaddr.Provider{
			{LocalName: "hashicup"}: addr.NewDefaultProvider("hashicup"),
		},
		ProviderRequirements: map[tfaddr.Provider]version.Constraints{
			addr.NewDefaultProvider("hashicup"): vc,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expectedMergedSchema_v015, mergedSchema, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("schema differs: %s", diff)
	}

	newMergedSchema, err := sm.SchemaForModule(&module.Meta{
		Path: "testdata",
		ProviderReferences: map[module.ProviderRef]tfaddr.Provider{
			{LocalName: "hcc"}: addr.NewDefaultProvider("hashicup"),
		},
		ProviderRequirements: map[tfaddr.Provider]version.Constraints{
			addr.NewDefaultProvider("hashicup"): vc,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expectedMergedSchema_v015_aliased, newMergedSchema, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("schema differs: %s", diff)
	}
}

func TestMergeWithJsonProviderSchemas_v012(t *testing.T) {
	sm := NewSchemaMerger(testCoreSchema())
	sr := testSchemaReader(t, filepath.Join("testdata", "provider-schemas-0.12.json"), true)
	sm.SetSchemaReader(sr)
	sm.SetTerraformVersion(v0_12_0)
	meta := testModuleMeta(t, "testdata/test-config-0.12.tf")
	mergedSchema, err := sm.SchemaForModule(meta)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expectedMergedSchema_v012, mergedSchema, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("schema differs: %s", diff)
	}
}

func TestMergeWithJsonProviderSchemas_v013(t *testing.T) {
	sm := NewSchemaMerger(testCoreSchema())
	sr := testSchemaReader(t, filepath.Join("testdata", "provider-schemas-0.13.json"), false)
	sm.SetSchemaReader(sr)
	sm.SetTerraformVersion(v0_13_0)
	meta := testModuleMeta(t, "testdata/test-config-0.13.tf")
	mergedSchema, err := sm.SchemaForModule(meta)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expectedMergedSchema_v013, mergedSchema, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("schema differs: %s", diff)
	}
}

func TestMergeWithJsonProviderSchemas_concurrencyBug(t *testing.T) {
	jsonSchema := &tfjson.ProviderSchemas{}
	b, err := ioutil.ReadFile(filepath.Join("testdata", "provider-schema-terraform.json"))
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(b, jsonSchema)
	if err != nil {
		t.Fatal(err)
	}

	pAddr := tfaddr.NewProvider(tfaddr.BuiltInProviderHost, tfaddr.BuiltInProviderNamespace, "terraform")
	ps := ProviderSchemaFromJson(jsonSchema.Schemas[pAddr.String()], pAddr)
	meta := &module.Meta{
		Path:      "testdir",
		Filenames: []string{"test.tf"},
		ProviderReferences: map[module.ProviderRef]tfaddr.Provider{
			{LocalName: "terraform"}: pAddr,
		},
		ProviderRequirements: module.ProviderRequirements{
			pAddr: version.Constraints{},
		},
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func(t *testing.T) {
		defer wg.Done()
		sm := NewSchemaMerger(testCoreSchema())
		sm.SetSchemaReader(exactSchemaReader{ps: ps})
		sm.SetTerraformVersion(v0_15_0)
		_, err := sm.SchemaForModule(meta)
		if err != nil {
			t.Error(err)
		}
	}(t)
	go func(t *testing.T) {
		defer wg.Done()
		sm := NewSchemaMerger(testCoreSchema())
		sm.SetSchemaReader(exactSchemaReader{ps: ps})
		sm.SetTerraformVersion(v0_15_0)
		_, err := sm.SchemaForModule(meta)
		if err != nil {
			t.Error(err)
		}
	}(t)
	wg.Wait()
}

type exactSchemaReader struct {
	ps *ProviderSchema
}

func (sr exactSchemaReader) ProviderSchema(modPath string, addr tfaddr.Provider, vc version.Constraints) (*ProviderSchema, error) {
	return sr.ps, nil
}

func TestMergeWithJsonProviderSchemas_v015(t *testing.T) {
	sm := NewSchemaMerger(testCoreSchema())
	sr := testSchemaReader(t, filepath.Join("testdata", "provider-schemas-0.15.json"), false)
	sm.SetSchemaReader(sr)
	sm.SetTerraformVersion(v0_15_0)
	meta := testModuleMeta(t, "testdata/test-config-0.15.tf")
	mergedSchema, err := sm.SchemaForModule(meta)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expectedMergedSchema_v015, mergedSchema, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("schema differs: %s", diff)
	}
}

func TestMergeWithJsonProviderSchemasAndModuleVariables_v015(t *testing.T) {
	sm := NewSchemaMerger(testCoreSchema())
	sr := testSchemaReader(t, filepath.Join("testdata", "provider-schemas-0.15.json"), false)
	sm.SetSchemaReader(sr)
	sm.SetModuleReader(testModuleReader())
	sm.SetTerraformVersion(v0_15_0)
	meta := testModuleMeta(t, "testdata/test-config-0.15.tf")
	mergedSchema, err := sm.SchemaForModule(meta)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expectedMergedSchemaWithModule_v015, mergedSchema, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("schema differs: %s", diff)
	}
}

func TestMergeWithJsonProviderSchemasAndModuleVariables_registryModule(t *testing.T) {
	sm := NewSchemaMerger(testCoreSchema())
	sm.SetModuleReader(testRegistryModuleReader())
	sm.SetTerraformVersion(v0_15_0)
	meta := testModuleMeta(t, "testdata/test-config-remote-module.tf")
	mergedSchema, err := sm.SchemaForModule(meta)
	if err != nil {
		t.Fatal(err)
	}

	moduleSchema := mergedSchema.Blocks["module"]

	if diff := cmp.Diff(expectedRemoteModuleSchema, moduleSchema, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("schema differs: %s", diff)
	}
}

func testModuleMeta(t *testing.T, path string) *module.Meta {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	filename := filepath.Base(path)

	f, diags := hclsyntax.ParseConfig(b, filename, hcl.InitialPos)
	if len(diags) > 0 {
		t.Fatal(diags)
	}
	meta, diags := earlydecoder.LoadModule("testdata", map[string]*hcl.File{
		filename: f,
	})
	if diags.HasErrors() {
		t.Fatal(diags)
	}
	return meta
}

func testSchemaReader(t *testing.T, jsonPath string, legacyStyle bool) SchemaReader {
	ps := &tfjson.ProviderSchemas{}
	b, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(b, ps)
	if err != nil {
		t.Fatal(err)
	}

	if legacyStyle {
		return &testJsonSchemaReader{
			ps:          ps,
			useTypeOnly: true,
			migrations: map[tfaddr.Provider]tfaddr.Provider{
				addr.NewLegacyProvider("null"):      addr.NewDefaultProvider("null"),
				addr.NewLegacyProvider("random"):    addr.NewDefaultProvider("random"),
				addr.NewLegacyProvider("terraform"): addr.NewBuiltInProvider("terraform"),
			},
		}
	}
	return &testJsonSchemaReader{
		ps: ps,
		migrations: map[tfaddr.Provider]tfaddr.Provider{
			// the builtin provider doesn't have entry in required_providers
			addr.NewLegacyProvider("terraform"): addr.NewBuiltInProvider("terraform"),
		},
	}
}

type testJsonSchemaReader struct {
	ps          *tfjson.ProviderSchemas
	useTypeOnly bool
	migrations  map[tfaddr.Provider]tfaddr.Provider
}

func testModuleReader() ModuleReader {
	return &testModuleReaderStruct{}
}

type testModuleReaderStruct struct {
}

func (m *testModuleReaderStruct) RegistryModuleMeta(addr tfaddr.Module, cons version.Constraints) (*registry.ModuleData, error) {
	return nil, nil
}

func (m *testModuleReaderStruct) ModuleCalls(modPath string) (module.ModuleCalls, error) {
	return module.ModuleCalls{
		Declared: map[string]module.DeclaredModuleCall{
			"example": {
				LocalName:  "example",
				SourceAddr: module.LocalSourceAddr("./source"),
			},
		},
	}, nil
}

func (m *testModuleReaderStruct) LocalModuleMeta(modPath string) (*module.Meta, error) {
	if modPath == filepath.Join("testdata", "source") {
		return &module.Meta{
			Path: "path",
			Variables: map[string]module.Variable{
				"test": {
					Type:        cty.String,
					Description: "test var",
				},
			},
		}, nil
	}
	return nil, fmt.Errorf("invalid source")
}

func testRegistryModuleReader() ModuleReader {
	return &testRegistryModuleReaderStruct{}
}

type testRegistryModuleReaderStruct struct {
}

func (m *testRegistryModuleReaderStruct) RegistryModuleMeta(addr tfaddr.Module, cons version.Constraints) (*registry.ModuleData, error) {
	return nil, nil
}

func (m *testRegistryModuleReaderStruct) ModuleCalls(modPath string) (module.ModuleCalls, error) {
	return module.ModuleCalls{
		Installed: map[string]module.InstalledModuleCall{
			"remote-example": {
				LocalName:  "remote-example",
				SourceAddr: tfaddr.MustParseModuleSource("registry.terraform.io/namespace/foo/bar"),
				Path:       ".terraform/modules/remote-example",
			},
		},
	}, nil
}

func (m *testRegistryModuleReaderStruct) LocalModuleMeta(modPath string) (*module.Meta, error) {
	if modPath == ".terraform/modules/remote-example" {
		return &module.Meta{
			Path: ".terraform/modules/remote-example",
			Variables: map[string]module.Variable{
				"test": {
					Type:        cty.String,
					Description: "test var",
				},
			},
		}, nil
	}
	return nil, fmt.Errorf("invalid source")
}

func (r *testJsonSchemaReader) ProviderSchema(_ string, pAddr tfaddr.Provider, _ version.Constraints) (*ProviderSchema, error) {
	if newAddr, ok := r.migrations[pAddr]; ok {
		pAddr = newAddr
	}

	addr := pAddr.String()
	if r.useTypeOnly {
		addr = pAddr.Type
	}

	jsonSchema, ok := r.ps.Schemas[addr]
	if !ok {
		return nil, fmt.Errorf("%s: schema not found", pAddr.String())
	}

	return ProviderSchemaFromJson(jsonSchema, pAddr), nil
}

func testCoreSchema() *schema.BodySchema {
	return &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"provider": {
				Labels: []*schema.LabelSchema{
					{
						Name:                   "name",
						SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name, lang.TokenModifierDependent},
					},
				},
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Provider},
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"alias": {Constraint: schema.LiteralType{Type: cty.String}, IsOptional: true},
					},
				},
			},
			"resource": {
				Labels: []*schema.LabelSchema{
					{
						Name:                   "type",
						SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
					},
					{
						Name:                   "name",
						SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
					},
				},
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Resource},
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"count": {
							Constraint: schema.AnyExpression{OfType: cty.Number},
							IsOptional: true,
						},
					},
				},
			},
			"data": {
				Labels: []*schema.LabelSchema{
					{
						Name:                   "type",
						SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
					},
					{
						Name:                   "name",
						SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
					},
				},
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Data},
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"count": {
							Constraint: schema.AnyExpression{OfType: cty.Number},
							IsOptional: true,
						},
					},
				},
			},
			"module": {
				Labels: []*schema.LabelSchema{
					{
						Name:                   "name",
						SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
					},
				},
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Module},
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"source": {
							Constraint:             schema.LiteralType{Type: cty.String},
							IsRequired:             true,
							IsDepKey:               true,
							SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent},
						},
						"version": {
							Constraint: schema.LiteralType{Type: cty.String},
							IsOptional: true,
						},
					},
				},
			},
		},
	}
}
