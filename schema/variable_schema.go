package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/terraform-schema/module"
)

func SchemaForVariables(vars map[string]module.Variable) (*schema.BodySchema, error) {
	varSchemas := make(map[string]*schema.AttributeSchema)

	for name, v := range vars {
		varSchemas[name] = &schema.AttributeSchema{
			Description: lang.MarkupContent{
				Value: v.Description,
				Kind:  lang.PlainTextKind,
			},
			Expr:        schema.ExprConstraints{schema.LiteralTypeExpr{Type: v.Type}},
			IsSensitive: v.IsSensitive,
			IsRequired:  v.IsRequired,
		}
	}

	return &schema.BodySchema{
		Attributes: varSchemas,
	}, nil
}
