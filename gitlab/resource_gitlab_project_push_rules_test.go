package gitlab

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	gitlab "github.com/xanzy/go-gitlab"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccGitlabProjectPushRules_basic(t *testing.T) {
	var pushRules gitlab.ProjectPushRules
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabProjectPushRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabProjectPushRulesConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectPushRulesExists("gitlab_project_push_rule.foo", &pushRules),
				),
			},
		},
	})
}

func testAccCheckGitlabProjectPushRulesExists(n string, pushRules *gitlab.ProjectPushRules) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		repoName := rs.Primary.Attributes["project"]
		if repoName == "" {
			return fmt.Errorf("No project ID is set")
		}
		conn := testAccProvider.Meta().(*gitlab.Client)
		gotPushRules, _, err := conn.Projects.GetProjectPushRules(repoName)
		if err != nil {
			return err
		}
		*pushRules = *gotPushRules
		return nil
	}
}

func testAccCheckGitlabProjectPushRulesDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*gitlab.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_project" {
			continue
		}

		gotRepo, resp, err := conn.Projects.GetProject(rs.Primary.ID, nil)
		if err == nil {
			if gotRepo != nil && fmt.Sprintf("%d", gotRepo.ID) == rs.Primary.ID {
				return fmt.Errorf("Repository still exists")
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func testAccGitlabProjectPushRulesConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance test - Push Rule"
  visibility_level = "public"
}

resource "gitlab_project_push_rules" "foo" {
  project = "${gitlab_project.foo.id}"
  commit_message_regex = "^(foo|bar).*"
}
`, rInt)
}
