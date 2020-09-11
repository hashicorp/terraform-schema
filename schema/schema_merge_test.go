package schema

import (
	"encoding/json"
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

func TestMergeCoreWithJsonProviderSchemas_v012(t *testing.T) {
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

	files := map[string]*hcl.File{
		"test.tf": f,
	}

	mergedSchema, err := MergeCoreWithJsonProviderSchemas(files, testCoreSchema, ps)
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

func TestMergeCoreWithJsonProviderSchemas_v013(t *testing.T) {
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

	files := map[string]*hcl.File{
		"test.tf": f,
	}

	mergedSchema, err := MergeCoreWithJsonProviderSchemas(files, testCoreSchema, ps)
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

var testCoreSchema = &schema.BodySchema{
	Blocks: map[string]*schema.BlockSchema{
		"provider": {
			Labels: []*schema.LabelSchema{
				{Name: "name"},
			},
			Body: &schema.BodySchema{
				Attributes: map[string]*schema.AttributeSchema{
					"alias": {ValueType: cty.String},
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
					"count": {ValueType: cty.Number},
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
					"count": {ValueType: cty.Number},
				},
			},
		},
	},
}
