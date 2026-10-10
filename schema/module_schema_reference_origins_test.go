// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl-lang/decoder"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/reference"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-schema/module"
	"github.com/zclconf/go-cty/cty"
)

type testDecoderPathReader struct {
	paths map[string]*decoder.PathContext
}

func (r *testDecoderPathReader) Paths(ctx context.Context) []lang.Path {
	paths := make([]lang.Path, 0, len(r.paths))
	for path := range r.paths {
		paths = append(paths, lang.Path{Path: path})
	}
	return paths
}

func (r *testDecoderPathReader) PathContext(path lang.Path) (*decoder.PathContext, error) {
	if ctx, ok := r.paths[path.Path]; ok {
		return ctx, nil
	}
	return nil, fmt.Errorf("path not found: %q", path.Path)
}

func TestUninstalledModuleBlock_collectsVarReferenceOrigins(t *testing.T) {
	source := "git::https://github.com/example/module.git"
	depSchema := schemaForUninstalledModuleBlock(module.DeclaredModuleCall{
		LocalName: "app",
	})

	bodySchema := &schema.BodySchema{
		Blocks: map[string]*schema.BlockSchema{
			"module": {
				Labels: []*schema.LabelSchema{{Name: "name"}},
				Body: &schema.BodySchema{
					Extensions: &schema.BodyExtensions{
						Count:   true,
						ForEach: true,
					},
					Attributes: map[string]*schema.AttributeSchema{
						"source": {
							Constraint: schema.LiteralType{Type: cty.String},
							IsRequired: true,
							IsDepKey:   true,
						},
					},
				},
				DependentBody: map[schema.SchemaKey]*schema.BodySchema{
					schema.NewSchemaKey(schema.DependencyKeys{
						Attributes: []schema.AttributeDependent{
							{
								Name: "source",
								Expr: schema.ExpressionValue{
									Static: cty.StringVal(source),
								},
							},
						},
					}): depSchema,
				},
			},
		},
	}

	cfg := `module "app" {
  source = "git::https://github.com/example/module.git"
  name   = var.location
}
`
	f, diags := hclsyntax.ParseConfig([]byte(cfg), "main.tf", hcl.InitialPos)
	if len(diags) > 0 {
		t.Fatal(diags)
	}

	dirPath := t.TempDir()
	d := decoder.NewDecoder(&testDecoderPathReader{
		paths: map[string]*decoder.PathContext{
			dirPath: {
				Schema: bodySchema,
				Files: map[string]*hcl.File{
					"main.tf": f,
				},
			},
		},
	})
	d.SetContext(decoder.NewDecoderContext())

	pathDecoder, err := d.Path(lang.Path{Path: dirPath})
	if err != nil {
		t.Fatal(err)
	}

	origins, err := pathDecoder.CollectReferenceOrigins()
	if err != nil {
		t.Fatal(err)
	}

	found := false
	expectedAddr := lang.Address{
		lang.RootStep{Name: "var"},
		lang.AttrStep{Name: "location"},
	}
	for _, origin := range origins {
		localOrigin, ok := origin.(reference.LocalOrigin)
		if !ok {
			continue
		}
		if cmp.Equal(localOrigin.Addr, expectedAddr) {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("expected var.location reference origin, got %#v", origins)
	}
}
