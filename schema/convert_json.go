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
		ps.Provider.HoverURL = urlForProvider(pAddr, nil)
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

func bodySchemaFromJson(schemaBlock *tfjson.SchemaBlock) *schema.BodySchema {
	if schemaBlock == nil {
		s := schema.NewBodySchema()
		return s
	}

	attributes := convertAttributesFromJson(schemaBlock.Attributes)

	// Attributes and block types of the same name should never occur
	// in providers which use the official plugin SDK but we give chance
	// for real blocks to override the "converted" ones just in case
	blocks := convertibleAttributesToBlocks(schemaBlock.Attributes)
	for name, block := range convertBlocksFromJson(schemaBlock.NestedBlocks) {
		blocks[name] = block
	}

	return &schema.BodySchema{
		Attributes:   attributes,
		Blocks:       blocks,
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
			IsSensitive:  attr.Sensitive,
			Expr:         exprConstraintsFromAttribute(attr),
		}
	}
	return cAttrs
}

// convertibleAttributesToBlocks is responsible for mimicking
// Terraform's builtin backwards-compatible logic where
// list(object) or set(object) attributes are also accessible
// as blocks.
// See https://github.com/hashicorp/terraform/blob/v1.0.3/internal/lang/blocktoattr/schema.go
func convertibleAttributesToBlocks(attributes map[string]*tfjson.SchemaAttribute) map[string]*schema.BlockSchema {
	blocks := make(map[string]*schema.BlockSchema, 0)

	for name, attr := range attributes {
		if typeCanBeBlocks(attr.AttributeType) {
			blockSchema, ok := blockSchemaForAttribute(attr)
			if !ok {
				continue
			}
			blocks[name] = blockSchema
		}
	}

	return blocks
}

// typeCanBeBlocks returns true if the given type is a list-of-object or
// set-of-object type, and would thus be subject to the blocktoattr fixup
// if used as an attribute type.
func typeCanBeBlocks(ty cty.Type) bool {
	return (ty.IsListType() || ty.IsSetType()) && ty.ElementType().IsObjectType()
}

func blockSchemaForAttribute(attr *tfjson.SchemaAttribute) (*schema.BlockSchema, bool) {
	if attr.AttributeType == cty.NilType {
		return nil, false
	}

	blockType := schema.BlockTypeNil
	switch {
	case attr.AttributeType.IsListType():
		blockType = schema.BlockTypeList
	case attr.AttributeType.IsSetType():
		blockType = schema.BlockTypeSet
	default:
		return nil, false
	}

	minItems := uint64(0)
	if attr.Required {
		minItems = 1
	}

	return &schema.BlockSchema{
		Description:  markupContent(attr.Description, attr.DescriptionKind),
		Type:         blockType,
		IsDeprecated: attr.Deprecated,
		MinItems:     minItems,
		Body:         bodySchemaForCtyObjectType(attr.AttributeType.ElementType()),
	}, true
}

func bodySchemaForCtyObjectType(typ cty.Type) *schema.BodySchema {
	if !typ.IsObjectType() {
		return nil
	}

	attrTypes := typ.AttributeTypes()
	ret := &schema.BodySchema{
		Attributes: make(map[string]*schema.AttributeSchema, len(attrTypes)),
	}
	for name, attrType := range attrTypes {
		ret.Attributes[name] = &schema.AttributeSchema{
			Expr:       convertAttributeTypeToExprConstraints(attrType),
			IsOptional: true,
		}
	}
	return ret
}

func exprConstraintsFromAttribute(attr *tfjson.SchemaAttribute) schema.ExprConstraints {
	var expr schema.ExprConstraints
	if attr.AttributeType != cty.NilType {
		return convertAttributeTypeToExprConstraints(attr.AttributeType)
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

func convertAttributeTypeToExprConstraints(attrType cty.Type) schema.ExprConstraints {
	ec := schema.ExprConstraints{
		schema.TraversalExpr{OfType: attrType},
		schema.LiteralTypeExpr{Type: attrType},
	}

	if attrType.IsListType() {
		ec = append(ec, schema.ListExpr{
			Elem: convertAttributeTypeToExprConstraints(attrType.ElementType()),
		})
	}
	if attrType.IsSetType() {
		ec = append(ec, schema.SetExpr{
			Elem: convertAttributeTypeToExprConstraints(attrType.ElementType()),
		})
	}
	if attrType.IsTupleType() {
		te := schema.TupleExpr{Elems: make([]schema.ExprConstraints, 0)}
		for _, elemType := range attrType.TupleElementTypes() {
			te.Elems = append(te.Elems, convertAttributeTypeToExprConstraints(elemType))
		}
		ec = append(ec, te)
	}
	if attrType.IsMapType() {
		ec = append(ec, schema.MapExpr{
			Elem: convertAttributeTypeToExprConstraints(attrType.ElementType()),
		})
	}
	if attrType.IsObjectType() {
		ec = append(ec, convertCtyObjectToObjectExprAttr(attrType))
	}

	return ec
}

func convertCtyObjectToObjectExprAttr(obj cty.Type) schema.ObjectExpr {
	attrTypes := obj.AttributeTypes()
	attributes := make(schema.ObjectExprAttributes, len(attrTypes))
	for name, attrType := range attrTypes {
		aSchema := &schema.AttributeSchema{
			Expr: convertAttributeTypeToExprConstraints(attrType),
		}

		if obj.AttributeOptional(name) {
			aSchema.IsOptional = true
		} else {
			aSchema.IsRequired = true
		}

		attributes[name] = aSchema
	}
	return schema.ObjectExpr{
		Attributes: attributes,
	}
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

	return &schema.DocsLink{
		URL:     urlForProvider(addr, v),
		Tooltip: fmt.Sprintf("%s Documentation", addr.ForDisplay()),
	}
}

func urlForProvider(addr tfaddr.Provider, v *version.Version) string {
	if !providerHasDocs(addr) {
		return ""
	}

	ver := "latest"
	if v != nil {
		ver = v.String()
	}

	return fmt.Sprintf("https://registry.terraform.io/providers/%s/%s/%s/docs",
		addr.Namespace, addr.Type, ver)
}

func providerHasDocs(addr tfaddr.Provider) bool {
	if addr.IsBuiltIn() {
		// Ideally this should point to versioned TF core docs
		// but there aren't any for the built-in provider yet
		return false
	}
	if addr.IsLegacy() {
		// The Registry does know where legacy providers live
		// but it doesn't provide stable (legacy) URLs
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
