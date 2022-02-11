package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabBranch_basic(t *testing.T) {

	var b gitlab.Branch
	rInt := acctest.RandInt()
	rSupportInt := acctest.RandInt()
	rModifyInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabBranchDestroy,
		Steps: []resource.TestStep{
			// Create the branch
			{
				Config: testAccGitlabBranchConfig(rInt, rSupportInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabBranchExists("gitlab_branch.this", &b),
					testAccCheckGitlabBranchAttributes(&b, &testAccGitlabBranchExpectedAttributes{
						Name: fmt.Sprintf("example-%[1]d", rInt),
					}),
				),
			},
			// Test ImportState
			{
				ResourceName:      "gitlab_branch.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update the branch to toggle all the values to their inverse
			{
				Config: testAccGitlabBranchConfig(rModifyInt, rSupportInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabBranchExists("gitlab_branch.this", &b),
					testAccCheckGitlabBranchAttributes(&b, &testAccGitlabBranchExpectedAttributes{
						Name: fmt.Sprintf("example-%[1]d", rModifyInt),
					}),
				),
			},
		},
	})
}

func testAccCheckGitlabBranchExists(n string, b *gitlab.Branch) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		project, branchName, err := parseTwoPartID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error in splitting project and branch IDs")
		}

		branch, _, err := testGitlabClient.Branches.GetBranch(project, branchName)
		if err != nil {
			return err
		}

		if branch.Name == branchName {
			*b = *branch
			return nil
		}
		return fmt.Errorf("Branch does not exist")
	}
}

type testAccGitlabBranchExpectedAttributes struct {
	Name string
}

func testAccCheckGitlabBranchAttributes(b *gitlab.Branch, want *testAccGitlabBranchExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if b.Name != want.Name {
			return fmt.Errorf("got name %q; want %q", b.Name, want.Name)
		}
		return nil
	}
}

func testAccCheckGitlabBranchDestroy(s *terraform.State) error {
	var project string
	var branchName string
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "gitlab_project" {
			project = rs.Primary.ID
		} else if rs.Type == "gitlab_branch" {
			branchName = rs.Primary.ID
		}
	}

	branch, _, err := testGitlabClient.Branches.GetBranch(project, branchName)
	if err == nil {
		if branch != nil {
			return fmt.Errorf("project branch protection %s still exists", branch)
		}
	}
	if !is404(err) {
		return err
	}
	return nil
}

func testAccGitlabBranchConfig(rInt int, rSupportInt int) string {
	return fmt.Sprintf(`
resource "gitlab_group" "this" {
  name        = "example-%[1]d"
  path        = "example-%[1]d"
  description = "An example group"
}
resource "gitlab_project" "this" {
  name                   = "example-%[1]d"
  namespace_id           = gitlab_group.this.id
  default_branch         = "main"
  initialize_with_readme = true
}
resource "gitlab_repository_file" "this" {
  project        = gitlab_project.this.id
  file_path      = "meow.txt"
  branch         = gitlab_project.this.default_branch
  content        = base64encode("Meow goes the cat")
  author_email   = "terraform@example.com"
  author_name    = "Terraform"
  commit_message = "feature: add meow file"
}
resource "gitlab_branch" "this" {
  project = gitlab_project.this.id
  branch  = "example-%[2]d"
  ref     = gitlab_project.this.default_branch
}
resource "gitlab_branch_protection" "this" {
  project            = gitlab_project.this.id
  branch             = gitlab_branch.this.branch
  push_access_level  = "maintainer"
  merge_access_level = "maintainer"
}
    `, rSupportInt, rInt)
}
