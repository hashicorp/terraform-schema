package backends

import (
	"sort"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

var (
	v0_12_2  = version.Must(version.NewVersion("0.12.2"))
	v0_12_4  = version.Must(version.NewVersion("0.12.4"))
	v0_12_6  = version.Must(version.NewVersion("0.12.6"))
	v0_12_8  = version.Must(version.NewVersion("0.12.8"))
	v0_12_10 = version.Must(version.NewVersion("0.12.10"))
	v0_12_14 = version.Must(version.NewVersion("0.12.14"))
	v0_13_0  = version.Must(version.NewVersion("0.13.0"))
	v0_13_1  = version.Must(version.NewVersion("0.13.1"))
	v0_14_0  = version.Must(version.NewVersion("0.14.0"))
	v0_15_0  = version.Must(version.NewVersion("0.15.0"))
)

func BackendTypesAsExprConstraints(tfVersion *version.Version) schema.ExprConstraints {
	ec := make(schema.ExprConstraints, 0)

	for backendType, bs := range backendBodySchemas(tfVersion) {
		lv := schema.LiteralValue{
			Val:          cty.StringVal(backendType),
			IsDeprecated: bs.IsDeprecated,
		}
		if bs != nil {
			lv.Description = bs.Description
		}
		ec = append(ec, lv)
	}

	sort.SliceStable(ec, func(i, j int) bool {
		leftVal := ec[i].(schema.LiteralValue)
		rightVal := ec[j].(schema.LiteralValue)
		return leftVal.Val.AsString() < rightVal.Val.AsString()
	})

	return ec
}

func ConfigsAsExprConstraints(tfVersion *version.Version) map[string]schema.ExprConstraints {
	ecs := make(map[string]schema.ExprConstraints, 0)

	for backendType, bs := range backendBodySchemas(tfVersion) {
		ecs[backendType] = schema.ExprConstraints{
			objectExprFromBodySchema(bs),
		}
	}

	return ecs
}

func ConfigsAsDependentBodies(tfVersion *version.Version) map[schema.SchemaKey]*schema.BodySchema {
	depBodies := make(map[schema.SchemaKey]*schema.BodySchema, 0)

	for backendType, bodySchema := range backendBodySchemas(tfVersion) {
		depBodies[labelKey(backendType)] = bodySchema
	}

	return depBodies
}

func labelKey(value string) schema.SchemaKey {
	return schema.NewSchemaKey(schema.DependencyKeys{
		Labels: []schema.LabelDependent{{Index: 0, Value: value}},
	})
}

func backendBodySchemas(v *version.Version) map[string]*schema.BodySchema {
	if v == nil {
		return map[string]*schema.BodySchema{}
	}

	v = v.Core()

	// https://github.com/hashicorp/terraform/blob/v0.12.0/backend/init/init.go
	backends := map[string]*schema.BodySchema{
		// Enhanced backends
		"local":  localBackend(v),
		"remote": remoteBackend(v),

		// Remote State backends
		"artifactory": artifactoryBackend(v),
		"azurerm":     azureRmBackend(v),
		"consul":      consulBackend(v),
		"etcd":        etcdv2Backend(v),
		"etcdv3":      etcdv3Backend(v),
		"gcs":         gcsBackend(v),
		"http":        httpBackend(v),
		"manta":       mantaBackend(v),
		"pg":          pgBackend(v),
		"s3":          s3Backend(v),
		"swift":       swiftBackend(v),

		// Deprecated backends
		"atlas": {
			IsDeprecated: true,
			Description:  lang.Markdown("`atlas` backend is **DEPRECATED**, please use `remote` instead"),
		},
		"azure": {
			IsDeprecated: true,
			Description:  lang.Markdown("`azure` name is **DEPRECATED**, please use `azurerm` instead"),
		},
	}

	if v.GreaterThanOrEqual(v0_12_2) {
		// https://github.com/hashicorp/terraform/commit/b887d447
		backends["oss"] = ossBackend(v)
	}

	if v.GreaterThanOrEqual(v0_13_0) {
		// https://github.com/hashicorp/terraform/commit/76e5b446
		backends["cos"] = cosBackend(v)
		// https://github.com/hashicorp/terraform/commit/23fb8f6d
		backends["kubernetes"] = kubernetesBackend(v)
	}

	if v.GreaterThanOrEqual(v0_15_0) {
		// https://github.com/hashicorp/terraform/commit/b8e3b803
		delete(backends, "atlas")
	}

	return backends
}
