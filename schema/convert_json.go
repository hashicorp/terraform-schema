// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"fmt"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	tfjson "github.com/hashicorp/terraform-json"
	tfaddr "github.com/hashicorp/terraform-registry-address"
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
			Constraint:   exprConstraintFromSchemaAttribute(attr),
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
			Constraint: convertAttributeTypeToConstraint(attrType),
			IsOptional: true,
		}
	}
	return ret
}

func exprConstraintFromSchemaAttribute(attr *tfjson.SchemaAttribute) schema.Constraint {
	if attr.AttributeType != cty.NilType {
		return convertAttributeTypeToConstraint(attr.AttributeType)
	}
	if attr.AttributeNestedType != nil {
		switch attr.AttributeNestedType.NestingMode {
		case tfjson.SchemaNestingModeSingle:
			return convertJsonAttributesToObjectConstraint(attr.AttributeNestedType.Attributes)
		case tfjson.SchemaNestingModeList:
			return schema.List{
				Elem:     convertJsonAttributesToObjectConstraint(attr.AttributeNestedType.Attributes),
				MinItems: attr.AttributeNestedType.MinItems,
				MaxItems: attr.AttributeNestedType.MaxItems,
			}
		case tfjson.SchemaNestingModeSet:
			return schema.Set{
				Elem:     convertJsonAttributesToObjectConstraint(attr.AttributeNestedType.Attributes),
				MinItems: attr.AttributeNestedType.MinItems,
				MaxItems: attr.AttributeNestedType.MaxItems,
			}
		case tfjson.SchemaNestingModeMap:
			return schema.Map{
				Elem:     convertJsonAttributesToObjectConstraint(attr.AttributeNestedType.Attributes),
				MinItems: attr.AttributeNestedType.MinItems,
				MaxItems: attr.AttributeNestedType.MaxItems,
			}
		}
	}
	return nil
}

func convertAttributeTypeToConstraint(attrType cty.Type) schema.Constraint {
	if attrType.IsPrimitiveType() || attrType == cty.DynamicPseudoType {
		return schema.AnyExpression{OfType: attrType}
	}

	cons := schema.OneOf{
		schema.AnyExpression{
			OfType:                  attrType,
			SkipLiteralComplexTypes: true,
		},
	}

	if attrType.IsListType() {
		cons = append(cons, schema.List{
			Elem: convertAttributeTypeToConstraint(attrType.ElementType()),
		})
	}
	if attrType.IsSetType() {
		cons = append(cons, schema.Set{
			Elem: convertAttributeTypeToConstraint(attrType.ElementType()),
		})
	}
	if attrType.IsTupleType() {
		te := schema.Tuple{Elems: make([]schema.Constraint, 0)}
		for _, elemType := range attrType.TupleElementTypes() {
			te.Elems = append(te.Elems, convertAttributeTypeToConstraint(elemType))
		}
		cons = append(cons, te)
	}
	if attrType.IsMapType() {
		cons = append(cons, schema.Map{
			Elem: convertAttributeTypeToConstraint(attrType.ElementType()),
		})
	}
	if attrType.IsObjectType() {
		cons = append(cons, convertCtyObjectToObjectCons(attrType))
	}

	return cons
}

func convertCtyObjectToObjectCons(obj cty.Type) schema.Object {
	attrTypes := obj.AttributeTypes()
	attributes := make(schema.ObjectAttributes, len(attrTypes))
	for name, attrType := range attrTypes {
		aSchema := &schema.AttributeSchema{
			Constraint: convertAttributeTypeToConstraint(attrType),
		}

		if obj.AttributeOptional(name) {
			aSchema.IsOptional = true
		} else {
			aSchema.IsRequired = true
		}

		attributes[name] = aSchema
	}
	return schema.Object{
		Attributes: attributes,
	}
}

func convertJsonAttributesToObjectConstraint(attrs map[string]*tfjson.SchemaAttribute) schema.Object {
	attributes := make(schema.ObjectAttributes, len(attrs))
	for name, attr := range attrs {
		attributes[name] = &schema.AttributeSchema{
			Description:  markupContent(attr.Description, attr.DescriptionKind),
			IsDeprecated: attr.Deprecated,
			IsComputed:   attr.Computed,
			IsOptional:   attr.Optional,
			IsRequired:   attr.Required,
			Constraint:   exprConstraintFromSchemaAttribute(attr),
		}
	}
	return schema.Object{
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
