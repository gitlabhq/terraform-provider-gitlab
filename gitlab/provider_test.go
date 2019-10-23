package gitlab

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
	"testing"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	if os.Getenv(resource.TestEnvVar) != "" {
		testAccProvider = Provider().(*schema.Provider)
		if err := testAccProvider.Configure(&terraform.ResourceConfig{}); err != nil {
			panic(err)
		}
		testAccProviders = map[string]terraform.ResourceProvider{
			"gitlab": testAccProvider,
		}
	}
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

func testGitLabLicensePreCheck(t *testing.T) {
	if v := os.Getenv("GITLAB_LICENSE_FILE"); v == "" {
		t.Skipf("GITLAB_LICENSE_FILE must be set to run EE tests.")
	}
}
