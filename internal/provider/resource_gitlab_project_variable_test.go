//go:build acceptance
// +build acceptance

package provider

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
)

func testAccCheckGitlabProjectVariableExists(name string) resource.TestCheckFunc {
	var (
		key              string
		value            string
		variableType     string
		protected        string
		masked           string
		environmentScope string
	)

	return resource.ComposeTestCheckFunc(
		// Load the real resource values using the GitLab API.
		func(state *terraform.State) error {
			attributes := state.RootModule().Resources[name].Primary.Attributes

			got, _, err := testGitlabClient.ProjectVariables.GetVariable(attributes["project"], attributes["key"], nil, gitlab.WithContext(context.Background()), withEnvironmentScopeFilter(context.Background(), attributes["environment_scope"]))
			if err != nil {
				return err
			}

			key = got.Key
			value = got.Value
			variableType = string(got.VariableType)
			protected = strconv.FormatBool(got.Protected)
			masked = strconv.FormatBool(got.Masked)
			environmentScope = got.EnvironmentScope

			return nil
		},

		// Check that the real values match what was configured in the resource.
		resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttrPtr(name, "key", &key),
			resource.TestCheckResourceAttrPtr(name, "value", &value),
			resource.TestCheckResourceAttrPtr(name, "variable_type", &variableType),
			resource.TestCheckResourceAttrPtr(name, "masked", &masked),
			resource.TestCheckResourceAttrPtr(name, "protected", &protected),
			resource.TestCheckResourceAttrPtr(name, "environment_scope", &environmentScope),
		),
	)
}

func testAccGitlabProjectVariableCheckAllVariablesDestroyed(ctx testAccGitlabProjectContext) func(state *terraform.State) error {
	return func(state *terraform.State) error {
		vars, _, err := testGitlabClient.ProjectVariables.ListVariables(ctx.project.ID, nil)
		if err != nil {
			return err
		}

		if len(vars) > 0 {
			return fmt.Errorf("expected no project variables but found %d variables %v", len(vars), vars)
		}

		return nil
	}
}

func TestAccGitlabProjectVariable_basic(t *testing.T) {
	ctx := testAccGitlabProjectStart(t)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccGitlabProjectVariableCheckAllVariablesDestroyed(ctx),
		Steps: []resource.TestStep{
			// Create a project variable from a project name.
			{
				Config: fmt.Sprintf(`
resource "gitlab_project_variable" "foo" {
  project = "%s"
  key = "my_key"
  value = "my_value"
}
`, ctx.project.PathWithNamespace),
				Check: testAccCheckGitlabProjectVariableExists("gitlab_project_variable.foo"),
			},
			{
				ResourceName:      "gitlab_project_variable.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Same, using the project id.
			{
				Config: fmt.Sprintf(`
resource "gitlab_project_variable" "foo" {
  project = %d
  key = "my_key"
  value = "my_value"
}
`, ctx.project.ID),
				Check: testAccCheckGitlabProjectVariableExists("gitlab_project_variable.foo"),
			},
			// Check that the variable is recreated if deleted out-of-band.
			{
				PreConfig: func() {
					if _, err := testGitlabClient.ProjectVariables.RemoveVariable(ctx.project.ID, "my_key", nil); err != nil {
						t.Fatalf("failed to remove variable: %v", err)
					}
				},
				Config: fmt.Sprintf(`
resource "gitlab_project_variable" "foo" {
  project = %d
  key = "my_key"
  value = "my_value"
}
`, ctx.project.ID),
				Check: testAccCheckGitlabProjectVariableExists("gitlab_project_variable.foo"),
			},
			// Update the variable_type.
			{
				Config: fmt.Sprintf(`
resource "gitlab_project_variable" "foo" {
  project = %d
  key = "my_key"
  value = "my_value"
  variable_type = "file"
}
`, ctx.project.ID),
				Check: testAccCheckGitlabProjectVariableExists("gitlab_project_variable.foo"),
			},
			// Update all other attributes.
			{
				Config: fmt.Sprintf(`
resource "gitlab_project_variable" "foo" {
  project = %d
  key = "my_key"
  value = "my_value_2"
  protected = true
  masked = true
}
`, ctx.project.ID),
				Check: testAccCheckGitlabProjectVariableExists("gitlab_project_variable.foo"),
			},
			{
				ResourceName:      "gitlab_project_variable.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Try to update with an illegal masked variable.
			// ref: https://docs.gitlab.com/ce/ci/variables/README.html#masked-variable-requirements
			{
				Config: fmt.Sprintf(`
resource "gitlab_project_variable" "foo" {
  project = %d
  key = "my_key"
  value = <<EOF
i am multiline
EOF
  masked = true
}
`, ctx.project.ID),
				ExpectError: regexp.MustCompile(regexp.QuoteMeta(
					"Invalid value for a masked variable. Check the masked variable requirements: https://docs.gitlab.com/ee/ci/variables/#masked-variable-requirements",
				)),
			},
		},
	})
}

