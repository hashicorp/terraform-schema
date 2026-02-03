// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/hashicorp/terraform-schema/internal/schema/backends"
	tfmod "github.com/hashicorp/terraform-schema/module"
	"github.com/hashicorp/terraform-schema/registry"
	"github.com/zclconf/go-cty/cty"
)

type SchemaMerger struct {
	coreSchema       *schema.BodySchema
	terraformVersion *version.Version
	stateReader      StateReader
}

// StateReader exposes a set of methods to read data from the internal language server state
type StateReader interface {
	// DeclaredModuleCalls returns a map of declared module calls for the given module
	// A declared module call refers to a module block in the configuration
	DeclaredModuleCalls(modPath string) (map[string]tfmod.DeclaredModuleCall, error)

	// InstalledModulePath checks if there is an installed module available for
	// the given normalized source address.
	InstalledModulePath(rootPath string, normalizedSource string) (string, bool)

	// LocalModuleMeta returns the module meta data for a local module. This is the result
	// of the [earlydecoder] when processing module files
	LocalModuleMeta(modPath string) (*tfmod.Meta, error)

	// RegistryModuleMeta returns the module meta data for public registry modules. We fetch this
	// data from the registry API.
	RegistryModuleMeta(addr tfaddr.Module, cons version.Constraints) (*registry.ModuleData, error)

	// ProviderSchema returns the schema for a provider we have stored in memory. The can come
	// from different sources.
	ProviderSchema(modPath string, addr tfaddr.Provider, vc version.Constraints) (*ProviderSchema, error)
}

func NewSchemaMerger(coreSchema *schema.BodySchema) *SchemaMerger {
	return &SchemaMerger{
		coreSchema: coreSchema,
	}
}

func (m *SchemaMerger) SetStateReader(mr StateReader) {
	m.stateReader = mr
}

func (m *SchemaMerger) SetTerraformVersion(v *version.Version) {
	m.terraformVersion = v
}

