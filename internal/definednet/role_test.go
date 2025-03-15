package definednet_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	. "github.com/onsi/gomega/gstruct"
	"github.com/sendsmaily/terraform-provider-definednet/internal/definednet"
)

var _ = Describe("creating roles", func() {
	Specify("roles are created on Defined.net", func(ctx SpecContext) {
		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest(http.MethodPost, "/v1/roles"),
			ghttp.VerifyJSONRepresenting(map[string]any{
				"name":        "test: Role",
				"description": "Role's description",
				"firewallRules": []map[string]any{
					{
						"protocol":      "TCP",
						"description":   "Allow SSH access",
						"allowedRoleID": "allowed-role-id",
						"portRange": map[string]int{
							"from": 22,
							"to":   22,
						},
					},
					{
						"protocol":    "ANY",
						"description": "Allow ephemeral ports",
						"allowedTags": []string{
							"tag:one",
							"tag:two",
						},
						"portRange": map[string]int{
							"from": 32768,
							"to":   65535,
						},
					},
				},
			}),
			ghttp.RespondWithJSONEncoded(http.StatusOK, map[string]any{}),
		))

		Expect(definednet.CreateRole(ctx, client, definednet.CreateRoleRequest{
			Name:        "test: Role",
			Description: "Role's description",
			FirewallRules: []definednet.FirewallRule{
				{
					Protocol:      "TCP",
					Description:   "Allow SSH access",
					AllowedRoleID: "allowed-role-id",
					PortRange: definednet.PortRange{
						From: 22,
						To:   22,
					},
				},
				{
					Protocol:    "ANY",
					Description: "Allow ephemeral ports",
					AllowedTags: []string{
						"tag:one",
						"tag:two",
					},
					PortRange: definednet.PortRange{
						From: 32768,
						To:   65535,
					},
				},
			},
		})).Error().NotTo(HaveOccurred())
	})
})

var _ = Describe("deleting roles", func() {
	Specify("roles are deleted from Defined.net", func(ctx SpecContext) {
		server.AppendHandlers(ghttp.VerifyRequest(http.MethodDelete, "/v1/roles/role-id"))
		Expect(definednet.DeleteRole(ctx, client, definednet.DeleteRoleRequest{
			ID: "role-id",
		})).To(Succeed())
		Expect(server.ReceivedRequests()).NotTo(BeEmpty(), "assert sanity")
	})
})

var _ = Describe("getting roles", func() {
	Specify("Defined.net responses are returned", func(ctx SpecContext) {
		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest(http.MethodGet, "/v1/roles/role-id"),
			ghttp.RespondWith(http.StatusOK, roleJSONResponse)),
		)

		Expect(definednet.GetRole(ctx, client, definednet.GetRoleRequest{
			ID: "role-id",
		})).
			To(PointTo(MatchAllFields(Fields{
				"ID":          Equal("role-id"),
				"Name":        Equal("test: Role"),
				"Description": Equal("Role's description"),
				"FirewallRules": HaveExactElements(
					MatchAllFields(Fields{
						"Protocol":      Equal("TCP"),
						"Description":   Equal("Allow SSH access"),
						"AllowedRoleID": Equal("allowed-role-id"),
						"AllowedTags":   BeEmpty(),
						"PortRange": MatchAllFields(Fields{
							"From": Equal(22),
							"To":   Equal(22),
						}),
					}),
					MatchAllFields(Fields{
						"Protocol":      Equal("ANY"),
						"Description":   Equal("Allow ephemeral ports"),
						"AllowedRoleID": BeEmpty(),
						"AllowedTags": HaveExactElements(
							"tag:one",
							"tag:two",
						),
						"PortRange": MatchAllFields(Fields{
							"From": Equal(32768),
							"To":   Equal(65535),
						}),
					}),
				),
			})))

		Expect(server.ReceivedRequests()).NotTo(BeEmpty(), "assert sanity")
	})
})

