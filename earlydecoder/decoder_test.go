package earlydecoder

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/hashicorp/terraform-schema/backend"
	"github.com/hashicorp/terraform-schema/module"
	"github.com/zclconf/go-cty-debug/ctydebug"
	"github.com/zclconf/go-cty/cty"
)

type testCase struct {
	name          string
	cfg           string
	expectedMeta  *module.Meta
	expectedError hcl.Diagnostics
}

var customComparer = []cmp.Option{
	cmp.Comparer(compareVersionConstraint),
	ctydebug.CmpOptions,
}

func TestLoadModule(t *testing.T) {
	path := t.TempDir()

	testCases := []testCase{
		{
			"empty config",
			``,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
			},
			nil,
		},
		{
			"core requirements only",
			`
terraform {
  required_version = "~> 0.12"
}`,
			&module.Meta{
				Path:                 path,
				CoreRequirements:     mustConstraints(t, "~> 0.12"),
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
			},
			nil,
		},
		{
			"legacy inferred provider requirements",
			`
provider "aws" {
  region = "eu-west-2"
}

resource "google_storage_bucket" "bucket" {
  name = "test-bucket"
}

data "blah_foobar" "test" {
  name = "something"
}

provider "grafana" {
  url    = "http://grafana.example.com/"
  org_id = 1
}
`,
			&module.Meta{
				Path: path,
				ProviderReferences: map[module.ProviderRef]tfaddr.Provider{
					{LocalName: "aws"}:     tfaddr.NewLegacyProvider("aws"),
					{LocalName: "blah"}:    tfaddr.NewLegacyProvider("blah"),
					{LocalName: "google"}:  tfaddr.NewLegacyProvider("google"),
					{LocalName: "grafana"}: tfaddr.NewLegacyProvider("grafana"),
				},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{
					tfaddr.NewLegacyProvider("aws"):     {},
					tfaddr.NewLegacyProvider("blah"):    {},
					tfaddr.NewLegacyProvider("google"):  {},
					tfaddr.NewLegacyProvider("grafana"): {},
				},
				Variables: map[string]module.Variable{},
				Outputs:   map[string]module.Output{},
			},
			nil,
		},
		{
			"simplified 0.12 provider requirements",
			`
terraform {
  required_providers {
    aws = "1.2.0"
    google = ">= 3.0.0"
  }
}
provider "aws" {
  region = "eu-west-2"
}

resource "google_storage_bucket" "bucket" {
  name = "test-bucket"
}

provider "grafana" {
  url    = "http://grafana.example.com/"
  org_id = 1
}
`,
			&module.Meta{
				Path: path,
				ProviderReferences: map[module.ProviderRef]tfaddr.Provider{
					{LocalName: "aws"}:     tfaddr.NewLegacyProvider("aws"),
					{LocalName: "google"}:  tfaddr.NewLegacyProvider("google"),
					{LocalName: "grafana"}: tfaddr.NewLegacyProvider("grafana"),
				},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{
					tfaddr.NewLegacyProvider("aws"):     mustConstraints(t, "1.2.0"),
					tfaddr.NewLegacyProvider("google"):  mustConstraints(t, ">= 3.0.0"),
					tfaddr.NewLegacyProvider("grafana"): {},
				},
				Variables: map[string]module.Variable{},
				Outputs:   map[string]module.Output{},
			},
			nil,
		},
		{
			"version-only 0.12 provider requirements",
			`
terraform {
  required_providers {
    aws = {
    	version = "1.2.0"
    }
    google = {
    	version = ">= 3.0.0"
    }
  }
}
provider "aws" {
  region = "eu-west-2"
}

resource "google_storage_bucket" "bucket" {
  name = "test-bucket"
}

provider "grafana" {
  url    = "http://grafana.example.com/"
  org_id = 1
}
`,
			&module.Meta{
				Path: path,
				ProviderReferences: map[module.ProviderRef]tfaddr.Provider{
					{LocalName: "aws"}:     tfaddr.NewLegacyProvider("aws"),
					{LocalName: "google"}:  tfaddr.NewLegacyProvider("google"),
					{LocalName: "grafana"}: tfaddr.NewLegacyProvider("grafana"),
				},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{
					tfaddr.NewLegacyProvider("aws"):     mustConstraints(t, "1.2.0"),
					tfaddr.NewLegacyProvider("google"):  mustConstraints(t, ">= 3.0.0"),
					tfaddr.NewLegacyProvider("grafana"): {},
				},
				Variables: map[string]module.Variable{},
				Outputs:   map[string]module.Output{},
			},
			nil,
		},
		{
			"0.13+ provider requirements",
			`
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "1.0.0"
    }
    google = {
      source  = "hashicorp/google"
      version = "2.0.0"
    }
    grafana = {
      source  = "grafana/grafana"
      version = "2.1.0"
    }
  }
}
provider "aws" {
  region = "eu-west-2"
}

resource "google_storage_bucket" "bucket" {
  name = "test-bucket"
}

provider "grafana" {
  url    = "http://grafana.example.com/"
  org_id = 1
}
`,
			&module.Meta{
				Path: path,
				ProviderReferences: map[module.ProviderRef]tfaddr.Provider{
					{LocalName: "aws"}: {
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "hashicorp",
						Type:      "aws",
					},
					{LocalName: "google"}: {
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "hashicorp",
						Type:      "google",
					},
					{LocalName: "grafana"}: {
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "grafana",
						Type:      "grafana",
					},
				},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{
					{
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "hashicorp",
						Type:      "aws",
					}: mustConstraints(t, "1.0.0"),
					{
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "hashicorp",
						Type:      "google",
					}: mustConstraints(t, "2.0.0"),
					{
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "grafana",
						Type:      "grafana",
					}: mustConstraints(t, "2.1.0"),
				},
				Variables: map[string]module.Variable{},
				Outputs:   map[string]module.Output{},
			},
			nil,
		},
		{
			"multiple valid version requirements",
			`
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 1.0.0"
    }
    google = {
      source  = "hashicorp/google"
      version = "2.0.0"
    }
  }
}

terraform {
  required_providers {
    aws = {
      version = "1.1.0"
    }
  }
}

provider "aws" {
  region = "eu-west-2"
}

resource "google_storage_bucket" "bucket" {
  name = "test-bucket"
}
`,
			&module.Meta{
				Path: path,
				ProviderReferences: map[module.ProviderRef]tfaddr.Provider{
					{LocalName: "aws"}: {
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "hashicorp",
						Type:      "aws",
					},
					{LocalName: "google"}: {
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "hashicorp",
						Type:      "google",
					},
				},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{
					{
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "hashicorp",
						Type:      "aws",
					}: mustConstraints(t, ">= 1.0.0,1.1.0"),
					{
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "hashicorp",
						Type:      "google",
					}: mustConstraints(t, "2.0.0"),
				},
				Variables: map[string]module.Variable{},
				Outputs:   map[string]module.Output{},
			},
			nil,
		},
		{
			"multiple invalid version requirements",
			`
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 1.0.0"
    }
    google = {
      source  = "hashicorp/google"
      version = "2.0.0"
    }
  }
}

terraform {
  required_providers {
    aws = {
    	source = "hashicorp/aws"
      version = "1.1.0"
    }
  }
}

provider "aws" {
  region = "eu-west-2"
}

resource "google_storage_bucket" "bucket" {
  name = "test-bucket"
}
`,
			&module.Meta{
				Path: path,
				ProviderReferences: map[module.ProviderRef]tfaddr.Provider{
					{LocalName: "aws"}: {
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "hashicorp",
						Type:      "aws",
					},
					{LocalName: "google"}: {
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "hashicorp",
						Type:      "google",
					},
				},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{
					{
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "hashicorp",
						Type:      "aws",
					}: mustConstraints(t, ">= 1.0.0,1.1.0"),
					{
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "hashicorp",
						Type:      "google",
					}: mustConstraints(t, "2.0.0"),
				},
				Variables: map[string]module.Variable{},
				Outputs:   map[string]module.Output{},
			},
			nil,
		},
		{
			"0.13+ provider aliases",
			`
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "1.0.0"
    }
    google = {
      source  = "hashicorp/google"
      version = "2.0.0"
    }
  }
}
provider "aws" {
  alias = "euwest"
  region = "eu-west-2"
}
`,
			&module.Meta{
				Path: path,
				ProviderReferences: map[module.ProviderRef]tfaddr.Provider{
					{LocalName: "aws"}: {
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "hashicorp",
						Type:      "aws",
					},
					{LocalName: "aws", Alias: "euwest"}: {
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "hashicorp",
						Type:      "aws",
					},
					{LocalName: "google"}: {
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "hashicorp",
						Type:      "google",
					},
				},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{
					{
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "hashicorp",
						Type:      "aws",
					}: mustConstraints(t, "1.0.0"),
					{
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "hashicorp",
						Type:      "google",
					}: mustConstraints(t, "2.0.0"),
				},
				Variables: map[string]module.Variable{},
				Outputs:   map[string]module.Output{},
			},
			nil,
		},
		{
			"0.15+ provider aliases",
			`
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "1.0.0"
      configuration_aliases = [aws.east]
    }
    google = {
      source  = "hashicorp/google"
      version = "2.0.0"
    }
  }
}
provider "aws" {
  alias = "west"
  region = "eu-west-2"
}
`,
			&module.Meta{
				Path: path,
				ProviderReferences: map[module.ProviderRef]tfaddr.Provider{
					{LocalName: "aws"}: {
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "hashicorp",
						Type:      "aws",
					},
					{LocalName: "aws", Alias: "east"}: {
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "hashicorp",
						Type:      "aws",
					},
					{LocalName: "aws", Alias: "west"}: {
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "hashicorp",
						Type:      "aws",
					},
					{LocalName: "google"}: {
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "hashicorp",
						Type:      "google",
					},
				},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{
					{
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "hashicorp",
						Type:      "aws",
					}: mustConstraints(t, "1.0.0"),
					{
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "hashicorp",
						Type:      "google",
					}: mustConstraints(t, "2.0.0"),
				},
				Variables: map[string]module.Variable{},
				Outputs:   map[string]module.Output{},
			},
			nil,
		},
		{
			"explicit provider association",
			`
terraform {
  required_providers {
    goo = {
      source  = "hashicorp/google-beta"
      version = "2.0.0"
    }
  }
}

resource "google_something" "test" {
	provider = goo
}
`,
			&module.Meta{
				Path: path,
				ProviderReferences: map[module.ProviderRef]tfaddr.Provider{
					{LocalName: "goo"}: {
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "hashicorp",
						Type:      "google-beta",
					},
				},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{
					{
						Hostname:  tfaddr.DefaultRegistryHost,
						Namespace: "hashicorp",
						Type:      "google-beta",
					}: mustConstraints(t, "2.0.0"),
				},
				Variables: map[string]module.Variable{},
				Outputs:   map[string]module.Output{},
			},
			nil,
		},
	}

	runTestCases(testCases, t, path)
}

