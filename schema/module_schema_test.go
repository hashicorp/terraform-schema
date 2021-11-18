package schema

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/module"
	"github.com/zclconf/go-cty-debug/ctydebug"
	"github.com/zclconf/go-cty/cty"
)

func TestSchemaForDependentModuleBlock_emptyMeta(t *testing.T) {
	meta := &module.Meta{}
	depSchema, err := schemaForDependentModuleBlock("refname", meta)
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
	depSchema, err := schemaForDependentModuleBlock("refname", meta)
	if err != nil {
		t.Fatal(err)
	}
	expectedDepSchema := &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{
			"example_var": {
				Expr: schema.ExprConstraints{
					schema.TraversalExpr{OfType: cty.String},
					schema.LiteralTypeExpr{Type: cty.String},
				},
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
				Expr: schema.ExprConstraints{
					schema.TraversalExpr{OfType: cty.String},
					schema.LiteralTypeExpr{Type: cty.String},
				},
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
						Type:    cty.String,
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
	}
	if diff := cmp.Diff(expectedDepSchema, depSchema, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("schema mismatch: %s", diff)
	}
}
