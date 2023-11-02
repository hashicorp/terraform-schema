// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/hcl/v2"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/module"
	"github.com/hashicorp/terraform-schema/registry"
	"github.com/zclconf/go-cty-debug/ctydebug"
	"github.com/zclconf/go-cty/cty"
)

func TestSchemaForDependentModuleBlock_emptyMeta(t *testing.T) {
	meta := &module.Meta{}
	module := module.InstalledModuleCall{
		LocalName: "refname",
	}
	depSchema, err := schemaForDependentModuleBlock(module, meta)
	if err != nil {
		t.Fatal(err)
	}
	expectedDepSchema := &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{},
		TargetableAs: []*schema.Targetable{
			{
				Address: lang.Address{
					lang.RootStep{Name: "module"},
					lang.AttrStep{Name: "refname"},
				},
				ScopeId:           refscope.ModuleScope,
				AsType:            cty.Object(map[string]cty.Type{}),
				NestedTargetables: []*schema.Targetable{},
			},
		},
		ImpliedOrigins: schema.ImpliedOrigins{},
	}
	if diff := cmp.Diff(expectedDepSchema, depSchema, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("schema mismatch: %s", diff)
	}
}

func TestSchemaForDependentModuleBlock_basic(t *testing.T) {
	meta := &module.Meta{
		Path: "./local",
		Variables: map[string]module.Variable{
			"example_var": {
				Description: "Test var",
				Type:        cty.String,
				IsSensitive: true,
			},
			"another_var": {
				DefaultValue: cty.StringVal("bar"),
			},
		},
		Outputs: map[string]module.Output{
			"first": {
				Description: "first output",
				IsSensitive: true,
				Value:       cty.BoolVal(true),
			},
			"second": {
				Description: "second output",
				Value:       cty.ListVal([]cty.Value{cty.StringVal("test")}),
			},
		},
	}
	module := module.InstalledModuleCall{
		LocalName: "refname",
	}
	depSchema, err := schemaForDependentModuleBlock(module, meta)
	if err != nil {
		t.Fatal(err)
	}
	expectedDepSchema := &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{
			"example_var": {
				Constraint:  schema.AnyExpression{OfType: cty.String},
				Description: lang.PlainText("Test var"),
				IsRequired:  true,
				IsSensitive: true,
				OriginForTarget: &schema.PathTarget{
					Address: schema.Address{
						schema.StaticStep{Name: "var"},
						schema.AttrNameStep{},
					},
					Path: lang.Path{
						Path:       "./local",
						LanguageID: "terraform",
					},
					Constraints: schema.Constraints{
						ScopeId: "variable",
						Type:    cty.String,
					},
				},
			},
			"another_var": {
				Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType},
				IsOptional: true,
				OriginForTarget: &schema.PathTarget{
					Address: schema.Address{
						schema.StaticStep{Name: "var"},
						schema.AttrNameStep{},
					},
					Path: lang.Path{
						Path:       "./local",
						LanguageID: "terraform",
					},
					Constraints: schema.Constraints{
						ScopeId: "variable",
						Type:    cty.DynamicPseudoType,
					},
				},
			},
		},
		TargetableAs: []*schema.Targetable{
			{
				Address: lang.Address{
					lang.RootStep{Name: "module"},
					lang.AttrStep{Name: "refname"},
				},
				ScopeId: refscope.ModuleScope,
				AsType: cty.Object(map[string]cty.Type{
					"first":  cty.Bool,
					"second": cty.List(cty.String),
				}),
				NestedTargetables: []*schema.Targetable{
					{
						Address: lang.Address{
							lang.RootStep{Name: "module"},
							lang.AttrStep{Name: "refname"},
							lang.AttrStep{Name: "first"},
						},
						ScopeId:     refscope.ModuleScope,
						AsType:      cty.Bool,
						IsSensitive: true,
						Description: lang.PlainText("first output"),
					},
					{
						Address: lang.Address{
							lang.RootStep{Name: "module"},
							lang.AttrStep{Name: "refname"},
							lang.AttrStep{Name: "second"},
						},
						ScopeId:     refscope.ModuleScope,
						AsType:      cty.List(cty.String),
						Description: lang.PlainText("second output"),
						NestedTargetables: []*schema.Targetable{
							{
								Address: lang.Address{
									lang.RootStep{Name: "module"},
									lang.AttrStep{Name: "refname"},
									lang.AttrStep{Name: "second"},
									lang.IndexStep{Key: cty.NumberIntVal(0)},
								},
								ScopeId: refscope.ModuleScope,
								AsType:  cty.String,
							},
						},
					},
				},
			},
		},
		ImpliedOrigins: schema.ImpliedOrigins{
			{
				OriginAddress: lang.Address{
					lang.RootStep{Name: "module"},
					lang.AttrStep{Name: "refname"},
					lang.AttrStep{Name: "first"},
				},
				TargetAddress: lang.Address{
					lang.RootStep{Name: "output"},
					lang.AttrStep{Name: "first"},
				},
				Path:        lang.Path{Path: "./local", LanguageID: "terraform"},
				Constraints: schema.Constraints{ScopeId: "output"},
			},
			{
				OriginAddress: lang.Address{
					lang.RootStep{Name: "module"},
					lang.AttrStep{Name: "refname"},
					lang.AttrStep{Name: "second"},
				},
				TargetAddress: lang.Address{
					lang.RootStep{Name: "output"},
					lang.AttrStep{Name: "second"},
				},
				Path:        lang.Path{Path: "./local", LanguageID: "terraform"},
				Constraints: schema.Constraints{ScopeId: "output"},
			},
		},
	}

	sort.Slice(depSchema.ImpliedOrigins, func(i, j int) bool {
		return depSchema.ImpliedOrigins[i].OriginAddress.String() < depSchema.ImpliedOrigins[j].OriginAddress.String()
	})

	if diff := cmp.Diff(expectedDepSchema, depSchema, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("schema mismatch: %s", diff)
	}
}

