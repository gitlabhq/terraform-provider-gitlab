package gitlab

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

func TestAccGitlabServicePipelinesEmail_basic(t *testing.T) {
	var pipelinesEmailService gitlab.PipelinesEmailService
	rInt := acctest.RandInt()
	pipelinesEmailResourceName := "gitlab_service_pipelines_email.email"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabServicePipelinesEmailDestroy,
		Steps: []resource.TestStep{
			// Create a project and a pipelines email service
			{
				Config: testAccGitlabServicePipelinesEmailConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServicePipelinesEmailExists(pipelinesEmailResourceName, &pipelinesEmailService),
					testRecipients(&pipelinesEmailService, []string{"test@example.com"}),
					resource.TestCheckResourceAttr(pipelinesEmailResourceName, "notify_only_broken_pipelines", "true"),
					resource.TestCheckResourceAttr(pipelinesEmailResourceName, "branches_to_be_notified", "default"),
				),
			},
			// Update the pipelinesEmail service
			{
				Config: testAccGitlabServicePipelinesEmailUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServicePipelinesEmailExists(pipelinesEmailResourceName, &pipelinesEmailService),
					testRecipients(&pipelinesEmailService, []string{"test@example.com", "test2@example.com"}),
					resource.TestCheckResourceAttr(pipelinesEmailResourceName, "notify_only_broken_pipelines", "false"),
					resource.TestCheckResourceAttr(pipelinesEmailResourceName, "branches_to_be_notified", "all"),
				),
			},
			// Update the pipelinesEmail service to get back to previous settings
			{
				Config: testAccGitlabServicePipelinesEmailConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServicePipelinesEmailExists(pipelinesEmailResourceName, &pipelinesEmailService),
					testRecipients(&pipelinesEmailService, []string{"test@example.com"}),
					resource.TestCheckResourceAttr(pipelinesEmailResourceName, "notify_only_broken_pipelines", "true"),
					resource.TestCheckResourceAttr(pipelinesEmailResourceName, "branches_to_be_notified", "default"),
				),
			},
		},
	})
}

// lintignore: AT002 // TODO: Resolve this tfproviderlint issue
func TestAccGitlabServicePipelinesEmail_import(t *testing.T) {
	pipelinesEmailResourceName := "gitlab_service_pipelines_email.email"
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabServicePipelinesEmailDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabServicePipelinesEmailConfig(rInt),
			},
			{
				ResourceName:      pipelinesEmailResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGitlabServicePipelinesEmailExists(n string, service *gitlab.PipelinesEmailService) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		project := rs.Primary.Attributes["project"]
		if project == "" {
			return fmt.Errorf("No project ID is set")
		}
		conn := testAccProvider.Meta().(*gitlab.Client)

		pipelinesEmailService, _, err := conn.Services.GetPipelinesEmailService(project)
		if err != nil {
			return fmt.Errorf("PipelinesEmail service does not exist in project %s: %v", project, err)
		}
		*service = *pipelinesEmailService

		return nil
	}
}

func testRecipients(service *gitlab.PipelinesEmailService, expected []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		res_string := service.Properties.Recipients
		res := strings.Split(res_string, ",")
		if len(res) != len(expected) {
			return fmt.Errorf("'recipients' field does not have the correct size expected: %d, found: %d", len(expected), len(res))
		}
		sort.Strings(res)
		sort.Strings(expected)
		for i, r := range res {
			e := expected[i]
			if r != e {
				return fmt.Errorf("expected: %s, found: %s", r, e)
			}

		}
		return nil
	}
}

func testAccCheckGitlabServicePipelinesEmailDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*gitlab.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_project" {
			continue
		}

		gotRepo, resp, err := conn.Projects.GetProject(rs.Primary.ID, nil)
		if err == nil {
			if gotRepo != nil && fmt.Sprintf("%d", gotRepo.ID) == rs.Primary.ID {
				if gotRepo.MarkedForDeletionAt == nil {
					return fmt.Errorf("Repository still exists")
				}
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func testAccGitlabServicePipelinesEmailConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
    name         = "foo-%d"
    description  = "Terraform acceptance tests"
}

resource "gitlab_service_pipelines_email" "email" {
    project                      = gitlab_project.foo.id
    recipients                   = ["test@example.com"]
}
`, rInt)
}

func testAccGitlabServicePipelinesEmailUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
    name         = "foo-%d"
    description  = "Terraform acceptance tests"
}

resource "gitlab_service_pipelines_email" "email" {
    project                      = gitlab_project.foo.id
    recipients                   = ["test@example.com", "test2@example.com"]
    notify_only_broken_pipelines = false
    branches_to_be_notified      = "all"
}
`, rInt)
}
