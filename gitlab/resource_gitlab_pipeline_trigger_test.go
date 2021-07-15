package gitlab

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

func TestAccGitlabPipelineTrigger_basic(t *testing.T) {
	var trigger gitlab.PipelineTrigger
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabPipelineTriggerDestroy,
		Steps: []resource.TestStep{
			// Create a project and pipeline trigger with default options
			{
				Config: testAccGitlabPipelineTriggerConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabPipelineTriggerExists("gitlab_pipeline_trigger.trigger", &trigger),
					testAccCheckGitlabPipelineTriggerAttributes(&trigger, &testAccGitlabPipelineTriggerExpectedAttributes{
						Description: "External Pipeline Trigger",
					}),
				),
			},
			// Update the pipeline trigger to change the parameters
			{
				Config: testAccGitlabPipelineTriggerUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabPipelineTriggerExists("gitlab_pipeline_trigger.trigger", &trigger),
					testAccCheckGitlabPipelineTriggerAttributes(&trigger, &testAccGitlabPipelineTriggerExpectedAttributes{
						Description: "Trigger",
					}),
				),
			},
			// Update the pipeline trigger to get back to initial settings
			{
				Config: testAccGitlabPipelineTriggerConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabPipelineTriggerExists("gitlab_pipeline_trigger.trigger", &trigger),
					testAccCheckGitlabPipelineTriggerAttributes(&trigger, &testAccGitlabPipelineTriggerExpectedAttributes{
						Description: "External Pipeline Trigger",
					}),
				),
			},
		},
	})
}

// lintignore: AT002 // TODO: Resolve this tfproviderlint issue
func TestAccGitlabPipelineTrigger_import(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "gitlab_pipeline_trigger.trigger"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabPipelineTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabPipelineTriggerConfig(rInt),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: getPipelineTriggerImportID(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func getPipelineTriggerImportID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("Not Found: %s", n)
		}

		pipelineTriggerID := rs.Primary.ID
		if pipelineTriggerID == "" {
			return "", fmt.Errorf("No pipeline trigger ID is set")
		}
		projectID := rs.Primary.Attributes["project"]
		if projectID == "" {
			return "", fmt.Errorf("No project ID is set")
		}

		return fmt.Sprintf("%s:%s", projectID, pipelineTriggerID), nil
	}
}

func testAccCheckGitlabPipelineTriggerExists(n string, trigger *gitlab.PipelineTrigger) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		triggerID := rs.Primary.ID
		repoName := rs.Primary.Attributes["project"]
		if repoName == "" {
			return fmt.Errorf("No project ID is set")
		}
		conn := testAccProvider.Meta().(*gitlab.Client)

		triggers, _, err := conn.PipelineTriggers.ListPipelineTriggers(repoName, nil)
		if err != nil {
			return err
		}
		for _, gotTrigger := range triggers {
			if strconv.Itoa(gotTrigger.ID) == triggerID {
				*trigger = *gotTrigger
				return nil
			}
		}
		return fmt.Errorf("Pipeline Trigger does not exist")
	}
}

type testAccGitlabPipelineTriggerExpectedAttributes struct {
	Description string
}

func testAccCheckGitlabPipelineTriggerAttributes(trigger *gitlab.PipelineTrigger, want *testAccGitlabPipelineTriggerExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if trigger.Description != want.Description {
			return fmt.Errorf("got description %q; want %q", trigger.Description, want.Description)
		}

		return nil
	}
}

func testAccCheckGitlabPipelineTriggerDestroy(s *terraform.State) error {
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

func testAccGitlabPipelineTriggerConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_pipeline_trigger" "trigger" {
	project = "${gitlab_project.foo.id}"
	description = "External Pipeline Trigger"
}
	`, rInt)
}

func testAccGitlabPipelineTriggerUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_pipeline_trigger" "trigger" {
  project = "${gitlab_project.foo.id}"
  description = "Trigger"
}
	`, rInt)
}
