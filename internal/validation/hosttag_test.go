package validation_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sendsmaily/terraform-provider-definednet/internal/validation"
)

var _ = Describe("validating host tags", func() {
	DescribeTable("valid tags pass validation",
		func(ctx SpecContext, tag string) {
			res := new(validator.StringResponse)
			validation.HostTag().ValidateString(ctx, validator.StringRequest{
				Path:        path.Empty().AtName("test"),
				ConfigValue: basetypes.NewStringValue("valid_host_key:valid_host-value"),
			}, res)

			Expect(res.Diagnostics.HasError()).To(BeFalse(), GetDiagnosticsMessage(res.Diagnostics))
		},
		Entry("flat tags", "key:value"),
		Entry("tags using underscore", "underscore_key:underscore_value"),
		Entry("tags using hyphens", "hyphen-key:hyphen-value"),
	)

	DescribeTable("invalid tags fail validation",
		func(ctx SpecContext, tag string) {
			res := new(validator.StringResponse)
			validation.HostTag().ValidateString(ctx, validator.StringRequest{
				Path:        path.Empty().AtName("test"),
				ConfigValue: basetypes.NewStringValue(tag),
			}, res)

			Expect(res.Diagnostics).To(ContainElement(diag.NewAttributeErrorDiagnostic(
				path.Empty().AtName("test"),
				"Invalid Attribute Value Match",
				fmt.Sprintf("Attribute test must be in the format KEY:VALUE, got: %s", tag),
			)))
		},
		Entry("invalid format", "tag"),
		Entry("invalid value", "key:"),
		Entry("missing key", ":value"),
		Entry("disallowed characters", "$key:$value"),
		Entry("disallowed integers", "key1:value2"),
		Entry("disallowed uppercase characters", "TAG:VALUE"),
		Entry("disallowed whitespace", "key:with value"),
	)
})
