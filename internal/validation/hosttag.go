package validation

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var tagPattern = regexp.MustCompile(`^[a-z_-]+?:[a-z_-]+?$`)

// HostTag validates the passed host tag conforms to Defined.net required format.
func HostTag() validator.String {
	return stringvalidator.RegexMatches(tagPattern, "must be in the format KEY:VALUE")
}
