//go:build acceptance
// +build acceptance

package provider

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabGroupSamlLink_basic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "gitlab_group_saml_link.foo"

	// PreCheck runs after Config so load test data here
	var samlLink gitlab.SAMLGroupLink
	testSamlLink := gitlab.SAMLGroupLink{
		Name: "test_saml_group",
	}

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabGroupSamlLinkDestroy,
		Steps: []resource.TestStep{

			// Create a group SAML link as a developer (uses testAccGitlabGroupLdapSamlCreateConfig for Config)
			{
				SkipFunc: isRunningInCE,
				Config: fmt.Sprintf(`
				resource "gitlab_group_saml_link" "foo" {
					group_id 		= "%d"
					access_level 	= "Developer"
					saml_group_name = "%s"
				
				}`, rInt, rInt, testSamlLink.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupSamlLinkExists(resourceName, &samlLink)),
			},

			// Import the group SAML link (re-uses testAccGitlabGroupSamlLinkCreateConfig for Config)
			{
				SkipFunc:          isRunningInCE,
				ResourceName:      resourceName,
				ImportStateIdFunc: getGitlabGroupSamlLinkImportID(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},

			// Update the group SAML link to change the access level (uses testAccGitlabGroupSamlLinkUpdateConfig for Config)
			{
				SkipFunc: isRunningInCE,
				Config: fmt.Sprintf(`
				resource "gitlab_group_saml_link" "foo" {
					group_id 		= "%d"
					access_level 	= "Maintainer"
					saml_group_name = "%s"
				}`, rInt, rInt, testSamlLink.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupSamlLinkExists(resourceName, &samlLink)),
			},
		},
	})
}

func getGitlabGroupSamlLinkImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not Found: %s", resourceName)
		}

		groupID := rs.Primary.Attributes["group_id"]
		if groupID == "" {
			return "", fmt.Errorf("No group ID is set")
		}
		samlGroupName := rs.Primary.Attributes["saml_group_name"]
		if samlGroupName == "" {
			return "", fmt.Errorf("No SAML group name is set")
		}

		return fmt.Sprintf("%s:%s", groupID, samlGroupName), nil
	}
}

func testAccCheckGitlabGroupSamlLinkExists(resourceName string, samlLink *gitlab.SAMLGroupLink) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Clear the "found" SAML link before checking for existence
		*samlLink = gitlab.SAMLGroupLink{}

		resourceState, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		err := testAccGetGitlabGroupSamlLink(samlLink, resourceState)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckGitlabGroupSamlLinkDestroy(s *terraform.State) error {
	// Can't check for links if the group is destroyed so make sure all groups are destroyed instead
	for _, resourceState := range s.RootModule().Resources {
		if resourceState.Type != "gitlab_group" {
			continue
		}

		group, _, err := testGitlabClient.Groups.GetGroup(resourceState.Primary.ID, nil)
		if err == nil {
			if group != nil && fmt.Sprintf("%d", group.ID) == resourceState.Primary.ID {
				if group.MarkedForDeletionOn == nil {
					return fmt.Errorf("Group still exists")
				}
			}
		}
		if !is404(err) {
			return err
		}
		return nil
	}
	return nil
}

func testAccGetGitlabGroupSamlLink(samlLink *gitlab.SAMLGroupLink, resourceState *terraform.ResourceState) error {
	groupId := resourceState.Primary.Attributes["group_id"]
	if groupId == "" {
		return fmt.Errorf("No group ID is set")
	}

	// Construct our desired SAML Link from the config values
	desiredSamlLink := gitlab.SAMLGroupLink{
		AccessLevel: resourceState.Primary.Attributes["access_level"],
		Name:        resourceState.Primary.Attributes["saml_group_name"],
	}

	desiredSamlLinkId := buildTwoPartID(&groupId, &desiredSamlLink.Name)

	// Try to fetch all group links from GitLab
	currentSamlLinks, _, err := testGitlabClient.Groups.ListGroupSamlLinks(groupId, nil)
	if err != nil {
		return err
	}

	found := false

	// Check if the SAML link exists in the returned list of links
	for _, currentSamlLink := range currentSamlLinks {
		if buildTwoPartID(&groupId, &currentSamlLink.Name) == desiredSamlLinkId {
			found = true
			*samlLink = *currentSamlLink
			break
		}
	}

	if !found {
		return errors.New(fmt.Sprintf("SamlLink %s does not exist.", desiredSamlLinkId)) // nolint // TODO: Resolve this golangci-lint issue: S1028: should use fmt.Errorf(...) instead of errors.New(fmt.Sprintf(...)) (gosimple)
	}

	return nil
}