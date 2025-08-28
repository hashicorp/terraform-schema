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
	tfjson "github.com/hashicorp/terraform-json"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/hashicorp/terraform-schema/internal/addr"
	tfmod "github.com/hashicorp/terraform-schema/module"
	tfschema "github.com/hashicorp/terraform-schema/schema"
	tfsearch "github.com/hashicorp/terraform-schema/search"
	"github.com/zclconf/go-cty-debug/ctydebug"
	"github.com/zclconf/go-cty/cty"
)

func TestSearchSchemaMerger_SchemaForSearch_noCoreSchema(t *testing.T) {
	sm := NewSearchSchemaMerger(nil)

	_, err := sm.SchemaForSearch(nil)
	if err == nil {
		t.Fatal("expected error for nil core schema")
	}

	if !errors.Is(err, tfschema.CoreSchemaRequiredErr{}) {
		t.Fatalf("unexpected error: %#v", err)
	}
}

func TestSearchSchemaMerger_SchemaForSearch_noProviderSchema(t *testing.T) {
	testCoreSchema := &schema.BodySchema{}

	sm := NewSearchSchemaMerger(testCoreSchema)

	_, err := sm.SchemaForSearch(&tfsearch.Meta{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSearchSchemaMerger_SchemaForSearch_providerNameMatch(t *testing.T) {
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
			"list": {},
		},
	}
	sm := NewSearchSchemaMerger(testCoreSchema)
	sm.SetStateReader(&testSearchSchemaReader{
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

	givenBodySchema, err := sm.SchemaForSearch(&tfsearch.Meta{
		Path: "local/path",
		ProviderReferences: map[tfsearch.ProviderRef]tfaddr.Provider{
			{LocalName: "data"}: addr.NewDefaultProvider("data"),
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
			"list": {
				DependentBody: map[schema.SchemaKey]*schema.BodySchema{},
			},
		},
	}

	if diff := cmp.Diff(expectedBodySchema, givenBodySchema, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("schema mismatch: %s", diff)
	}
}

func TestSchemaMerger_SchemaForSearch_twiceMerged(t *testing.T) {
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
			"list": {},
		},
	}
	sm := NewSearchSchemaMerger(testCoreSchema)
	sm.SetStateReader(&testSearchSchemaReader{
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

	givenBodySchema, err := sm.SchemaForSearch(&tfsearch.Meta{
		Path: "local/path",
		ProviderReferences: map[tfsearch.ProviderRef]tfaddr.Provider{
			{LocalName: "data"}: addr.NewDefaultProvider("data"),
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
			"list": {
				DependentBody: map[schema.SchemaKey]*schema.BodySchema{},
			},
		},
	}

	if diff := cmp.Diff(expectedBodySchema, givenBodySchema, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("schema mismatch: %s", diff)
	}

	// now merge again with different local name
	givenBodySchema, err = sm.SchemaForSearch(&tfsearch.Meta{
		Path: "local/path",
		ProviderReferences: map[tfsearch.ProviderRef]tfaddr.Provider{
			{LocalName: "data"}: addr.NewDefaultProvider("data"),
		},
	})

	if err != nil {
		t.Fatal(err)
	}
	expectedBodySchema = &schema.BodySchema{
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
			"list": {
				DependentBody: map[schema.SchemaKey]*schema.BodySchema{},
			},
		},
	}

	if diff := cmp.Diff(expectedBodySchema, givenBodySchema, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("schema mismatch: %s", diff)
	}
}

func TestSearchSchemaMerger_SchemaForSearch_variables(t *testing.T) {
	testCoreSchema := &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"provider": {},
			"list":     {},
			"variable": {},
		},
	}
	sm := NewSearchSchemaMerger(testCoreSchema)
	sm.SetStateReader(&testSearchSchemaReader{})

	givenBodySchema, err := sm.SchemaForSearch(&tfsearch.Meta{
		Variables: map[string]tfsearch.Variable{
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
			"list": {
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

func TestSearchSchemaMerger_SchemaForSearch_lists(t *testing.T) {
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
			"list": {
				Labels: []*schema.LabelSchema{
					{Name: "name"},
				},
			},
		},
	}
	sm := NewSearchSchemaMerger(testCoreSchema)
	sm.SetStateReader(&testSearchSchemaReader{
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
					ListResourceSchemas: map[string]*tfjson.Schema{
						"dummy_resource": {
							Block: &tfjson.SchemaBlock{
								NestedBlocks: map[string]*tfjson.SchemaBlockType{
									"config": {
										Block: &tfjson.SchemaBlock{
											Attributes: map[string]*tfjson.SchemaAttribute{
												"count": {
													AttributeType: cty.Number,
													Optional:      true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	givenBodySchema, err := sm.SchemaForSearch(&tfsearch.Meta{
		Path: "local/path",
		ProviderReferences: map[tfsearch.ProviderRef]tfaddr.Provider{
			{LocalName: "data"}: addr.NewDefaultProvider("data"),
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
			"list": {
				Labels: []*schema.LabelSchema{
					{Name: "name"},
				},
				DependentBody: map[schema.SchemaKey]*schema.BodySchema{
					`{"labels":[{"index":0,"value":"dummy_resource"}],"attrs":[{"name":"provider","expr":{"addr":"data"}}]}`: {
						Detail: "hashicorp/data",
						Blocks: map[string]*schema.BlockSchema{
							"config": {
								Description: lang.Markdown("Filters specific to the list type"),
								MaxItems:    1,
								Body: &schema.BodySchema{
									Blocks: map[string]*schema.BlockSchema{
										"config": {
											Labels: []*schema.LabelSchema{},
											Body: &schema.BodySchema{
												Blocks: map[string]*schema.BlockSchema{},
												Attributes: map[string]*schema.AttributeSchema{
													"count": {
														IsOptional: true,
														Constraint: schema.AnyExpression{OfType: cty.Number},
													},
												},
											},
										},
									},
									Attributes: map[string]*schema.AttributeSchema{},
									Detail:     "hashicorp/data",
								},
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

type testSearchSchemaReader struct {
	ps *tfjson.ProviderSchemas
}

func (r *testSearchSchemaReader) InstalledModulePath(rootPath string, normalizedSource string) (string, bool) {
	if normalizedSource == "git::https://example.com/vpc.git" {
		return "fake/git/path", true
	}
	if normalizedSource == "registry.terraform.io/registry/source/test" {
		return "fake/registry/path", true
	}

	return "", false
}

func (r *testSearchSchemaReader) ProviderSchema(modPath string, addr tfaddr.Provider, vc version.Constraints) (*tfschema.ProviderSchema, error) {
	jsonSchema, ok := r.ps.Schemas[addr.String()]
	if !ok {
		return nil, fmt.Errorf("%s: schema not found", addr.String())
	}

	return tfschema.ProviderSchemaFromJson(jsonSchema, addr), nil
}

func (r *testSearchSchemaReader) LocalModuleMeta(modPath string) (*tfmod.Meta, error) {
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
			ProviderRequirements: tfmod.ProviderRequirements{
				addr.NewDefaultProvider("data"): version.MustConstraints(version.NewConstraint("1.0")),
			},
			Filenames: []string{"main.tf"},
		}, nil
	}

	return nil, fmt.Errorf("invalid source")
}
