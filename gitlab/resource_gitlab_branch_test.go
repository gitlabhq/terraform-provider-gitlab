package gitlab

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

func TestAccGitlabBranch_basic(t *testing.T) {
	var branch gitlab.Branch
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabBranchDestroy,
		Steps: []resource.TestStep{
			// Create a group
			{
				Config: testAccGitlabBranchConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabBranchExists("gitlab_branch.foo", &branch, rInt),
				),
			},
		},
	})
}

func testAccCheckGitlabBranchDestroy(s *terraform.State) error {
	// conn := testAccProvider.Meta().(*gitlab.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_branch" {
			continue
		}
		// branch, resp, err := conn.Branches.GetBranch(rs.Primary.ID ,rs.Primary.Name)
		// // if err == nil {
		// // 	if group != nil && fmt.Sprintf("%d", group.ID) == rs.Primary.ID {
		// // 		if group.MarkedForDeletionOn == nil {
		// // 			return fmt.Errorf("Group still exists")
		// // 		}
		// // 	}
		// // }
		// if resp.StatusCode != 404 {
		// 	return err
		// }
		// return nil
	}
	log.Println("destroy method")
	return nil
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
		branch = gotBranch
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
