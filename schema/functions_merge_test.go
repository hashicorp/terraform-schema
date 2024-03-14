// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	tfjson "github.com/hashicorp/terraform-json"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/hashicorp/terraform-schema/internal/addr"
	tfmod "github.com/hashicorp/terraform-schema/module"
	"github.com/zclconf/go-cty-debug/ctydebug"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

func TestFunctionsMerger_FunctionsForModule_noCoreFunctions(t *testing.T) {
	sm := NewFunctionsMerger(nil)

	_, err := sm.FunctionsForModule(nil)
	if err == nil {
		t.Fatal("expected error for nil core schema")
	}

	if !errors.Is(err, coreFunctionsRequiredErr{}) {
		t.Fatalf("unexpected error: %#v", err)
	}
}

func TestFunctionsMerger_FunctionsForModule_noMeta(t *testing.T) {
	coreFunctions := map[string]schema.FunctionSignature{
		"foo": {
			Params: []function.Parameter{
				{Name: "bar", Type: cty.String, Description: "bar function"},
			},
		},
	}

	sm := NewFunctionsMerger(coreFunctions)

	functions, err := sm.FunctionsForModule(nil)
	if err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}

	if len(functions) != 1 {
		t.Fatalf("unexpected functions: %#v", functions)
	}

	if functions["foo"].Params[0].Name != "bar" {
		t.Fatalf("unexpected function: %#v", functions)
	}
}

var providerSchemaWithFunctions = map[string]*tfjson.ProviderSchema{
	"registry.terraform.io/hashicorp/test": {
		ConfigSchema:      &tfjson.Schema{},
		DataSourceSchemas: map[string]*tfjson.Schema{},
		ResourceSchemas:   map[string]*tfjson.Schema{},
		Functions: map[string]*tfjson.FunctionSignature{
			"bar": {
				Parameters: []*tfjson.FunctionParameter{
					{Name: "baz", Type: cty.String, Description: "baz param"},
				},
				Description: "bar function",
				ReturnType:  cty.String,
			},
			"alleven": {
				VariadicParameter: &tfjson.FunctionParameter{
					Name: "numbers",
					Type: cty.List(cty.Number),
				},
				Description: "Returns true if all passed arguments are even numbers",
				ReturnType:  cty.Bool,
			},
		},
	},
}

func TestFunctionsMerger_FunctionsForModule_18(t *testing.T) {
	coreFunctions := map[string]schema.FunctionSignature{
		"foo": {
			Params: []function.Parameter{
				{Name: "bar", Type: cty.String, Description: "bar param"},
			},
			Description: "foo function",
		},
	}

	fm := NewFunctionsMerger(coreFunctions)
	fm.SetStateReader(&testJsonSchemaReader{
		ps: &tfjson.ProviderSchemas{
			FormatVersion: "1.0",
			Schemas:       providerSchemaWithFunctions,
		},
	})
	fm.SetTerraformVersion(version.Must(version.NewVersion("1.8")))

	testProvider := addr.NewDefaultProvider("test")
	versionConstraints := version.MustConstraints(version.NewConstraint("1.0.0"))
	meta := &tfmod.Meta{
		ProviderReferences: map[tfmod.ProviderRef]tfaddr.Provider{
			{LocalName: "localtest"}: testProvider,
		},
		ProviderRequirements: tfmod.ProviderRequirements{
			testProvider: versionConstraints,
		},
	}

	expectedFunctions := map[string]schema.FunctionSignature{
		"foo": {
			Params: []function.Parameter{
				{Name: "bar", Type: cty.String, Description: "bar param"},
			},
			Description: "foo function",
		},
		"provider::localtest::bar": {
			Params: []function.Parameter{
				{Name: "baz", Type: cty.String, Description: "baz param"},
			},
			Description: "bar function",
			Detail:      "hashicorp/test",
			ReturnType:  cty.String,
		},
		"provider::localtest::alleven": {
			Params: []function.Parameter{},
			VarParam: &function.Parameter{
				Name: "numbers", Type: cty.List(cty.Number),
			},
			Description: "Returns true if all passed arguments are even numbers",
			Detail:      "hashicorp/test",
			ReturnType:  cty.Bool,
		},
	}

	givenFunctions, err := fm.FunctionsForModule(meta)
	if err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}

	if diff := cmp.Diff(expectedFunctions, givenFunctions, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("functions mismatch: %s", diff)
	}
}

func TestFunctionsMerger_FunctionsForModule_17(t *testing.T) {
	fm := NewFunctionsMerger(map[string]schema.FunctionSignature{})
	fm.SetStateReader(&testJsonSchemaReader{
		ps: &tfjson.ProviderSchemas{
			FormatVersion: "1.0",
			Schemas:       providerSchemaWithFunctions,
		},
	})
	fm.SetTerraformVersion(version.Must(version.NewVersion("1.7")))

	testProvider := addr.NewDefaultProvider("test")
	versionConstraints := version.MustConstraints(version.NewConstraint("1.0.0"))
	meta := &tfmod.Meta{
		ProviderReferences: map[tfmod.ProviderRef]tfaddr.Provider{
			{LocalName: "localtest"}: testProvider,
		},
		ProviderRequirements: tfmod.ProviderRequirements{
			testProvider: versionConstraints,
		},
	}

	expectedFunctions := map[string]schema.FunctionSignature{}

	givenFunctions, err := fm.FunctionsForModule(meta)
	if err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}

	if diff := cmp.Diff(expectedFunctions, givenFunctions, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("functions mismatch: %s", diff)
	}
}