var _ = Describe("updating roles", func() {
	Specify("roles are updated on Defined.net", func(ctx SpecContext) {
		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest(http.MethodPut, "/v1/roles/role-id"),
			ghttp.VerifyJSONRepresenting(map[string]any{
				"name":        "test: Role",
				"description": "Role's description",
				"firewallRules": []map[string]any{
					{
						"protocol":      "TCP",
						"description":   "Allow SSH access",
						"allowedRoleID": "allowed-role-id",
						"portRange": map[string]int{
							"from": 22,
							"to":   22,
						},
					},
					{
						"protocol":    "ANY",
						"description": "Allow ephemeral ports",
						"allowedTags": []string{
							"tag:one",
							"tag:two",
						},
						"portRange": map[string]int{
							"from": 32768,
							"to":   65535,
						},
					},
				},
			}),
			ghttp.RespondWithJSONEncoded(http.StatusOK, map[string]any{}),
		))

		Expect(definednet.UpdateRole(ctx, client, definednet.UpdateRoleRequest{
			ID:          "role-id",
			Name:        "test: Role",
			Description: "Role's description",
			FirewallRules: []definednet.FirewallRule{
				{
					Protocol:      "TCP",
					Description:   "Allow SSH access",
					AllowedRoleID: "allowed-role-id",
					PortRange: definednet.PortRange{
						From: 22,
						To:   22,
					},
				},
				{
					Protocol:    "ANY",
					Description: "Allow ephemeral ports",
					AllowedTags: []string{
						"tag:one",
						"tag:two",
					},
					PortRange: definednet.PortRange{
						From: 32768,
						To:   65535,
					},
				},
			},
		})).Error().NotTo(HaveOccurred())
		Expect(server.ReceivedRequests()).NotTo(BeEmpty(), "assert sanity")
	})

	Specify("Defined.net responses are returned", func(ctx SpecContext) {
		server.AppendHandlers(ghttp.RespondWith(http.StatusOK, roleJSONResponse))

		Expect(definednet.UpdateRole(ctx, client, definednet.UpdateRoleRequest{})).
			To(PointTo(MatchAllFields(Fields{
				"ID":          Equal("role-id"),
				"Name":        Equal("test: Role"),
				"Description": Equal("Role's description"),
				"FirewallRules": HaveExactElements(
					MatchAllFields(Fields{
						"Protocol":      Equal("TCP"),
						"Description":   Equal("Allow SSH access"),
						"AllowedRoleID": Equal("allowed-role-id"),
						"AllowedTags":   BeEmpty(),
						"PortRange": MatchAllFields(Fields{
							"From": Equal(22),
							"To":   Equal(22),
						}),
					}),
					MatchAllFields(Fields{
						"Protocol":      Equal("ANY"),
						"Description":   Equal("Allow ephemeral ports"),
						"AllowedRoleID": BeEmpty(),
						"AllowedTags": HaveExactElements(
							"tag:one",
							"tag:two",
						),
						"PortRange": MatchAllFields(Fields{
							"From": Equal(32768),
							"To":   Equal(65535),
						}),
					}),
				),
			})))

		Expect(server.ReceivedRequests()).NotTo(BeEmpty(), "assert sanity")
	})
})

var roleJSONResponse = `{
  "data": {
    "id": "role-id",
    "name": "test: Role",
    "description": "Role's description",
    "createdAt": "2023-02-15T13:59:09Z",
    "modifiedAt": "2023-02-15T13:59:09Z",
    "firewallRules": [
      {
        "protocol": "TCP",
        "description": "Allow SSH access",
        "allowedRoleID": "allowed-role-id",
        "portRange": {
          "from": 22,
          "to": 22
        }
      },
      {
        "protocol": "ANY",
        "description": "Allow ephemeral ports",
        "allowedTags": [
          "tag:one",
          "tag:two"
        ],
        "portRange": {
          "from": 32768,
          "to": 65535
        }
      }
    ]
  },
  "metadata": {}
}`
