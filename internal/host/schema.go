package host

import (
	_ "embed"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Schema is the host resource's schema.
var Schema = schema.Schema{
	MarkdownDescription: resourceDescription,
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "Host's ID",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			Description: "Host's name",
			Required:    true,
		},
		"network_id": schema.StringAttribute{
			Description: "Enrolled Network ID",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"role_id": schema.StringAttribute{
			Description: "Host's role ID on Defined.net",
			Optional:    true,
		},
		"ip_address": schema.StringAttribute{
			Description: "Host's IP address on Defined.net overlay network",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"tags": schema.ListAttribute{
			Description: "Host's tags on Defined.net",
			ElementType: types.StringType,
			Optional:    true,
		},
		"enrollment_code": schema.StringAttribute{
			Description: "Host's enrollment code",
			Sensitive:   true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	},
}

//go:embed docs/resource.md
var resourceDescription string
