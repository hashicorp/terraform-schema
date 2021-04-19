package earlydecoder

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-registry-address"
	"github.com/hashicorp/terraform-schema/module"
)

func TestLoadModule(t *testing.T) {
	path := t.TempDir()

	testCases := []struct {
		name         string
		cfg          string
		expectedMeta *module.Meta
	}{
		{
			"empty config",
			``,
			&module.Meta{
				Path:                 path,
				ProviderReferences:   map[module.ProviderRef]tfaddr.Provider{},
				ProviderRequirements: map[tfaddr.Provider]version.Constraints{},
			},
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
			},
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
			},
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
			},
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
			},
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
			},
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
			},
		},
	}

	opts := cmp.Comparer(compareVersionConstraint)

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
			if len(diags) > 0 {
				t.Fatal(diags)
			}

			if diff := cmp.Diff(tc.expectedMeta, meta, opts); diff != "" {
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
