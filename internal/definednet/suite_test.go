package definednet_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/samber/lo"
	"github.com/sendsmaily/terraform-provider-definednet/internal/definednet"
)

var (
	server *ghttp.Server
	client definednet.Client
)

var _ = BeforeEach(func() {
	server = ghttp.NewServer()
	DeferCleanup(server.Close)

	client = lo.Must(definednet.NewClient(server.URL(), "supersecret", "test"))
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "internal/definednet")
}
