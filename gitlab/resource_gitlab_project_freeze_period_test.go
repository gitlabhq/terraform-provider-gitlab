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

func TestAccGitlabFreezePeriod_basic(t *testing.T) {
	var schedule gitlab.FreezePeriod
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabFreezePeriodDestroy,
		Steps: []resource.TestStep{
			// Create a project and freeze period with default options
			{
				Config: testAccGitlabFreezePeriodConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabFreezePeriodExists("gitlab_project_freeze_period.schedule", &schedule),
					testAccCheckGitlabFreezePeriodAttributes(&schedule, &testAccGitlabFreezePeriodExpectedAttributes{
						FreezeStart:  "0 23 * * 5",
						FreezeEnd:    "0 7 * * 1",
						CronTimezone: "UTC",
					}),
				),
			},
			// Update the freeze period to change the parameters
			{
				Config: testAccGitlabFreezePeriodUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabFreezePeriodExists("gitlab_project_freeze_period.schedule", &schedule),
					testAccCheckGitlabFreezePeriodAttributes(&schedule, &testAccGitlabFreezePeriodExpectedAttributes{
						FreezeStart:  "0 23 * * 5",
						FreezeEnd:    "0 7 * * 3",
						CronTimezone: "UTC",
					}),
				),
			},
			// Update the freeze period to get back to initial settings
			{
				Config: testAccGitlabFreezePeriodConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabFreezePeriodExists("gitlab_project_freeze_period.schedule", &schedule),
					testAccCheckGitlabFreezePeriodAttributes(&schedule, &testAccGitlabFreezePeriodExpectedAttributes{
						FreezeStart:  "0 23 * * 5",
						FreezeEnd:    "0 7 * * 1",
						CronTimezone: "UTC",
					}),
				),
			},
		},
	})
}

func testAccCheckGitlabFreezePeriodExists(n string, freezePeriod *gitlab.FreezePeriod) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		scheduleID := rs.Primary.ID
		repoName := rs.Primary.Attributes["project_id"]
		if repoName == "" {
			return fmt.Errorf("No project ID is set")
		}
		conn := testAccProvider.Meta().(*gitlab.Client)

		freezePeriods, _, err := conn.FreezePeriods.ListFreezePeriods(repoName, nil)
		if err != nil {
			return err
		}
		for _, gotFreezePeriod := range freezePeriods {
			if strconv.Itoa(gotFreezePeriod.ID) == scheduleID {
				*freezePeriod = *gotFreezePeriod
				return nil
			}
		}
		return fmt.Errorf("Freeze Period does not exist")
	}
}

type testAccGitlabFreezePeriodExpectedAttributes struct {
	FreezeStart  string
	FreezeEnd    string
	CronTimezone string
}

func testAccCheckGitlabFreezePeriodAttributes(freezePeriod *gitlab.FreezePeriod, want *testAccGitlabFreezePeriodExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if freezePeriod.FreezeStart != want.FreezeStart {
			return fmt.Errorf("got freeze_start %q; want %q", freezePeriod.FreezeStart, want.FreezeStart)
		}
		if freezePeriod.FreezeEnd != want.FreezeEnd {
			return fmt.Errorf("got freeze_end %q; want %q", freezePeriod.FreezeEnd, want.FreezeEnd)
		}

		if freezePeriod.CronTimezone != want.CronTimezone {
			return fmt.Errorf("got cron_timezone %q; want %q", freezePeriod.CronTimezone, want.CronTimezone)
		}

		return nil
	}
}

func testAccCheckGitlabFreezePeriodDestroy(s *terraform.State) error {
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

func testAccGitlabFreezePeriodConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_project_freeze_period" "schedule" {
	project_id = gitlab_project.foo.id
	freeze_start = "0 23 * * 5"
	freeze_end =  "0 7 * * 1"
	cron_timezone = "UTC"
}
	`, rInt)
}

func testAccGitlabFreezePeriodUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_project_freeze_period" "schedule" {
  project_id = gitlab_project.foo.id
  freeze_start = "0 23 * * 5"
  freeze_end =  "0 7 * * 3"
  cron_timezone = "UTC"
}
	`, rInt)
}
