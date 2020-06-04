package gitlab

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/xanzy/go-gitlab"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"gitlab": testAccProvider,
	}
}

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

func isRunningInCE() (bool, error) {
	isEE, err := isRunningInEE()
	return !isEE, err
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("GITLAB_TOKEN"); v == "" {
		t.Fatal("GITLAB_TOKEN must be set for acceptance tests")
	}
}
