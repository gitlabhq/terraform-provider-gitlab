package gitlab

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabInstanceVariable_basic(t *testing.T) {
	var instanceVariable gitlab.InstanceVariable
	rString := acctest.RandString(5)

	// lintignore: AT001 // TODO: Resolve this tfproviderlint issue
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a variable with default options
			{
				Config: testAccGitlabInstanceVariableConfig(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabInstanceVariableExists("gitlab_instance_variable.foo", &instanceVariable),
					testAccCheckGitlabInstanceVariableAttributes(&instanceVariable, &testAccGitlabInstanceVariableExpectedAttributes{
						Key:   fmt.Sprintf("key_%s", rString),
						Value: fmt.Sprintf("value-%s", rString),
					}),
				),
			},
			// Update the instance variable to toggle all the values to their inverse
			{
				Config: testAccGitlabInstanceVariableUpdateConfig(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabInstanceVariableExists("gitlab_instance_variable.foo", &instanceVariable),
					testAccCheckGitlabInstanceVariableAttributes(&instanceVariable, &testAccGitlabInstanceVariableExpectedAttributes{
						Key:       fmt.Sprintf("key_%s", rString),
						Value:     fmt.Sprintf("value-inverse-%s", rString),
						Protected: true,
					}),
				),
			},
			// Update the instance variable to toggle the options back
			{
				Config: testAccGitlabInstanceVariableConfig(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabInstanceVariableExists("gitlab_instance_variable.foo", &instanceVariable),
					testAccCheckGitlabInstanceVariableAttributes(&instanceVariable, &testAccGitlabInstanceVariableExpectedAttributes{
						Key:       fmt.Sprintf("key_%s", rString),
						Value:     fmt.Sprintf("value-%s", rString),
						Protected: false,
					}),
				),
			},
		},
	})
}

func testAccCheckGitlabInstanceVariableExists(n string, instanceVariable *gitlab.InstanceVariable) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		key := rs.Primary.Attributes["key"]
		if key == "" {
			return fmt.Errorf("No variable key is set")
		}
		conn := testAccProvider.Meta().(*gitlab.Client)

		gotVariable, _, err := conn.InstanceVariables.GetVariable(key)
		if err != nil {
			return err
		}
		*instanceVariable = *gotVariable
		return nil
	}
}

type testAccGitlabInstanceVariableExpectedAttributes struct {
	Key       string
	Value     string
	Protected bool
	Masked    bool
}

func testAccCheckGitlabInstanceVariableAttributes(variable *gitlab.InstanceVariable, want *testAccGitlabInstanceVariableExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if variable.Key != want.Key {
			return fmt.Errorf("got key %s; want %s", variable.Key, want.Key)
		}

		if variable.Value != want.Value {
			return fmt.Errorf("got value %s; value %s", variable.Value, want.Value)
		}

		if variable.Protected != want.Protected {
			return fmt.Errorf("got protected %t; want %t", variable.Protected, want.Protected)
		}

		if variable.Masked != want.Masked {
			return fmt.Errorf("got masked %t; want %t", variable.Masked, want.Masked)
		}

		return nil
	}
}

func testAccGitlabInstanceVariableConfig(rString string) string {
	return fmt.Sprintf(`
resource "gitlab_instance_variable" "foo" {
  key = "key_%s"
  value = "value-%s"
  variable_type = "file"
  masked = false
}
	`, rString, rString)
}

func testAccGitlabInstanceVariableUpdateConfig(rString string) string {
	return fmt.Sprintf(`
resource "gitlab_instance_variable" "foo" {
  key = "key_%s"
  value = "value-inverse-%s"
  protected = true
  masked = false
}
	`, rString, rString)
}
