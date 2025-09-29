package validation

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PortRange is a port range validator.
func PortRange(start, end int32) validator.Object {
	return portRangeValidator{
		start: start,
		end:   end,
	}
}

type portRangeValidator struct {
	start, end int32
}

func (v portRangeValidator) Description(context.Context) string {
	return fmt.Sprintf("port value must be between %d and %d", v.start, v.end)
}

func (v portRangeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v portRangeValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	val, d := req.ConfigValue.ToObjectValue(ctx)
	resp.Diagnostics.Append(d...)

	if resp.Diagnostics.HasError() || val.IsNull() || val.IsUnknown() {
		return
	}

	attrs := val.Attributes()
	from, ok := attrs["from"].(types.Int32)
	if !ok {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeTypeDiagnostic(
			req.Path.AtName("from"),
			"must be int32",
			attrs["from"].Type(ctx).String(),
		))
	}

	to, ok := attrs["to"].(types.Int32)
	if !ok {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeTypeDiagnostic(
			req.Path.AtName("to"),
			"must be int32",
			attrs["to"].Type(ctx).String(),
		))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if from.ValueInt32() < v.start || from.ValueInt32() > v.end {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.Path.AtName("from"),
			v.Description(ctx),
			from.String(),
		))
	}

	if to.ValueInt32() < v.start || to.ValueInt32() > v.end {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.Path.AtName("to"),
			v.Description(ctx),
			to.String(),
		))
	}

	if to.ValueInt32() <= from.ValueInt32() {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
			req.Path,
			`"to" port must be greater than "from" port`,
		))
	}
}