func (m *SchemaMerger) SchemaForModule(meta *tfmod.Meta) (*schema.BodySchema, error) {
	if m.coreSchema == nil {
		return nil, CoreSchemaRequiredErr{}
	}

	if meta == nil {
		return m.coreSchema, nil
	}

	if m.stateReader == nil {
		return m.coreSchema, nil
	}

	mergedSchema := m.coreSchema.Copy()

	if mergedSchema.Blocks["provider"].DependentBody == nil {
		mergedSchema.Blocks["provider"].DependentBody = make(map[schema.SchemaKey]*schema.BodySchema)
	}
	if mergedSchema.Blocks["resource"].DependentBody == nil {
		mergedSchema.Blocks["resource"].DependentBody = make(map[schema.SchemaKey]*schema.BodySchema)
	}
	if ephemeralBlock, ok := mergedSchema.Blocks["ephemeral"]; ok && ephemeralBlock.DependentBody == nil {
		mergedSchema.Blocks["ephemeral"].DependentBody = make(map[schema.SchemaKey]*schema.BodySchema)
	}
	if mergedSchema.Blocks["data"].DependentBody == nil {
		mergedSchema.Blocks["data"].DependentBody = make(map[schema.SchemaKey]*schema.BodySchema)
	}
	if mergedSchema.Blocks["module"].DependentBody == nil {
		mergedSchema.Blocks["module"].DependentBody = make(map[schema.SchemaKey]*schema.BodySchema)
	}
	if checkBlock, ok := mergedSchema.Blocks["check"]; ok && checkBlock.Body.Blocks["data"].DependentBody == nil {
		mergedSchema.Blocks["check"].Body.Blocks["data"].DependentBody = make(map[schema.SchemaKey]*schema.BodySchema)
	}
	if actionBlock, ok := mergedSchema.Blocks["action"]; ok && actionBlock.DependentBody == nil {
		mergedSchema.Blocks["action"].DependentBody = make(map[schema.SchemaKey]*schema.BodySchema)
	}

	providerRefs := ProviderReferences(meta.ProviderReferences)

	for pAddr, pVersionCons := range meta.ProviderRequirements {
		pSchema, err := m.stateReader.ProviderSchema(meta.Path, pAddr, pVersionCons)
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
				if TypeBelongsToProvider(rName, localRef) {
					depKeys := schema.DependencyKeys{
						Labels: []schema.LabelDependent{
							{Index: 0, Value: rName},
						},
					}
					mergedSchema.Blocks["resource"].DependentBody[schema.NewSchemaKey(depKeys)] = rSchema
				}
			}

			if m.terraformVersion.GreaterThanOrEqual(v1_14) {
				for arName, arSchema := range pSchema.ActionResources {
					// Create a BodySchema that ensures a config block exists
					actionBodySchema := &schema.BodySchema{
						HoverURL:     arSchema.HoverURL,
						DocsLink:     arSchema.DocsLink,
						Detail:       arSchema.Detail,
						Description:  arSchema.Description,
						IsDeprecated: arSchema.IsDeprecated,
						Blocks: map[string]*schema.BlockSchema{
							"config": {
								Description: lang.Markdown("Provider specific action configuration"),
								MaxItems:    1,
								Body:        arSchema,
							},
						},
					}

					depKeys := schema.DependencyKeys{
						Labels: []schema.LabelDependent{
							{Index: 0, Value: arName},
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
					mergedSchema.Blocks["action"].DependentBody[schema.NewSchemaKey(depKeys)] = actionBodySchema

					if TypeBelongsToProvider(arName, localRef) {
						depKeys := schema.DependencyKeys{
							Labels: []schema.LabelDependent{
								{Index: 0, Value: arName},
							},
						}
						mergedSchema.Blocks["action"].DependentBody[schema.NewSchemaKey(depKeys)] = actionBodySchema
					}
				}
			}

			// Ephemeral resources were introduced in Terraform 1.10, so we don't need to
			// merge them for older versions
			if m.terraformVersion.GreaterThanOrEqual(v1_10) {
				for erName, erSchema := range pSchema.EphemeralResources {
					depKeys := schema.DependencyKeys{
						Labels: []schema.LabelDependent{
							{Index: 0, Value: erName},
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
					mergedSchema.Blocks["ephemeral"].DependentBody[schema.NewSchemaKey(depKeys)] = erSchema

					// No explicit association is required
					// if the ephemeral resource prefix matches provider name
					if TypeBelongsToProvider(erName, localRef) {
						depKeys := schema.DependencyKeys{
							Labels: []schema.LabelDependent{
								{Index: 0, Value: erName},
							},
						}
						mergedSchema.Blocks["ephemeral"].DependentBody[schema.NewSchemaKey(depKeys)] = erSchema
					}
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
					remoteStateDs := dsSchema.Copy()

					remoteStateDs.Attributes["backend"].IsDepKey = true
					remoteStateDs.Attributes["backend"].SemanticTokenModifiers = lang.SemanticTokenModifiers{lang.TokenModifierDependent}
					remoteStateDs.Attributes["backend"].Constraint = backends.BackendTypesAsOneOfConstraint(m.terraformVersion)
					delete(remoteStateDs.Attributes, "config")

					depBodies := m.dependentBodyForRemoteStateDataSource(remoteStateDs, providerAddr, localRef)
					for key, depBody := range depBodies {
						mergedSchema.Blocks["data"].DependentBody[key] = depBody
						if _, ok := mergedSchema.Blocks["check"]; ok {
							mergedSchema.Blocks["check"].Body.Blocks["data"].DependentBody[key] = depBody
						}
					}

					dsSchema = remoteStateDs
				}

				mergedSchema.Blocks["data"].DependentBody[schema.NewSchemaKey(depKeys)] = dsSchema

				if _, ok := mergedSchema.Blocks["check"]; ok {
					mergedSchema.Blocks["check"].Body.Blocks["data"].DependentBody[schema.NewSchemaKey(depKeys)] = dsSchema
				}

				// No explicit association is required
				// if the resource prefix matches provider name
				if TypeBelongsToProvider(dsName, localRef) {
					depKeys := schema.DependencyKeys{
						Labels: []schema.LabelDependent{
							{Index: 0, Value: dsName},
						},
					}
					mergedSchema.Blocks["data"].DependentBody[schema.NewSchemaKey(depKeys)] = dsSchema
					if _, ok := mergedSchema.Blocks["check"]; ok {
						mergedSchema.Blocks["check"].Body.Blocks["data"].DependentBody[schema.NewSchemaKey(depKeys)] = dsSchema
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

	declared, err := m.stateReader.DeclaredModuleCalls(meta.Path)
	if err != nil {
		return mergedSchema, nil
	}

	for _, module := range declared {
		depKeys := schema.DependencyKeys{
			// Fetching based only on the source can cause conflicts for multiple versions of the same module
			// specially if they have different versions or the source of those modules have been modified
			// inside the .terraform folder. This is a compromise that we made in this moment since it would impact only auto completion
			Attributes: []schema.AttributeDependent{
				{
					Name: "source",
					Expr: schema.ExpressionValue{
						Static: cty.StringVal(module.RawSourceAddr),
					},
				},
			},
		}

		switch sourceAddr := module.SourceAddr.(type) {
		case tfaddr.Module:
			// 1. See if we have a local installation of the module available
			installedDir, ok := m.stateReader.InstalledModulePath(meta.Path, sourceAddr.String())
			if ok {
				path := filepath.Join(meta.Path, installedDir)

				modMeta, err := m.stateReader.LocalModuleMeta(path)
				if err == nil {
					depSchema, err := schemaForDependentModuleBlock(module, modMeta)
					if err == nil {
						mergedSchema.Blocks["module"].DependentBody[schema.NewSchemaKey(depKeys)] = depSchema
					}

					// We continue here, so we don't end up overwriting the schema with one from the registry
					continue
				}
			}

			// 2. See if we have fetched the module schema from the registry
			modMeta, err := m.stateReader.RegistryModuleMeta(sourceAddr, module.Version)
			if err != nil {
				continue
			}

			depSchema, err := schemaForDependentRegistryModuleBlock(module, modMeta)
			if err == nil {
				mergedSchema.Blocks["module"].DependentBody[schema.NewSchemaKey(depKeys)] = depSchema
			}

		case tfmod.RemoteSourceAddr:
			installedDir, ok := m.stateReader.InstalledModulePath(meta.Path, sourceAddr.String())
			if !ok {
				continue
			}
			path := filepath.Join(meta.Path, installedDir)

			modMeta, err := m.stateReader.LocalModuleMeta(path)
			if err == nil {
				depSchema, err := schemaForDependentModuleBlock(module, modMeta)
				if err == nil {
					mergedSchema.Blocks["module"].DependentBody[schema.NewSchemaKey(depKeys)] = depSchema
				}
			}

		case tfmod.LocalSourceAddr:
			path := filepath.Join(meta.Path, sourceAddr.String())

			modMeta, err := m.stateReader.LocalModuleMeta(path)
			if err == nil {
				depSchema, err := schemaForDependentModuleBlock(module, modMeta)
				if err == nil {
					mergedSchema.Blocks["module"].DependentBody[schema.NewSchemaKey(depKeys)] = depSchema
				}
			}
		}
	}

	return mergedSchema, nil
}

// TypeBelongsToProvider returns true if the given type
// (resource or data source) name belongs to a particular provider.
//
// This reflects internal implementation in Terraform at
// https://github.com/hashicorp/terraform/blob/488bbd80/internal/addrs/resource.go#L68-L77
func TypeBelongsToProvider(typeName string, pRef tfmod.ProviderRef) bool {
	return typeName == pRef.LocalName || strings.HasPrefix(typeName, pRef.LocalName+"_")
}

func variableDependentBody(vars map[string]tfmod.Variable) map[schema.SchemaKey]*schema.BodySchema {
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
					Constraint:  schema.LiteralType{Type: mVar.Type},
					IsOptional:  true,
					Description: lang.Markdown("Default value to use when variable is not explicitly set"),
				},
			},
		}
	}

	return depBodies
}

type ProviderReferences map[tfmod.ProviderRef]tfaddr.Provider

func (pr ProviderReferences) ReferencesOfProvider(addr tfaddr.Provider) []tfmod.ProviderRef {
	refs := make([]tfmod.ProviderRef, 0)

	for ref, pAddr := range pr {
		if pAddr.Equals(addr) {
			refs = append(refs, ref)
		}
	}

	return refs
}
