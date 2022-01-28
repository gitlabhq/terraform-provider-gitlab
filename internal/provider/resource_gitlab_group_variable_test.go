package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabGroupVariable_basic(t *testing.T) {
	var groupVariable gitlab.GroupVariable
	rString := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabGroupVariableDestroy,
		Steps: []resource.TestStep{
			// Create a group and variable with default options
			{
				Config: testAccGitlabGroupVariableConfig(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupVariableExists("gitlab_group_variable.foo", &groupVariable),
					testAccCheckGitlabGroupVariableAttributes(&groupVariable, &testAccGitlabGroupVariableExpectedAttributes{
						Key:              fmt.Sprintf("key_%s", rString),
						Value:            fmt.Sprintf("value-%s", rString),
						EnvironmentScope: "*",
					}),
				),
			},
			// Update the group variable to toggle all the values to their inverse
			{
				Config: testAccGitlabGroupVariableUpdateConfig(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupVariableExists("gitlab_group_variable.foo", &groupVariable),
					testAccCheckGitlabGroupVariableAttributes(&groupVariable, &testAccGitlabGroupVariableExpectedAttributes{
						Key:              fmt.Sprintf("key_%s", rString),
						Value:            fmt.Sprintf("value-inverse-%s", rString),
						Protected:        true,
						EnvironmentScope: "*",
					}),
				),
			},
			// Update the group variable to toggle the options back
			{
				Config: testAccGitlabGroupVariableConfig(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupVariableExists("gitlab_group_variable.foo", &groupVariable),
					testAccCheckGitlabGroupVariableAttributes(&groupVariable, &testAccGitlabGroupVariableExpectedAttributes{
						Key:              fmt.Sprintf("key_%s", rString),
						Value:            fmt.Sprintf("value-%s", rString),
						Protected:        false,
						EnvironmentScope: "*",
					}),
				),
			},
		},
	})
}

func TestAccGitlabGroupVariable_scope(t *testing.T) {
	var groupVariableA, groupVariableB gitlab.GroupVariable
	rString := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabGroupVariableDestroy,
		Steps: []resource.TestStep{
			// Create a group and variables with same keys, different scopes
			{
				Config:   testAccGitlabGroupVariableScopeConfig(rString, "*", "review/*"),
				SkipFunc: isRunningInCE,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupVariableExists("gitlab_group_variable.a", &groupVariableA),
					testAccCheckGitlabGroupVariableExists("gitlab_group_variable.b", &groupVariableB),
					testAccCheckGitlabGroupVariableAttributes(&groupVariableA, &testAccGitlabGroupVariableExpectedAttributes{
						Key:              fmt.Sprintf("key_%s", rString),
						Value:            fmt.Sprintf("value-%s-a", rString),
						EnvironmentScope: "*",
					}),
					testAccCheckGitlabGroupVariableAttributes(&groupVariableB, &testAccGitlabGroupVariableExpectedAttributes{
						Key:              fmt.Sprintf("key_%s", rString),
						Value:            fmt.Sprintf("value-%s-b", rString),
						EnvironmentScope: "review/*",
					}),
				),
			},
			// Change a variable's scope
			{
				Config:   testAccGitlabGroupVariableScopeConfig(rString, "my-new-scope", "review/*"),
				SkipFunc: isRunningInCE,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupVariableExists("gitlab_group_variable.a", &groupVariableA),
					testAccCheckGitlabGroupVariableExists("gitlab_group_variable.b", &groupVariableB),
					testAccCheckGitlabGroupVariableAttributes(&groupVariableA, &testAccGitlabGroupVariableExpectedAttributes{
						Key:              fmt.Sprintf("key_%s", rString),
						Value:            fmt.Sprintf("value-%s-a", rString),
						EnvironmentScope: "my-new-scope",
					}),
					testAccCheckGitlabGroupVariableAttributes(&groupVariableB, &testAccGitlabGroupVariableExpectedAttributes{
						Key:              fmt.Sprintf("key_%s", rString),
						Value:            fmt.Sprintf("value-%s-b", rString),
						EnvironmentScope: "review/*",
					}),
				),
			},
			// Change both variables scopes at the same time
			{
				Config:   testAccGitlabGroupVariableScopeConfig(rString, "my-new-new-scope", "review/hello-world"),
				SkipFunc: isRunningInCE,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupVariableExists("gitlab_group_variable.a", &groupVariableA),
					testAccCheckGitlabGroupVariableExists("gitlab_group_variable.b", &groupVariableB),
					testAccCheckGitlabGroupVariableAttributes(&groupVariableA, &testAccGitlabGroupVariableExpectedAttributes{
						Key:              fmt.Sprintf("key_%s", rString),
						Value:            fmt.Sprintf("value-%s-a", rString),
						EnvironmentScope: "my-new-new-scope",
					}),
					testAccCheckGitlabGroupVariableAttributes(&groupVariableB, &testAccGitlabGroupVariableExpectedAttributes{
						Key:              fmt.Sprintf("key_%s", rString),
						Value:            fmt.Sprintf("value-%s-b", rString),
						EnvironmentScope: "review/hello-world",
					}),
				),
			},
		},
	})
}

