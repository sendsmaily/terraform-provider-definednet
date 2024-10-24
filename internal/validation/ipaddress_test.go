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

var _ = Describe("validating IPv4 addresses", func() {
	Specify("valid values pass validation", func(ctx SpecContext) {
		res := new(validator.StringResponse)
		validation.IPAddress().ValidateString(ctx, validator.StringRequest{
			Path:        path.Empty().AtName("test"),
			ConfigValue: basetypes.NewStringValue("172.16.0.1"),
		}, res)

		Expect(res.Diagnostics.HasError()).To(BeFalse(), GetDiagnosticsMessage(res.Diagnostics))
	})

	DescribeTable("invalid values fail validation",
		func(ctx SpecContext, addr string) {
			res := new(validator.StringResponse)
			validation.IPAddress().ValidateString(ctx, validator.StringRequest{
				Path:        path.Empty().AtName("test"),
				ConfigValue: basetypes.NewStringValue(addr),
			}, res)

			Expect(res.Diagnostics).To(ContainElement(diag.NewAttributeErrorDiagnostic(
				path.Empty().AtName("test"),
				"Invalid Attribute Value Match",
				fmt.Sprintf("Attribute test value must be an IP address, got: %s", addr),
			)))
		},
		Entry("invalid address", "172.16.0.256"),
		Entry("host-port", "172.16.0.1:4242"),
	)
})

var _ = Describe("validating IPv6 addresses", func() {
	Specify("valid values pass validation", func(ctx SpecContext) {
		res := new(validator.StringResponse)
		validation.IPAddress().ValidateString(ctx, validator.StringRequest{
			Path:        path.Empty().AtName("test"),
			ConfigValue: basetypes.NewStringValue("fd:beef::1"),
		}, res)

		Expect(res.Diagnostics.HasError()).To(BeFalse(), GetDiagnosticsMessage(res.Diagnostics))
	})

	DescribeTable("invalid values fail validation",
		func(ctx SpecContext, addr string) {
			res := new(validator.StringResponse)
			validation.IPAddress().ValidateString(ctx, validator.StringRequest{
				Path:        path.Empty().AtName("test"),
				ConfigValue: basetypes.NewStringValue(addr),
			}, res)

			Expect(res.Diagnostics).To(ContainElement(diag.NewAttributeErrorDiagnostic(
				path.Empty().AtName("test"),
				"Invalid Attribute Value Match",
				fmt.Sprintf("Attribute test value must be an IP address, got: %s", addr),
			)))
		},
		Entry("invalid address", "fd:beef::g"),
		Entry("host-port", "[fd:beef::1]:4242"),
	)
})

var _ = Describe("validating null values", func() {
	Specify("null values pass validation", func(ctx SpecContext) {
		res := new(validator.StringResponse)
		validation.IPAddress().ValidateString(ctx, validator.StringRequest{
			Path:        path.Empty().AtName("test"),
			ConfigValue: basetypes.NewStringNull(),
		}, res)

		Expect(res.Diagnostics.HasError()).To(BeFalse(), GetDiagnosticsMessage(res.Diagnostics))
	})
})

var _ = Describe("validating unknown values", func() {
	Specify("unknown values pass validation", func(ctx SpecContext) {
		res := new(validator.StringResponse)
		validation.IPAddress().ValidateString(ctx, validator.StringRequest{
			Path:        path.Empty().AtName("test"),
			ConfigValue: basetypes.NewStringUnknown(),
		}, res)

		Expect(res.Diagnostics.HasError()).To(BeFalse(), GetDiagnosticsMessage(res.Diagnostics))
	})
})
