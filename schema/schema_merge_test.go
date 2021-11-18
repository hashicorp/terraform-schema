package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	tfjson "github.com/hashicorp/terraform-json"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/hashicorp/terraform-schema/earlydecoder"
	"github.com/hashicorp/terraform-schema/module"
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
						"alias": {Expr: schema.LiteralTypeOnly(cty.String), IsOptional: true},
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
						"count": {Expr: schema.LiteralTypeOnly(cty.Number), IsOptional: true},
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
						"count": {Expr: schema.LiteralTypeOnly(cty.Number), IsOptional: true},
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
							Expr:       schema.LiteralTypeOnly(cty.String),
							IsRequired: true,
							IsDepKey:   true,
						},
						"version": {
							Expr:       schema.LiteralTypeOnly(cty.String),
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

func TestSchemaMerger_SchemaForModule_twiceMerged(t *testing.T) {
	testCoreSchema := &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"provider": {
				Labels: []*schema.LabelSchema{
					{Name: "name"},
				},
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"alias": {
							Expr:       schema.LiteralTypeOnly(cty.String),
							IsOptional: true,
						},
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
						"count": {
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
							IsOptional: true,
						},
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
						"count": {
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
							IsOptional: true,
						},
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
							Expr:       schema.LiteralTypeOnly(cty.String),
							IsRequired: true,
							IsDepKey:   true,
						},
						"version": {
							Expr:       schema.LiteralTypeOnly(cty.String),
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

	vc, err := version.NewConstraint("0.0.0")
	if err != nil {
		t.Fatal(err)
	}
	mergedSchema, err := sm.SchemaForModule(&module.Meta{
		Path: "testdata",
		ProviderReferences: map[module.ProviderRef]tfaddr.Provider{
			{LocalName: "hashicup"}: tfaddr.NewDefaultProvider("hashicup"),
		},
		ProviderRequirements: map[tfaddr.Provider]version.Constraints{
			tfaddr.NewDefaultProvider("hashicup"): vc,
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
			{LocalName: "hcc"}: tfaddr.NewDefaultProvider("hashicup"),
		},
		ProviderRequirements: map[tfaddr.Provider]version.Constraints{
			tfaddr.NewDefaultProvider("hashicup"): vc,
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
				tfaddr.NewLegacyProvider("null"):      tfaddr.NewDefaultProvider("null"),
				tfaddr.NewLegacyProvider("random"):    tfaddr.NewDefaultProvider("random"),
				tfaddr.NewLegacyProvider("terraform"): tfaddr.NewBuiltInProvider("terraform"),
			},
		}
	}
	return &testJsonSchemaReader{
		ps: ps,
		migrations: map[tfaddr.Provider]tfaddr.Provider{
			// the builtin provider doesn't have entry in required_providers
			tfaddr.NewLegacyProvider("terraform"): tfaddr.NewBuiltInProvider("terraform"),
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

func (m *testModuleReaderStruct) ModuleCalls(modPath string) ([]module.ModuleCall, error) {
	return []module.ModuleCall{
		{
			LocalName:  "example",
			SourceAddr: "source",
			Path:       "path",
		},
	}, nil
}

func (m *testModuleReaderStruct) ModuleMeta(modPath string) (*module.Meta, error) {
	if modPath == "path" {
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
					{Name: "name"},
				},
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"alias": {Expr: schema.LiteralTypeOnly(cty.String), IsOptional: true},
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
						"count": {
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
							IsOptional: true,
						},
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
						"count": {
							Expr: schema.ExprConstraints{
								schema.TraversalExpr{OfType: cty.Number},
								schema.LiteralTypeExpr{Type: cty.Number},
							},
							IsOptional: true,
						},
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
							Expr:       schema.LiteralTypeOnly(cty.String),
							IsRequired: true,
							IsDepKey:   true,
						},
						"version": {
							Expr:       schema.LiteralTypeOnly(cty.String),
							IsOptional: true,
						},
					},
				},
			},
		},
	}
}
