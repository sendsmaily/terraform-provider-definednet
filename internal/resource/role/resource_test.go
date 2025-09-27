package role_test

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

var _ = DescribeTable("role resource management",
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
	Entry("assert role is created in expected configuration",
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/role.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("test: Role"),
				"description": config.StringVariable("Role's description"),
			},
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					// Assert sanity.
					plancheck.ExpectResourceAction("definednet_role.test", plancheck.ResourceActionCreate),
				},
			},
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					"definednet_role.test",
					tfjsonpath.New("id"),
					knownvalue.StringRegexp(regexp.MustCompile(`^role-[A-Z0-9]+$`)),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("definednet_role.test", "name", "test: Role"),
				resource.TestCheckResourceAttr("definednet_role.test", "description", "Role's description"),
			),
		},
	),
	Entry("assert role is updated with expected configuration",
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/role.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("test: Role"),
				"description": config.StringVariable("Role's description"),
			},
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					// Assert sanity.
					plancheck.ExpectResourceAction("definednet_role.test", plancheck.ResourceActionCreate),
				},
			},
		},
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/role.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("test: Updated role"),
				"description": config.StringVariable("Updated role's description"),
			},
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					// Assert sanity.
					plancheck.ExpectResourceAction("definednet_role.test", plancheck.ResourceActionUpdate),
				},
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("definednet_role.test", "name", "test: Updated role"),
				resource.TestCheckResourceAttr("definednet_role.test", "description", "Updated role's description"),
			),
		},
	),
	Entry("assert importing the role populates the state",
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/role.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("test: Role"),
				"description": config.StringVariable("Role's description"),
			},
		},
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/role.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("test: Role"),
				"description": config.StringVariable("Role's description"),
			},
			ResourceName:      "definednet_role.test",
			ImportState:       true,
			ImportStateVerify: true,
		},
	),
)

