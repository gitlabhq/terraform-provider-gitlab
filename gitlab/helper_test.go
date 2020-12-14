package gitlab

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Fourcast/go-gitlab"
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

// Returns true if the acceptance test is running Gitlab EE.
// Meant to be used as SkipFunc to skip tests that work only on Gitlab CE.
func isRunningInEE() (bool, error) {
	if conn, ok := testAccProvider.Meta().(*gitlab.Client); ok {
		version, _, err := conn.Version.GetVersion()
		if err != nil {
			return false, err
		}
		if strings.Contains(version.String(), "-ee") {
			return true, nil
		}
	} else {
		return false, errors.New("Provider not initialized, unable to get GitLab connection")
	}
	return false, nil
}

// Returns true if the acceptance test is running Gitlab CE.
// Meant to be used as SkipFunc to skip tests that work only on Gitlab EE.
func isRunningInCE() (bool, error) {
	isEE, err := isRunningInEE()
	return !isEE, err
}
