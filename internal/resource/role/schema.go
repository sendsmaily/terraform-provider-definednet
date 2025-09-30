package role

import (
	_ "embed"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sendsmaily/terraform-provider-definednet/internal/validation"
)

// Schema is the host resource's schema.
var Schema = schema.Schema{
	MarkdownDescription: resourceDescription,
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "Role's ID",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			Description: "Role's name",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(50),
			},
		},
		"description": schema.StringAttribute{
			Description: "Role's description",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(255),
			},
		},
	},
	Blocks: map[string]schema.Block{
		"rule": schema.SetNestedBlock{
			Description: "Role's firewall rule",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"protocol": schema.StringAttribute{
						Description: "Network protocol. One of `ANY`, `TCP`, `UDP`, or `ICMP`.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("ANY", "TCP", "UDP", "ICMP"),
						},
					},
					"description": schema.StringAttribute{
						Description: "Role's description",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtMost(255),
						},
					},
					"allowed_role_id": schema.StringAttribute{
						Description: "Allowed role's ID",
						Optional:    true,
					},
					"allowed_tags": schema.SetAttribute{
						Description: "Allowed hosts' tags",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.ValueStringsAre(validation.HostTag()),
						},
					},
					"port": schema.Int32Attribute{
						Description: "Allowed port",
						Optional:    true,
						Validators: []validator.Int32{
							int32validator.Between(1, 65535),
						},
					},
					"port_range": schema.SingleNestedAttribute{
						Description: "Allowed port range",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"from": schema.Int32Attribute{
								Description: "Start of the allowed port range",
								Required:    true,
								Validators: []validator.Int32{
									int32validator.Between(1, 65535),
								},
							},
							"to": schema.Int32Attribute{
								Description: "End of the allowed port range",
								Required:    true,
								Validators: []validator.Int32{
									int32validator.Between(1, 65535),
								},
							},
						},
					},
				},
				Validators: []validator.Object{
					objectvalidator.AtLeastOneOf(
						path.MatchRelative().AtName("allowed_role_id"),
						path.MatchRelative().AtName("allowed_tags"),
					),
				},
			},
		},
	},
}

//go:embed docs/resource.md
var resourceDescription string
