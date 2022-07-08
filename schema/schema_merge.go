package schema

import (
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/hashicorp/terraform-schema/internal/schema/backends"
	"github.com/hashicorp/terraform-schema/module"
	"github.com/hashicorp/terraform-schema/registry"
	"github.com/zclconf/go-cty/cty"
)

type SchemaMerger struct {
	coreSchema       *schema.BodySchema
	schemaReader     SchemaReader
	terraformVersion *version.Version
	moduleReader     ModuleReader
}

type ModuleReader interface {
	ModuleCalls(modPath string) (module.ModuleCalls, error)
	LocalModuleMeta(modPath string) (*module.Meta, error)
	RegistryModuleMeta(addr tfaddr.Module, cons version.Constraints) (*registry.ModuleData, error)
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

func (m *SchemaMerger) SetModuleReader(mr ModuleReader) {
	m.moduleReader = mr
}

func (m *SchemaMerger) SetTerraformVersion(v *version.Version) {
	m.terraformVersion = v
}

func (m *SchemaMerger) SchemaForModule(meta *module.Meta) (*schema.BodySchema, error) {
	if m.coreSchema == nil {
		return nil, coreSchemaRequiredErr{}
	}

	if meta == nil {
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
	if mergedSchema.Blocks["module"].DependentBody == nil {
		mergedSchema.Blocks["module"].DependentBody = make(map[schema.SchemaKey]*schema.BodySchema)
	}

	providerRefs := ProviderReferences(meta.ProviderReferences)

	if m.schemaReader != nil {
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
						dsSchema.Attributes["backend"].SemanticTokenModifiers = lang.SemanticTokenModifiers{lang.TokenModifierDependent}
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

	if m.moduleReader != nil {
		reader := m.moduleReader
		mc, err := reader.ModuleCalls(meta.Path)
		if err != nil {
			return mergedSchema, nil
		}

		for _, module := range mc.Declared {
			sourceAddr, ok := module.SourceAddr.(tfaddr.Module)
			if !ok {
				// TODO: local sources (See https://github.com/hashicorp/terraform-ls/issues/598)
				continue
			}

			modMeta, err := reader.RegistryModuleMeta(sourceAddr, module.Version)
			if err != nil {
				continue
			}

			// Fetching based only on the source can cause conflicts for multiple versions of the same module
			// specially if they have different versions or the source of those modules have been modified
			// inside the .terraform folder. This is a compromise that we made in this moment since it would impact only auto completion
			depKeys := schema.DependencyKeys{
				Attributes: []schema.AttributeDependent{
					{
						Name: "source",
						Expr: schema.ExpressionValue{
							Static: cty.StringVal(sourceAddr.String()),
						},
					},
				},
			}
			// There's likely more edge cases with how source address can be represented in config
			// vs in module manifest, but for now we at least account for the common case of external registries
			depKeysAddr := schema.DependencyKeys{
				Attributes: []schema.AttributeDependent{
					{
						Name: "source",
						Expr: schema.ExpressionValue{
							Static: cty.StringVal(sourceAddr.Package.ForRegistryProtocol()),
						},
					},
				},
			}

			depSchema, err := schemaForDeclaredDependentModuleBlock(module, modMeta)
			if err == nil {
				mergedSchema.Blocks["module"].DependentBody[schema.NewSchemaKey(depKeys)] = depSchema
				mergedSchema.Blocks["module"].DependentBody[schema.NewSchemaKey(depKeysAddr)] = depSchema
			}
		}

		for _, module := range mc.Installed {
			if module.SourceAddr == nil {
				// This should never happen for installed modules, but to
				// be safe we skip all modules with an empty source address
				continue
			}
			modMeta, err := reader.LocalModuleMeta(module.Path)
			if err != nil {
				continue
			}

			depKeys := schema.DependencyKeys{
				// Fetching based only on the source can cause conflicts for multiple versions of the same module
				// specially if they have different versions or the source of those modules have been modified
				// inside the .terraform folder. This is a compromise that we made in this moment since it would impact only auto completion
				Attributes: []schema.AttributeDependent{
					{
						Name: "source",
						Expr: schema.ExpressionValue{
							Static: cty.StringVal(module.SourceAddr.String()),
						},
					},
				},
			}

			depSchema, err := schemaForDependentModuleBlock(module, modMeta)
			if err == nil {
				mergedSchema.Blocks["module"].DependentBody[schema.NewSchemaKey(depKeys)] = depSchema
			}

			// There's likely more edge cases with how source address can be represented in config
			// vs in module manifest, but for now we at least account for the common case of external registries
			registryAddr, ok := module.SourceAddr.(tfaddr.Module)
			if err == nil && ok {
				depKeys := schema.DependencyKeys{
					// Fetching based only on the source can cause conflicts for multiple versions of the same module
					// specially if they have different versions or the source of those modules have been modified
					// inside the .terraform folder. This is a compromise that we made in this moment since it would impact only auto completion
					Attributes: []schema.AttributeDependent{
						{
							Name: "source",
							Expr: schema.ExpressionValue{
								Static: cty.StringVal(registryAddr.Package.ForRegistryProtocol()),
							},
						},
					},
				}

				mergedSchema.Blocks["module"].DependentBody[schema.NewSchemaKey(depKeys)] = depSchema
			}
		}
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
