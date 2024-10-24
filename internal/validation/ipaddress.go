package validation

import (
	"context"
	"net"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// IPAddress validates the value is an IP address.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func IPAddress() ipaddressValidator {
	return ipaddressValidator{}
}

type ipaddressValidator struct{}

var _ validator.String = ipaddressValidator{}

func (validator ipaddressValidator) Description(_ context.Context) string {
	return "value must be an IP address"
}

func (validator ipaddressValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (v ipaddressValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if net.ParseIP(value) == nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value,
		))
	}
}
