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

func TestAccGitlabPipelineScheduleVariable_basic(t *testing.T) {
	var variable gitlab.PipelineVariable
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabPipelineScheduleVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabPipelineScheduleVariableConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabPipelineScheduleVariableExists("gitlab_pipeline_schedule_variable.schedule_var", &variable),
					testAccCheckGitlabPipelineScheduleVariableAttributes(&variable, &testAccGitlabPipelineScheduleVariableExpectedAttributes{
						Key:   "TERRAFORMED_TEST_VALUE",
						Value: "test",
					}),
				),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_pipeline_schedule_variable.schedule_var",
				ImportState:       true,
				ImportStateIdFunc: getPipelineScheduleVariableID("gitlab_pipeline_schedule_variable.schedule_var"),
				ImportStateVerify: true,
			},
			{
				Config: testAccGitlabPipelineScheduleVariableUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabPipelineScheduleVariableExists("gitlab_pipeline_schedule_variable.schedule_var", &variable),
					testAccCheckGitlabPipelineScheduleVariableAttributes(&variable, &testAccGitlabPipelineScheduleVariableExpectedAttributes{
						Key:   "TERRAFORMED_TEST_VALUE",
						Value: "test_updated",
					}),
				),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_pipeline_schedule_variable.schedule_var",
				ImportState:       true,
				ImportStateIdFunc: getPipelineScheduleVariableID("gitlab_pipeline_schedule_variable.schedule_var"),
				ImportStateVerify: true,
			},
			{
				Config: testAccGitlabPipelineScheduleVariableConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabPipelineScheduleVariableExists("gitlab_pipeline_schedule_variable.schedule_var", &variable),
					testAccCheckGitlabPipelineScheduleVariableAttributes(&variable, &testAccGitlabPipelineScheduleVariableExpectedAttributes{
						Key:   "TERRAFORMED_TEST_VALUE",
						Value: "test",
					}),
				),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_pipeline_schedule_variable.schedule_var",
				ImportState:       true,
				ImportStateIdFunc: getPipelineScheduleVariableID("gitlab_pipeline_schedule_variable.schedule_var"),
				ImportStateVerify: true,
			},
		},
	})
}

func getPipelineScheduleVariableID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("not found: %s", n)
		}

		pipelineScheduleVariableID := rs.Primary.ID
		if pipelineScheduleVariableID == "" {
			return "", fmt.Errorf("no pipeline schedule variable ID is set")
		}
		projectID := rs.Primary.Attributes["project"]
		if projectID == "" {
			return "", fmt.Errorf("no project ID is set")
		}

		return fmt.Sprintf("%s:%s", projectID, pipelineScheduleVariableID), nil
	}
}

func testAccCheckGitlabPipelineScheduleVariableExists(n string, variable *gitlab.PipelineVariable) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		project := rs.Primary.Attributes["project"]
		scheduleID, err := strconv.Atoi(rs.Primary.Attributes["pipeline_schedule_id"])
		if err != nil {
			return fmt.Errorf("failed to convert PipelineSchedule.ID to int")
		}

		pipelineSchedule, _, err := testGitlabClient.PipelineSchedules.GetPipelineSchedule(project, scheduleID)
		if err != nil {
			return err
		}

		for _, pipelineVariable := range pipelineSchedule.Variables {
			if pipelineVariable.Key == rs.Primary.Attributes["key"] {
				*variable = *pipelineVariable
				return nil
			}
		}
		return fmt.Errorf("PipelineScheduleVariable %s does not exist", variable.Key)
	}
}

type testAccGitlabPipelineScheduleVariableExpectedAttributes struct {
	Key   string
	Value string
}

func testAccCheckGitlabPipelineScheduleVariableAttributes(variable *gitlab.PipelineVariable, want *testAccGitlabPipelineScheduleVariableExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if variable.Key != want.Key {
			return fmt.Errorf("got key %q; want %q", variable.Key, want.Key)
		}

		return nil
	}
}

func testAccGitlabPipelineScheduleVariableConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_pipeline_schedule" "schedule" {
	project = "${gitlab_project.foo.id}"
	description = "Pipeline Schedule"
	ref = "master"
	cron = "0 1 * * *"
}

resource "gitlab_pipeline_schedule_variable" "schedule_var" {
	project = "${gitlab_project.foo.id}"
	pipeline_schedule_id = "${gitlab_pipeline_schedule.schedule.id}"
	key = "TERRAFORMED_TEST_VALUE"
	value = "test"
}
	`, rInt)
}

func testAccGitlabPipelineScheduleVariableUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_pipeline_schedule" "schedule" {
	project = "${gitlab_project.foo.id}"
	description = "Pipeline Schedule"
	ref = "master"
	cron = "0 1 * * *"
}

resource "gitlab_pipeline_schedule_variable" "schedule_var" {
	project = "${gitlab_project.foo.id}"
	pipeline_schedule_id = "${gitlab_pipeline_schedule.schedule.id}"
	key = "TERRAFORMED_TEST_VALUE"
	value = "test_updated"
}
	`, rInt)
}

func testAccCheckGitlabPipelineScheduleVariableDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_pipeline_schedule_variable" {
			continue
		}

		psidString := rs.Primary.Attributes["pipeline_schedule_id"]
		psid, err := strconv.Atoi(psidString)
		if err != nil {
			return fmt.Errorf("could not convert pipeline schedule id to integer: %s", err)
		}

		gotPS, _, err := testGitlabClient.PipelineSchedules.GetPipelineSchedule(rs.Primary.Attributes["project"], psid)
		if err == nil {
			for _, v := range gotPS.Variables {
				if buildTwoPartID(&psidString, &v.Key) == rs.Primary.ID {
					return fmt.Errorf("pipeline schedule variable still exists")
				}
			}
		}
	}
	return nil
}
