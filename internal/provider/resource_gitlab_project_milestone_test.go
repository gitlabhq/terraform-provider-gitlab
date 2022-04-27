package provider

import (
    "errors"
    "fmt"
    "time"
    "testing"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
    "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
    gitlab "github.com/xanzy/go-gitlab"
)

func TestAccGitlabProjectMilestone_basic(t *testing.T) {
    testAccCheck(t)

    var milestone gitlab.Milestone
    var milestoneUpdate gitlab.Milestone
    rInt1, rInt2 := acctest.RandInt(), acctest.RandInt()
    project := testAccCreateProject(t)

    resource.Test(t, resource.TestCase{
        PreCheck:          func() { testAccPreCheck(t) },
        ProviderFactories: providerFactories,
        CheckDestroy:      testAccCheckGitlabProjectMilestoneDestroy,
        Steps: []resource.TestStep{
            {
                // create Milestone with required values only
                Config: testAccGitlabProjectMilestoneConfigRequiredOnly(project.ID, rInt1, ""),
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckGitlabProjectMilestoneExists("this", &milestone),
                    testAccCheckGitlabProjectMilestoneAttributes("this", &milestone, &testAccGitlabProjectMilestoneExpectedAttributes{
                        Title:       fmt.Sprintf("test-%d", rInt1),
                        ProjectID:   project.ID,
                        Description: "",
                        StartDate:   gitlab.ISOTime{},
                        DueDate:     gitlab.ISOTime{},
                        State:       "active",
                        Expired:     false,
                    }),
                ),
            },
            {
                // verify import
                ResourceName:      "gitlab_project_milestone.this",
                ImportState:       true,
                ImportStateVerify: true,
            },
            {
                // update some Milestone attributes
                Config: testAccGitlabProjectMilestoneConfigAll(project.ID, rInt2, "2022-04-10", "2022-04-15", "closed"),
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckGitlabProjectMilestoneExists("this", &milestoneUpdate),
                    testAccCheckGitlabProjectMilestoneAttributes("this", &milestoneUpdate, &testAccGitlabProjectMilestoneExpectedAttributes{
                        Title:       fmt.Sprintf("test-%d", rInt2),
                        ProjectID:   project.ID,
                        Description: fmt.Sprintf("test-%d", rInt2),
                        StartDate:   gitlab.ISOTime(time.Date(2022, time.April, 10, 0, 0, 0, 0, time.UTC)),
                        DueDate:     gitlab.ISOTime(time.Date(2022, time.April, 15, 0, 0, 0, 0, time.UTC)),
                        State:       "closed",
                        Expired:     true,
                    }),
                ),
            },
        },
    })
}

func testAccCheckGitlabProjectMilestoneDestroy(s *terraform.State) error {
    for _, rs := range s.RootModule().Resources {
        if rs.Type != "gitlab_project_milestone" {
            continue
        }
        project, milestoneID, err := resourceGitLabProjectMilestoneParseId(rs.Primary.ID)
        if err != nil {
            return err
        }

        milestone, _, err := testGitlabClient.Milestones.GetMilestone(project, milestoneID)
        if err == nil && milestone != nil {
            return errors.New("Milestone still exists")
        }
        if !is404(err) {
            return err
        }
        return nil
    }
    return nil
}

func testAccCheckGitlabProjectMilestoneAttributes(n string, milestone *gitlab.Milestone, want *testAccGitlabProjectMilestoneExpectedAttributes) resource.TestCheckFunc {
    return func(s *terraform.State) error {
        if milestone.Title != want.Title {
            return fmt.Errorf("Got milestone title '%s'; want '%s'", milestone.Title, want.Title)
        }
        if milestone.ProjectID != want.ProjectID {
            return fmt.Errorf("Got milestone project_id '%d'; want '%d'", milestone.ProjectID, want.ProjectID)
        }
        if milestone.Description != want.Description {
            return fmt.Errorf("Got milestone description '%s'; want '%s'", milestone.Description, want.Description)
        }
        startDate := gitlab.ISOTime(time.Date(0001, time.January, 1, 0, 0, 0, 0, time.UTC))
        if milestone.StartDate != nil {
            startDate = *milestone.StartDate
        }
        if startDate != want.StartDate {
            return fmt.Errorf("Got milestone start_date '%s'; want '%s'", milestone.StartDate, want.StartDate)
        }
        dueDate := gitlab.ISOTime(time.Date(0001, time.January, 1, 0, 0, 0, 0, time.UTC))
        if milestone.DueDate != nil {
            dueDate = *milestone.DueDate
        }
        if dueDate != want.DueDate {
            return fmt.Errorf("Got milestone due_date '%s'; want '%s'", milestone.DueDate, want.DueDate)
        }
        if milestone.State != want.State {
            return fmt.Errorf("Got milestone state '%s'; want '%s'", milestone.State, want.State)
        }
        expired := false
        if milestone.Expired != nil {
            expired = *milestone.Expired
        }
        if expired != want.Expired {
            return fmt.Errorf("Got milestone expired '%v'; want '%v'", milestone.Expired, want.Expired)
        }
        return nil
    }
}

func testAccCheckGitlabProjectMilestoneExists(n string, milestone *gitlab.Milestone) resource.TestCheckFunc {
    return func(s *terraform.State) error {
        rs, ok := s.RootModule().Resources[fmt.Sprintf("gitlab_project_milestone.%s", n)]
        if !ok {
            return fmt.Errorf("Not Found: %s", n)
        }
        project, milestoneID, err := resourceGitLabProjectMilestoneParseId(rs.Primary.ID)
        if err != nil {
            return fmt.Errorf("Error in splitting project and milestoneID")
        }
        gotMilestone, _, err := testGitlabClient.Milestones.GetMilestone(project, milestoneID)
        if err != nil {
            return err
        }
        *milestone = *gotMilestone
        return err
    }
}

func testAccGitlabProjectMilestoneConfigRequiredOnly(project int, rInt int, additinalOptions string) string {
    return fmt.Sprintf(`
    resource "gitlab_project_milestone" "this" {
        project_id  = "%d"
        title       = "test-%d"
        %s
    }
  `, project, rInt, additinalOptions)
}

func testAccGitlabProjectMilestoneConfigAll(project int, rInt int, startDate string, dueDate string, state string) string {
    additinalOptions := fmt.Sprintf(`
        description = "test-%d"
        start_date  = "%s"
        due_date    = "%s"
        state       = "%s"
  `, rInt, startDate, dueDate, state)
  return testAccGitlabProjectMilestoneConfigRequiredOnly(project, rInt, additinalOptions)
}

type testAccGitlabProjectMilestoneExpectedAttributes struct {
    Title       string
    ProjectID   int
    Description string
    StartDate   gitlab.ISOTime
    DueDate     gitlab.ISOTime
    State       string
    Expired     bool
}
