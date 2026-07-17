// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	tfjson "github.com/hashicorp/terraform-json"
	tfpolicy "github.com/hashicorp/terraform-schema/policy"
	"github.com/zclconf/go-cty-debug/ctydebug"
	"github.com/zclconf/go-cty/cty"
)

func TestPolicySchemaMerger_SchemaForSearch_inputs(t *testing.T) {
	testCoreSchema := &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"module_policy":   {},
			"resource_policy": {},
			"provider_policy": {},
			"input":           {},
		},
	}
	sm := NewSchemaMerger(testCoreSchema)
	sm.SetStateReader(&testSearchSchemaReader{})

	givenBodySchema, err := sm.SchemaForPolicy(&tfpolicy.Meta{
		Inputs: map[string]tfpolicy.Input{
			"region": {Type: cty.String, Description: "AWS region"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	expectedBodySchema := &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"module_policy":   {},
			"resource_policy": {},
			"provider_policy": {},
			"input": {
				Labels: []*schema.LabelSchema{{Name: "name", IsDepKey: true, Description: lang.MarkupContent{Value: "Input name", Kind: lang.PlainTextKind}}},
				DependentBody: map[schema.SchemaKey]*schema.BodySchema{
					`{"labels":[{"index":0,"value":"region"}]}`: {
						Attributes: map[string]*schema.AttributeSchema{
							"default": {
								Constraint:  schema.LiteralType{Type: cty.String},
								Description: lang.MarkupContent{Value: "Default value to use when input is not explicitly set", Kind: lang.MarkdownKind},
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

type testSearchSchemaReader struct {
	ps *tfjson.ProviderSchemas
}
