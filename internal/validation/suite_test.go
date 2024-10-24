package validation_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "internal/validation")
}

func GetDiagnosticsMessage(diags diag.Diagnostics) string {
	GinkgoHelper()
	return strings.Join(lo.Map(diags, func(d diag.Diagnostic, _ int) string {
		return fmt.Sprintf("%s: %s", d.Summary(), d.Detail())
	}), "; ")
}
