package backends

import "github.com/hashicorp/hcl-lang/schema"

func objectExprFromBodySchema(bs *schema.BodySchema) schema.ObjectExpr {
	if bs == nil {
		return schema.ObjectExpr{}
	}

	oe := schema.ObjectExpr{
		Description: bs.Description,
		Attributes:  bs.Attributes,
	}

	for bType, block := range bs.Blocks {
		oe.Attributes[bType] = &schema.AttributeSchema{
			Description:  block.Description,
			IsDeprecated: block.IsDeprecated,
		}

		if block.MinItems > 0 {
			oe.Attributes[bType].IsRequired = true
		} else {
			oe.Attributes[bType].IsOptional = true
		}

		switch block.Type {
		case schema.BlockTypeObject:
			oe.Attributes[bType].Expr = schema.ExprConstraints{
				objectExprFromBodySchema(block.Body),
			}
		case schema.BlockTypeList:
			oe.Attributes[bType].Expr = schema.ExprConstraints{
				schema.ListExpr{
					Elem: schema.ExprConstraints{
						objectExprFromBodySchema(block.Body),
					},
					MinItems: block.MinItems,
					MaxItems: block.MaxItems,
				},
			}
		case schema.BlockTypeSet:
			oe.Attributes[bType].Expr = schema.ExprConstraints{
				schema.SetExpr{
					Elem: schema.ExprConstraints{
						objectExprFromBodySchema(block.Body),
					},
					MinItems: block.MinItems,
					MaxItems: block.MaxItems,
				},
			}
		case schema.BlockTypeMap:
			oe.Attributes[bType].Expr = schema.ExprConstraints{
				schema.MapExpr{
					Elem: schema.ExprConstraints{
						objectExprFromBodySchema(block.Body),
					},
					MinItems: block.MinItems,
					MaxItems: block.MaxItems,
				},
			}
		}
	}

	return oe
}
