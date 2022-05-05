package schema

import (
	"fmt"
	"sort"

	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/hcl/v2"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/hashicorp/terraform-schema/internal/schema/refscope"
	"github.com/hashicorp/terraform-schema/module"
	"github.com/zclconf/go-cty/cty"
)

func schemaForDependentModuleBlock(module module.InstalledModuleCall, modMeta *module.Meta) (*schema.BodySchema, error) {
	attributes := make(map[string]*schema.AttributeSchema, 0)

	for name, modVar := range modMeta.Variables {
		aSchema := moduleVarToAttribute(modVar)
		varType := typeOfModuleVar(modVar)
		aSchema.Expr = convertAttributeTypeToExprConstraints(varType)
		aSchema.OriginForTarget = &schema.PathTarget{
			Address: schema.Address{
				schema.StaticStep{Name: "var"},
				schema.AttrNameStep{},
			},
			Path: lang.Path{
				Path:       modMeta.Path,
				LanguageID: ModuleLanguageID,
			},
			Constraints: schema.Constraints{
				ScopeId: refscope.VariableScope,
				Type:    varType,
			},
		}

		attributes[name] = aSchema
	}

	bodySchema := &schema.BodySchema{
		Attributes: attributes,
	}

	if module.LocalName == "" {
		// avoid creating output refs if we don't have reference name
		return bodySchema, nil
	}

	modOutputTypes := make(map[string]cty.Type, 0)
	modOutputVals := make(map[string]cty.Value, 0)
	targetableOutputs := make(schema.Targetables, 0)

	for name, output := range modMeta.Outputs {
		addr := lang.Address{
			lang.RootStep{Name: "module"},
			lang.AttrStep{Name: module.LocalName},
			lang.AttrStep{Name: name},
		}

		typ := cty.DynamicPseudoType
		if !output.Value.IsNull() {
			typ = output.Value.Type()
		}

		targetable := &schema.Targetable{
			Address:           addr,
			ScopeId:           refscope.ModuleScope,
			AsType:            typ,
			IsSensitive:       output.IsSensitive,
			NestedTargetables: schema.NestedTargetablesForValue(addr, refscope.ModuleScope, output.Value),
		}
		if output.Description != "" {
			targetable.Description = lang.PlainText(output.Description)
		}

		targetableOutputs = append(targetableOutputs, targetable)

		modOutputTypes[name] = typ
		modOutputVals[name] = output.Value
	}

	sort.Sort(targetableOutputs)

	addr := lang.Address{
		lang.RootStep{Name: "module"},
		lang.AttrStep{Name: module.LocalName},
	}
	bodySchema.TargetableAs = append(bodySchema.TargetableAs, &schema.Targetable{
		Address:           addr,
		ScopeId:           refscope.ModuleScope,
		AsType:            cty.Object(modOutputTypes),
		NestedTargetables: targetableOutputs,
	})

	if len(modMeta.Filenames) > 0 {
		filename := modMeta.Filenames[0]

		// Prioritize main.tf based on best practices as documented at
		// https://learn.hashicorp.com/tutorials/terraform/module-create
		if sliceContains(modMeta.Filenames, "main.tf") {
			filename = "main.tf"
		}

		bodySchema.Targets = &schema.Target{
			Path: lang.Path{
				Path:       modMeta.Path,
				LanguageID: "terraform",
			},
			Range: hcl.Range{
				Filename: filename,
				Start:    hcl.InitialPos,
				End:      hcl.InitialPos,
			},
		}
	}

	moduleSourceRegistry, err := tfaddr.ParseRawModuleSourceRegistry(module.SourceAddr)
	if err == nil && moduleSourceRegistry.PackageAddr.Host == "registry.terraform.io" {
		versionStr := ""
		if module.Version == nil {
			versionStr = "latest"
		} else {
			versionStr = module.Version.String()
		}

		bodySchema.DocsLink = &schema.DocsLink{
			URL: fmt.Sprintf(
				`https://registry.terraform.io/modules/%s/%s`,
				moduleSourceRegistry.PackageAddr.ForRegistryProtocol(),
				versionStr,
			),
		}
	}

	return bodySchema, nil
}

func sliceContains(slice []string, value string) bool {
	for _, val := range slice {
		if val == value {
			return true
		}
	}
	return false
}
