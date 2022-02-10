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
    rmodifyInt := acctest.RandInt()

    resource.Test(t, resource.TestCase{
        PreCheck:          func() { testAccPreCheck(t) },
        ProviderFactories: providerFactories,
        CheckDestroy:      testAccCheckGitlabBranchDestroy,
        Steps: []resource.TestStep{
            // Create the branch
            {
                Config: testAccGitlabBranchConfig(rInt),
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckGitlabBranchExists("gitlab_branch.this", &b),
                    testAccCheckGitlabBranchAttributes(&b, &testAccGitlabBranchExpectedAttributes{
                        Name: fmt.Sprintf("example-%[1]d", rInt),
                    }),
                ),
            },
            // Update the branch to toggle all the values to their inverse
            {
                Config: testAccGitlabBranchConfig(rmodifyInt),
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckGitlabBranchExists("gitlab_branch.this", &b),
                    testAccCheckGitlabBranchAttributes(&b, &testAccGitlabBranchExpectedAttributes{
                        Name: fmt.Sprintf("example-%[1]d", rmodifyInt),
                    }),
                ),
            },
            // Update the branch to toggle the options back
            {
                Config: testAccGitlabBranchConfig(rInt),
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckGitlabBranchExists("gitlab_branch.this", &b),
                    testAccCheckGitlabBranchAttributes(&b, &testAccGitlabBranchExpectedAttributes{
                        Name: fmt.Sprintf("example-%[1]d", rInt),
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

        project, branch_name, err := parseTwoPartID(rs.Primary.ID)
        if err != nil {
            return fmt.Errorf("Error in Splitting Project and Branch Ids")
        }

        branch, _, err := testGitlabClient.Branches.GetBranch(project, branch_name)
        if err != nil {
            return err
        }

        if branch.Name == branch_name {
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
    var branch_name string
    for _, rs := range s.RootModule().Resources {
        if rs.Type == "gitlab_project" {
            project = rs.Primary.ID
        } else if rs.Type == "gitlab_branch" {
            branch_name = rs.Primary.ID
        }
    }

    branch, _, err := testGitlabClient.Branches.GetBranch(project, branch_name)
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

func testAccGitlabBranchConfig(rInt int) string {
    return fmt.Sprintf(`
resource "gitlab_group" "this" {
  name        = "example-%[1]d"
  path        = "example"
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
  name    = "example-%[1]d"
  ref     = gitlab_project.this.default_branch
}
resource "gitlab_branch_protection" "this" {
  project            = gitlab_project.this.id
  branch             = gitlab_branch.this.name
  push_access_level  = "maintainer"
  merge_access_level = "maintainer"
}
    `, rInt)
}
