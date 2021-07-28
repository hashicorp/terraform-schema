package module

import (
	"github.com/zclconf/go-cty/cty"
)

type Output struct {
	Description string
	IsSensitive bool
	Value       cty.Value
}
