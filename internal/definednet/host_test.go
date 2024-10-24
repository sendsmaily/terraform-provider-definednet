package definednet_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	. "github.com/onsi/gomega/gstruct"
	"github.com/sendsmaily/terraform-provider-definednet/internal/definednet"
)

var _ = Describe("deleting hosts", func() {
	Specify("hosts are deleted from Defined.net", func(ctx SpecContext) {
		server.AppendHandlers(ghttp.VerifyRequest(http.MethodDelete, "/v1/hosts/host-id"))
		Expect(definednet.DeleteHost(ctx, client, definednet.DeleteHostRequest{
			ID: "host-id",
		})).To(Succeed())
		Expect(server.ReceivedRequests()).NotTo(BeEmpty(), "assert sanity")
	})
})

var _ = Describe("getting hosts", func() {
	Specify("Defined.net responses are returned", func(ctx SpecContext) {
		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest(http.MethodGet, "/v1/hosts/host-id"),
			ghttp.RespondWith(http.StatusOK, hostJSONResponse)),
		)

		Expect(definednet.GetHost(ctx, client, definednet.GetHostRequest{
			ID: "host-id",
		})).To(PointTo(MatchAllFields(Fields{
			"ID":              Equal("host-id"),
			"NetworkID":       Equal("network-id"),
			"RoleID":          Equal("role-id"),
			"Name":            Equal("host.defined.test"),
			"IPAddress":       Equal("10.0.0.1"),
			"StaticAddresses": HaveExactElements("127.0.0.1:8484", "172.16.0.1:8484"),
			"ListenPort":      Equal(8484),
			"IsLighthouse":    BeTrue(),
			"IsRelay":         BeTrue(),
			"Tags":            HaveExactElements("tag:one", "tag:two"),
		})))
		Expect(server.ReceivedRequests()).NotTo(BeEmpty(), "assert sanity")
	})
})

var _ = Describe("updating hosts", func() {
	Specify("hosts are updated on Defined.net", func(ctx SpecContext) {
		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest(http.MethodPut, "/v2/hosts/host-id"),
			ghttp.VerifyJSONRepresenting(map[string]any{
				"roleID":          "role-id",
				"name":            "host.defined.test",
				"staticAddresses": []string{"127.0.0.1:8484", "172.16.0.1:8484"},
				"listenPort":      8484,
				"tags":            []string{"tag:one", "tag:two"},
			}),
			ghttp.RespondWithJSONEncoded(http.StatusOK, map[string]any{}),
		))

		Expect(definednet.UpdateHost(ctx, client, definednet.UpdateHostRequest{
			ID:              "host-id",
			RoleID:          "role-id",
			Name:            "host.defined.test",
			StaticAddresses: []string{"127.0.0.1:8484", "172.16.0.1:8484"},
			ListenPort:      8484,
			Tags:            []string{"tag:one", "tag:two"},
		})).Error().NotTo(HaveOccurred())
		Expect(server.ReceivedRequests()).NotTo(BeEmpty(), "assert sanity")
	})

	Specify("Defined.net responses are returned", func(ctx SpecContext) {
		server.AppendHandlers(ghttp.RespondWith(http.StatusOK, hostJSONResponse))

		Expect(definednet.UpdateHost(ctx, client, definednet.UpdateHostRequest{})).
			To(PointTo(MatchAllFields(Fields{
				"ID":              Equal("host-id"),
				"NetworkID":       Equal("network-id"),
				"RoleID":          Equal("role-id"),
				"Name":            Equal("host.defined.test"),
				"IPAddress":       Equal("10.0.0.1"),
				"StaticAddresses": HaveExactElements("127.0.0.1:8484", "172.16.0.1:8484"),
				"ListenPort":      Equal(8484),
				"IsLighthouse":    BeTrue(),
				"IsRelay":         BeTrue(),
				"Tags":            HaveExactElements("tag:one", "tag:two"),
			})))

		Expect(server.ReceivedRequests()).NotTo(BeEmpty(), "assert sanity")
	})
})

var hostJSONResponse = `{
  "data": {
    "createdAt": "2024-10-18T08:37:30Z",
    "id": "host-id",
    "ipAddress": "10.0.0.1",
    "isBlocked": false,
    "isLighthouse": true,
    "isRelay": true,
    "listenPort": 8484,
    "name": "host.defined.test",
    "networkID": "network-id",
    "organizationID": "org-id",
    "roleID": "role-id",
    "staticAddresses": [
      "127.0.0.1:8484",
      "172.16.0.1:8484"
    ],
	"tags": [
	  "tag:one",
	  "tag:two"
	],
    "metadata": {
      "lastSeenAt": "2023-01-25T18:15:27Z",
      "platform": "dnclient",
      "updateAvailable": false,
      "version": "0.1.9"
    }
  },
  "metadata": {}
}`
