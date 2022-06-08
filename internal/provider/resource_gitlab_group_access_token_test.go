//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

func TestAccGitlabGroupAccessToken_basic(t *testing.T) {
	var gat testAccGitlabGroupAccessTokenWrapper
	var groupVariable gitlab.GroupVariable

	testGroup := testAccCreateGroups(t, 1)[0]

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabGroupAccessTokenDestroy,
		Steps: []resource.TestStep{
			// Create a Group and a Group Access Token
			{
				Config: testAccGitlabGroupAccessTokenConfig(testGroup.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupAccessTokenExists("gitlab_group_access_token.this", &gat),
					testAccCheckGitlabGroupAccessTokenAttributes(&gat, &testAccGitlabGroupAccessTokenExpectedAttributes{
						name:        "my group token",
						scopes:      map[string]bool{"read_repository": true, "api": true, "write_repository": true, "read_api": true},
						expiresAt:   "2099-01-01",
						accessLevel: gitlab.AccessLevelValue(gitlab.DeveloperPermissions),
					}),
				),
			},
			// Update the Group Access Token to change the parameters
			{
				Config: testAccGitlabGroupAccessTokenUpdateConfig(testGroup.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupAccessTokenExists("gitlab_group_access_token.this", &gat),
					testAccCheckGitlabGroupAccessTokenAttributes(&gat, &testAccGitlabGroupAccessTokenExpectedAttributes{
						name:        "my new group token",
						scopes:      map[string]bool{"read_repository": false, "api": true, "write_repository": false, "read_api": false},
						expiresAt:   "2099-05-01",
						accessLevel: gitlab.AccessLevelValue(gitlab.MaintainerPermissions),
					}),
				),
			},
			// Update the Group Access Token Access Level to Owner
			{
				Config: testAccGitlabGroupAccessTokenUpdateAccessLevel(testGroup.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupAccessTokenExists("gitlab_group_access_token.this", &gat),
					testAccCheckGitlabGroupAccessTokenAttributes(&gat, &testAccGitlabGroupAccessTokenExpectedAttributes{
						name:        "my new group token",
						scopes:      map[string]bool{"read_repository": false, "api": true, "write_repository": false, "read_api": false},
						expiresAt:   "2099-05-01",
						accessLevel: gitlab.AccessLevelValue(gitlab.OwnerPermissions),
					}),
				),
			},
			// Add a CICD variable with Group Access Token value
			{
				Config: testAccGitlabGroupAccessTokenUpdateConfigWithCICDvar(testGroup.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupAccessTokenExists("gitlab_group_access_token.this", &gat),
					testAccCheckGitlabGroupVariableExists("gitlab_group_variable.var", &groupVariable),
					testAccCheckGitlabGroupAccessTokenAttributes(&gat, &testAccGitlabGroupAccessTokenExpectedAttributes{
						name:        "my new group token",
						scopes:      map[string]bool{"read_repository": false, "api": true, "write_repository": false, "read_api": false},
						expiresAt:   "2099-05-01",
						accessLevel: gitlab.AccessLevelValue(gitlab.MaintainerPermissions),
					}),
				),
			},
			//Restore Group Access Token initial parameters
			{
				Config: testAccGitlabGroupAccessTokenConfig(testGroup.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupAccessTokenExists("gitlab_group_access_token.this", &gat),
					testAccCheckGitlabGroupAccessTokenAttributes(&gat, &testAccGitlabGroupAccessTokenExpectedAttributes{
						name:        "my group token",
						scopes:      map[string]bool{"read_repository": true, "api": true, "write_repository": true, "read_api": true},
						expiresAt:   "2099-01-01",
						accessLevel: gitlab.AccessLevelValue(gitlab.DeveloperPermissions),
					}),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_group_access_token.this",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// the token is only known during creating. We explicitly mention this limitation in the docs.
					"token",
				},
			},
		},
	})
}