func TestLoadModule_Variables(t *testing.T) {
	path := t.TempDir()

	testCases := []testCase{
		{
			"no name variables",
			`
variable "" {
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
			},
			nil,
		},
		{
			"no name variables",
			`
variable {
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
			},
			hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Missing name for variable",
					Detail:   "All variable blocks must have 1 labels (name).",
					Subject: &hcl.Range{
						Filename: "test.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 10,
							Byte:   10,
						},
						End: hcl.Pos{
							Line:   2,
							Column: 11,
							Byte:   11,
						},
					},
					Context: &hcl.Range{
						Filename: "test.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 1,
							Byte:   1,
						},
						End: hcl.Pos{
							Line:   2,
							Column: 11,
							Byte:   11,
						},
					},
				},
			},
		},
		{
			"double label variables",
			`
variable "one" "two" {
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
			},
			hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Extraneous label for variable",
					Detail:   "Only 1 labels (name) are expected for variable blocks.",
					Subject: &hcl.Range{
						Filename: "test.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 16,
							Byte:   16,
						},
						End: hcl.Pos{
							Line:   2,
							Column: 21,
							Byte:   21,
						},
					},
					Context: &hcl.Range{
						Filename: "test.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 1,
							Byte:   1,
						},
						End: hcl.Pos{
							Line:   2,
							Column: 23,
							Byte:   23,
						},
					},
				},
			},
		},
		{
			"empty variables",
			`
variable "name" {
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables: map[string]module.Variable{
					"name": {
						Type: cty.DynamicPseudoType,
					},
				},
				Outputs: map[string]module.Output{},
			},
			nil,
		},
		{
			"variables with type",
			`
variable "name" {
	type = string
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables: map[string]module.Variable{
					"name": {
						Type: cty.String,
					},
				},
				Outputs: map[string]module.Output{},
			},
			nil,
		},
		{
			"variables with description",
			`
variable "name" {
	description = "description"
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables: map[string]module.Variable{
					"name": {
						Type:        cty.DynamicPseudoType,
						Description: "description",
					},
				},
				Outputs: map[string]module.Output{},
			},
			nil,
		},
		{
			"variables with sensitive",
			`
variable "name" {
	sensitive = true
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables: map[string]module.Variable{
					"name": {
						Type:        cty.DynamicPseudoType,
						IsSensitive: true,
					},
				},
				Outputs: map[string]module.Output{},
			},
			nil,
		},
		{
			"variables with type and description and sensitive",
			`
variable "name" {
	type = string
	description = "description"
	sensitive = true
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables: map[string]module.Variable{
					"name": {
						Type:        cty.String,
						Description: "description",
						IsSensitive: true,
					},
				},
				Outputs: map[string]module.Output{},
			},
			nil,
		},
		{
			"variables with default",
			`
variable "name" {
  	default = {}
}`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables: map[string]module.Variable{
					"name": {
						Type:         cty.DynamicPseudoType,
						DefaultValue: cty.EmptyObjectVal,
					},
				},
				Outputs: map[string]module.Output{},
			},
			nil,
		},
		{
			"empty output",
			`
output "name" {
}
`,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs: map[string]module.Output{
					"name": {Value: cty.NilVal},
				},
			},
			nil,
		},
	}

	runTestCases(testCases, t, path)
}

