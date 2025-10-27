package lighthouse_test

import (
	"testing"

	tfprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sendsmaily/terraform-provider-definednet/internal/definednet"
	"github.com/sendsmaily/terraform-provider-definednet/internal/provider"
	fakeserver "github.com/sendsmaily/terraform-provider-definednet/internal/testing/server"
)

var providerFactory func() tfprovider.Provider

var _ = BeforeEach(func() {
	server := fakeserver.New()
	DeferCleanup(server.Close)

	providerFactory = provider.New(
		func(string, string, string) (definednet.Client, error) {
			return server.Client(), nil
		},
		"test",
	)
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "internal/resource/lighthouse")
}
