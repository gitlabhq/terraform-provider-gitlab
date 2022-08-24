//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccGitlabGroupSamlLink_basic(t *testing.T) {
	testAccCheckEE(t)
	testAccRequiresAtLeast(t, "15.3")

	testGroup := testAccCreateGroups(t, 1)[0]

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabGroupSamlLinkDestroy,
		Steps: []resource.TestStep{

			// Create a group SAML link as a developer
			{
				Config: fmt.Sprintf(`
					resource "gitlab_group_saml_link" "this" {
						group   		= "%d"
						access_level 	= "developer"
						saml_group_name = "test_saml_group"

					}
				`, testGroup.ID),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_group_saml_link.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update the group SAML link to change the access level
			{
				Config: fmt.Sprintf(`
					resource "gitlab_group_saml_link" "this" {
						group   		= "%d"
						access_level 	= "maintainer"
						saml_group_name = "test_saml_group"

					}
				`, testGroup.ID),
			},
		},
	})
}

func testAccCheckGitlabGroupSamlLinkDestroy(s *terraform.State) error {
	for _, resourceState := range s.RootModule().Resources {
		if resourceState.Type != "gitlab_group_saml_link" {
			continue
		}

		group, samlGroupName, err := parseTwoPartID(resourceState.Primary.ID)
		if err != nil {
			return err
		}

		samlGroupLink, _, err := testGitlabClient.Groups.GetGroupSAMLLink(group, samlGroupName)
		if err == nil {
			if samlGroupLink != nil {
				return fmt.Errorf("SAML Group Link still exists")
			}
		}
		if !is404(err) {
			return err
		}
		return nil
	}
	return nil
}