func TestLoadModule_backend(t *testing.T) {
	path := t.TempDir()

	testCases := []testCase{
		{
			"no backend",
			`
terraform {

}`,
			&module.Meta{
				Path:                 path,
				Backend:              nil,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
			},
			nil,
		},
		{
			"s3 backend",
			`
terraform {
  backend "s3" {
  	blah = "test"
  }
}`,
			&module.Meta{
				Path: path,
				Backend: &module.Backend{
					Type: "s3",
					Data: &backend.UnknownBackendData{},
				},
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
			},
			nil,
		},
		{
			"empty remote backend",
			`
terraform {
  backend "remote" {}
}`,
			&module.Meta{
				Path: path,
				Backend: &module.Backend{
					Type: "remote",
					Data: &backend.Remote{},
				},
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
			},
			nil,
		},
		{
			"remote backend with hostname",
			`
terraform {
  backend "remote" {
  	hostname = "app.terraform.io"
  }
}`,
			&module.Meta{
				Path: path,
				Backend: &module.Backend{
					Type: "remote",
					Data: &backend.Remote{Hostname: "app.terraform.io"},
				},
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
			},
			nil,
		},
		{
			"remote backend with hostname and more attributes",
			`
terraform {
  backend "remote" {
    hostname = "app.terraform.io"
    organization = "test"

    workspaces {
      name = "test"
    }
  }
}`,
			&module.Meta{
				Path: path,
				Backend: &module.Backend{
					Type: "remote",
					Data: &backend.Remote{Hostname: "app.terraform.io"},
				},
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
				Variables:            map[string]module.Variable{},
				Outputs:              map[string]module.Output{},
			},
			nil,
		},
	}

	runTestCases(testCases, t, path)
}

func runTestCases(testCases []testCase, t *testing.T, path string) {
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.name), func(t *testing.T) {
			f, diags := hclsyntax.ParseConfig([]byte(tc.cfg), "test.tf", hcl.InitialPos)
			if len(diags) > 0 {
				t.Fatal(diags)
			}
			files := map[string]*hcl.File{
				"test.tf": f,
			}

			meta, diags := LoadModule(path, files)

			if diff := cmp.Diff(tc.expectedError, diags, customComparer...); diff != "" {
				t.Fatalf("expected errors doesn't match: %s", diff)
			}

			if diff := cmp.Diff(tc.expectedMeta, meta, customComparer...); diff != "" {
				t.Fatalf("module meta doesn't match: %s", diff)
			}
		})
	}
}

func mustConstraints(t *testing.T, vc string) version.Constraints {
	c, err := version.NewConstraint(vc)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func compareVersionConstraint(x, y version.Constraint) bool {
	return x.String() == y.String()
}
