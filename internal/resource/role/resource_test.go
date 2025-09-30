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
				"name": config.StringVariable("test: Role"),
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
				"rules": config.SetVariable(
					config.ObjectVariable(map[string]config.Variable{
						"port":            config.IntegerVariable(22),
						"protocol":        config.StringVariable("TCP"),
						"description":     config.StringVariable("SSH access"),
						"allowed_role_id": config.StringVariable("role:abcdef"),
						"allowed_tags": config.SetVariable(
							config.StringVariable("tag:one"),
							config.StringVariable("tag:two"),
						),
					}),
					config.ObjectVariable(map[string]config.Variable{
						"port":            config.IntegerVariable(443),
						"protocol":        config.StringVariable("TCP"),
						"description":     config.StringVariable("HTTPS access"),
						"allowed_role_id": config.StringVariable("role:123456"),
						"allowed_tags": config.SetVariable(
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
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					"definednet_role.test",
					tfjsonpath.New("rule"),
					knownvalue.SetExact([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"port":            knownvalue.Int32Exact(22),
							"port_range":      knownvalue.Null(),
							"protocol":        knownvalue.StringExact("TCP"),
							"description":     knownvalue.StringExact("SSH access"),
							"allowed_role_id": knownvalue.StringExact("role:abcdef"),
							"allowed_tags": knownvalue.SetExact([]knownvalue.Check{
								knownvalue.StringExact("tag:one"),
								knownvalue.StringExact("tag:two"),
							}),
						}),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"port":            knownvalue.Int32Exact(443),
							"port_range":      knownvalue.Null(),
							"protocol":        knownvalue.StringExact("TCP"),
							"description":     knownvalue.StringExact("HTTPS access"),
							"allowed_role_id": knownvalue.StringExact("role:123456"),
							"allowed_tags": knownvalue.SetExact([]knownvalue.Check{
								knownvalue.StringExact("tag:https_one"),
								knownvalue.StringExact("tag:https_two"),
							}),
						}),
					}),
				),
			},
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
				"rules": config.SetVariable(
					config.ObjectVariable(map[string]config.Variable{
						"port":            config.IntegerVariable(22),
						"protocol":        config.StringVariable("TCP"),
						"description":     config.StringVariable("SSH access"),
						"allowed_role_id": config.StringVariable("role:abcdef"),
						"allowed_tags": config.SetVariable(
							config.StringVariable("tag:one"),
							config.StringVariable("tag:two"),
						),
					}),
					config.ObjectVariable(map[string]config.Variable{
						"port":            config.IntegerVariable(443),
						"protocol":        config.StringVariable("TCP"),
						"description":     config.StringVariable("HTTPS access"),
						"allowed_role_id": config.StringVariable("role:123456"),
						"allowed_tags": config.SetVariable(
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
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					"definednet_role.test",
					tfjsonpath.New("rule"),
					knownvalue.SetExact([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"port":            knownvalue.Int32Exact(22),
							"port_range":      knownvalue.Null(),
							"protocol":        knownvalue.StringExact("TCP"),
							"description":     knownvalue.StringExact("SSH access"),
							"allowed_role_id": knownvalue.StringExact("role:abcdef"),
							"allowed_tags": knownvalue.SetExact([]knownvalue.Check{
								knownvalue.StringExact("tag:one"),
								knownvalue.StringExact("tag:two"),
							}),
						}),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"port":            knownvalue.Int32Exact(443),
							"port_range":      knownvalue.Null(),
							"protocol":        knownvalue.StringExact("TCP"),
							"description":     knownvalue.StringExact("HTTPS access"),
							"allowed_role_id": knownvalue.StringExact("role:123456"),
							"allowed_tags": knownvalue.SetExact([]knownvalue.Check{
								knownvalue.StringExact("tag:https_one"),
								knownvalue.StringExact("tag:https_two"),
							}),
						}),
					}),
				),
			},
		},
	),
	Entry("assert importing the role populates the state",
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/role_port.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("test: Role"),
				"description": config.StringVariable("Role's description"),
				"rules": config.SetVariable(
					config.ObjectVariable(map[string]config.Variable{
						"port":            config.IntegerVariable(22),
						"protocol":        config.StringVariable("TCP"),
						"description":     config.StringVariable("SSH access"),
						"allowed_role_id": config.StringVariable("role:abcdef"),
						"allowed_tags": config.SetVariable(
							config.StringVariable("tag:one"),
							config.StringVariable("tag:two"),
						),
					}),
					config.ObjectVariable(map[string]config.Variable{
						"port":            config.IntegerVariable(443),
						"protocol":        config.StringVariable("TCP"),
						"description":     config.StringVariable("HTTPS access"),
						"allowed_role_id": config.StringVariable("role:123456"),
						"allowed_tags": config.SetVariable(
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
				"name": config.StringVariable("test: Role"),
			},
			ResourceName:      "definednet_role.test",
			ImportState:       true,
			ImportStateVerify: true,
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
				"rules": config.SetVariable(
					config.ObjectVariable(map[string]config.Variable{
						"port_from":       config.IntegerVariable(1024),
						"port_to":         config.IntegerVariable(2048),
						"protocol":        config.StringVariable("TCP"),
						"description":     config.StringVariable("Range one"),
						"allowed_role_id": config.StringVariable("role:abcdef"),
						"allowed_tags": config.SetVariable(
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
						"allowed_tags": config.SetVariable(
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
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					"definednet_role.test",
					tfjsonpath.New("rule"),
					knownvalue.SetExact([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"port": knownvalue.Null(),
							"port_range": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"from": knownvalue.Int32Exact(1024),
								"to":   knownvalue.Int32Exact(2048),
							}),
							"protocol":        knownvalue.StringExact("TCP"),
							"description":     knownvalue.StringExact("Range one"),
							"allowed_role_id": knownvalue.StringExact("role:abcdef"),
							"allowed_tags": knownvalue.SetExact([]knownvalue.Check{
								knownvalue.StringExact("tag:one"),
								knownvalue.StringExact("tag:two"),
							}),
						}),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"port": knownvalue.Null(),
							"port_range": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"from": knownvalue.Int32Exact(4096),
								"to":   knownvalue.Int32Exact(8192),
							}),
							"protocol":        knownvalue.StringExact("TCP"),
							"description":     knownvalue.StringExact("Range two"),
							"allowed_role_id": knownvalue.StringExact("role:123456"),
							"allowed_tags": knownvalue.SetExact([]knownvalue.Check{
								knownvalue.StringExact("tag:https_one"),
								knownvalue.StringExact("tag:https_two"),
							}),
						}),
					}),
				),
			},
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
				"rules": config.SetVariable(
					config.ObjectVariable(map[string]config.Variable{
						"port_from":       config.IntegerVariable(1024),
						"port_to":         config.IntegerVariable(2048),
						"protocol":        config.StringVariable("TCP"),
						"description":     config.StringVariable("Range one"),
						"allowed_role_id": config.StringVariable("role:abcdef"),
						"allowed_tags": config.SetVariable(
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
						"allowed_tags": config.SetVariable(
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
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					"definednet_role.test",
					tfjsonpath.New("rule"),
					knownvalue.SetExact([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"port": knownvalue.Null(),
							"port_range": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"from": knownvalue.Int32Exact(1024),
								"to":   knownvalue.Int32Exact(2048),
							}),
							"protocol":        knownvalue.StringExact("TCP"),
							"description":     knownvalue.StringExact("Range one"),
							"allowed_role_id": knownvalue.StringExact("role:abcdef"),
							"allowed_tags": knownvalue.SetExact([]knownvalue.Check{
								knownvalue.StringExact("tag:one"),
								knownvalue.StringExact("tag:two"),
							}),
						}),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"port": knownvalue.Null(),
							"port_range": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"from": knownvalue.Int32Exact(4096),
								"to":   knownvalue.Int32Exact(8192),
							}),
							"protocol":        knownvalue.StringExact("TCP"),
							"description":     knownvalue.StringExact("Range two"),
							"allowed_role_id": knownvalue.StringExact("role:123456"),
							"allowed_tags": knownvalue.SetExact([]knownvalue.Check{
								knownvalue.StringExact("tag:https_one"),
								knownvalue.StringExact("tag:https_two"),
							}),
						}),
					}),
				),
			},
		},
	),
	Entry("assert importing the role populates the state",
		resource.TestStep{
			ConfigFile: config.StaticFile("testdata/role_port_range.tf"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable("test: Role"),
				"description": config.StringVariable("Role's description"),
				"rules": config.SetVariable(
					config.ObjectVariable(map[string]config.Variable{
						"port_from":       config.IntegerVariable(1024),
						"port_to":         config.IntegerVariable(2048),
						"protocol":        config.StringVariable("TCP"),
						"description":     config.StringVariable("Range one"),
						"allowed_role_id": config.StringVariable("role:abcdef"),
						"allowed_tags": config.SetVariable(
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
						"allowed_tags": config.SetVariable(
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
				"name": config.StringVariable("test: Role"),
			},
			ResourceName:      "definednet_role.test",
			ImportState:       true,
			ImportStateVerify: true,
		},
	),
)
