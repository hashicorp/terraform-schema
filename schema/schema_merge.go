package schema

import (
	"fmt"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/hcl/v2"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-schema/internal/addrs"
	"github.com/hashicorp/terraform-schema/internal/refdecoder"
)

type SchemaMerger struct {
	coreSchema  *schema.BodySchema
	parsedFiles map[string]*hcl.File

	coreVersion      *version.Version
	providerVersions map[addrs.Provider]*version.Version
}

func NewSchemaMerger(coreSchema *schema.BodySchema) *SchemaMerger {
	return &SchemaMerger{
		coreSchema:       coreSchema,
		parsedFiles:      make(map[string]*hcl.File, 0),
		providerVersions: make(map[addrs.Provider]*version.Version, 0),
	}
}

// SetParsedFiles sets a map of parsed files where key is a filename
func (m *SchemaMerger) SetParsedFiles(files map[string]*hcl.File) {
	m.parsedFiles = files
}

// SetCoreVersion sets version of Terraform (core) to help identify core schema
// and schema of the builtin terraform provider
func (m *SchemaMerger) SetCoreVersion(v *version.Version) {
	m.coreVersion = v
}

// SetProviderVersions sets versions of providers to help identify
// where the provider schemas came from
func (m *SchemaMerger) SetProviderVersions(versions map[string]*version.Version) error {
	versionMap := make(map[addrs.Provider]*version.Version, 0)

	for addr, ver := range versions {
		srcAddr, err := addrs.ParseProviderSourceString(addr)
		if err != nil {
			return err
		}
		versionMap[srcAddr] = ver
	}

	m.providerVersions = versionMap

	return nil
}

// MergeWithJsonProviderSchemas provides a merged schema based on
// terraform-json formatted provider schema and any other data
// provided via setters
func (m *SchemaMerger) MergeWithJsonProviderSchemas(ps *tfjson.ProviderSchemas) (*schema.BodySchema, error) {
	if m.coreSchema == nil {
		return nil, coreSchemaRequiredErr{}
	}

	if ps == nil {
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

	refs, diags := refdecoder.DecodeProviderReferences(m.parsedFiles)
	if diags.HasErrors() && len(refs) == 0 {
		return m.coreSchema, nil
	}

	for sourceString, provider := range ps.Schemas {
		srcAddr, err := addrs.ParseProviderSourceString(sourceString)
		if err != nil {
			return m.coreSchema, nil
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

		detail := m.detailForSrcAddr(srcAddr)

		for _, localRef := range localRefs {
			pSchema := convertBodySchemaFromJson(detail, providerSchema)
			pSchema.DocsLink = m.docsLinkForProvider(srcAddr)

			mergedSchema.Blocks["provider"].DependentBody[schema.NewSchemaKey(schema.DependencyKeys{
				Labels: []schema.LabelDependent{
					{Index: 0, Value: localRef.LocalName},
				},
			})] = pSchema

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

func (m *SchemaMerger) docsLinkForProvider(addr addrs.Provider) *schema.DocsLink {
	if !providerHasDocs(addr) {
		return nil
	}

	version := "latest"
	for pAddr, ver := range m.providerVersions {
		if addr.Equals(pAddr) {
			version = ver.String()
		}
	}

	return &schema.DocsLink{
		URL: fmt.Sprintf("https://registry.terraform.io/providers/%s/%s/%s/docs",
			addr.Namespace, addr.Type, version),
		Tooltip: fmt.Sprintf("%s Documentation", addr.ForDisplay()),
	}
}

func providerHasDocs(addr addrs.Provider) bool {
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

func (m *SchemaMerger) detailForSrcAddr(addr addrs.Provider) string {
	if addr.IsBuiltIn() {
		if m.coreVersion == nil {
			return "(builtin)"
		}
		return fmt.Sprintf("(builtin %s)", m.coreVersion.String())
	}

	detail := addr.ForDisplay()
	for pAddr, ver := range m.providerVersions {
		if addr.Equals(pAddr) {
			detail += " " + ver.String()
			break
		}
	}

	return detail
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
	}

	// backwards compatibility with v0.12
	return lang.PlainText(value)
}