func TestAccGitlabProjectVariable_scoped(t *testing.T) {
	ctx := testAccGitlabProjectStart(t)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy: func(state *terraform.State) error {
			// Destroy behavior is nondeterministic for variables with scopes in GitLab versions prior to 13.4
			// ref: https://gitlab.com/gitlab-org/gitlab/-/merge_requests/39209
			if isAtLeast134, err := isGitLabVersionAtLeast(context.Background(), testGitlabClient, "13.4")(); err != nil {
				return err
			} else if isAtLeast134 {
				return testAccGitlabProjectVariableCheckAllVariablesDestroyed(ctx)(state)
			}
			return nil
		},
		Steps: []resource.TestStep{
			// Create a project variable from a project id
			{
				Config: fmt.Sprintf(`
resource "gitlab_project_variable" "foo" {
  project = %d
  key = "my_key"
  value = "my_value"
}
`, ctx.project.ID),
				Check: testAccCheckGitlabProjectVariableExists("gitlab_project_variable.foo"),
			},
			// Update the scope.
			{
				Config: fmt.Sprintf(`
resource "gitlab_project_variable" "foo" {
  project = %d
  key = "my_key"
  value = "my_value"
  environment_scope = "foo"
}
`, ctx.project.ID),
				Check: testAccCheckGitlabProjectVariableExists("gitlab_project_variable.foo"),
			},
			// Add a second variable with the same key and different scope.
			{
				Config: fmt.Sprintf(`
resource "gitlab_project_variable" "foo" {
  project = %[1]d
  key = "my_key"
  value = "my_value"
  environment_scope = "foo"
}

resource "gitlab_project_variable" "bar" {
  project = %[1]d
  key = "my_key"
  value = "my_value_2"
  environment_scope = "bar"
}
`, ctx.project.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckGitlabProjectVariableExists("gitlab_project_variable.foo"),
					testAccCheckGitlabProjectVariableExists("gitlab_project_variable.bar"),
				),
			},
			{
				ResourceName:      "gitlab_project_variable.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "gitlab_project_variable.bar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update an attribute on one of the variables.
			// Updating a variable with a non-unique key only works reliably on GitLab v13.4+.
			{
				SkipFunc: isGitLabVersionLessThan(context.Background(), testGitlabClient, "13.4"),
				Config: fmt.Sprintf(`
resource "gitlab_project_variable" "foo" {
  project = %[1]d
  key = "my_key"
  value = "my_value"
  environment_scope = "foo"
}

resource "gitlab_project_variable" "bar" {
  project = %[1]d
  key = "my_key"
  value = "my_value_2_updated"
  environment_scope = "bar"
}
`, ctx.project.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckGitlabProjectVariableExists("gitlab_project_variable.foo"),
					testAccCheckGitlabProjectVariableExists("gitlab_project_variable.bar"),
				),
			},
			// Try to have two variables with the same keys and scopes.
			// On versions of GitLab < 13.4 this can sometimes result in an inconsistent state instead of an error.
			{
				SkipFunc: isGitLabVersionLessThan(context.Background(), testGitlabClient, "13.4"),
				Config: fmt.Sprintf(`
resource "gitlab_project_variable" "foo" {
  project = %[1]d
  key = "my_key"
  value = "my_value"
  environment_scope = "foo"
}

resource "gitlab_project_variable" "bar" {
  project = %[1]d
  key = "my_key"
  value = "my_value_2"
  environment_scope = "foo"
}
`, ctx.project.ID),
				ExpectError: regexp.MustCompile(regexp.QuoteMeta("(my_key) has already been taken")),
			},
		},
	})
}
