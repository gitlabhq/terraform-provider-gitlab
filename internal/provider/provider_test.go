package gitlab

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	if os.Getenv(resource.TestEnvVar) != "" {
		testAccProvider = Provider()
		if err := testAccProvider.Configure(context.TODO(), &terraform.ResourceConfig{}); err != nil {
			panic(fmt.Sprintf("%#v", err)) // lintignore: R009 // TODO: Resolve this tfproviderlint issue
		}
		testAccProviders = map[string]*schema.Provider{
			"gitlab": testAccProvider,
		}
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("GITLAB_TOKEN"); v == "" {
		t.Fatal("GITLAB_TOKEN must be set for acceptance tests")
	}
}