func testAccCheckGitlabGroupAccessTokenExists(n string, gat *testAccGitlabGroupAccessTokenWrapper) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		group, tokenString, err := parseTwoPartID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing ID: %s", rs.Primary.ID)
		}
		groupAccessTokenID, err := strconv.Atoi(tokenString)
		if err != nil {
			return fmt.Errorf("%s cannot be converted to int", tokenString)
		}

		groupId := rs.Primary.Attributes["group"]
		if groupId == "" {
			return fmt.Errorf("No group ID is set")
		}
		if groupId != group {
			return fmt.Errorf("Group [%s] in group identifier [%s] it's different from group stored into the state [%s]", group, rs.Primary.ID, groupId)
		}

		tokens, _, err := testGitlabClient.GroupAccessTokens.ListGroupAccessTokens(groupId, nil)
		if err != nil {
			return err
		}

		for _, token := range tokens {
			if token.ID == groupAccessTokenID {
				gat.groupAccessToken = token
				gat.group = groupId
				gat.token = rs.Primary.Attributes["token"]
				return nil
			}
		}
		return fmt.Errorf("Group Access Token does not exist")
	}
}

type testAccGitlabGroupAccessTokenExpectedAttributes struct {
	name        string
	scopes      map[string]bool
	expiresAt   string
	accessLevel gitlab.AccessLevelValue
}

type testAccGitlabGroupAccessTokenWrapper struct {
	groupAccessToken *gitlab.GroupAccessToken
	group            string
	token            string
}

func testAccCheckGitlabGroupAccessTokenAttributes(gatWrap *testAccGitlabGroupAccessTokenWrapper, want *testAccGitlabGroupAccessTokenExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		gat := gatWrap.groupAccessToken
		if gat.Name != want.name {
			return fmt.Errorf("got Name %q; want %q", gat.Name, want.name)
		}

		if gat.ExpiresAt.String() != want.expiresAt {
			return fmt.Errorf("got ExpiresAt %q; want %q", gat.ExpiresAt.String(), want.expiresAt)
		}

		if gat.AccessLevel != want.accessLevel {
			return fmt.Errorf("got AccessLevel %q; want %q", gat.AccessLevel, want.accessLevel)
		}

		for _, scope := range gat.Scopes {
			if !want.scopes[scope] {
				return fmt.Errorf("got a not wanted Scope %q, received %v", scope, gat.Scopes)
			}
			want.scopes[scope] = false
		}
		for k, v := range want.scopes {
			if v {
				return fmt.Errorf("not got a wanted Scope %q, received %v", k, gat.Scopes)
			}
		}

		git, err := gitlab.NewClient(gatWrap.token, gitlab.WithBaseURL(testGitlabClient.BaseURL().String()))
		if err != nil {
			return fmt.Errorf("Cannot use the token to instantiate a new client %s", err)
		}
		_, _, err = git.Groups.GetGroup(gatWrap.group, nil)
		if err != nil {
			return fmt.Errorf("Cannot use the token to perform an API call %s", err)
		}

		return nil
	}
}

func testAccCheckGitlabGroupAccessTokenDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_group" {
			continue
		}

		group, resp, err := testGitlabClient.Groups.GetGroup(rs.Primary.ID, nil)
		if err == nil {
			if group != nil && fmt.Sprintf("%d", group.ID) == rs.Primary.ID {
				if group.MarkedForDeletionOn == nil {
					return fmt.Errorf("Group still exists")
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

func testAccGitlabGroupAccessTokenConfig(groupId int) string {
	return fmt.Sprintf(`
resource "gitlab_group_access_token" "this" {
  name = "my group token"
  group = %d
  expires_at = "2099-01-01"
  access_level = "developer"
  scopes = ["read_repository" , "api", "write_repository", "read_api"]
}
	`, groupId)
}

func testAccGitlabGroupAccessTokenUpdateConfig(groupId int) string {
	return fmt.Sprintf(`
resource "gitlab_group_access_token" "this" {
  name = "my new group token"
  group = %d
  expires_at = "2099-05-01"
  access_level = "maintainer"
  scopes = ["api"]
}
	`, groupId)
}

func testAccGitlabGroupAccessTokenUpdateAccessLevel(groupId int) string {
	return fmt.Sprintf(`
resource "gitlab_group_access_token" "this" {
  name = "my new group token"
  group = %d
  expires_at = "2099-05-01"
  access_level = "owner"
  scopes = ["api"]
}
	`, groupId)
}

func testAccGitlabGroupAccessTokenUpdateConfigWithCICDvar(groupId int) string {
	return fmt.Sprintf(`
resource "gitlab_group_access_token" "this" {
  name = "my new group token"
  group = %d
  expires_at = "2099-05-01"
  access_level = "maintainer"
  scopes = ["api"]
}

resource "gitlab_group_variable" "var" {
  group   = %d
  key     = "my_grp_access_token"
  value   = gitlab_group_access_token.this.token
 }

	`, groupId, groupId)
}
