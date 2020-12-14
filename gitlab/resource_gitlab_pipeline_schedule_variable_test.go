package gitlab

import (
	"fmt"
	"strconv"
	"testing"

	gitlab "github.com/Fourcast/go-gitlab"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccGitlabPipelineScheduleVariable_basic(t *testing.T) {
	var variable gitlab.PipelineVariable
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabProjectDestroy,
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
		},
	})
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

		conn := testAccProvider.Meta().(*gitlab.Client)
		pipelineSchedule, _, err := conn.PipelineSchedules.GetPipelineSchedule(project, scheduleID)
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
