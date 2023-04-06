// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/module"
	"github.com/zclconf/go-cty-debug/ctydebug"
	"github.com/zclconf/go-cty/cty"
)

func TestSchemaForVariables(t *testing.T) {
	testCases := []struct {
		name           string
		variables      map[string]module.Variable
		expectedSchema *schema.BodySchema
	}{
		{
			"empty schema",
			make(map[string]module.Variable),
			&schema.BodySchema{Attributes: make(map[string]*schema.AttributeSchema)},
		},
		{
			"one attribute schema",
			map[string]module.Variable{
				"name": {
					Description: "name of the module",
					Type:        cty.String,
				},
			},
			&schema.BodySchema{Attributes: map[string]*schema.AttributeSchema{
				"name": {
					Description: lang.MarkupContent{
						Value: "name of the module",
						Kind:  lang.PlainTextKind,
					},
					IsRequired: true,
					Constraint: schema.LiteralType{Type: cty.String},
					OriginForTarget: &schema.PathTarget{
						Address:     schema.Address{schema.StaticStep{Name: "var"}, schema.AttrNameStep{}},
						Path:        lang.Path{Path: "./local", LanguageID: "terraform"},
						Constraints: schema.Constraints{ScopeId: "variable", Type: cty.String},
					},
				},
			}},
		},
		{
			"two attribute schema",
			map[string]module.Variable{
				"name": {
					Description:  "name of the module",
					Type:         cty.String,
					DefaultValue: cty.StringVal("default"),
				},
				"id": {
					Description: "id of the module",
					Type:        cty.Number,
					IsSensitive: true,
				},
			},
			&schema.BodySchema{Attributes: map[string]*schema.AttributeSchema{
				"name": {
					Description: lang.MarkupContent{
						Value: "name of the module",
						Kind:  lang.PlainTextKind,
					},
					IsOptional: true,
					Constraint: schema.LiteralType{Type: cty.String},
					OriginForTarget: &schema.PathTarget{
						Address:     schema.Address{schema.StaticStep{Name: "var"}, schema.AttrNameStep{}},
						Path:        lang.Path{Path: "./local", LanguageID: "terraform"},
						Constraints: schema.Constraints{ScopeId: "variable", Type: cty.String},
					},
				},
				"id": {
					Description: lang.MarkupContent{
						Value: "id of the module",
						Kind:  lang.PlainTextKind,
					},
					Constraint:  schema.LiteralType{Type: cty.Number},
					IsSensitive: true,
					IsRequired:  true,
					OriginForTarget: &schema.PathTarget{
						Address:     schema.Address{schema.StaticStep{Name: "var"}, schema.AttrNameStep{}},
						Path:        lang.Path{Path: "./local", LanguageID: "terraform"},
						Constraints: schema.Constraints{ScopeId: "variable", Type: cty.Number},
					},
				},
			}},
		},
		{
			"attribute with type from default value",
			map[string]module.Variable{
				"name": {
					Description:  "name of the module",
					Type:         cty.DynamicPseudoType,
					DefaultValue: cty.StringVal("default"),
				},
				"id": {
					Description:  "id of the module",
					Type:         cty.NilType,
					DefaultValue: cty.NumberIntVal(42),
				},
			},
			&schema.BodySchema{Attributes: map[string]*schema.AttributeSchema{
				"name": {
					Description: lang.MarkupContent{
						Value: "name of the module",
						Kind:  lang.PlainTextKind,
					},
					IsOptional: true,
					Constraint: schema.LiteralType{Type: cty.String},
					OriginForTarget: &schema.PathTarget{
						Address:     schema.Address{schema.StaticStep{Name: "var"}, schema.AttrNameStep{}},
						Path:        lang.Path{Path: "./local", LanguageID: "terraform"},
						Constraints: schema.Constraints{ScopeId: "variable", Type: cty.String},
					},
				},
				"id": {
					Description: lang.MarkupContent{
						Value: "id of the module",
						Kind:  lang.PlainTextKind,
					},
					Constraint: schema.LiteralType{Type: cty.Number},
					IsOptional: true,
					OriginForTarget: &schema.PathTarget{
						Address:     schema.Address{schema.StaticStep{Name: "var"}, schema.AttrNameStep{}},
						Path:        lang.Path{Path: "./local", LanguageID: "terraform"},
						Constraints: schema.Constraints{ScopeId: "variable", Type: cty.Number},
					},
				},
			}},
		},
	}

	modPath := "./local"
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.name), func(t *testing.T) {
			actualSchema, err := SchemaForVariables(tc.variables, modPath)

			if err != nil {
				t.Fatal(err)
			}

			diff := cmp.Diff(tc.expectedSchema, actualSchema, ctydebug.CmpOptions)
			if diff != "" {
				t.Fatalf("unexpected schema %s", diff)
			}
		})
	}
}
