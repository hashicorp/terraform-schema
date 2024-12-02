// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package schema

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	tfschema "github.com/hashicorp/terraform-schema/schema"
	"github.com/hashicorp/terraform-schema/stack"
	"github.com/zclconf/go-cty-debug/ctydebug"
	"github.com/zclconf/go-cty/cty"
)

func TestDeploySchemaMerger_SchemaForDeployment_noCoreSchema(t *testing.T) {
	sm := NewDeploySchemaMerger(nil)

	_, err := sm.SchemaForDeployment(nil)
	if err == nil {
		t.Fatal("expected error for nil core schema")
	}

	if !errors.Is(err, tfschema.CoreSchemaRequiredErr{}) {
		t.Fatalf("unexpected error: %#v", err)
	}
}

func TestDeploySchemaMerger_SchemaForDeployment_no_inputs(t *testing.T) {
	testCoreSchema := &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"deployment": {
				Labels: []*schema.LabelSchema{
					{Name: "name"},
				},
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"inputs": {
							Description: lang.Markdown("A mapping of stack variable names to values for this deployment. The keys of this map must correspond to the names of variables defined for the stack. The values must be valid HCL literals meeting the type constraint of those variables. Values are also expressions, currently with access to identity token references only"),
							IsOptional:  true,
							Constraint: schema.Map{
								Name: "map of variable references",
								Elem: schema.AnyExpression{OfType: cty.DynamicPseudoType},
							},
						},
					},
				},
			},
		},
	}

	sm := NewDeploySchemaMerger(testCoreSchema)

	givenBodySchema, err := sm.SchemaForDeployment(&stack.Meta{
		Variables: map[string]stack.Variable{},
	})
	if err != nil {
		t.Fatal(err)
	}

	expectedBodySchema := &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"deployment": {
				Labels: []*schema.LabelSchema{
					{Name: "name"},
				},
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"inputs": {
							Description: lang.Markdown("A mapping of stack variable names to values for this deployment. The keys of this map must correspond to the names of variables defined for the stack. The values must be valid HCL literals meeting the type constraint of those variables. Values are also expressions, currently with access to identity token references only"),
							IsOptional:  true,
							Constraint:  schema.Object{Attributes: schema.ObjectAttributes{}},
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

func TestDeploySchemaMerger_SchemaForDeployment_inputs(t *testing.T) {
	testCoreSchema := &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"deployment": {
				Labels: []*schema.LabelSchema{
					{Name: "name"},
				},
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"inputs": {
							Description: lang.Markdown("A mapping of stack variable names to values for this deployment. The keys of this map must correspond to the names of variables defined for the stack. The values must be valid HCL literals meeting the type constraint of those variables. Values are also expressions, currently with access to identity token references only"),
							IsOptional:  true,
							Constraint: schema.Map{
								Name: "map of variable references",
								Elem: schema.AnyExpression{OfType: cty.DynamicPseudoType},
							},
						},
					},
				},
			},
		},
	}

	sm := NewDeploySchemaMerger(testCoreSchema)

	givenBodySchema, err := sm.SchemaForDeployment(&stack.Meta{
		Variables: map[string]stack.Variable{
			"foo": {
				Type: cty.String,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	expectedBodySchema := &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"deployment": {
				Labels: []*schema.LabelSchema{
					{Name: "name"},
				},
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"inputs": {
							Description: lang.Markdown("A mapping of stack variable names to values for this deployment. The keys of this map must correspond to the names of variables defined for the stack. The values must be valid HCL literals meeting the type constraint of those variables. Values are also expressions, currently with access to identity token references only"),
							IsOptional:  true,
							Constraint: schema.Object{Attributes: schema.ObjectAttributes{
								"foo": {
									IsRequired: true,
									Constraint: schema.AnyExpression{OfType: cty.String},
									OriginForTarget: &schema.PathTarget{
										Address:     schema.Address{schema.StaticStep{Name: "var"}, schema.AttrNameStep{}},
										Path:        lang.Path{LanguageID: "terraform-stack"},
										Constraints: schema.Constraints{ScopeId: "variable", Type: cty.String},
									},
								},
							}},
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
