package validation_test

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sendsmaily/terraform-provider-definednet/internal/validation"
)

var _ = Describe("validating port ranges", func() {
	Specify("valid values pass validation", func(ctx SpecContext) {
		res := new(validator.ObjectResponse)
		validation.PortRange(1, 65535).ValidateObject(ctx, validator.ObjectRequest{
			Path: path.Empty().AtName("test"),
			ConfigValue: basetypes.NewObjectValueMust(
				map[string]attr.Type{
					"from": types.Int32Type,
					"to":   types.Int32Type,
				},
				map[string]attr.Value{
					"from": types.Int32Value(1024),
					"to":   types.Int32Value(2048),
				},
			),
		}, res)

		Expect(res.Diagnostics.HasError()).To(BeFalse(), GetDiagnosticsMessage(res.Diagnostics))
	})

	DescribeTable("invalid values fail validation",
		func(ctx SpecContext, from, to int, summary, detail string) {
			res := new(validator.ObjectResponse)
			validation.PortRange(2048, 4096).ValidateObject(ctx, validator.ObjectRequest{
				Path: path.Empty().AtName("test"),
				ConfigValue: basetypes.NewObjectValueMust(
					map[string]attr.Type{
						"from": types.Int32Type,
						"to":   types.Int32Type,
					},
					map[string]attr.Value{
						"from": types.Int32Value(int32(from)),
						"to":   types.Int32Value(int32(to)),
					},
				),
			}, res)

			Expect(res.Diagnostics.Errors()).To(ContainElement(SatisfyAll(
				HaveField("Summary()", summary),
				HaveField("Detail()", detail),
			)))
		},
		Entry("From port value undercuts the allowed range",
			1024, 4096,
			"Invalid Attribute Value",
			"Attribute test.from port value must be between 2048 and 4096, got: 1024",
		),
		Entry("From port value exceeds the allowed range",
			8192, 16384,
			"Invalid Attribute Value",
			"Attribute test.from port value must be between 2048 and 4096, got: 8192",
		),
		Entry("To port value undercuts the allowed range",
			512, 1024,
			"Invalid Attribute Value",
			"Attribute test.to port value must be between 2048 and 4096, got: 1024",
		),
		Entry("To port value exceeds the allowed range",
			2048, 8192,
			"Invalid Attribute Value",
			"Attribute test.to port value must be between 2048 and 4096, got: 8192",
		),
		Entry("To port value is greater than from port value",
			4096, 2048,
			"Invalid Attribute Combination",
			`"to" port must be greater than "from" port`,
		),
	)
})

var _ = Describe("validating null values", func() {
	Specify("null values pass validation", func(ctx SpecContext) {
		res := new(validator.ObjectResponse)
		validation.PortRange(1024, 2048).ValidateObject(ctx, validator.ObjectRequest{
			Path:        path.Empty().AtName("test"),
			ConfigValue: basetypes.NewObjectNull(nil),
		}, res)

		Expect(res.Diagnostics.HasError()).To(BeFalse(), GetDiagnosticsMessage(res.Diagnostics))
	})
})

var _ = Describe("validating unknown values", func() {
	Specify("unknown values pass validation", func(ctx SpecContext) {
		res := new(validator.ObjectResponse)
		validation.PortRange(1024, 2048).ValidateObject(ctx, validator.ObjectRequest{
			Path:        path.Empty().AtName("test"),
			ConfigValue: basetypes.NewObjectUnknown(nil),
		}, res)

		Expect(res.Diagnostics.HasError()).To(BeFalse(), GetDiagnosticsMessage(res.Diagnostics))
	})
})
