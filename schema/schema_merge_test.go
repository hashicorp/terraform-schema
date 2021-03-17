package schema

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/zclconf/go-cty/cty"
)

func TestMergeWithJsonProviderSchemas_noCoreSchema(t *testing.T) {
	sm := NewSchemaMerger(nil)

	_, err := sm.MergeWithJsonProviderSchemas(nil)
	if err == nil {
		t.Fatal("expected error for nil core schema")
	}

	if !errors.Is(err, coreSchemaRequiredErr{}) {
		t.Fatalf("unexpected error: %#v", err)
	}
}

func TestMergeWithJsonProviderSchemas_noProviderSchema(t *testing.T) {
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
		},
	}
	sm := NewSchemaMerger(testCoreSchema)

	_, err := sm.MergeWithJsonProviderSchemas(nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMergeWithJsonProviderSchemas_v012(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/test-config-0.12.tf")
	if err != nil {
		t.Fatal(err)
	}
	f, diags := hclsyntax.ParseConfig(b, "test.tf", hcl.InitialPos)
	if len(diags) > 0 {
		t.Fatal(diags)
	}

	ps := &tfjson.ProviderSchemas{}
	b, err = ioutil.ReadFile("testdata/provider-schemas-0.12.json")
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(b, ps)
	if err != nil {
		t.Fatal(err)
	}

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
		},
	}
	sm := NewSchemaMerger(testCoreSchema)
	sm.SetParsedFiles(map[string]*hcl.File{
		"test.tf": f,
	})

	mergedSchema, err := sm.MergeWithJsonProviderSchemas(ps)
	if err != nil {
		t.Fatal(err)
	}

	opts := cmp.Options{
		cmpopts.IgnoreUnexported(cty.Type{}),
	}

	if diff := cmp.Diff(expectedMergedSchema_v012, mergedSchema, opts); diff != "" {
		t.Fatalf("schema differs: %s", diff)
	}
}

func TestMergeWithJsonProviderSchemas_v013(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/test-config-0.13.tf")
	if err != nil {
		t.Fatal(err)
	}
	f, diags := hclsyntax.ParseConfig(b, "test.tf", hcl.InitialPos)
	if len(diags) > 0 {
		t.Fatal(diags)
	}

	ps := &tfjson.ProviderSchemas{}
	b, err = ioutil.ReadFile("testdata/provider-schemas-0.13.json")
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(b, ps)
	if err != nil {
		t.Fatal(err)
	}

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
		},
	}
	sm := NewSchemaMerger(testCoreSchema)
	sm.SetParsedFiles(map[string]*hcl.File{
		"test.tf": f,
	})

	mergedSchema, err := sm.MergeWithJsonProviderSchemas(ps)
	if err != nil {
		t.Fatal(err)
	}

	opts := cmp.Options{
		cmpopts.IgnoreUnexported(cty.Type{}),
	}

	if diff := cmp.Diff(expectedMergedSchema_v013, mergedSchema, opts); diff != "" {
		t.Fatalf("schema differs: %s", diff)
	}
}