func TestSchemaForDependentModuleBlock_Target(t *testing.T) {
	type testCase struct {
		name           string
		meta           *module.Meta
		expectedSchema *schema.BodySchema
	}

	testCases := []testCase{
		{
			"no target",
			&module.Meta{
				Path:      "./local",
				Variables: map[string]module.Variable{},
				Outputs:   map[string]module.Output{},
				Filenames: nil,
			},
			&schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{},
				TargetableAs: []*schema.Targetable{
					{
						Address: lang.Address{
							lang.RootStep{Name: "module"},
							lang.AttrStep{Name: "refname"},
						},
						ScopeId:           refscope.ModuleScope,
						AsType:            cty.Object(map[string]cty.Type{}),
						NestedTargetables: []*schema.Targetable{},
					},
				},
				Targets:        nil,
				ImpliedOrigins: schema.ImpliedOrigins{},
			},
		},
		{
			"without main.tf",
			&module.Meta{
				Path:      "./local",
				Variables: map[string]module.Variable{},
				Outputs:   map[string]module.Output{},
				Filenames: []string{"a_file.tf", "b_file.tf"},
			},
			&schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{},
				TargetableAs: []*schema.Targetable{
					{
						Address: lang.Address{
							lang.RootStep{Name: "module"},
							lang.AttrStep{Name: "refname"},
						},
						ScopeId:           refscope.ModuleScope,
						AsType:            cty.Object(map[string]cty.Type{}),
						NestedTargetables: []*schema.Targetable{},
					},
				},
				Targets: &schema.Target{
					Path: lang.Path{Path: "./local", LanguageID: "terraform"},
					Range: hcl.Range{
						Filename: "a_file.tf",
						Start:    hcl.InitialPos,
						End:      hcl.InitialPos,
					},
				},
				ImpliedOrigins: schema.ImpliedOrigins{},
			},
		},
		{
			"with main.tf",
			&module.Meta{
				Path:      "./local",
				Variables: map[string]module.Variable{},
				Outputs:   map[string]module.Output{},
				Filenames: []string{"a_file.tf", "main.tf"},
			},
			&schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{},
				TargetableAs: []*schema.Targetable{
					{
						Address: lang.Address{
							lang.RootStep{Name: "module"},
							lang.AttrStep{Name: "refname"},
						},
						ScopeId:           refscope.ModuleScope,
						AsType:            cty.Object(map[string]cty.Type{}),
						NestedTargetables: []*schema.Targetable{},
					},
				},
				Targets: &schema.Target{
					Path: lang.Path{Path: "./local", LanguageID: "terraform"},
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.InitialPos,
						End:      hcl.InitialPos,
					},
				},
				ImpliedOrigins: schema.ImpliedOrigins{},
			},
		},
	}
	module := module.InstalledModuleCall{
		LocalName: "refname",
	}

	for _, tc := range testCases {
		depSchema, err := schemaForDependentModuleBlock(module, tc.meta)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(tc.expectedSchema, depSchema, ctydebug.CmpOptions); diff != "" {
			t.Fatalf("schema mismatch: %s", diff)
		}
	}
}

