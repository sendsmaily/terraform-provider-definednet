package definednet_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	. "github.com/onsi/gomega/gstruct"
	"github.com/sendsmaily/terraform-provider-definednet/internal/definednet"
)

var _ = Describe("creating enrollment codes", func() {
	Specify("enrollment codes are created on Defined.net", func(ctx SpecContext) {
		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest(http.MethodPost, "/v1/hosts/host-id/enrollment-code"),
			ghttp.RespondWithJSONEncoded(http.StatusOK, map[string]any{}),
		))

		Expect(definednet.CreateEnrollmentCode(ctx, client, definednet.CreateEnrollmentCodeRequest{
			ID: "host-id",
		})).Error().NotTo(HaveOccurred())
		Expect(server.ReceivedRequests()).NotTo(BeEmpty(), "assert sanity")
	})

	Specify("Defined.net responses are returned", func(ctx SpecContext) {
		server.AppendHandlers(ghttp.RespondWith(http.StatusOK, enrollmentCodeJSONResponse))
		Expect(definednet.CreateEnrollmentCode(ctx, client, definednet.CreateEnrollmentCodeRequest{})).
			To(PointTo(MatchAllFields(Fields{
				"Code":            Equal("supersecret"),
				"LifetimeSeconds": Equal(300),
			})))
		Expect(server.ReceivedRequests()).NotTo(BeEmpty(), "assert sanity")
	})
})

var enrollmentCodeJSONResponse = `{
	"data": {
	  "code": "supersecret",
	  "lifetimeSeconds": 300
	},
	"metadata": {}
  }`
