package earlydecoder

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/terraform-schema/backend"
	"github.com/zclconf/go-cty/cty"
)

func decodeBackendsBlock(block *hcl.Block) (backend.BackendData, hcl.Diagnostics) {
	bType := block.Labels[0]
	attrs, diags := block.Body.JustAttributes()

	switch bType {
	case "remote":
		if attr, ok := attrs["hostname"]; ok {
			val, vDiags := attr.Expr.Value(nil)
			diags = append(diags, vDiags...)
			if val.IsWhollyKnown() && val.Type() == cty.String {
				return &backend.Remote{
					Hostname: val.AsString(),
				}, nil
			}
		}

		return &backend.Remote{}, nil
	}

	return &backend.UnknownBackendData{}, diags
}