func TestSchemaForDependentModuleBlock_DocsLink(t *testing.T) {
	type testCase struct {
		name           string
		meta           *module.Meta
		module         module.InstalledModuleCall
		expectedSchema *schema.BodySchema
	}

	testCases := []testCase{
		{
			"local module",
			&module.Meta{
				Path:      "./local",
				Variables: map[string]module.Variable{},
				Outputs:   map[string]module.Output{},
				Filenames: nil,
			},
			module.InstalledModuleCall{
				LocalName:  "refname",
				SourceAddr: module.LocalSourceAddr("./local"),
			},
			&schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{},
				TargetableAs: []*schema.Targetable{
					{
						Address: lang.Address{
							lang.RootStep{Name: "module"},
							lang.AttrStep{Name: "refname"},
						},
						ScopeId:           refscope.ModuleScope,
						AsType:            cty.Object(map[string]cty.Type{}),
						NestedTargetables: []*schema.Targetable{},
					},
				},
				ImpliedOrigins: schema.ImpliedOrigins{},
				Targets:        nil,
			},
		},
		{
			"remote module",
			&module.Meta{
				Path:      "registry.terraform.io/terraform-aws-modules/vpc/aws",
				Variables: map[string]module.Variable{},
				Outputs:   map[string]module.Output{},
				Filenames: nil,
			},
			module.InstalledModuleCall{
				LocalName:  "vpc",
				SourceAddr: tfaddr.MustParseModuleSource("registry.terraform.io/terraform-aws-modules/vpc/aws"),
			},
			&schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{},
				TargetableAs: []*schema.Targetable{
					{
						Address: lang.Address{
							lang.RootStep{Name: "module"},
							lang.AttrStep{Name: "vpc"},
						},
						ScopeId:           refscope.ModuleScope,
						AsType:            cty.Object(map[string]cty.Type{}),
						NestedTargetables: []*schema.Targetable{},
					},
				},
				ImpliedOrigins: schema.ImpliedOrigins{},
				DocsLink: &schema.DocsLink{
					URL: "https://registry.terraform.io/modules/terraform-aws-modules/vpc/aws/latest",
				},
			},
		},
		{
			"remote module with version",
			&module.Meta{
				Path:      "registry.terraform.io/terraform-aws-modules/vpc/aws",
				Variables: map[string]module.Variable{},
				Outputs:   map[string]module.Output{},
				Filenames: nil,
			},
			module.InstalledModuleCall{
				LocalName:  "vpc",
				SourceAddr: tfaddr.MustParseModuleSource("registry.terraform.io/terraform-aws-modules/vpc/aws"),
				Version:    version.Must(version.NewVersion("1.33.7")),
			},
			&schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{},
				TargetableAs: []*schema.Targetable{
					{
						Address: lang.Address{
							lang.RootStep{Name: "module"},
							lang.AttrStep{Name: "vpc"},
						},
						ScopeId:           refscope.ModuleScope,
						AsType:            cty.Object(map[string]cty.Type{}),
						NestedTargetables: []*schema.Targetable{},
					},
				},
				ImpliedOrigins: schema.ImpliedOrigins{},
				DocsLink: &schema.DocsLink{
					URL: "https://registry.terraform.io/modules/terraform-aws-modules/vpc/aws/1.33.7",
				},
			},
		},
		{
			"remote module on unknown registry",
			&module.Meta{
				Path:      "example.com/terraform-aws-modules/vpc/aws",
				Variables: map[string]module.Variable{},
				Outputs:   map[string]module.Output{},
				Filenames: nil,
			},
			module.InstalledModuleCall{
				LocalName:  "vpc",
				SourceAddr: tfaddr.MustParseModuleSource("example.com/terraform-aws-modules/vpc/aws"),
				Version:    version.Must(version.NewVersion("1.33.7")),
			},
			&schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{},
				TargetableAs: []*schema.Targetable{
					{
						Address: lang.Address{
							lang.RootStep{Name: "module"},
							lang.AttrStep{Name: "vpc"},
						},
						ScopeId:           refscope.ModuleScope,
						AsType:            cty.Object(map[string]cty.Type{}),
						NestedTargetables: []*schema.Targetable{},
					},
				},
				ImpliedOrigins: schema.ImpliedOrigins{},
			},
		},
	}

	for _, tc := range testCases {
		depSchema, err := schemaForDependentModuleBlock(tc.module, tc.meta)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(tc.expectedSchema, depSchema, ctydebug.CmpOptions); diff != "" {
			t.Fatalf("schema mismatch: %s", diff)
		}
	}
}

