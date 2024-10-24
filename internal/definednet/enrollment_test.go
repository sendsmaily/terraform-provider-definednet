package definednet_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	. "github.com/onsi/gomega/gstruct"
	"github.com/sendsmaily/terraform-provider-definednet/internal/definednet"
)

var _ = Describe("creating host enrollments", func() {
	Specify("enrollments are created on Defined.net", func(ctx SpecContext) {
		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest(http.MethodPost, "/v1/host-and-enrollment-code"),
			ghttp.VerifyJSONRepresenting(map[string]any{
				"networkID":       "network-id",
				"roleID":          "role-id",
				"name":            "host.defined.test",
				"ipAddress":       "10.0.0.1",
				"staticAddresses": []string{"127.0.0.1:8484", "172.16.0.1:8484"},
				"listenPort":      8484,
				"isLighthouse":    true,
				"isRelay":         true,
				"tags":            []string{"tag:one", "tag:two"},
			}),
			ghttp.RespondWithJSONEncoded(http.StatusOK, map[string]any{}),
		))

		Expect(definednet.CreateEnrollment(ctx, client, definednet.CreateEnrollmentRequest{
			NetworkID:       "network-id",
			RoleID:          "role-id",
			Name:            "host.defined.test",
			IPAddress:       "10.0.0.1",
			StaticAddresses: []string{"127.0.0.1:8484", "172.16.0.1:8484"},
			ListenPort:      8484,
			IsLighthouse:    true,
			IsRelay:         true,
			Tags:            []string{"tag:one", "tag:two"},
		})).Error().NotTo(HaveOccurred())

		Expect(server.ReceivedRequests()).NotTo(BeEmpty(), "assert sanity")
	})

	Specify("Defined.net responses are returned", func(ctx SpecContext) {
		server.AppendHandlers(ghttp.RespondWith(http.StatusOK, enrollmentJSONResponse))

		Expect(definednet.CreateEnrollment(ctx, client, definednet.CreateEnrollmentRequest{})).
			To(PointTo(MatchAllFields(Fields{
				"Host": MatchAllFields(Fields{
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
				}),
				"EnrollmentCode": MatchAllFields(Fields{
					"Code":            Equal("supersecret"),
					"LifetimeSeconds": Equal(300),
				}),
			})))

		Expect(server.ReceivedRequests()).NotTo(BeEmpty(), "assert sanity")
	})
})

var enrollmentJSONResponse = `{
  "data": {
    "host": {
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
    "enrollmentCode": {
      "code": "supersecret",
      "lifetimeSeconds": 300
    }
  },
  "metadata": {}
}`