func testAccCheckGitlabGroupVariableExists(n string, groupVariable *gitlab.GroupVariable) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		repoName := rs.Primary.Attributes["group"]
		if repoName == "" {
			return fmt.Errorf("No group ID is set")
		}
		key := rs.Primary.Attributes["key"]
		if key == "" {
			return fmt.Errorf("No variable key is set")
		}
		gotVariable, _, err := testGitlabClient.GroupVariables.GetVariable(repoName, key, modifyRequestAddEnvironmentFilter(rs.Primary.Attributes["environment_scope"]))
		if err != nil {
			return err
		}
		*groupVariable = *gotVariable
		return nil
	}
}

type testAccGitlabGroupVariableExpectedAttributes struct {
	Key              string
	Value            string
	Protected        bool
	Masked           bool
	EnvironmentScope string
}

func testAccCheckGitlabGroupVariableAttributes(variable *gitlab.GroupVariable, want *testAccGitlabGroupVariableExpectedAttributes) resource.TestCheckFunc {
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

		if variable.EnvironmentScope != want.EnvironmentScope {
			return fmt.Errorf("got environment_scope %s; want %s", variable.EnvironmentScope, want.EnvironmentScope)
		}

		return nil
	}
}

func testAccCheckGitlabGroupVariableDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_group" {
			continue
		}

		_, _, err := testGitlabClient.Groups.GetGroup(rs.Primary.ID, nil)
		if err == nil { // nolint // TODO: Resolve this golangci-lint issue: SA9003: empty branch (staticcheck)
			//if gotRepo != nil && fmt.Sprintf("%d", gotRepo.ID) == rs.Primary.ID {
			//	if gotRepo.MarkedForDeletionAt == nil {
			//		return fmt.Errorf("Repository still exists")
			//	}
			//}
		}
		if !is404(err) {
			return err
		}
		return nil
	}
	return nil
}

func testAccGitlabGroupVariableConfig(rString string) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
name = "foo%v"
path = "foo%v"
}

resource "gitlab_group_variable" "foo" {
  group = "${gitlab_group.foo.id}"
  key = "key_%s"
  value = "value-%s"
  variable_type = "file"
  masked = false
}
	`, rString, rString, rString, rString)
}

func testAccGitlabGroupVariableUpdateConfig(rString string) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
name = "foo%v"
path = "foo%v"
}

resource "gitlab_group_variable" "foo" {
  group = "${gitlab_group.foo.id}"
  key = "key_%s"
  value = "value-inverse-%s"
  protected = true
  masked = false
}
	`, rString, rString, rString, rString)
}

func testAccGitlabGroupVariableScopeConfig(rString, scopeA, scopeB string) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
  name = "foo%v"
  path = "foo%v"
}

resource "gitlab_group_variable" "a" {
  group             = "${gitlab_group.foo.id}"
  key               = "key_%s"
  value             = "value-%s-a"
  environment_scope = "%s"
}

resource "gitlab_group_variable" "b" {
  group             = "${gitlab_group.foo.id}"
  key               = "key_%s"
  value             = "value-%s-b"
  environment_scope = "%s"
}
	`, rString, rString, rString, rString, scopeA, rString, rString, scopeB)
}
