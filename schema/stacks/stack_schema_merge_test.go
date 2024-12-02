// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package schema

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/hcl/v2"
	tfjson "github.com/hashicorp/terraform-json"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	tfmod "github.com/hashicorp/terraform-schema/module"
	tfschema "github.com/hashicorp/terraform-schema/schema"
	"github.com/hashicorp/terraform-schema/stack"
	"github.com/zclconf/go-cty-debug/ctydebug"
	"github.com/zclconf/go-cty/cty"
)

func TestStackSchemaMerger_SchemaForStack_noCoreSchema(t *testing.T) {
	sm := NewStackSchemaMerger(nil)

	_, err := sm.SchemaForStack(nil)
	if err == nil {
		t.Fatal("expected error for nil core schema")
	}

	if !errors.Is(err, tfschema.CoreSchemaRequiredErr{}) {
		t.Fatalf("unexpected error: %#v", err)
	}
}

func TestStackSchemaMerger_SchemaForStack_noProviderSchema(t *testing.T) {
	testCoreSchema := &schema.BodySchema{}

	sm := NewStackSchemaMerger(testCoreSchema)

	_, err := sm.SchemaForStack(&stack.Meta{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestStackSchemaMerger_SchemaForStack_providerNameMatch(t *testing.T) {
	testCoreSchema := &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"provider": {
				Labels: []*schema.LabelSchema{
					{Name: "type"},
					{Name: "name"},
				},
				Body: &schema.BodySchema{
					Blocks: map[string]*schema.BlockSchema{
						"config": {},
					},
				},
			},
			"component": {},
		},
	}
	sm := NewStackSchemaMerger(testCoreSchema)
	sm.SetStateReader(&testStackSchemaReader{
		ps: &tfjson.ProviderSchemas{
			FormatVersion: "1.0",
			Schemas: map[string]*tfjson.ProviderSchema{
				"registry.terraform.io/hashicorp/data": {
					ConfigSchema: &tfjson.Schema{
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
	})

	givenBodySchema, err := sm.SchemaForStack(&stack.Meta{
		ProviderRequirements: map[string]stack.ProviderRequirement{
			"data": {Source: tfaddr.MustParseProviderSource("hashicorp/data"), VersionConstraints: version.MustConstraints(version.NewConstraint("1.0"))},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	expectedBodySchema := &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"provider": {
				Labels: []*schema.LabelSchema{
					{Name: "type"},
					{Name: "name"},
				},
				Body: &schema.BodySchema{
					Blocks: map[string]*schema.BlockSchema{
						"config": {},
					},
				},
				DependentBody: map[schema.SchemaKey]*schema.BodySchema{
					`{"labels":[{"index":0,"value":"data"}]}`: {
						Blocks: map[string]*schema.BlockSchema{
							"config": {
								Body: &schema.BodySchema{
									Blocks: map[string]*schema.BlockSchema{},
									Attributes: map[string]*schema.AttributeSchema{
										"foobar": {
											IsOptional: true,
											Constraint: schema.AnyExpression{OfType: cty.Bool},
										},
									},
									Detail: "hashicorp/data",
									DocsLink: &schema.DocsLink{
										URL:     "https://registry.terraform.io/providers/hashicorp/data/latest/docs",
										Tooltip: "hashicorp/data Documentation",
									},
									HoverURL: "https://registry.terraform.io/providers/hashicorp/data/latest/docs",
								},
							},
						},
						Detail: "hashicorp/data",
						DocsLink: &schema.DocsLink{
							URL:     "https://registry.terraform.io/providers/hashicorp/data/latest/docs",
							Tooltip: "hashicorp/data Documentation",
						},
						HoverURL: "https://registry.terraform.io/providers/hashicorp/data/latest/docs",
					},
				},
			},
			"component": {
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
					{Name: "type"},
					{Name: "name"},
				},
				Body: &schema.BodySchema{
					Blocks: map[string]*schema.BlockSchema{
						"config": {},
					},
				},
			},
			"component": {},
		},
	}
	sm := NewStackSchemaMerger(testCoreSchema)
	sm.SetStateReader(&testStackSchemaReader{
		ps: &tfjson.ProviderSchemas{
			FormatVersion: "1.0",
			Schemas: map[string]*tfjson.ProviderSchema{
				"registry.terraform.io/hashicorp/data": {
					ConfigSchema: &tfjson.Schema{
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
	})

	givenBodySchema, err := sm.SchemaForStack(&stack.Meta{
		ProviderRequirements: map[string]stack.ProviderRequirement{
			"data": {Source: tfaddr.MustParseProviderSource("hashicorp/data"), VersionConstraints: version.MustConstraints(version.NewConstraint("1.0"))},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	expectedBodySchema := &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"provider": {
				Labels: []*schema.LabelSchema{
					{Name: "type"},
					{Name: "name"},
				},
				Body: &schema.BodySchema{
					Blocks: map[string]*schema.BlockSchema{
						"config": {},
					},
				},
				DependentBody: map[schema.SchemaKey]*schema.BodySchema{
					`{"labels":[{"index":0,"value":"data"}]}`: {
						Blocks: map[string]*schema.BlockSchema{
							"config": {
								Body: &schema.BodySchema{
									Blocks: map[string]*schema.BlockSchema{},
									Attributes: map[string]*schema.AttributeSchema{
										"foobar": {
											IsOptional: true,
											Constraint: schema.AnyExpression{OfType: cty.Bool},
										},
									},
									Detail: "hashicorp/data",
									DocsLink: &schema.DocsLink{
										URL:     "https://registry.terraform.io/providers/hashicorp/data/latest/docs",
										Tooltip: "hashicorp/data Documentation",
									},
									HoverURL: "https://registry.terraform.io/providers/hashicorp/data/latest/docs",
								},
							},
						},
						Detail: "hashicorp/data",
						DocsLink: &schema.DocsLink{
							URL:     "https://registry.terraform.io/providers/hashicorp/data/latest/docs",
							Tooltip: "hashicorp/data Documentation",
						},
						HoverURL: "https://registry.terraform.io/providers/hashicorp/data/latest/docs",
					},
				},
			},
			"component": {
				DependentBody: map[schema.SchemaKey]*schema.BodySchema{},
			},
		},
	}

	if diff := cmp.Diff(expectedBodySchema, givenBodySchema, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("schema mismatch: %s", diff)
	}

	// now merge again with different local name
	givenBodySchema, err = sm.SchemaForStack(&stack.Meta{
		ProviderRequirements: map[string]stack.ProviderRequirement{
			"daataa": {Source: tfaddr.MustParseProviderSource("hashicorp/data"), VersionConstraints: version.MustConstraints(version.NewConstraint("1.0"))},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	expectedBodySchema = &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"provider": {
				Labels: []*schema.LabelSchema{
					{Name: "type"},
					{Name: "name"},
				},
				Body: &schema.BodySchema{
					Blocks: map[string]*schema.BlockSchema{
						"config": {},
					},
				},
				DependentBody: map[schema.SchemaKey]*schema.BodySchema{
					`{"labels":[{"index":0,"value":"daataa"}]}`: {
						Blocks: map[string]*schema.BlockSchema{
							"config": {
								Body: &schema.BodySchema{
									Blocks: map[string]*schema.BlockSchema{},
									Attributes: map[string]*schema.AttributeSchema{
										"foobar": {
											IsOptional: true,
											Constraint: schema.AnyExpression{OfType: cty.Bool},
										},
									},
									Detail: "hashicorp/data",
									DocsLink: &schema.DocsLink{
										URL:     "https://registry.terraform.io/providers/hashicorp/data/latest/docs",
										Tooltip: "hashicorp/data Documentation",
									},
									HoverURL: "https://registry.terraform.io/providers/hashicorp/data/latest/docs",
								},
							},
						},
						Detail: "hashicorp/data",
						DocsLink: &schema.DocsLink{
							URL:     "https://registry.terraform.io/providers/hashicorp/data/latest/docs",
							Tooltip: "hashicorp/data Documentation",
						},
						HoverURL: "https://registry.terraform.io/providers/hashicorp/data/latest/docs",
					},
				},
			},
			"component": {
				DependentBody: map[schema.SchemaKey]*schema.BodySchema{},
			},
		},
	}

	if diff := cmp.Diff(expectedBodySchema, givenBodySchema, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("schema mismatch: %s", diff)
	}
}

func TestStackSchemaMerger_SchemaForStack_variables(t *testing.T) {
	testCoreSchema := &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"provider":  {},
			"component": {},
			"variable":  {},
		},
	}
	sm := NewStackSchemaMerger(testCoreSchema)
	sm.SetStateReader(&testStackSchemaReader{})

	givenBodySchema, err := sm.SchemaForStack(&stack.Meta{
		Variables: map[string]stack.Variable{
			"foo": {Type: cty.String, Description: "A foo variable", IsSensitive: true, DefaultValue: cty.StringVal("bar")},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	expectedBodySchema := &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"provider": {
				DependentBody: map[schema.SchemaKey]*schema.BodySchema{},
			},
			"component": {
				DependentBody: map[schema.SchemaKey]*schema.BodySchema{},
			},
			"variable": {
				Labels: []*schema.LabelSchema{{Name: "name", IsDepKey: true, Description: lang.MarkupContent{Value: "Variable name", Kind: lang.PlainTextKind}}},
				DependentBody: map[schema.SchemaKey]*schema.BodySchema{
					`{"labels":[{"index":0,"value":"foo"}]}`: {
						Attributes: map[string]*schema.AttributeSchema{
							"default": {
								Constraint:  schema.LiteralType{Type: cty.String},
								Description: lang.MarkupContent{Value: "Default value to use when variable is not explicitly set", Kind: lang.MarkdownKind},
								IsOptional:  true,
							},
						},
					},
				},
			},
		},
	}

	if diff := cmp.Diff(expectedBodySchema, givenBodySchema, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("schema mismatch: %s", diff)
	}
}

func TestStackSchemaMerger_SchemaForStack_components(t *testing.T) {
	testCoreSchema := &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"provider":  {},
			"component": {},
		},
	}
	sm := NewStackSchemaMerger(testCoreSchema)
	sm.SetStateReader(&testStackSchemaReader{})

	givenBodySchema, err := sm.SchemaForStack(&stack.Meta{
		Components: map[string]stack.Component{
			"reg": {
				Source:     "registry/source/test",
				Version:    version.MustConstraints(version.NewConstraint("1.0")),
				SourceAddr: tfmod.ParseModuleSourceAddr("registry/source/test"),
			},
			"local": {
				Source:     "./local/path",
				SourceAddr: tfmod.ParseModuleSourceAddr("./local/path"),
			},
			"git": {
				Source:     "git::https://example.com/vpc.git",
				SourceAddr: tfmod.ParseModuleSourceAddr("git::https://example.com/vpc.git"),
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	expectedBodySchema := &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"provider": {
				DependentBody: map[schema.SchemaKey]*schema.BodySchema{},
			},
			"component": {
				DependentBody: map[schema.SchemaKey]*schema.BodySchema{
					`{"attrs":[{"name":"source","expr":{"static":"./local/path"}}]}`: {
						Attributes: map[string]*schema.AttributeSchema{
							"inputs": {Constraint: schema.Object{Attributes: schema.ObjectAttributes{}}},
							"providers": {Constraint: schema.Object{Attributes: schema.ObjectAttributes{
								"test": {Constraint: schema.Reference{OfScopeId: "provider"}},
							}}},
						},
						TargetableAs: []*schema.Targetable{
							{
								Address:           lang.Address{lang.RootStep{Name: "component"}, lang.AttrStep{Name: "local"}},
								ScopeId:           "module",
								AsType:            cty.Object(map[string]cty.Type{}),
								NestedTargetables: schema.Targetables{},
							},
						},
						Targets: &schema.Target{
							Path: lang.Path{LanguageID: "terraform"},
							Range: hcl.Range{
								Filename: "main.tf",
								Start:    hcl.Pos{Line: 1, Column: 1, Byte: 0},
								End:      hcl.Pos{Line: 1, Column: 1, Byte: 0},
							},
						},
						ImpliedOrigins: schema.ImpliedOrigins{},
					},
					`{"attrs":[{"name":"source","expr":{"static":"git::https://example.com/vpc.git"}}]}`: {
						Attributes: map[string]*schema.AttributeSchema{
							"inputs": {Constraint: schema.Object{Attributes: schema.ObjectAttributes{
								"foo": {
									IsRequired: true,
									Constraint: schema.AnyExpression{OfType: cty.String},
									OriginForTarget: &schema.PathTarget{
										Address:     schema.Address{schema.StaticStep{Name: "var"}, schema.AttrNameStep{}},
										Path:        lang.Path{LanguageID: "terraform"},
										Constraints: schema.Constraints{ScopeId: "variable", Type: cty.String},
									},
								},
							}}},
							"providers": {Constraint: schema.Object{Attributes: schema.ObjectAttributes{}}},
						},
						TargetableAs: schema.Targetables{
							{
								Address:           lang.Address{lang.RootStep{Name: "component"}, lang.AttrStep{Name: "git"}},
								ScopeId:           "module",
								AsType:            cty.Object(map[string]cty.Type{}),
								NestedTargetables: schema.Targetables{},
							},
						},
						ImpliedOrigins: schema.ImpliedOrigins{},
					},
					`{"attrs":[{"name":"source","expr":{"static":"registry/source/test"}}]}`: {
						Attributes: map[string]*schema.AttributeSchema{
							"inputs":    {Constraint: schema.Object{Attributes: schema.ObjectAttributes{}}},
							"providers": {Constraint: schema.Object{Attributes: schema.ObjectAttributes{}}},
							"version": {
								Description: lang.MarkupContent{
									Value: "Accepts a comma-separated list of version constraints for registry modules. Required for registry modules",
									Kind:  lang.MarkdownKind,
								},
								IsRequired: true,
								Constraint: schema.LiteralType{
									Type: cty.String,
								},
							},
						},
						TargetableAs: schema.Targetables{
							{
								Address: lang.Address{lang.RootStep{Name: "component"}, lang.AttrStep{Name: "reg"}},
								ScopeId: "module",
								AsType: cty.Object(map[string]cty.Type{
									"bar": cty.DynamicPseudoType,
								}),
								NestedTargetables: schema.Targetables{
									{
										Address: lang.Address{lang.RootStep{Name: "component"}, lang.AttrStep{Name: "reg"}, lang.AttrStep{Name: "bar"}},
										ScopeId: "component",
										AsType:  cty.DynamicPseudoType,
									},
								},
							},
						},
						ImpliedOrigins: schema.ImpliedOrigins{
							{
								OriginAddress: lang.Address{lang.RootStep{Name: "component"}, lang.AttrStep{Name: "reg"}, lang.AttrStep{Name: "bar"}},
								TargetAddress: lang.Address{lang.RootStep{Name: "output"}, lang.AttrStep{Name: "bar"}},
								Path:          lang.Path{LanguageID: "terraform"},
								Constraints:   schema.Constraints{ScopeId: "output"},
							},
						},
					},
				},
			},
		},
	}

	if diff := cmp.Diff(expectedBodySchema, givenBodySchema, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("schema mismatch: %s", diff)
	}
}

type testStackSchemaReader struct {
	ps *tfjson.ProviderSchemas
}

func (r *testStackSchemaReader) InstalledModulePath(rootPath string, normalizedSource string) (string, bool) {
	if normalizedSource == "git::https://example.com/vpc.git" {
		return "fake/git/path", true
	}
	if normalizedSource == "registry.terraform.io/registry/source/test" {
		return "fake/registry/path", true
	}

	return "", false
}

func (r *testStackSchemaReader) ProviderSchema(modPath string, addr tfaddr.Provider, vc version.Constraints) (*tfschema.ProviderSchema, error) {
	jsonSchema, ok := r.ps.Schemas[addr.String()]
	if !ok {
		return nil, fmt.Errorf("%s: schema not found", addr.String())
	}

	return tfschema.ProviderSchemaFromJson(jsonSchema, addr), nil
}

func (r *testStackSchemaReader) LocalModuleMeta(modPath string) (*tfmod.Meta, error) {
	switch filepath.ToSlash(modPath) {
	case "fake/git/path":
		return &tfmod.Meta{
			Variables: map[string]tfmod.Variable{
				"foo": {Type: cty.String},
			},
		}, nil
	case "fake/registry/path":
		return &tfmod.Meta{
			Outputs: map[string]tfmod.Output{
				"bar": {Value: cty.DynamicVal},
			},
		}, nil
	case "local/path":
		return &tfmod.Meta{
			ProviderReferences: map[tfmod.ProviderRef]tfaddr.Provider{
				{LocalName: "test"}: tfaddr.NewProvider("registry.terraform.io", "hashicorp", "test"),
			},
			Filenames: []string{"main.tf"},
		}, nil
	}

	return nil, fmt.Errorf("invalid source")
}
