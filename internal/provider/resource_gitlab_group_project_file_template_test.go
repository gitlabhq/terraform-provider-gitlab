//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabGroupProjectFileTemplate_basic(t *testing.T) {
	// Since we do some manual setup in this test, we need to handle the test skip first.
	baseGroup := testAccCreateGroups(t, 1)[0]
	firstProject := testAccCreateProjectWithNamespace(t, baseGroup.ID)
	secondProject := testAccCreateProjectWithNamespace(t, baseGroup.ID)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckProjectFileTemplateDestroy,
		Steps: []resource.TestStep{
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGroupProjectFileTemplateConfig(baseGroup.ID, firstProject.ID),
				Check: resource.ComposeTestCheckFunc(
					// Note - we can't use the testAccCheckGitlabGroupAttributes, because that checks the TF
					// state attributes, and file project template explicitly doesn't exist there.
					testAccCheckGitlabGroupFileTemplateValue(baseGroup, firstProject),
					resource.TestCheckResourceAttr("gitlab_group_project_file_template.linking_template", "group_id", strconv.Itoa(baseGroup.ID)),
					resource.TestCheckResourceAttr("gitlab_group_project_file_template.linking_template", "file_template_project_id", strconv.Itoa(firstProject.ID)),
				),
			},
			{
				//Test that when we update the project name, it re-links the group to the new project
				SkipFunc: isRunningInCE,
				Config:   testAccGroupProjectFileTemplateConfig(baseGroup.ID, secondProject.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupFileTemplateValue(baseGroup, secondProject),
					resource.TestCheckResourceAttr("gitlab_group_project_file_template.linking_template", "group_id", strconv.Itoa(baseGroup.ID)),
					resource.TestCheckResourceAttr("gitlab_group_project_file_template.linking_template", "file_template_project_id", strconv.Itoa(secondProject.ID)),
				),
			},
		},
	},
	)
}

func testAccCheckGitlabGroupFileTemplateValue(g *gitlab.Group, p *gitlab.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		//Re-retrieve the group to ensure we have the most up-to-date group info
		g, _, err := testGitlabClient.Groups.GetGroup(g.ID, &gitlab.GetGroupOptions{})
		if is404(err) {
			return fmt.Errorf("Group no longer exists, expected group to exist with a file_template_project_id")
		}

		if g.FileTemplateProjectID == p.ID {
			return nil
		}
		return fmt.Errorf("Group file_template_project_id doesn't match. Wanted %d, received %d", p.ID, g.FileTemplateProjectID)
	}
}

func testAccCheckProjectFileTemplateDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "gitlab_group_project_file_template" {
			continue
		}

		// To test if the resource was destroyed, we need to retrieve the group.
		gid := rs.Primary.ID
		group, _, err := testGitlabClient.Groups.GetGroup(gid, nil)
		if err != nil {
			return err
		}

		// the test should succeed if the group is still present and has a 0 file_template_project_id value
		if group != nil && group.FileTemplateProjectID != 0 {
			return fmt.Errorf("Group still has a template project attached")
		}
		return nil
	}
	return nil
}

func testAccGroupProjectFileTemplateConfig(groupID int, projectID int) string {
	return fmt.Sprintf(
		`
resource "gitlab_group_project_file_template" "linking_template" {
 group_id = %d
 file_template_project_id = %d
}
`, groupID, projectID)
}
