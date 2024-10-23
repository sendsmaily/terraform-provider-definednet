package lighthouse_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sendsmaily/terraform-provider-definednet/internal/definednet"
	fakeprovider "github.com/sendsmaily/terraform-provider-definednet/internal/testing/provider"
	fakeserver "github.com/sendsmaily/terraform-provider-definednet/internal/testing/server"
)

var (
	provider *fakeprovider.Provider

	server *fakeserver.Server
	client definednet.Client
)

var _ = BeforeEach(func() {
	server = fakeserver.New()
	DeferCleanup(server.Close)

	client = server.Client()

	provider = fakeprovider.New(client)
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "internal/resource/lighthouse")
}