func TestSchemaForDeclaredDependentModuleBlock_basic(t *testing.T) {
	meta := &registry.ModuleData{
		Version: version.Must(version.NewVersion("1.0.0")),
		Inputs: []registry.Input{
			{
				Name:        "example_var",
				Type:        cty.String,
				Description: lang.PlainText("Test var"),
				Required:    true,
			},
			{
				Name:    "foo_var",
				Type:    cty.DynamicPseudoType,
				Default: cty.NumberIntVal(42),
			},
			{
				Name: "another_var",
				Type: cty.DynamicPseudoType,
			},
		},
		Outputs: []registry.Output{
			{
				Name:        "first",
				Description: lang.PlainText("first output"),
			},
			{
				Name:        "second",
				Description: lang.PlainText("second output"),
			},
		},
	}
	module := module.DeclaredModuleCall{
		LocalName:  "refname",
		SourceAddr: tfaddr.MustParseModuleSource("terraform-aws-modules/eks/aws"),
	}
	depSchema, err := schemaForDeclaredDependentModuleBlock(module, meta)
	if err != nil {
		t.Fatal(err)
	}
	expectedDepSchema := &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{
			"example_var": {
				Constraint:  schema.AnyExpression{OfType: cty.String},
				Description: lang.PlainText("Test var"),
				IsRequired:  true,
			},
			"foo_var": {
				Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType},
				IsOptional: true,
			},
			"another_var": {
				Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType},
				IsOptional: true,
			},
		},
		TargetableAs: []*schema.Targetable{
			{
				Address: lang.Address{
					lang.RootStep{Name: "module"},
					lang.AttrStep{Name: "refname"},
				},
				ScopeId: refscope.ModuleScope,
				AsType: cty.Object(map[string]cty.Type{
					"first":  cty.DynamicPseudoType,
					"second": cty.DynamicPseudoType,
				}),
				NestedTargetables: []*schema.Targetable{
					{
						Address: lang.Address{
							lang.RootStep{Name: "module"},
							lang.AttrStep{Name: "refname"},
							lang.AttrStep{Name: "first"},
						},
						ScopeId:     refscope.ModuleScope,
						AsType:      cty.DynamicPseudoType,
						Description: lang.PlainText("first output"),
					},
					{
						Address: lang.Address{
							lang.RootStep{Name: "module"},
							lang.AttrStep{Name: "refname"},
							lang.AttrStep{Name: "second"},
						},
						ScopeId:     refscope.ModuleScope,
						AsType:      cty.DynamicPseudoType,
						Description: lang.PlainText("second output"),
					},
				},
			},
		},
		DocsLink: &schema.DocsLink{
			URL: "https://registry.terraform.io/modules/terraform-aws-modules/eks/aws/1.0.0",
		},
	}
	if diff := cmp.Diff(expectedDepSchema, depSchema, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("schema mismatch: %s", diff)
	}
}
