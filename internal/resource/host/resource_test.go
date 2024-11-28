package host_test

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
)

var _ = DescribeTable("host resource management",
	func(steps ...resource.TestStep) {
		resource.Test(GinkgoT(), resource.TestCase{
			Steps: lo.Map(steps, func(step resource.TestStep, _ int) resource.TestStep {
				step.ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
					"definednet": providerserver.NewProtocol6WithError(providerFactory()),
				}

				return step
			}),
		})
	},
	Entry("assert host is created in expected configuration",
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/host.tf"),
			ConfigVariables: config.Variables{
				"name":       config.StringVariable("host.defined.test"),
				"network_id": config.StringVariable("network-id"),
				"role_id":    config.StringVariable("role-id"),
				"tags": config.ListVariable(
					config.StringVariable("tag:one"),
					config.StringVariable("tag:two"),
				),
			},
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					// Assert sanity.
					plancheck.ExpectResourceAction("definednet_host.test", plancheck.ResourceActionCreate),
				},
			},
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					"definednet_host.test",
					tfjsonpath.New("id"),
					knownvalue.StringRegexp(regexp.MustCompile(`^host-[A-Z0-9]+$`)),
				),
				statecheck.ExpectKnownValue(
					"definednet_host.test",
					tfjsonpath.New("ip_address"),
					knownvalue.StringRegexp(regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)),
				),
				statecheck.ExpectSensitiveValue("definednet_host.test", tfjsonpath.New("enrollment_code")),
				statecheck.ExpectKnownValue(
					"definednet_host.test",
					tfjsonpath.New("enrollment_code"),
					knownvalue.StringRegexp(regexp.MustCompile(`^.{32}$`)),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("definednet_host.test", "name", "host.defined.test"),
				resource.TestCheckResourceAttr("definednet_host.test", "network_id", "network-id"),
				resource.TestCheckResourceAttr("definednet_host.test", "role_id", "role-id"),
				resource.TestCheckResourceAttr("definednet_host.test", "tags.0", "tag:one"),
				resource.TestCheckResourceAttr("definednet_host.test", "tags.1", "tag:two"),
			),
		},
	),
	Entry("assert simple updates are executed in-place",
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/host.tf"),
			ConfigVariables: config.Variables{
				"name":       config.StringVariable("host.defined.test"),
				"network_id": config.StringVariable("network-id"),
				"role_id":    config.StringVariable("role-id"),
				"tags": config.ListVariable(
					config.StringVariable("tag:one"),
					config.StringVariable("tag:two"),
				),
			},
		},
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/host.tf"),
			ConfigVariables: config.Variables{
				"name":       config.StringVariable("updated-host.defined.test"),
				"network_id": config.StringVariable("network-id"),
				"role_id":    config.StringVariable("updated-role-id"),
				"tags": config.ListVariable(
					config.StringVariable("tag:one"),
					config.StringVariable("tag:three"),
				),
			},
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction("definednet_host.test", plancheck.ResourceActionUpdate),
				},
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("definednet_host.test", "name", "updated-host.defined.test"),
				resource.TestCheckResourceAttr("definednet_host.test", "network_id", "network-id"),
				resource.TestCheckResourceAttr("definednet_host.test", "role_id", "updated-role-id"),
				resource.TestCheckResourceAttr("definednet_host.test", "tags.0", "tag:one"),
				resource.TestCheckResourceAttr("definednet_host.test", "tags.1", "tag:three"),
			),
		},
	),
	Entry("assert updating network_id replaces the host",
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/host.tf"),
			ConfigVariables: config.Variables{
				"name":       config.StringVariable("host.defined.test"),
				"network_id": config.StringVariable("network-id"),
				"role_id":    config.StringVariable("role-id"),
				"tags": config.ListVariable(
					config.StringVariable("tag:one"),
					config.StringVariable("tag:two"),
				),
			},
		},
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/host.tf"),
			ConfigVariables: config.Variables{
				"name":       config.StringVariable("host.defined.test"),
				"network_id": config.StringVariable("updated-network-id"),
				"role_id":    config.StringVariable("role-id"),
				"tags": config.ListVariable(
					config.StringVariable("tag:one"),
					config.StringVariable("tag:two"),
				),
			},
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					// Assert the resource is replaced.
					plancheck.ExpectResourceAction("definednet_host.test", plancheck.ResourceActionReplace),
				},
			},
		},
	),
	Entry("assert importing host populates the host",
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/host.tf"),
			ConfigVariables: config.Variables{
				"name":       config.StringVariable("host.defined.test"),
				"network_id": config.StringVariable("network-id"),
				"role_id":    config.StringVariable("role-id"),
				"tags": config.ListVariable(
					config.StringVariable("tag:one"),
					config.StringVariable("tag:two"),
				),
			},
		},
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/host.tf"),
			ConfigVariables: config.Variables{
				"name":       config.StringVariable("host.defined.test"),
				"network_id": config.StringVariable("network-id"),
				"role_id":    config.StringVariable("role-id"),
				"tags": config.ListVariable(
					config.StringVariable("tag:one"),
					config.StringVariable("tag:two"),
				),
			},
			ResourceName:            "definednet_host.test",
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"enrollment_code"},
		},
	),
	Entry("assert optional fields are optional",
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/host_minimal.tf"),
			ConfigVariables: config.Variables{
				"name":       config.StringVariable("host.defined.test"),
				"network_id": config.StringVariable("network-id"),
			},
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue("definednet_host.minimal_test", tfjsonpath.New("role_id"), knownvalue.Null()),
				statecheck.ExpectKnownValue("definednet_host.minimal_test", tfjsonpath.New("tags"), knownvalue.Null()),
			},
		},
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/host_minimal.tf"),
			ConfigVariables: config.Variables{
				"name":       config.StringVariable("host.defined.test"),
				"network_id": config.StringVariable("network-id"),
			},
			ResourceName:            "definednet_host.minimal_test",
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"enrollment_code"},
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue("definednet_host.minimal_test", tfjsonpath.New("role_id"), knownvalue.Null()),
				statecheck.ExpectKnownValue("definednet_host.minimal_test", tfjsonpath.New("tags"), knownvalue.Null()),
			},
		},
	),
)
