package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabRunner_basic(t *testing.T) {
	group := testAccCreateGroups(t, 1)[0]
	//The runner token is not populated on the return from the group create, so re-retrieve it to get the token.
	group, _, err := testGitlabClient.Groups.GetGroup(group.ID, &gitlab.GetGroupOptions{})
	if err != nil {
		t.Fail()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() {},
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckManagedLicenseDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "gitlab_runner" "this" {
					token = "%s"
					description = "Lorem Ipsum"
				}
				`, group.RunnersToken),
			},
			{
				ResourceName:      "gitlab_runner.this",
				ImportState:       true,
				ImportStateVerify: true,
				//These need to be ignored since they don't come back in the "get" command
				ImportStateVerifyIgnore: []string{"authentication_token", "token"},
			},
			{
				Config: fmt.Sprintf(`
				resource "gitlab_runner" "this" {
					token = "%s"
					description = "Lorem Ipsum Dolor Sit Amet"
				}
				`, group.RunnersToken),
			},
			{
				ResourceName:      "gitlab_runner.this",
				ImportState:       true,
				ImportStateVerify: true,
				//These need to be ignored since they don't come back in the "get" command
				ImportStateVerifyIgnore: []string{"authentication_token", "token"},
			},
		},
	})
}

func TestAccGitlabRunner_comprehensive(t *testing.T) {
	group := testAccCreateGroups(t, 1)[0]
	//The runner token is not populated on the return from the group create, so re-retrieve it to get the token.
	group, _, err := testGitlabClient.Groups.GetGroup(group.ID, &gitlab.GetGroupOptions{})
	if err != nil {
		t.Fail()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() {},
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckManagedLicenseDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "gitlab_runner" "this" {
					token = "%s"
					description = "Lorem Ipsum"

					paused = false
					locked = false
					run_untagged = false
					tag_list = ["tag_one", "tag_two"]
					access_level = "ref_protected"
					maximum_timeout = 3600					
				}
				`, group.RunnersToken),
			},
			{
				ResourceName:      "gitlab_runner.this",
				ImportState:       true,
				ImportStateVerify: true,
				//These need to be ignored since they don't come back in the "get" command
				ImportStateVerifyIgnore: []string{"authentication_token", "token"},
			},
			{
				Config: fmt.Sprintf(`
				resource "gitlab_runner" "this" {
					token = "%s"
					description = "Lorem Ipsum Dolor Sit Amet"

					paused = true
					locked = true
					run_untagged = true
					tag_list = ["tag_one", "tag_two", "tag_three"]
					access_level = "not_protected"
					maximum_timeout = 4200
				}
				`, group.RunnersToken),
			},
			{
				ResourceName:      "gitlab_runner.this",
				ImportState:       true,
				ImportStateVerify: true,
				//These need to be ignored since they don't come back in the "get" command
				ImportStateVerifyIgnore: []string{"authentication_token", "token"},
			},
		},
	})
}
