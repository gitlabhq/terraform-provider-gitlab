package gitlab

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// testAccCompareGitLabAttribute compares an attribute in two ResourceData's for
// equivalency.
func testAccCompareGitLabAttribute(attr string, expected, received *schema.ResourceData) error {
	e := expected.Get(attr)
	r := received.Get(attr)
	switch e.(type) {
	case *schema.Set:
		if !e.(*schema.Set).Equal(r) {
			return fmt.Errorf(`attribute set %s expected "%+v" received "%+v"`, attr, e, r)
		}
	default:
		// Stringify to check because of type differences
		if fmt.Sprintf("%v", e) != fmt.Sprintf("%v", r) {
			return fmt.Errorf(`attribute %s expected "%+v" received "%+v"`, attr, e, r)
		}
	}
	return nil
}

func testAccIsSkippedAttribute(needle string, haystack []string) bool {
	if needle != "" {
		for _, hay := range haystack {
			if hay == needle {
				return true
			}
		}
	}

	return false
}
