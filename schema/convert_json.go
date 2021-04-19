package schema

import (
	"fmt"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-registry-address"
	"github.com/zclconf/go-cty/cty"
)

func ProviderSchemaFromJson(jsonSchema *tfjson.ProviderSchema, pAddr tfaddr.Provider) *ProviderSchema {
	ps := &ProviderSchema{
		Resources:   map[string]*schema.BodySchema{},
		DataSources: map[string]*schema.BodySchema{},
	}

	if jsonSchema.ConfigSchema != nil {
		ps.Provider = bodySchemaFromJson(jsonSchema.ConfigSchema.Block)
		ps.Provider.Detail = detailForSrcAddr(pAddr, nil)
		ps.Provider.DocsLink = docsLinkForProvider(pAddr, nil)
	}

	for rName, rSchema := range jsonSchema.ResourceSchemas {
		ps.Resources[rName] = bodySchemaFromJson(rSchema.Block)
		ps.Resources[rName].Detail = detailForSrcAddr(pAddr, nil)
	}

	for dsName, dsSchema := range jsonSchema.DataSourceSchemas {
		ps.DataSources[dsName] = bodySchemaFromJson(dsSchema.Block)
		ps.DataSources[dsName].Detail = detailForSrcAddr(pAddr, nil)
	}

	return ps
}

func (ps *ProviderSchema) SetProviderVersion(pAddr tfaddr.Provider, v *version.Version) {
	if ps.Provider != nil {
		ps.Provider.Detail = detailForSrcAddr(pAddr, v)
		ps.Provider.DocsLink = docsLinkForProvider(pAddr, v)
	}
	for _, rSchema := range ps.Resources {
		rSchema.Detail = detailForSrcAddr(pAddr, v)
	}
	for _, dsSchema := range ps.DataSources {
		dsSchema.Detail = detailForSrcAddr(pAddr, v)
	}
}

func bodySchemaFromJson(schemaBlock *tfjson.SchemaBlock) *schema.BodySchema {
	if schemaBlock == nil {
		s := schema.NewBodySchema()
		return s
	}

	return &schema.BodySchema{
		Attributes:   convertAttributesFromJson(schemaBlock.Attributes),
		Blocks:       convertBlocksFromJson(schemaBlock.NestedBlocks),
		IsDeprecated: schemaBlock.Deprecated,
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
		case tfjson.SchemaNestingModeSingle:
			blockType = schema.BlockTypeObject
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
			Body:         bodySchemaFromJson(block),
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

func docsLinkForProvider(addr tfaddr.Provider, v *version.Version) *schema.DocsLink {
	if !providerHasDocs(addr) {
		return nil
	}

	ver := "latest"
	if v != nil {
		ver = v.String()
	}

	return &schema.DocsLink{
		URL: fmt.Sprintf("https://registry.terraform.io/providers/%s/%s/%s/docs",
			addr.Namespace, addr.Type, ver),
		Tooltip: fmt.Sprintf("%s Documentation", addr.ForDisplay()),
	}
}

func providerHasDocs(addr tfaddr.Provider) bool {
	if addr.IsBuiltIn() {
		// Ideally this should point to versioned TF core docs
		// but there aren't any for the built-in provider yet
		return false
	}
	if addr.Hostname != "registry.terraform.io" {
		// docs URLs outside of the official Registry aren't standardized yet
		return false
	}
	return true
}

func detailForSrcAddr(addr tfaddr.Provider, v *version.Version) string {
	if addr.IsBuiltIn() {
		if v == nil {
			return "(builtin)"
		}
		return fmt.Sprintf("(builtin %s)", v.String())
	}

	detail := addr.ForDisplay()
	if v != nil {
		detail += " " + v.String()
	}

	return detail
}
