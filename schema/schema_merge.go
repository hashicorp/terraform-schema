package schema

import (
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/hashicorp/terraform-schema/internal/schema/backends"
	"github.com/hashicorp/terraform-schema/module"
)

type SchemaMerger struct {
	coreSchema       *schema.BodySchema
	schemaReader     SchemaReader
	terraformVersion *version.Version
}

type ProviderSchema struct {
	Provider    *schema.BodySchema
	Resources   map[string]*schema.BodySchema
	DataSources map[string]*schema.BodySchema
}

func (ps *ProviderSchema) Copy() *ProviderSchema {
	if ps == nil {
		return nil
	}

	newPs := &ProviderSchema{
		Provider: ps.Provider.Copy(),
	}

	if ps.Resources != nil {
		newPs.Resources = make(map[string]*schema.BodySchema, len(ps.Resources))
		for name, rSchema := range ps.Resources {
			newPs.Resources[name] = rSchema.Copy()
		}
	}

	if ps.DataSources != nil {
		newPs.DataSources = make(map[string]*schema.BodySchema, len(ps.DataSources))
		for name, rSchema := range ps.DataSources {
			newPs.DataSources[name] = rSchema.Copy()
		}
	}

	return newPs
}

type SchemaReader interface {
	ProviderSchema(modPath string, addr tfaddr.Provider, vc version.Constraints) (*ProviderSchema, error)
}

func NewSchemaMerger(coreSchema *schema.BodySchema) *SchemaMerger {
	return &SchemaMerger{
		coreSchema: coreSchema,
	}
}

func (m *SchemaMerger) SetSchemaReader(sr SchemaReader) {
	m.schemaReader = sr
}

func (m *SchemaMerger) SetTerraformVersion(v *version.Version) {
	m.terraformVersion = v
}

func (m *SchemaMerger) SchemaForModule(meta *module.Meta) (*schema.BodySchema, error) {
	if m.coreSchema == nil {
		return nil, coreSchemaRequiredErr{}
	}

	if meta == nil || m.schemaReader == nil {
		return m.coreSchema, nil
	}

	mergedSchema := m.coreSchema.Copy()

	if mergedSchema.Blocks["provider"].DependentBody == nil {
		mergedSchema.Blocks["provider"].DependentBody = make(map[schema.SchemaKey]*schema.BodySchema)
	}
	if mergedSchema.Blocks["resource"].DependentBody == nil {
		mergedSchema.Blocks["resource"].DependentBody = make(map[schema.SchemaKey]*schema.BodySchema)
	}
	if mergedSchema.Blocks["data"].DependentBody == nil {
		mergedSchema.Blocks["data"].DependentBody = make(map[schema.SchemaKey]*schema.BodySchema)
	}

	providerRefs := ProviderReferences(meta.ProviderReferences)

	for pAddr, pVersionCons := range meta.ProviderRequirements {
		pSchema, err := m.schemaReader.ProviderSchema(meta.Path, pAddr, pVersionCons)
		if err != nil {
			continue
		}

		refs := providerRefs.ReferencesOfProvider(pAddr)
		for _, localRef := range refs {
			if pSchema.Provider != nil {
				mergedSchema.Blocks["provider"].DependentBody[schema.NewSchemaKey(schema.DependencyKeys{
					Labels: []schema.LabelDependent{
						{Index: 0, Value: localRef.LocalName},
					},
				})] = pSchema.Provider
			}

			providerAddr := lang.Address{
				lang.RootStep{Name: localRef.LocalName},
			}
			if localRef.Alias != "" {
				providerAddr = append(providerAddr, lang.AttrStep{Name: localRef.Alias})
			}

			for rName, rSchema := range pSchema.Resources {
				depKeys := schema.DependencyKeys{
					Labels: []schema.LabelDependent{
						{Index: 0, Value: rName},
					},
					Attributes: []schema.AttributeDependent{
						{
							Name: "provider",
							Expr: schema.ExpressionValue{
								Address: providerAddr,
							},
						},
					},
				}
				mergedSchema.Blocks["resource"].DependentBody[schema.NewSchemaKey(depKeys)] = rSchema

				// No explicit association is required
				// if the resource prefix matches provider name
				if strings.HasPrefix(rName, localRef.LocalName+"_") {
					depKeys := schema.DependencyKeys{
						Labels: []schema.LabelDependent{
							{Index: 0, Value: rName},
						},
					}
					mergedSchema.Blocks["resource"].DependentBody[schema.NewSchemaKey(depKeys)] = rSchema
				}
			}

			for dsName, dsSchema := range pSchema.DataSources {
				depKeys := schema.DependencyKeys{
					Labels: []schema.LabelDependent{
						{Index: 0, Value: dsName},
					},
					Attributes: []schema.AttributeDependent{
						{
							Name: "provider",
							Expr: schema.ExpressionValue{
								Address: providerAddr,
							},
						},
					},
				}

				// Add backend-related core bits of schema
				if isRemoteStateDataSource(pAddr, dsName) {
					dsSchema.Attributes["backend"].IsDepKey = true
					dsSchema.Attributes["backend"].Expr = backends.BackendTypesAsExprConstraints(m.terraformVersion)

					delete(dsSchema.Attributes, "config")
					depBodies := m.dependentBodyForRemoteStateDataSource(providerAddr, localRef)
					for key, depBody := range depBodies {
						mergedSchema.Blocks["data"].DependentBody[key] = depBody
					}
				}

				mergedSchema.Blocks["data"].DependentBody[schema.NewSchemaKey(depKeys)] = dsSchema

				// No explicit association is required
				// if the resource prefix matches provider name
				if strings.HasPrefix(dsName, localRef.LocalName+"_") {
					depKeys := schema.DependencyKeys{
						Labels: []schema.LabelDependent{
							{Index: 0, Value: dsName},
						},
					}
					mergedSchema.Blocks["data"].DependentBody[schema.NewSchemaKey(depKeys)] = dsSchema
				}
			}
		}
	}

	if _, ok := mergedSchema.Blocks["variable"]; ok {
		mergedSchema.Blocks["variable"].Labels = []*schema.LabelSchema{
			{
				Name:        "name",
				IsDepKey:    true,
				Description: lang.PlainText("Variable name"),
			},
		}
		mergedSchema.Blocks["variable"].DependentBody = variableDependentBody(meta.Variables)
	}
	return mergedSchema, nil
}

func variableDependentBody(vars map[string]module.Variable) map[schema.SchemaKey]*schema.BodySchema {
	depBodies := make(map[schema.SchemaKey]*schema.BodySchema)

	for name, mVar := range vars {
		depKeys := schema.DependencyKeys{
			Labels: []schema.LabelDependent{
				{Index: 0, Value: name},
			},
		}
		depBodies[schema.NewSchemaKey(depKeys)] = &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"default": {
					Expr:        schema.ExprConstraints{schema.LiteralTypeExpr{Type: mVar.Type}},
					IsOptional:  true,
					Description: lang.Markdown("Default value to use when variable is not explicitly set"),
				},
			},
		}
	}

	return depBodies
}

type ProviderReferences map[module.ProviderRef]tfaddr.Provider

func (pr ProviderReferences) ReferencesOfProvider(addr tfaddr.Provider) []module.ProviderRef {
	refs := make([]module.ProviderRef, 0)

	for ref, pAddr := range pr {
		if pAddr.Equals(addr) {
			refs = append(refs, ref)
		}
	}

	return refs
}
