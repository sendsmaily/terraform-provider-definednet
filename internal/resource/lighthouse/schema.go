package lighthouse

import (
	_ "embed"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sendsmaily/terraform-provider-definednet/internal/validation"
)

// Schema is the lighthouse resource's schema.
var Schema = schema.Schema{
	MarkdownDescription: resourceDescription,
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description: "Lighthouse's name",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(255),
			},
		},
		"network_id": schema.StringAttribute{
			Description: "Enrolled Network ID",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"role_id": schema.StringAttribute{
			Description: "Lighthouse's role ID on Defined.net",
			Optional:    true,
		},
		"static_addresses": schema.ListAttribute{
			Description: "Lighthouse's static IP addresses",
			ElementType: types.StringType,
			Required:    true,
			Validators: []validator.List{
				listvalidator.UniqueValues(),
				listvalidator.ValueStringsAre(validation.IPAddress()),
			},
		},
		"listen_port": schema.Int32Attribute{
			Description: "Lighthouse's listen port",
			Required:    true,
		},
		"tags": schema.ListAttribute{
			Description: "Lighthouse's tags on Defined.net",
			ElementType: types.StringType,
			Optional:    true,
			Validators: []validator.List{
				listvalidator.UniqueValues(),
				listvalidator.ValueStringsAre(validation.HostTag()),
			},
		},
		"id": schema.StringAttribute{
			Description: "Lighthouse's ID",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"ip_address": schema.StringAttribute{
			Description: "Lighthouse's IP address on Defined.net overlay network",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"enrollment_code": schema.StringAttribute{
			Description: "Lighthouse's enrollment code",
			Sensitive:   true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	},
	Blocks: map[string]schema.Block{
		"metrics": schema.SingleNestedBlock{
			Description: "Host's metrics exporter configuration",
			Attributes: map[string]schema.Attribute{
				"enabled": schema.BoolAttribute{
					Description: "Enable metrics exporter",
					Optional:    true,
				},
				"listen": schema.StringAttribute{
					Description: "Host-port for Prometheus metrics exporter listener",
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString("127.0.0.1:8080"),
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"path": schema.StringAttribute{
					Description: "Prometheus metrics exporter's HTTP path",
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString("/metrics"),
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"namespace": schema.StringAttribute{
					Description: "Prometheus metrics' namespace",
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString("nebula"),
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"subsystem": schema.StringAttribute{
					Description: "Prometheus metrics' subsystem",
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString("lighthouse"),
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"enable_extra_metrics": schema.BoolAttribute{
					Description: "Enable extra metrics",
					Optional:    true,
					Computed:    true,
					PlanModifiers: []planmodifier.Bool{
						boolplanmodifier.UseStateForUnknown(),
					},
				},
			},
		},
	},
}

//go:embed docs/resource.md
var resourceDescription string
