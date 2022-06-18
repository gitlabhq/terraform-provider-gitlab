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

func TestAccGitlabPipelineSchedule_basic(t *testing.T) {
	var schedule gitlab.PipelineSchedule
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabPipelineScheduleDestroy,
		Steps: []resource.TestStep{
			// Create a project and pipeline schedule with default options
			{
				Config: testAccGitlabPipelineScheduleConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabPipelineScheduleExists("gitlab_pipeline_schedule.schedule", &schedule),
					testAccCheckGitlabPipelineScheduleAttributes(&schedule, &testAccGitlabPipelineScheduleExpectedAttributes{
						Description:  "Pipeline Schedule",
						Ref:          "master",
						Cron:         "0 1 * * *",
						CronTimezone: "UTC",
						Active:       true,
					}),
				),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_pipeline_schedule.schedule",
				ImportStateIdFunc: getPipelineScheduleImportID("gitlab_pipeline_schedule.schedule"),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update the pipeline schedule to change the parameters
			{
				Config: testAccGitlabPipelineScheduleUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabPipelineScheduleExists("gitlab_pipeline_schedule.schedule", &schedule),
					testAccCheckGitlabPipelineScheduleAttributes(&schedule, &testAccGitlabPipelineScheduleExpectedAttributes{
						Description:  "Schedule",
						Ref:          "master",
						Cron:         "0 4 * * *",
						CronTimezone: "UTC",
						Active:       false,
					}),
				),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_pipeline_schedule.schedule",
				ImportStateIdFunc: getPipelineScheduleImportID("gitlab_pipeline_schedule.schedule"),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update the pipeline schedule to get back to initial settings
			{
				Config: testAccGitlabPipelineScheduleConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabPipelineScheduleExists("gitlab_pipeline_schedule.schedule", &schedule),
					testAccCheckGitlabPipelineScheduleAttributes(&schedule, &testAccGitlabPipelineScheduleExpectedAttributes{
						Description:  "Pipeline Schedule",
						Ref:          "master",
						Cron:         "0 1 * * *",
						CronTimezone: "UTC",
						Active:       true,
					}),
				),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_pipeline_schedule.schedule",
				ImportStateIdFunc: getPipelineScheduleImportID("gitlab_pipeline_schedule.schedule"),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func getPipelineScheduleImportID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("Not Found: %s", n)
		}

		pipelineScheduleID := rs.Primary.ID
		if pipelineScheduleID == "" {
			return "", fmt.Errorf("No pipeline schedule ID is set")
		}
		projectID := rs.Primary.Attributes["project"]
		if projectID == "" {
			return "", fmt.Errorf("No project ID is set")
		}

		return fmt.Sprintf("%s:%s", projectID, pipelineScheduleID), nil
	}
}

func testAccCheckGitlabPipelineScheduleExists(n string, schedule *gitlab.PipelineSchedule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		scheduleID := rs.Primary.ID
		repoName := rs.Primary.Attributes["project"]
		if repoName == "" {
			return fmt.Errorf("No project ID is set")
		}

		schedules, _, err := testGitlabClient.PipelineSchedules.ListPipelineSchedules(repoName, nil)
		if err != nil {
			return err
		}
		for _, gotSchedule := range schedules {
			if strconv.Itoa(gotSchedule.ID) == scheduleID {
				*schedule = *gotSchedule
				return nil
			}
		}
		return fmt.Errorf("Pipeline Schedule does not exist")
	}
}

type testAccGitlabPipelineScheduleExpectedAttributes struct {
	Description  string
	Ref          string
	Cron         string
	CronTimezone string
	Active       bool
}

func testAccCheckGitlabPipelineScheduleAttributes(schedule *gitlab.PipelineSchedule, want *testAccGitlabPipelineScheduleExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if schedule.Description != want.Description {
			return fmt.Errorf("got description %q; want %q", schedule.Description, want.Description)
		}
		if schedule.Ref != want.Ref {
			return fmt.Errorf("got ref %q; want %q", schedule.Ref, want.Ref)
		}

		if schedule.Cron != want.Cron {
			return fmt.Errorf("got cron %q; want %q", schedule.Cron, want.Cron)
		}

		if schedule.CronTimezone != want.CronTimezone {
			return fmt.Errorf("got cron_timezone %q; want %q", schedule.CronTimezone, want.CronTimezone)
		}

		if schedule.Active != want.Active {
			return fmt.Errorf("got active %t; want %t", schedule.Active, want.Active)
		}

		return nil
	}
}

func testAccCheckGitlabPipelineScheduleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_pipeline_schedule" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("could not convert pipeline schedule id to integer: %s", err)
		}

		gotPS, _, err := testGitlabClient.PipelineSchedules.GetPipelineSchedule(rs.Primary.Attributes["project"], id)
		if err == nil {
			if gotPS != nil && fmt.Sprintf("%d", gotPS.ID) == rs.Primary.ID {
				return fmt.Errorf("pipeline schedule still exists")
			}
		}
		if !is404(err) {
			return err
		}
		return nil
	}
	return nil
}

func testAccGitlabPipelineScheduleConfig(rInt int) string {
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
	`, rInt)
}

func testAccGitlabPipelineScheduleUpdateConfig(rInt int) string {
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
  description = "Schedule"
  ref = "master"
  cron = "0 4 * * *"
  active = false
}
	`, rInt)
}
