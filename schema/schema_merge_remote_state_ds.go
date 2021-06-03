package schema

import (
	"strings"

	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	tfaddr "github.com/hashicorp/terraform-registry-address"
	"github.com/hashicorp/terraform-schema/internal/schema/backends"
	"github.com/hashicorp/terraform-schema/module"
	"github.com/zclconf/go-cty/cty"
)

const remoteStateDsName = "terraform_remote_state"

func isRemoteStateDataSource(pAddr tfaddr.Provider, dsName string) bool {
	return (pAddr.Equals(tfaddr.NewBuiltInProvider("terraform")) ||
		pAddr.Equals(tfaddr.NewDefaultProvider("terraform")) ||
		pAddr.Equals(tfaddr.NewLegacyProvider("terraform"))) &&
		dsName == remoteStateDsName
}

func (sm *SchemaMerger) dependentBodyForRemoteStateDataSource(providerAddr lang.Address, localRef module.ProviderRef) map[schema.SchemaKey]*schema.BodySchema {
	m := make(map[schema.SchemaKey]*schema.BodySchema, 0)
	bExprConstraints := backends.ConfigsAsExprConstraints(sm.terraformVersion)

	for backendType, exprConstraints := range bExprConstraints {
		depKeys := schema.DependencyKeys{
			Labels: []schema.LabelDependent{
				{Index: 0, Value: remoteStateDsName},
			},
			Attributes: []schema.AttributeDependent{
				{
					Name: "provider",
					Expr: schema.ExpressionValue{
						Address: providerAddr,
					},
				},
				{
					Name: "backend",
					Expr: schema.ExpressionValue{
						Static: cty.StringVal(backendType),
					},
				},
			},
		}

		dsSchema := &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"config": {
					Expr:       exprConstraints,
					IsOptional: true,
				},
			},
		}

		m[schema.NewSchemaKey(depKeys)] = dsSchema

		// No explicit association is required
		// if the resource prefix matches provider name
		if strings.HasPrefix(remoteStateDsName, localRef.LocalName+"_") {
			depKeys := schema.DependencyKeys{
				Labels: []schema.LabelDependent{
					{Index: 0, Value: remoteStateDsName},
				},
				Attributes: []schema.AttributeDependent{
					{
						Name: "backend",
						Expr: schema.ExpressionValue{
							Static: cty.StringVal(backendType),
						},
					},
				},
			}
			m[schema.NewSchemaKey(depKeys)] = dsSchema
		}
	}

	return m
}
