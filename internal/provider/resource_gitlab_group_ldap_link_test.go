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

func TestAccGitlabGroupLdapLink_basic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "gitlab_group_ldap_link.foo"

	// PreCheck runs after Config so load test data here
	var ldapLink gitlab.LDAPGroupLink
	testLdapLink := gitlab.LDAPGroupLink{
		CN:       "default",
		Provider: "default",
	}

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabGroupLdapLinkDestroy,
		Steps: []resource.TestStep{

			// Create a group LDAP link as a developer (uses testAccGitlabGroupLdapLinkCreateConfig for Config)
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabGroupLdapLinkCreateConfig(rInt, &testLdapLink),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupLdapLinkExists(resourceName, &ldapLink),
					testAccCheckGitlabGroupLdapLinkAttributes(&ldapLink, &testAccGitlabGroupLdapLinkExpectedAttributes{
						accessLevel: "developer",
					})),
			},

			// Import the group LDAP link (re-uses testAccGitlabGroupLdapLinkCreateConfig for Config)
			{
				SkipFunc:          isRunningInCE,
				ResourceName:      resourceName,
				ImportStateIdFunc: getGitlabGroupLdapLinkImportID(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},

			// Update the group LDAP link to change the access level (uses testAccGitlabGroupLdapLinkUpdateConfig for Config)
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabGroupLdapLinkUpdateConfig(rInt, &testLdapLink),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupLdapLinkExists(resourceName, &ldapLink),
					testAccCheckGitlabGroupLdapLinkAttributes(&ldapLink, &testAccGitlabGroupLdapLinkExpectedAttributes{
						accessLevel: "maintainer",
					})),
			},
		},
	})
}

func getGitlabGroupLdapLinkImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not Found: %s", resourceName)
		}

		groupID := rs.Primary.Attributes["group_id"]
		if groupID == "" {
			return "", fmt.Errorf("No group ID is set")
		}
		ldapProvider := rs.Primary.Attributes["ldap_provider"]
		if ldapProvider == "" {
			return "", fmt.Errorf("No LDAP provider is set")
		}
		ldapCN := rs.Primary.Attributes["cn"]
		if ldapCN == "" {
			return "", fmt.Errorf("No LDAP CN is set")
		}

		return fmt.Sprintf("%s:%s:%s", groupID, ldapProvider, ldapCN), nil
	}
}

func testAccCheckGitlabGroupLdapLinkExists(resourceName string, ldapLink *gitlab.LDAPGroupLink) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Clear the "found" LDAP link before checking for existence
		*ldapLink = gitlab.LDAPGroupLink{}

		resourceState, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		err := testAccGetGitlabGroupLdapLink(ldapLink, resourceState)
		if err != nil {
			return err
		}

		return nil
	}
}

type testAccGitlabGroupLdapLinkExpectedAttributes struct {
	accessLevel string
}

func testAccCheckGitlabGroupLdapLinkAttributes(ldapLink *gitlab.LDAPGroupLink, want *testAccGitlabGroupLdapLinkExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		accessLevelId, ok := accessLevelValueToName[ldapLink.GroupAccess]
		if !ok {
			return fmt.Errorf("Invalid access level '%s'", accessLevelId)
		}
		if accessLevelId != want.accessLevel {
			return fmt.Errorf("Has access level %s; want %s", accessLevelId, want.accessLevel)
		}
		return nil
	}
}

func testAccCheckGitlabGroupLdapLinkDestroy(s *terraform.State) error {
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

func testAccGetGitlabGroupLdapLink(ldapLink *gitlab.LDAPGroupLink, resourceState *terraform.ResourceState) error {
	groupId := resourceState.Primary.Attributes["group_id"]
	if groupId == "" {
		return fmt.Errorf("No group ID is set")
	}

	// Construct our desired LDAP Link from the config values
	desiredLdapLink := gitlab.LDAPGroupLink{
		CN:          resourceState.Primary.Attributes["cn"],
		GroupAccess: accessLevelNameToValue[resourceState.Primary.Attributes["group_access"]],
		Provider:    resourceState.Primary.Attributes["ldap_provider"],
	}

	desiredLdapLinkId := buildTwoPartID(&desiredLdapLink.Provider, &desiredLdapLink.CN)

	// Try to fetch all group links from GitLab
	currentLdapLinks, _, err := testGitlabClient.Groups.ListGroupLDAPLinks(groupId, nil)
	if err != nil {
		// The read/GET API wasn't implemented in GitLab until version 12.8 (March 2020, well after the add and delete APIs).
		// If we 404, assume GitLab is at an older version and take things on faith.
		switch err.(type) { // nolint // TODO: Resolve this golangci-lint issue: S1034: assigning the result of this type assertion to a variable (switch err := err.(type)) could eliminate type assertions in switch cases (gosimple)
		case *gitlab.ErrorResponse:
			if err.(*gitlab.ErrorResponse).Response.StatusCode == 404 { // nolint // TODO: Resolve this golangci-lint issue: S1034(related information): could eliminate this type assertion (gosimple)
				// Do nothing
			} else {
				return err
			}
		default:
			return err
		}
	}

	// If we got here and don't have links, assume GitLab is below version 12.8 and skip the check
	if currentLdapLinks != nil {
		found := false

		// Check if the LDAP link exists in the returned list of links
		for _, currentLdapLink := range currentLdapLinks {
			if buildTwoPartID(&currentLdapLink.Provider, &currentLdapLink.CN) == desiredLdapLinkId {
				found = true
				*ldapLink = *currentLdapLink
				break
			}
		}

		if !found {
			return errors.New(fmt.Sprintf("LdapLink %s does not exist.", desiredLdapLinkId)) // nolint // TODO: Resolve this golangci-lint issue: S1028: should use fmt.Errorf(...) instead of errors.New(fmt.Sprintf(...)) (gosimple)
		}
	} else {
		*ldapLink = desiredLdapLink
	}

	return nil
}

func testAccGitlabGroupLdapLinkCreateConfig(rInt int, testLdapLink *gitlab.LDAPGroupLink) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
    name = "foo%d"
	path = "foo%d"
	description = "Terraform acceptance test - Group LDAP Links 1"
}

resource "gitlab_group_ldap_link" "foo" {
    group_id 		= "${gitlab_group.foo.id}"
    cn				= "%s"
	group_access 	= "developer"
	ldap_provider   = "%s"

}`, rInt, rInt, testLdapLink.CN, testLdapLink.Provider)
}

func testAccGitlabGroupLdapLinkUpdateConfig(rInt int, testLdapLink *gitlab.LDAPGroupLink) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
    name = "foo%d"
	path = "foo%d"
	description = "Terraform acceptance test - Group LDAP Links 2"
}

resource "gitlab_group_ldap_link" "foo" {
    group_id 		= "${gitlab_group.foo.id}"
    cn				= "%s"
	group_access 	= "maintainer"
	ldap_provider   = "%s"
}`, rInt, rInt, testLdapLink.CN, testLdapLink.Provider)
}
