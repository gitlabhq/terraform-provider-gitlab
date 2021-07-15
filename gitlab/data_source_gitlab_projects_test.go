package gitlab

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// lintignore: AT003 // TODO: Resolve this tfproviderlint issue
func TestAccDataGitlabProjectsSearch(t *testing.T) {
	projectName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataGitlabProjectsConfigGetProjectSearch(projectName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDataSourceGitlabProjects(
						"gitlab_project.search",
						"data.gitlab_projects.search",
					),
					resource.TestCheckResourceAttr(
						"data.gitlab_projects.search",
						"projects.0.owner.0.id",
						"1",
					),
					resource.TestCheckResourceAttr(
						"data.gitlab_projects.search",
						"projects.0.permissions.0.project_access.access_level",
						"40",
					),
					resource.TestCheckNoResourceAttr(
						"data.gitlab_projects.search",
						"projects.0.permissions.0.project_access.group_level",
					),
					resource.TestCheckResourceAttr(
						"data.gitlab_projects.search",
						"projects.0.namespace.0.kind",
						"user",
					),
				),
			},
		},
	})
}

// lintignore: AT003 // TODO: Resolve this tfproviderlint issue
func TestAccDataGitlabProjectsGroups(t *testing.T) {
	projectName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	groupName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	parentGroupName := fmt.Sprintf("tf-%s", acctest.RandString(5))
	subGroupName1 := fmt.Sprintf("tf-%s", acctest.RandString(5))
	subGroupName2 := fmt.Sprintf("tf-%s", acctest.RandString(5))
	subGroupProjectName1 := fmt.Sprintf("tf-%s", acctest.RandString(5))
	subGroupProjectName2 := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataGitlabProjectsConfigGetGroupProjectsByGroupId(groupName, projectName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDataSourceGitlabProjects(
						"gitlab_project.testProject",
						"data.gitlab_projects.group",
					),
					resource.TestCheckResourceAttr(
						"data.gitlab_projects.group",
						"projects.0.namespace.0.kind",
						"group",
					),
				),
			},
			{
				Config: testAccDataGitlabProjectsConfigGetNestedProjectsByParentGroupId(parentGroupName, subGroupName1, subGroupName2, subGroupProjectName1, subGroupProjectName2),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGitlabProjects(
						"gitlab_project.testProject1",
						"data.gitlab_projects.subGroups",
					),
					testAccDataSourceGitlabProjects(
						"gitlab_project.testProject2",
						"data.gitlab_projects.subGroups",
					),
				),
			},
		},
	})
}

func testAccDataSourceGitlabProjects(src string, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		project := s.RootModule().Resources[src]
		projectResource := project.Primary.Attributes

		search := s.RootModule().Resources[n]
		searchResource := search.Primary.Attributes

		if searchResource["id"] == "" {
			return fmt.Errorf("expected to get a project ID from Gitlab")
		}
		if searchResource["projects.#"] == "0" {
			return fmt.Errorf("expected to find at least one matching project from the datasource")
		}

		projectsNumber, err := strconv.Atoi(searchResource["projects.#"])
		if err != nil {
			return fmt.Errorf("the datasource returned no 'projects' attribute, got: %s", searchResource)
		}

		testAttributes := []string{
			"id",
			"name",
			"path",
		}

		for i := 0; i < projectsNumber; i++ {
			for j, attribute := range testAttributes {
				if searchResource[fmt.Sprintf("projects.%d.%s", i, attribute)] != projectResource[attribute] {
					break
				}
				if j == len(testAttributes)-1 {
					// Found an exact match
					return nil
				}
			}
		}

		var errorMessageExpected strings.Builder
		for _, attr := range testAttributes {
			errorMessageExpected.WriteString(fmt.Sprintf("%s=%v, ", attr, projectResource[fmt.Sprintf("%s", attr)])) // nolint // TODO: Resolve this golangci-lint issue: S1025: the argument is already a string, there's no need to use fmt.Sprintf (gosimple)
		}

		var errorMessageGot strings.Builder
		for i := 0; i < projectsNumber; i++ {
			errorMessageGot.WriteString(fmt.Sprintf("project_%d: ", i))
			for _, attr := range testAttributes {
				errorMessageGot.WriteString(fmt.Sprintf("%s_%d=%v, ", attr, i, searchResource[fmt.Sprintf("projects.%d.%s", i, attr)]))
			}
			errorMessageGot.WriteString("\n")
		}

		return fmt.Errorf("datasource did not return any match.\nExpected: %s\nGot:\n  %s", errorMessageExpected.String(), errorMessageGot.String())
	}
}

func testAccDataGitlabProjectsConfigGetProjectSearch(projectName string) string {
	return fmt.Sprintf(`

resource "gitlab_project" "search" {
  name = "%s"
  path = "%s"
}

data "gitlab_projects" "search" {
  search = gitlab_project.search.name
}
	`, projectName, projectName)
}

func testAccDataGitlabProjectsConfigGetGroupProjectsByGroupId(groupName string, projectName string) string {
	return fmt.Sprintf(`
resource "gitlab_group" "testGroup" {
  name = "%s"
  path = "%s"
  description = "Terraform acceptance tests"
}

resource "gitlab_project" "testProject"{
  name = "%s"
  namespace_id = gitlab_group.testGroup.id
}

data "gitlab_projects" "group" {
  group_id = gitlab_project.testProject.namespace_id
}
	`, groupName, groupName, projectName)
}

func testAccDataGitlabProjectsConfigGetNestedProjectsByParentGroupId(parentGroupName string, subGroupName1 string, subGroupName2 string, projectName1 string, projectName2 string) string {
	return fmt.Sprintf(`
resource "gitlab_group" "testGroup" {
  name = "%s"
  path = "%s"
}

resource "gitlab_group" "testSubGroup1" {
  name = "%s"
  path = "%s"
  parent_id = gitlab_group.testGroup.id
}

resource "gitlab_group" "testSubGroup2" {
  name = "%s"
  path = "%s"
  parent_id = gitlab_group.testGroup.id
}

resource "gitlab_project" "testProject1"{
  name = "%s"
  namespace_id = gitlab_group.testSubGroup1.id
  description = gitlab_group.testGroup.id
}

resource "gitlab_project" "testProject2"{
  name = "%s"
  namespace_id = gitlab_group.testSubGroup2.id
  // This is all just to avoid using explicit depends_on on the datasource
  // since it seems to break the acceptance tests
  description = gitlab_project.testProject1.description
}

data "gitlab_projects" "subGroups" {
  // This is to ensure the projects have been created before running the datasource
  group_id = gitlab_project.testProject2.description
  include_subgroups = true
}
	`, parentGroupName, parentGroupName, subGroupName1, subGroupName1, subGroupName2, subGroupName2, projectName1, projectName2)
}
