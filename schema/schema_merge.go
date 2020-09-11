package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-schema/internal/addrs"
	"github.com/hashicorp/terraform-schema/internal/refdecoder"
)

// MergeCoreWithJsonProviderSchemas provides a merged schema based on provided
// parsed files, core schema and terraform-json formatted provider schema
func MergeCoreWithJsonProviderSchemas(m map[string]*hcl.File, coreSchema *schema.BodySchema, ps *tfjson.ProviderSchemas) (
	*schema.BodySchema, error) {

	mergedSchema := coreSchema

	if mergedSchema.Blocks["provider"].DependentBody == nil {
		mergedSchema.Blocks["provider"].DependentBody = make(map[schema.SchemaKey]*schema.BodySchema)
	}
	if mergedSchema.Blocks["resource"].DependentBody == nil {
		mergedSchema.Blocks["resource"].DependentBody = make(map[schema.SchemaKey]*schema.BodySchema)
	}
	if mergedSchema.Blocks["data"].DependentBody == nil {
		mergedSchema.Blocks["data"].DependentBody = make(map[schema.SchemaKey]*schema.BodySchema)
	}

	refs, err := refdecoder.DecodeProviderReferences(m)
	if err != nil {
		return coreSchema, err
	}

	for sourceString, provider := range ps.Schemas {
		srcAddr, err := addrs.ParseProviderSourceString(sourceString)
		if err != nil {
			return coreSchema, err
		}

		localRefs := refs.LocalNamesByAddr(srcAddr)

		if len(localRefs) == 0 && (srcAddr.IsBuiltIn() || srcAddr.IsLegacy() || srcAddr.IsDefault()) {
			// Assume this provider does not have alias
			localRefs = append(localRefs, addrs.LocalProviderConfig{
				LocalName: srcAddr.Type,
			})
		}

		var providerSchema *tfjson.SchemaBlock
		if provider.ConfigSchema != nil {
			providerSchema = provider.ConfigSchema.Block
		}

		detail := srcAddr.ForDisplay()
		if srcAddr.IsBuiltIn() {
			detail = "(builtin)"
		}

		for _, localRef := range localRefs {
			mergedSchema.Blocks["provider"].DependentBody[schema.NewSchemaKey(schema.DependencyKeys{
				Labels: []schema.LabelDependent{
					{Index: 0, Value: localRef.LocalName},
				},
			})] = convertBodySchemaFromJson(detail, providerSchema)

			for rName, rJsonSchema := range provider.ResourceSchemas {
				rSchema := convertBodySchemaFromJson(detail, rJsonSchema.Block)

				depKeys := schema.DependencyKeys{
					Labels: []schema.LabelDependent{
						{Index: 0, Value: rName},
					},
				}
				if localRef.Alias != "" {
					depKeys.Attributes = append(depKeys.Attributes, schema.AttributeDependent{
						Name: "provider",
						Expr: schema.ExpressionValue{
							Reference: lang.Reference{
								lang.RootStep{Name: localRef.LocalName},
								lang.AttrStep{Name: localRef.Alias},
							},
						},
					})
				}

				mergedSchema.Blocks["resource"].DependentBody[schema.NewSchemaKey(depKeys)] = rSchema
			}

			for dsName, dsJsonSchema := range provider.DataSourceSchemas {
				dsSchema := convertBodySchemaFromJson(detail, dsJsonSchema.Block)

				depKeys := schema.DependencyKeys{
					Labels: []schema.LabelDependent{
						{Index: 0, Value: dsName},
					},
				}
				if localRef.Alias != "" {
					depKeys.Attributes = append(depKeys.Attributes, schema.AttributeDependent{
						Name: "provider",
						Expr: schema.ExpressionValue{
							Reference: lang.Reference{
								lang.RootStep{Name: localRef.LocalName},
								lang.AttrStep{Name: localRef.Alias},
							},
						},
					})
				}

				mergedSchema.Blocks["data"].DependentBody[schema.NewSchemaKey(depKeys)] = dsSchema
			}
		}
	}

	return mergedSchema, nil
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
			IsReadOnly:   (attr.Computed && !attr.Optional && !attr.Required),
			IsRequired:   attr.Required,
			ValueType:    attr.AttributeType,
		}
	}
	return cAttrs
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
	case "plain": // discrepancy between terraform-json & Terraform core
		return lang.PlainText(value)
	}

	// backwards compatibility with v0.12
	return lang.PlainText(value)
}
