package lighthouse_test

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	. "github.com/onsi/ginkgo/v2"
	"github.com/samber/lo"
	"github.com/sendsmaily/terraform-provider-definednet/internal/lighthouse"
)

var _ = DescribeTable("lighthouse resource management",
	func(steps ...resource.TestStep) {
		resource.Test(GinkgoT(), resource.TestCase{
			Steps: lo.Map(steps, func(step resource.TestStep, _ int) resource.TestStep {
				step.ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
					"definednet": providerserver.NewProtocol6WithError(provider.WithResource(lighthouse.NewResource())),
				}

				return step
			}),
		})
	},
	Entry("assert lighthouse is created in expected configuration",
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/lighthouse.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("lighthouse.defined.test"),
				"network_id":  config.StringVariable("network-id"),
				"role_id":     config.StringVariable("role-id"),
				"listen_port": config.IntegerVariable(8484),
				"static_addresses": config.ListVariable(
					config.StringVariable("127.0.0.1"),
					config.StringVariable("172.16.0.1"),
				),
				"tags": config.ListVariable(
					config.StringVariable("tag:one"),
					config.StringVariable("tag:two"),
				),
			},
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					// Assert sanity.
					plancheck.ExpectResourceAction("definednet_lighthouse.test", plancheck.ResourceActionCreate),
				},
			},
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					"definednet_lighthouse.test",
					tfjsonpath.New("id"),
					knownvalue.StringRegexp(regexp.MustCompile(`^host-[A-Z0-9]+$`)),
				),
				statecheck.ExpectKnownValue(
					"definednet_lighthouse.test",
					tfjsonpath.New("ip_address"),
					knownvalue.StringRegexp(regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)),
				),
				statecheck.ExpectSensitiveValue("definednet_lighthouse.test", tfjsonpath.New("enrollment_code")),
				statecheck.ExpectKnownValue(
					"definednet_lighthouse.test",
					tfjsonpath.New("enrollment_code"),
					knownvalue.StringRegexp(regexp.MustCompile(`^\w+$`)),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("definednet_lighthouse.test", "name", "lighthouse.defined.test"),
				resource.TestCheckResourceAttr("definednet_lighthouse.test", "network_id", "network-id"),
				resource.TestCheckResourceAttr("definednet_lighthouse.test", "role_id", "role-id"),
				resource.TestCheckResourceAttr("definednet_lighthouse.test", "listen_port", "8484"),
				resource.TestCheckResourceAttr("definednet_lighthouse.test", "static_addresses.0", "127.0.0.1"),
				resource.TestCheckResourceAttr("definednet_lighthouse.test", "static_addresses.1", "172.16.0.1"),
				resource.TestCheckResourceAttr("definednet_lighthouse.test", "tags.0", "tag:one"),
				resource.TestCheckResourceAttr("definednet_lighthouse.test", "tags.1", "tag:two"),
			),
		},
	),
	Entry("assert simple updates are executed in-place",
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/lighthouse.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("lighthouse.defined.test"),
				"network_id":  config.StringVariable("network-id"),
				"role_id":     config.StringVariable("role-id"),
				"listen_port": config.IntegerVariable(8484),
				"static_addresses": config.ListVariable(
					config.StringVariable("127.0.0.1"),
					config.StringVariable("172.16.0.1"),
				),
				"tags": config.ListVariable(
					config.StringVariable("tag:one"),
					config.StringVariable("tag:two"),
				),
			},
		},
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/lighthouse.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("updated-lighthouse.defined.test"),
				"network_id":  config.StringVariable("network-id"),
				"role_id":     config.StringVariable("updated-role-id"),
				"listen_port": config.IntegerVariable(6363),
				"static_addresses": config.ListVariable(
					config.StringVariable("127.0.0.1"),
					config.StringVariable("172.16.0.1"),
				),
				"tags": config.ListVariable(
					config.StringVariable("tag:one"),
					config.StringVariable("tag:three"),
				),
			},
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction("definednet_lighthouse.test", plancheck.ResourceActionUpdate),
				},
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("definednet_lighthouse.test", "listen_port", "6363"),
				resource.TestCheckResourceAttr("definednet_lighthouse.test", "static_addresses.0", "127.0.0.1"),
				resource.TestCheckResourceAttr("definednet_lighthouse.test", "static_addresses.1", "172.16.0.1"),
				resource.TestCheckResourceAttr("definednet_lighthouse.test", "name", "updated-lighthouse.defined.test"),
				resource.TestCheckResourceAttr("definednet_lighthouse.test", "network_id", "network-id"),
				resource.TestCheckResourceAttr("definednet_lighthouse.test", "role_id", "updated-role-id"),
				resource.TestCheckResourceAttr("definednet_lighthouse.test", "tags.0", "tag:one"),
				resource.TestCheckResourceAttr("definednet_lighthouse.test", "tags.1", "tag:three"),
			),
		},
	),
	Entry("assert updating network_id replaces the lighthouse",
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/lighthouse.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("lighthouse.defined.test"),
				"network_id":  config.StringVariable("network-id"),
				"role_id":     config.StringVariable("role-id"),
				"listen_port": config.IntegerVariable(8484),
				"static_addresses": config.ListVariable(
					config.StringVariable("127.0.0.1"),
					config.StringVariable("172.16.0.1"),
				),
				"tags": config.ListVariable(
					config.StringVariable("tag:one"),
					config.StringVariable("tag:two"),
				),
			},
		},
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/lighthouse.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("lighthouse.defined.test"),
				"network_id":  config.StringVariable("updated-network-id"),
				"role_id":     config.StringVariable("role-id"),
				"listen_port": config.IntegerVariable(8484),
				"static_addresses": config.ListVariable(
					config.StringVariable("127.0.0.1"),
					config.StringVariable("172.16.0.1"),
				),
				"tags": config.ListVariable(
					config.StringVariable("tag:one"),
					config.StringVariable("tag:two"),
				),
			},
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					// Assert the resource is replaced.
					plancheck.ExpectResourceAction("definednet_lighthouse.test", plancheck.ResourceActionReplace),
				},
			},
		},
	),
	Entry("assert importing lighthouse populates the lighthouse",
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/lighthouse.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("lighthouse.defined.test"),
				"network_id":  config.StringVariable("network-id"),
				"role_id":     config.StringVariable("role-id"),
				"listen_port": config.IntegerVariable(8484),
				"static_addresses": config.ListVariable(
					config.StringVariable("127.0.0.1"),
					config.StringVariable("172.16.0.1"),
				),
				"tags": config.ListVariable(
					config.StringVariable("tag:one"),
					config.StringVariable("tag:two"),
				),
			},
		},
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/lighthouse.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("lighthouse.defined.test"),
				"network_id":  config.StringVariable("network-id"),
				"role_id":     config.StringVariable("role-id"),
				"listen_port": config.IntegerVariable(8484),
				"static_addresses": config.ListVariable(
					config.StringVariable("127.0.0.1"),
					config.StringVariable("172.16.0.1"),
				),
				"tags": config.ListVariable(
					config.StringVariable("tag:one"),
					config.StringVariable("tag:two"),
				),
			},
			ResourceName:            "definednet_lighthouse.test",
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"enrollment_code"},
		},
	),
)
