//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabGroupShareGroup_basic(t *testing.T) {
	mainGroup := testAccCreateGroups(t, 1)[0]
	sharedGroup := testAccCreateGroups(t, 1)[0]

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabShareGroupDestroy,
		Steps: []resource.TestStep{
			// Share a new group with another group
			{
				Config: testAccGitlabGroupShareGroupConfig(mainGroup.ID, sharedGroup.ID,
					`
					group_access 	 = "guest"
					expires_at     = "2099-01-01"
					`,
				),
				Check: testAccCheckGitlabGroupSharedWithGroup(mainGroup.Name, sharedGroup.Name, "2099-01-01", gitlab.GuestPermissions),
			},
			{
				// Verify Import
				ResourceName:      "gitlab_group_share_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update the share group
			{
				Config: testAccGitlabGroupShareGroupConfig(mainGroup.ID, sharedGroup.ID, `group_access = "reporter"`),
				Check:  testAccCheckGitlabGroupSharedWithGroup(mainGroup.Name, sharedGroup.Name, "", gitlab.ReporterPermissions),
			},
			{
				// Verify Import
				ResourceName:      "gitlab_group_share_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update share group back to initial settings
			{
				Config: testAccGitlabGroupShareGroupConfig(mainGroup.ID, sharedGroup.ID,
					`
					group_access 	 = "guest"
					expires_at     = "2099-01-01"
					`,
				),
				Check: testAccCheckGitlabGroupSharedWithGroup(mainGroup.Name, sharedGroup.Name, "2099-01-01", gitlab.GuestPermissions),
			},
			{
				// Verify Import
				ResourceName:      "gitlab_group_share_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGitlabGroupSharedWithGroup(mainGroupName string, sharedGroupName string, expireTime string, accessLevel gitlab.AccessLevelValue) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		mainGroup, _, err := testGitlabClient.Groups.GetGroup(mainGroupName, nil)
		if err != nil {
			return err
		}

		sharedGroupsCount := len(mainGroup.SharedWithGroups)
		if sharedGroupsCount != 1 {
			return fmt.Errorf("Number of shared groups was %d (wanted %d)", sharedGroupsCount, 1)
		}

		sharedGroup := mainGroup.SharedWithGroups[0]

		if sharedGroup.GroupName != sharedGroupName {
			return fmt.Errorf("group name was %s (wanted %s)", sharedGroup.GroupName, sharedGroupName)
		}

		if gitlab.AccessLevelValue(sharedGroup.GroupAccessLevel) != accessLevel {
			return fmt.Errorf("groupAccessLevel was %d (wanted %d)", sharedGroup.GroupAccessLevel, accessLevel)
		}

		if sharedGroup.ExpiresAt == nil && expireTime != "" {
			return fmt.Errorf("expire time was nil (wanted %s)", expireTime)
		} else if sharedGroup.ExpiresAt != nil && sharedGroup.ExpiresAt.String() != expireTime {
			return fmt.Errorf("expire time was %s (wanted %s)", sharedGroup.ExpiresAt.String(), expireTime)
		}

		return nil
	}
}

func testAccCheckGitlabShareGroupDestroy(s *terraform.State) error {
	var groupId string
	var sharedGroupId int
	var err error

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "gitlab_group_share_group" {
			groupId, sharedGroupId, err = groupIdsFromId(rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("[ERROR] cannot get Group ID and ShareGroupId from input: %v", rs.Primary.ID)
			}

			// Get Main Group
			group, _, err := testGitlabClient.Groups.GetGroup(groupId, nil)
			if err != nil {
				return err
			}

			// Make sure that SharedWithGroups attribute on the main group does not contain the shared group id at all
			for _, sharedGroup := range group.SharedWithGroups {
				if sharedGroupId == sharedGroup.GroupID {
					return fmt.Errorf("GitLab Group Share %d still exists", sharedGroupId)
				}
			}
		}
	}

	return nil
}

func testAccGitlabGroupShareGroupConfig(mainGroupId int, shareGroupId int, shareGroupSettings string) string {
	return fmt.Sprintf(
		`
		resource "gitlab_group_share_group" "test" {
		  group_id       = %[1]d
			share_group_id = %[2]d
			%[3]s
		}
		`,
		mainGroupId,
		shareGroupId,
		shareGroupSettings,
	)
}
