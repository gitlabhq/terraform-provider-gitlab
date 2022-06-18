//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataGitlabProject_basic(t *testing.T) {
	projectname := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataGitlabProjectConfigByPathWithNamespace(projectname),
				Check: testAccDataSourceGitlabProject("gitlab_project.test", "data.gitlab_project.foo",
					[]string{"id", "name", "path", "visibility", "description"}),
			},
			{
				Config: testAccDataGitlabProjectConfig(projectname),
				Check: testAccDataSourceGitlabProject("gitlab_project.test", "data.gitlab_project.foo",
					[]string{"id", "name", "path", "visibility", "description"}),
			},
			{
				SkipFunc: isRunningInCE,
				Config:   testAccDataGitlabProjectConfigPushRules(projectname),
				Check: testAccDataSourceGitlabProject("gitlab_project.test", "data.gitlab_project.foo",
					[]string{"push_rules.0.author_email_regex"}),
			},
		},
	})
}

func testAccDataSourceGitlabProject(resourceName, dataSourceName string, testAttributes []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		project := s.RootModule().Resources[resourceName]
		projectResource := project.Primary.Attributes

		search := s.RootModule().Resources[dataSourceName]
		searchResource := search.Primary.Attributes

		if searchResource["id"] == "" {
			return fmt.Errorf("Expected to get a project ID from Gitlab")
		}

		for _, attribute := range testAttributes {
			if searchResource[attribute] != projectResource[attribute] {
				return fmt.Errorf("Expected the project %s to be: %s, but got: %s", attribute, projectResource[attribute], searchResource[attribute])
			}
		}
		return nil
	}
}

func testAccDataGitlabProjectConfig(projectname string) string {
	return fmt.Sprintf(`
resource "gitlab_project" "test"{
	name = "%s"
	path = "%s"
	description = "Terraform acceptance tests"
	visibility_level = "public"
}

data "gitlab_project" "foo" {
	id = "${gitlab_project.test.id}"
}
	`, projectname, projectname)
}

func testAccDataGitlabProjectConfigByPathWithNamespace(projectname string) string {
	return fmt.Sprintf(`
resource "gitlab_project" "test"{
	name = "%s"
	path = "%s"
	description = "Terraform acceptance tests"
	visibility_level = "public"
}

data "gitlab_project" "foo" {
	path_with_namespace = gitlab_project.test.path_with_namespace
}
	`, projectname, projectname)
}

func testAccDataGitlabProjectConfigPushRules(projectName string) string {
	return fmt.Sprintf(`
resource "gitlab_project" "test"{
	name = "%[1]s"
	path = "%[1]s"
	description = "Terraform acceptance tests"
	visibility_level = "public"
    push_rules {
        author_email_regex = "foo"
    }
}

data "gitlab_project" "foo" {
	id = gitlab_project.test.id
}
	`, projectName)
}
