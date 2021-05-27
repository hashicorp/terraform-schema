package module

import (
	"github.com/zclconf/go-cty/cty"
)

type Variable struct {
	Description  string
	Type         cty.Type
	DefaultValue cty.Value

	// In case the version it is before 0.14 sensitive will always be false
	// that was actually the default value for prior versions
	IsSensitive bool
}
