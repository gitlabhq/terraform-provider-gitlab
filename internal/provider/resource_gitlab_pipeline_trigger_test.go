//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

func TestAccGitlabPipelineTrigger_basic(t *testing.T) {
	var trigger gitlab.PipelineTrigger
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabPipelineTriggerDestroy,
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
			// Verify Import
			{
				ResourceName:      "gitlab_pipeline_trigger.trigger",
				ImportStateIdFunc: getPipelineTriggerImportID("gitlab_pipeline_trigger.trigger"),
				ImportState:       true,
				ImportStateVerify: true,
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
			// Verify Import
			{
				ResourceName:      "gitlab_pipeline_trigger.trigger",
				ImportStateIdFunc: getPipelineTriggerImportID("gitlab_pipeline_trigger.trigger"),
				ImportState:       true,
				ImportStateVerify: true,
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
			// Verify Import
			{
				ResourceName:      "gitlab_pipeline_trigger.trigger",
				ImportStateIdFunc: getPipelineTriggerImportID("gitlab_pipeline_trigger.trigger"),
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

		triggers, _, err := testGitlabClient.PipelineTriggers.ListPipelineTriggers(repoName, nil)
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
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_pipeline_trigger" {
			continue
		}

		project := rs.Primary.Attributes["project"]
		pipelineTriggerID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, _, err = testGitlabClient.PipelineTriggers.GetPipelineTrigger(project, pipelineTriggerID)
		if err == nil {
			return fmt.Errorf("Pipeline Trigger %d in project %s still exists", pipelineTriggerID, project)
		}
		if !is404(err) {
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
