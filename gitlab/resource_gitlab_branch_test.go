package gitlab

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	gitlab "github.com/xanzy/go-gitlab"
	"testing"
)

func TestAccGitlabBranch_basic(t *testing.T) {
	var branch gitlab.Branch
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabBranchDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabBranchConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabBranchExists("gitlab_branch.foo", "test", &branch, rInt),
					testAccCheckGitlabBranchAttributes("gitlab_branch.foo", "test", &branch, &testAccGitlabBranchExpectedAttributes{
						Name: fmt.Sprintf("testbranch-%d", rInt),
					}),
				),
			},
		},
	})
}

func testAccCheckGitlabBranchDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_branch" {
			continue
		}
		name := rs.Primary.Attributes["name"]
		project := rs.Primary.Attributes["project"]
		conn := testAccProvider.Meta().(*gitlab.Client)
		branch, resp, err := conn.Branches.GetBranch(project, name)
		if err == nil {
			if branch != nil && fmt.Sprintf("%s", branch.Name) == name {
				return fmt.Errorf("Branch still exists")
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func testAccCheckGitlabBranchAttributes(n, p string, branch *gitlab.Branch, want *testAccGitlabBranchExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if branch.WebURL == "" {
			return errors.New("got empty web url")
		}
		projectID := s.RootModule().Resources[fmt.Sprintf("gitlab_project.%s", p)].Primary.ID
		if s.RootModule().Resources[n].Primary.ID != fmt.Sprintf("%s-%s", projectID, want.Name) {
			return fmt.Errorf("Got ID: %s expected: %s-%s", s.RootModule().Resources[n].Primary.ID, projectID, want.Name)
		}
		if branch.Commit.ID == "" {
			return errors.New("Empty commit")
		}
		if branch.Name != want.Name {
			return fmt.Errorf("got name %s; want %s", branch.Name, want.Name)
		}
		if !branch.CanPush {
			return fmt.Errorf("can push %t; want %t", branch.CanPush, want.CanPush)
		}
		if branch.DevelopersCanPush {
			return errors.New("Developers can push expected output to be false")
		}
		if branch.DevelopersCanMerge {
			return errors.New("Developers can merge expected output to be false")
		}
		if branch.Default {
			return errors.New("Default branch set to true")
		}
		if branch.Merged {
			return errors.New("Merged set to true")
		}
		return nil
	}
}

func testAccCheckGitlabBranchExists(n, p string, branch *gitlab.Branch, rInt int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}
		name := rs.Primary.Attributes["name"]
		pid := s.RootModule().Resources[fmt.Sprintf("gitlab_project.%s", p)].Primary.ID
		conn := testAccProvider.Meta().(*gitlab.Client)
		gotBranch, _, err := conn.Branches.GetBranch(pid, name)
		*branch = *gotBranch
		return err
	}
}

func testAccGitlabBranchConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "gitlab_project" "test" {
		name = "foo-%d"
		description = "Terraform acceptance tests"
	  
		# So that acceptance tests can be run in a gitlab organization
		# with no billing
		visibility_level = "public"
	}
	resource "gitlab_branch" "foo" {
		name = "testbranch-%d"
		ref = "master"
		project = gitlab_project.test.id
	}
  `, rInt, rInt)
}

type testAccGitlabBranchExpectedAttributes struct {
	Name               string
	WebURL             string
	CanPush            bool
	Default            bool
	Merged             bool
	Protected          bool
	Ref                string
	Project            string
	DevelopersCanPush  bool
	DevelopersCanMerge bool
	Commit             *gitlab.Commit
}