var _ = DescribeTable("port-based firewall management",
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
	Entry("assert roles with single port rules can be created",
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/role_port.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("test: Role"),
				"description": config.StringVariable("Role's description"),
				"rules": config.ListVariable(
					config.ObjectVariable(map[string]config.Variable{
						"port":            config.IntegerVariable(22),
						"protocol":        config.StringVariable("TCP"),
						"description":     config.StringVariable("SSH access"),
						"allowed_role_id": config.StringVariable("role:abcdef"),
						"allowed_tags": config.ListVariable(
							config.StringVariable("tag:one"),
							config.StringVariable("tag:two"),
						),
					}),
					config.ObjectVariable(map[string]config.Variable{
						"port":            config.IntegerVariable(443),
						"protocol":        config.StringVariable("TCP"),
						"description":     config.StringVariable("HTTPS access"),
						"allowed_role_id": config.StringVariable("role:123456"),
						"allowed_tags": config.ListVariable(
							config.StringVariable("tag:https_one"),
							config.StringVariable("tag:https_two"),
						),
					}),
				),
			},
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					// Assert sanity.
					plancheck.ExpectResourceAction("definednet_role.test", plancheck.ResourceActionCreate),
				},
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("definednet_role.test", "rule.0.port", "22"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.0.protocol", "TCP"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.0.description", "SSH access"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.0.allowed_role_id", "role:abcdef"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.0.allowed_tags.0", "tag:one"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.0.allowed_tags.1", "tag:two"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.1.port", "443"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.1.protocol", "TCP"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.1.description", "HTTPS access"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.1.allowed_role_id", "role:123456"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.1.allowed_tags.0", "tag:https_one"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.1.allowed_tags.1", "tag:https_two"),
			),
		},
	),
	Entry("assert roles with single port rules can be updated",
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/role_port.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("test: Role"),
				"description": config.StringVariable("Role's description"),
			},
		},
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/role_port.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("test: Role"),
				"description": config.StringVariable("Role's description"),
				"rules": config.ListVariable(
					config.ObjectVariable(map[string]config.Variable{
						"port":            config.IntegerVariable(22),
						"protocol":        config.StringVariable("TCP"),
						"description":     config.StringVariable("SSH access"),
						"allowed_role_id": config.StringVariable("role:abcdef"),
						"allowed_tags": config.ListVariable(
							config.StringVariable("tag:one"),
							config.StringVariable("tag:two"),
						),
					}),
					config.ObjectVariable(map[string]config.Variable{
						"port":            config.IntegerVariable(443),
						"protocol":        config.StringVariable("TCP"),
						"description":     config.StringVariable("HTTPS access"),
						"allowed_role_id": config.StringVariable("role:123456"),
						"allowed_tags": config.ListVariable(
							config.StringVariable("tag:https_one"),
							config.StringVariable("tag:https_two"),
						),
					}),
				),
			},
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					// Assert sanity.
					plancheck.ExpectResourceAction("definednet_role.test", plancheck.ResourceActionUpdate),
				},
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("definednet_role.test", "rule.0.port", "22"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.0.protocol", "TCP"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.0.description", "SSH access"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.0.allowed_role_id", "role:abcdef"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.0.allowed_tags.0", "tag:one"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.0.allowed_tags.1", "tag:two"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.1.port", "443"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.1.protocol", "TCP"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.1.description", "HTTPS access"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.1.allowed_role_id", "role:123456"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.1.allowed_tags.0", "tag:https_one"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.1.allowed_tags.1", "tag:https_two"),
			),
		},
	),
	Entry("assert importing the role populates the state",
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/role_port.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("test: Role"),
				"description": config.StringVariable("Role's description"),
				"rules": config.ListVariable(
					config.ObjectVariable(map[string]config.Variable{
						"port":            config.IntegerVariable(22),
						"protocol":        config.StringVariable("TCP"),
						"description":     config.StringVariable("SSH access"),
						"allowed_role_id": config.StringVariable("role:abcdef"),
						"allowed_tags": config.ListVariable(
							config.StringVariable("tag:one"),
							config.StringVariable("tag:two"),
						),
					}),
					config.ObjectVariable(map[string]config.Variable{
						"port":            config.IntegerVariable(443),
						"protocol":        config.StringVariable("TCP"),
						"description":     config.StringVariable("HTTPS access"),
						"allowed_role_id": config.StringVariable("role:123456"),
						"allowed_tags": config.ListVariable(
							config.StringVariable("tag:https_one"),
							config.StringVariable("tag:https_two"),
						),
					}),
				),
			},
		},
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/role_port.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("test: Role"),
				"description": config.StringVariable("Role's description"),
				"rules": config.ListVariable(config.ObjectVariable(map[string]config.Variable{
					"port":            config.IntegerVariable(22),
					"protocol":        config.StringVariable("TCP"),
					"description":     config.StringVariable("SSH access"),
					"allowed_role_id": config.StringVariable("role:abcdef"),
					"allowed_tags": config.ListVariable(
						config.StringVariable("tag:one"),
						config.StringVariable("tag:two"),
					),
				})),
			},
			ResourceName:      "definednet_role.test",
			ImportState:       true,
			ImportStateVerify: true,
		},
	),
	Entry("assert port value sanity is enforced",
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/role_port.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("test: Role"),
				"description": config.StringVariable("Test port must not be zero."),
				"rules": config.ListVariable(
					config.ObjectVariable(map[string]config.Variable{
						"port":            config.IntegerVariable(0),
						"protocol":        config.StringVariable("TCP"),
						"description":     config.StringVariable("Port"),
						"allowed_role_id": config.StringVariable("role:abcdef"),
						"allowed_tags": config.ListVariable(
							config.StringVariable("tag:one"),
							config.StringVariable("tag:two"),
						),
					}),
				),
			},
			ExpectError: regexp.MustCompile(".+port value must be between 1 and 65535, got: 0"),
		},
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/role_port.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("test: Role"),
				"description": config.StringVariable("Test port must not be out of bounds."),
				"rules": config.ListVariable(
					config.ObjectVariable(map[string]config.Variable{
						"port":            config.IntegerVariable(65536),
						"protocol":        config.StringVariable("TCP"),
						"description":     config.StringVariable("Port"),
						"allowed_role_id": config.StringVariable("role:abcdef"),
						"allowed_tags": config.ListVariable(
							config.StringVariable("tag:one"),
							config.StringVariable("tag:two"),
						),
					}),
				),
			},
			ExpectError: regexp.MustCompile(".+port value must be between 1 and 65535, got: 65536"),
		},
	),
)

