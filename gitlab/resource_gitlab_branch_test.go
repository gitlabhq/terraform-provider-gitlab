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
					testAccCheckGitlabBranchExists("gitlab_branch.foo", &branch, rInt),
					testAccCheckGitlabBranchAttributes(&branch, &testAccGitlabBranchExpectedAttributes{
						Name:               fmt.Sprintf("testbranch-%d", rInt),
						CanPush:            true,
						DevelopersCanMerge: false,
						DevelopersCanPush:  false,
						Default:            false,
						Merged:             false,
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
		pid := s.RootModule().Resources["gitlab_project.test"].Primary.ID
		conn := testAccProvider.Meta().(*gitlab.Client)
		branch, resp, err := conn.Branches.GetBranch(pid, name)
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

func testAccCheckGitlabBranchAttributes(branch *gitlab.Branch, want *testAccGitlabBranchExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if branch.WebURL == "" {
			return errors.New("got empty web url")
		}
		if branch.Name != want.Name {
			return fmt.Errorf("got name %s; want %s", branch.Name, want.Name)
		}
		if branch.CanPush != want.CanPush {
			return fmt.Errorf("can push %t; want %t", branch.CanPush, want.CanPush)
		}
		if branch.DevelopersCanPush != want.DevelopersCanPush {
			return fmt.Errorf("Developers can push %t; want %t", branch.DevelopersCanPush, want.DevelopersCanPush)
		}
		if branch.DevelopersCanMerge != want.DevelopersCanMerge {
			return fmt.Errorf("Developers can merge %t; want %t", branch.DevelopersCanMerge, want.DevelopersCanMerge)
		}
		if branch.Default != want.Default {
			return fmt.Errorf("Default set %t; want %t", branch.CanPush, want.CanPush)
		}
		if branch.Merged != want.Merged {
			return fmt.Errorf("Merged %t; want %t", branch.CanPush, want.CanPush)
		}
		return nil
	}
}

func testAccCheckGitlabBranchExists(n string, branch *gitlab.Branch, rInt int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}
		name := rs.Primary.Attributes["name"]
		pid := s.RootModule().Resources["gitlab_project.test"].Primary.ID
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
