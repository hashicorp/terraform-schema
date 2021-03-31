package schema

import (
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-registry-address"
	"github.com/hashicorp/terraform-schema/module"
	"github.com/zclconf/go-cty/cty"
)

type SchemaMerger struct {
	coreSchema   *schema.BodySchema
	schemaReader SchemaReader
}

type ProviderSchema struct {
	Provider    *schema.BodySchema
	Resources   map[string]*schema.BodySchema
	DataSources map[string]*schema.BodySchema
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

func (m *SchemaMerger) SchemaForModule(meta *module.Meta) (*schema.BodySchema, error) {
	if m.coreSchema == nil {
		return nil, coreSchemaRequiredErr{}
	}

	if meta == nil || m.schemaReader == nil {
		return m.coreSchema, nil
	}

	mergedSchema := m.coreSchema

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

	return mergedSchema, nil
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

func convertBodySchemaFromJson(detail string, schemaBlock *tfjson.SchemaBlock) *schema.BodySchema {
	if schemaBlock == nil {
		s := schema.NewBodySchema()
		s.Detail = detail
		return s
	}

	return &schema.BodySchema{
		Attributes:   convertAttributesFromJson(schemaBlock.Attributes),
		Blocks:       convertBlocksFromJson(schemaBlock.NestedBlocks),
		IsDeprecated: schemaBlock.Deprecated,
		Detail:       detail,
		Description:  markupContent(schemaBlock.Description, schemaBlock.DescriptionKind),
	}
}

func convertBlocksFromJson(blocks map[string]*tfjson.SchemaBlockType) map[string]*schema.BlockSchema {
	cBlocks := make(map[string]*schema.BlockSchema, len(blocks))
	for name, jsonSchema := range blocks {
		block := jsonSchema.Block

		blockType := schema.BlockTypeNil
		labels := []*schema.LabelSchema{}

		switch jsonSchema.NestingMode {
		case tfjson.SchemaNestingModeMap:
			labels = []*schema.LabelSchema{
				{Name: "name"},
			}
			blockType = schema.BlockTypeMap
		case tfjson.SchemaNestingModeList:
			blockType = schema.BlockTypeList
		case tfjson.SchemaNestingModeSet:
			blockType = schema.BlockTypeSet
		}

		cBlocks[name] = &schema.BlockSchema{
			Description:  markupContent(block.Description, block.DescriptionKind),
			Type:         blockType,
			IsDeprecated: block.Deprecated,
			MinItems:     jsonSchema.MinItems,
			MaxItems:     jsonSchema.MaxItems,
			Labels:       labels,
			Body:         convertBodySchemaFromJson("", block),
		}
	}
	return cBlocks
}

func convertAttributesFromJson(attributes map[string]*tfjson.SchemaAttribute) map[string]*schema.AttributeSchema {
	cAttrs := make(map[string]*schema.AttributeSchema, len(attributes))
	for name, attr := range attributes {
		cAttrs[name] = &schema.AttributeSchema{
			Description:  markupContent(attr.Description, attr.DescriptionKind),
			IsDeprecated: attr.Deprecated,
			IsComputed:   attr.Computed,
			IsOptional:   attr.Optional,
			IsRequired:   attr.Required,
			Expr:         exprConstraintsFromAttribute(attr),
		}
	}
	return cAttrs
}

func exprConstraintsFromAttribute(attr *tfjson.SchemaAttribute) schema.ExprConstraints {
	var expr schema.ExprConstraints
	if attr.AttributeType != cty.NilType {
		expr = schema.LiteralTypeOnly(attr.AttributeType)
	}
	if attr.AttributeNestedType != nil {
		switch attr.AttributeNestedType.NestingMode {
		case tfjson.SchemaNestingModeSingle:
			return schema.ExprConstraints{
				convertJsonAttributesToObjectExprAttr(attr.AttributeNestedType.Attributes),
			}
		case tfjson.SchemaNestingModeList:
			return schema.ExprConstraints{
				schema.ListExpr{
					Elem: schema.ExprConstraints{
						convertJsonAttributesToObjectExprAttr(attr.AttributeNestedType.Attributes),
					},
					MinItems: attr.AttributeNestedType.MinItems,
					MaxItems: attr.AttributeNestedType.MaxItems,
				},
			}
		case tfjson.SchemaNestingModeSet:
			return schema.ExprConstraints{
				schema.SetExpr{
					Elem: schema.ExprConstraints{
						convertJsonAttributesToObjectExprAttr(attr.AttributeNestedType.Attributes),
					},
					MinItems: attr.AttributeNestedType.MinItems,
					MaxItems: attr.AttributeNestedType.MaxItems,
				},
			}
		case tfjson.SchemaNestingModeMap:
			return schema.ExprConstraints{
				schema.MapExpr{
					Elem: schema.ExprConstraints{
						convertJsonAttributesToObjectExprAttr(attr.AttributeNestedType.Attributes),
					},
					MinItems: attr.AttributeNestedType.MinItems,
					MaxItems: attr.AttributeNestedType.MaxItems,
				},
			}
		}
	}
	return expr
}

func convertJsonAttributesToObjectExprAttr(attrs map[string]*tfjson.SchemaAttribute) schema.ObjectExpr {
	attributes := make(schema.ObjectExprAttributes, len(attrs))
	for name, attr := range attrs {
		attributes[name] = &schema.AttributeSchema{
			Description:  markupContent(attr.Description, attr.DescriptionKind),
			IsDeprecated: attr.Deprecated,
			IsComputed:   attr.Computed,
			IsOptional:   attr.Optional,
			IsRequired:   attr.Required,
			Expr:         exprConstraintsFromAttribute(attr),
		}
	}
	return schema.ObjectExpr{
		Attributes: attributes,
	}
}

func markupContent(value string, kind tfjson.SchemaDescriptionKind) lang.MarkupContent {
	if value == "" {
		return lang.MarkupContent{}
	}
	switch kind {
	case tfjson.SchemaDescriptionKindMarkdown:
		return lang.Markdown(value)
	case tfjson.SchemaDescriptionKindPlain:
		return lang.PlainText(value)
	}

	// backwards compatibility with v0.12
	return lang.PlainText(value)
}