var _ = DescribeTable("port range-based firewall management",
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
	Entry("assert roles with port range rules can be created",
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/role_port_range.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("test: Role"),
				"description": config.StringVariable("Role's description"),
				"rules": config.ListVariable(
					config.ObjectVariable(map[string]config.Variable{
						"port_from":       config.IntegerVariable(1024),
						"port_to":         config.IntegerVariable(2048),
						"protocol":        config.StringVariable("TCP"),
						"description":     config.StringVariable("Range one"),
						"allowed_role_id": config.StringVariable("role:abcdef"),
						"allowed_tags": config.ListVariable(
							config.StringVariable("tag:one"),
							config.StringVariable("tag:two"),
						),
					}),
					config.ObjectVariable(map[string]config.Variable{
						"port_from":       config.IntegerVariable(4096),
						"port_to":         config.IntegerVariable(8192),
						"protocol":        config.StringVariable("TCP"),
						"description":     config.StringVariable("Range two"),
						"allowed_role_id": config.StringVariable("role:123456"),
						"allowed_tags": config.ListVariable(
							config.StringVariable("tag:https_one"),
							config.StringVariable("tag:https_two"),
						),
					}),
				),
			},
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					// Assert sanity.
					plancheck.ExpectResourceAction("definednet_role.test", plancheck.ResourceActionCreate),
				},
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("definednet_role.test", "rule.0.port_range.from", "1024"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.0.port_range.to", "2048"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.0.protocol", "TCP"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.0.description", "Range one"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.0.allowed_role_id", "role:abcdef"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.0.allowed_tags.0", "tag:one"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.0.allowed_tags.1", "tag:two"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.1.port_range.from", "4096"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.1.port_range.to", "8192"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.1.protocol", "TCP"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.1.description", "Range two"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.1.allowed_role_id", "role:123456"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.1.allowed_tags.0", "tag:https_one"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.1.allowed_tags.1", "tag:https_two"),
			),
		},
	),
	Entry("assert roles with port range rules can be updated",
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/role_port_range.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("test: Role"),
				"description": config.StringVariable("Role's description"),
			},
		},
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/role_port_range.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("test: Role"),
				"description": config.StringVariable("Role's description"),
				"rules": config.ListVariable(
					config.ObjectVariable(map[string]config.Variable{
						"port_from":       config.IntegerVariable(1024),
						"port_to":         config.IntegerVariable(2048),
						"protocol":        config.StringVariable("TCP"),
						"description":     config.StringVariable("Range one"),
						"allowed_role_id": config.StringVariable("role:abcdef"),
						"allowed_tags": config.ListVariable(
							config.StringVariable("tag:one"),
							config.StringVariable("tag:two"),
						),
					}),
					config.ObjectVariable(map[string]config.Variable{
						"port_from":       config.IntegerVariable(4096),
						"port_to":         config.IntegerVariable(8192),
						"protocol":        config.StringVariable("TCP"),
						"description":     config.StringVariable("Range two"),
						"allowed_role_id": config.StringVariable("role:123456"),
						"allowed_tags": config.ListVariable(
							config.StringVariable("tag:https_one"),
							config.StringVariable("tag:https_two"),
						),
					}),
				),
			},
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					// Assert sanity.
					plancheck.ExpectResourceAction("definednet_role.test", plancheck.ResourceActionUpdate),
				},
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("definednet_role.test", "rule.0.port_range.from", "1024"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.0.port_range.to", "2048"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.0.protocol", "TCP"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.0.description", "Range one"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.0.allowed_role_id", "role:abcdef"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.0.allowed_tags.0", "tag:one"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.0.allowed_tags.1", "tag:two"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.1.port_range.from", "4096"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.1.port_range.to", "8192"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.1.protocol", "TCP"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.1.description", "Range two"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.1.allowed_role_id", "role:123456"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.1.allowed_tags.0", "tag:https_one"),
				resource.TestCheckResourceAttr("definednet_role.test", "rule.1.allowed_tags.1", "tag:https_two"),
			),
		},
	),
	Entry("assert importing the role populates the state",
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/role_port_range.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("test: Role"),
				"description": config.StringVariable("Role's description"),
				"rules": config.ListVariable(
					config.ObjectVariable(map[string]config.Variable{
						"port_from":       config.IntegerVariable(1024),
						"port_to":         config.IntegerVariable(2048),
						"protocol":        config.StringVariable("TCP"),
						"description":     config.StringVariable("Range one"),
						"allowed_role_id": config.StringVariable("role:abcdef"),
						"allowed_tags": config.ListVariable(
							config.StringVariable("tag:one"),
							config.StringVariable("tag:two"),
						),
					}),
					config.ObjectVariable(map[string]config.Variable{
						"port_from":       config.IntegerVariable(4096),
						"port_to":         config.IntegerVariable(8192),
						"protocol":        config.StringVariable("TCP"),
						"description":     config.StringVariable("Range two"),
						"allowed_role_id": config.StringVariable("role:123456"),
						"allowed_tags": config.ListVariable(
							config.StringVariable("tag:https_one"),
							config.StringVariable("tag:https_two"),
						),
					}),
				),
			},
		},
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/role_port_range.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("test: Role"),
				"description": config.StringVariable("Role's description"),
			},
			ResourceName:      "definednet_role.test",
			ImportState:       true,
			ImportStateVerify: true,
		},
	),
	Entry("assert port range sanity is enforced",
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/role_port_range.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("test: Role"),
				"description": config.StringVariable("Test from port's zero value is rejected."),
				"rules": config.ListVariable(config.ObjectVariable(map[string]config.Variable{
					"port_from":       config.IntegerVariable(0),
					"port_to":         config.IntegerVariable(1024),
					"protocol":        config.StringVariable("TCP"),
					"description":     config.StringVariable("Range"),
					"allowed_role_id": config.StringVariable("role:abcdef"),
					"allowed_tags": config.ListVariable(
						config.StringVariable("tag:one"),
						config.StringVariable("tag:two"),
					),
				})),
			},
			ExpectError: regexp.MustCompile("Port must be >= 1"),
		},
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/role_port_range.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("test: Role"),
				"description": config.StringVariable("Test to port's zero value is rejected."),
				"rules": config.ListVariable(config.ObjectVariable(map[string]config.Variable{
					"port_from":       config.IntegerVariable(1024),
					"port_to":         config.IntegerVariable(0),
					"protocol":        config.StringVariable("TCP"),
					"description":     config.StringVariable("Range"),
					"allowed_role_id": config.StringVariable("role:abcdef"),
					"allowed_tags": config.ListVariable(
						config.StringVariable("tag:one"),
						config.StringVariable("tag:two"),
					),
				})),
			},
			ExpectError: regexp.MustCompile("Port must be >= 1"),
		},
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/role_port_range.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("test: Role"),
				"description": config.StringVariable("Test from port's out of bound value is rejected."),
				"rules": config.ListVariable(config.ObjectVariable(map[string]config.Variable{
					"port_from":       config.IntegerVariable(65536),
					"port_to":         config.IntegerVariable(1024),
					"protocol":        config.StringVariable("TCP"),
					"description":     config.StringVariable("Range"),
					"allowed_role_id": config.StringVariable("role:abcdef"),
					"allowed_tags": config.ListVariable(
						config.StringVariable("tag:one"),
						config.StringVariable("tag:two"),
					),
				})),
			},
			ExpectError: regexp.MustCompile("Port must be <= 65535"),
		},
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/role_port_range.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("test: Role"),
				"description": config.StringVariable("Test to port's out of bound value is rejected."),
				"rules": config.ListVariable(config.ObjectVariable(map[string]config.Variable{
					"port_from":       config.IntegerVariable(1024),
					"port_to":         config.IntegerVariable(65536),
					"protocol":        config.StringVariable("TCP"),
					"description":     config.StringVariable("Range"),
					"allowed_role_id": config.StringVariable("role:abcdef"),
					"allowed_tags": config.ListVariable(
						config.StringVariable("tag:one"),
						config.StringVariable("tag:two"),
					),
				})),
			},
			ExpectError: regexp.MustCompile("Port must be <= 65535"),
		},
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/role_port_range.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("test: Role"),
				"description": config.StringVariable("Test to port must be greater than from port."),
				"rules": config.ListVariable(config.ObjectVariable(map[string]config.Variable{
					"port_from":       config.IntegerVariable(2048),
					"port_to":         config.IntegerVariable(1024),
					"protocol":        config.StringVariable("TCP"),
					"description":     config.StringVariable("Range"),
					"allowed_role_id": config.StringVariable("role:abcdef"),
					"allowed_tags": config.ListVariable(
						config.StringVariable("tag:one"),
						config.StringVariable("tag:two"),
					),
				})),
			},
			ExpectError: regexp.MustCompile("To port must be greater than from port"),
		},
	),
)
