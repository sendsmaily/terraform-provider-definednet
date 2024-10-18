package definednet_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/onsi/gomega/gstruct"
	"github.com/sendsmaily/terraform-provider-definednet/internal/definednet"
)

var _ = Describe("API client's invariants", func() {
	Specify("client requires a valid endpoint URL", func() {
		Expect(definednet.NewClient("http://localhost:8000/api/v1", "supersecret")).Error().
			NotTo(HaveOccurred(), "assert sanity")

		Expect(definednet.NewClient("", "supersecret")).Error().
			To(MatchError("endpoint URL must be set"), "missing endpoint URL")

		Expect(definednet.NewClient("  ", "supersecret")).Error().
			To(MatchError("endpoint URL must be set"), "empty endpoint URL")

		Expect(definednet.NewClient("invalid.url", "supersecret")).Error().
			To(MatchError(`parse "invalid.url": invalid URI for request`), "invalid endpoint URL")
	})

	Specify("client requires a non-zero authorization token", func() {
		Expect(definednet.NewClient("http://localhost:8000/api/v1", "supersecret")).Error().
			NotTo(HaveOccurred(), "assert sanity")

		Expect(definednet.NewClient("http://localhost:8000/api/v1", "")).Error().
			To(MatchError("authorization token must be set"), "missing token")

		Expect(definednet.NewClient("http://localhost:8000/api/v1", "  ")).Error().
			To(MatchError("authorization token must be set"), "empty token")
	})
})

var _ = Describe("executing API requests", func() {
	Context("request headers", func() {
		Specify("all requests are executed with common HTTP headers", func(ctx SpecContext) {
			server.AppendHandlers(ghttp.RespondWith(http.StatusOK, nil))
			Expect(client.Do(ctx, http.MethodGet, []string{}, nil, nil)).To(Succeed())

			Expect(server.ReceivedRequests()).
				To(HaveExactElements(
					HaveField("Header", gstruct.MatchKeys(gstruct.IgnoreExtras, gstruct.Keys{
						"Accept":        HaveExactElements("application/json"),
						"Authorization": HaveExactElements("Bearer supersecret"),
						"User-Agent":    HaveExactElements("Terraform-smaily-definednet/0.1.0"),
					})),
				))
		})

		When("request does not specify body payload", func() {
			Specify("Content-Type header is not set on the request", func(ctx SpecContext) {
				server.AppendHandlers(ghttp.RespondWith(http.StatusOK, nil))
				Expect(client.Do(ctx, http.MethodPost, []string{}, nil, nil)).To(Succeed())

				Expect(server.ReceivedRequests()).
					To(HaveExactElements(HaveField("Header", Not(HaveKey("Content-Type")))))
			})
		})

		When("request specifies body payload", func() {
			Specify("Content-Type header is set on the request", func(ctx SpecContext) {
				server.AppendHandlers(ghttp.RespondWith(http.StatusOK, nil))
				Expect(client.Do(ctx, http.MethodPost, []string{}, map[string]any{}, nil)).To(Succeed())

				Expect(server.ReceivedRequests()).
					To(HaveExactElements(HaveField("Header",
						HaveKeyWithValue("Content-Type", HaveExactElements("application/json")),
					)))
			})
		})
	})

	Context("request path", func() {
		Specify("paths are compiled from passed components", func(ctx SpecContext) {
			components := []string{"path", "compiled", "from", "components"}
			server.AppendHandlers(ghttp.VerifyRequest(http.MethodGet, "/path/compiled/from/components"))
			Expect(client.Do(ctx, http.MethodGet, components, nil, nil)).To(Succeed())
			Expect(server.ReceivedRequests()).NotTo(BeEmpty(), "assert sanity")
		})

		Specify("paths are protected from URL injection", func(ctx SpecContext) {
			components := []string{"spaced value", "path/../../traversal/attempt"}
			server.AppendHandlers(ghttp.VerifyRequest(http.MethodGet, "/spaced value/path/../../traversal/attempt"))
			Expect(client.Do(ctx, http.MethodGet, components, nil, nil)).To(Succeed())
			Expect(server.ReceivedRequests()).NotTo(BeEmpty(), "assert sanity")
		})
	})

	Context("request body", func() {
		type (
			nested struct {
				Field string `json:"field"`
			}

			request struct {
				Field  string `json:"field"`
				Nested nested `json:"nested"`
			}
		)

		Specify("request payload containers are JSON-encoded into HTTP body", func(ctx SpecContext) {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest(http.MethodPost, "/"),
				ghttp.VerifyJSON(`{"field":"value","nested":{"field":"nested_value"}}`),
				ghttp.RespondWith(http.StatusOK, nil),
			))

			Expect(client.Do(
				ctx,
				http.MethodPost,
				[]string{},
				request{
					Field: "value",
					Nested: nested{
						Field: "nested_value",
					},
				}, nil,
			)).To(Succeed())

			Expect(server.ReceivedRequests()).NotTo(BeEmpty(), "assert sanity")
		})
	})
})

var _ = Describe("handling HTTP API success responses", func() {
	type response struct {
		Field  string `json:"field"`
		Nested struct {
			Field string `json:"field"`
		} `json:"nested"`
	}

	Specify("API responses are deserialized into passed response payload container", func(ctx SpecContext) {
		var container response

		server.AppendHandlers(ghttp.RespondWithJSONEncoded(http.StatusOK, map[string]any{
			"field": "value",
			"nested": map[string]any{
				"field": "nested_value",
			},
		}))

		Expect(client.Do(ctx, http.MethodGet, []string{}, nil, &container)).To(Succeed())
		Expect(container).To(gstruct.MatchAllFields(gstruct.Fields{
			"Field": Equal("value"),
			"Nested": gstruct.MatchAllFields(gstruct.Fields{
				"Field": Equal("nested_value"),
			}),
		}))
	})

	Specify("unexpected keys in responses are ignored", func(ctx SpecContext) {
		var container response

		server.AppendHandlers(ghttp.RespondWithJSONEncoded(http.StatusOK, map[string]any{
			"field": "value",
			"nested": map[string]any{
				"field": "nested_value",
			},
			"unexpected": "field",
		}))

		Expect(client.Do(ctx, http.MethodGet, []string{}, nil, &container)).To(Succeed())
		Expect(container).To(gstruct.MatchAllFields(gstruct.Fields{
			"Field": Equal("value"),
			"Nested": gstruct.MatchAllFields(gstruct.Fields{
				"Field": Equal("nested_value"),
			}),
		}), "assert sanity")
	})

	When("response payload container is not provided", func() {
		Specify("response payload decoding is skipped", func(ctx SpecContext) {
			server.AppendHandlers(ghttp.RespondWithJSONEncoded(http.StatusOK, map[string]any{
				"field": "value",
				"nested": map[string]any{
					"field": "nested_value",
				},
			}))

			Expect(client.Do(ctx, http.MethodGet, []string{}, nil, nil)).To(Succeed())
		})
	})
})

var _ = Describe("handling HTTP API error responses", func() {
	Specify("API errors are reported to the user", func(ctx SpecContext) {
		server.AppendHandlers(ghttp.RespondWithJSONEncoded(
			http.StatusBadRequest,
			map[string]any{
				"server": "response",
			},
		))

		Expect(client.Do(ctx, http.MethodGet, []string{}, nil, nil)).
			To(MatchError(`code=400 reason={"server":"response"}`))
	})
})
