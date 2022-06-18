//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
)

var testRSAPubKey string = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCi+ErxScCKIVqg2ZRJ6Mx2Yd/RTsh2DGyhUR8z8Iey4rpi1YOBlpTgjxxnSLy26J++Un/iWYDP8wMvEjXElkWz3z4I+Z3mfF3dv039FTCu+O17Mw20Ek4DJxdrKvOgul040sUG/ABVHo6DjqjokjoVJwzUrUmoOtbeMMD8hFN9bWdEVyTj18XQO8nvEe/VkbhCRhAlZC1l60fM07/7Tw83SV5UNAnBtOB+nfa3b24baO+Ijc4+PqYcBuUAF6DvhXW2gZPqf5wjDBJqlDlRTYDdHarMXZAKBpWfWj0gntbtEOM+Fnp6hS1HajaeveNSs6yQwgQEDN2boQnDuvXJ8Y7zW3YQKZp8z0uqWYJSIrYRKVEVYL7gDWL9NvdRV52d/RKPnE/BlL2chiAWBRCT8buQdjVtEPPoYbA1667PXZg6PI9yhCGEIjCj71XzPssA6VL/R7yUafsmNLsirWz9Uyh3HJWCcgNuO9mglP5nfFHIXSHQVhEUEYMfzv1iX5FrenU= terraform@gitlab.com"
var testRSAPubKeyUpdatedComment string = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCi+ErxScCKIVqg2ZRJ6Mx2Yd/RTsh2DGyhUR8z8Iey4rpi1YOBlpTgjxxnSLy26J++Un/iWYDP8wMvEjXElkWz3z4I+Z3mfF3dv039FTCu+O17Mw20Ek4DJxdrKvOgul040sUG/ABVHo6DjqjokjoVJwzUrUmoOtbeMMD8hFN9bWdEVyTj18XQO8nvEe/VkbhCRhAlZC1l60fM07/7Tw83SV5UNAnBtOB+nfa3b24baO+Ijc4+PqYcBuUAF6DvhXW2gZPqf5wjDBJqlDlRTYDdHarMXZAKBpWfWj0gntbtEOM+Fnp6hS1HajaeveNSs6yQwgQEDN2boQnDuvXJ8Y7zW3YQKZp8z0uqWYJSIrYRKVEVYL7gDWL9NvdRV52d/RKPnE/BlL2chiAWBRCT8buQdjVtEPPoYbA1667PXZg6PI9yhCGEIjCj71XzPssA6VL/R7yUafsmNLsirWz9Uyh3HJWCcgNuO9mglP5nfFHIXSHQVhEUEYMfzv1iX5FrenU= terraform2@foo.com"
var updatedRSAPubKey string = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDStVqW58VZ5afXFphIvu2JahndXslJZMkgWsNiYCNdk/NvrEbc4i7yZVoDPFQsbS9I6Ty1RMW7qy3KxJalMsVHcw8arCQFDxs/ka1NHGCUPl68t5ZxUOl900KRQ0lOzGnDQMqG/UUZdPw4CCmigTr6Z9ZBcD1fXAiUwbXR4tWrr5z9KWXC2HgF4WkIJUTIct7ilY1m9W0y79dI/+K8bZrurn3q2QK83pxqqWkLwvUsCxtlhMpwuyflyzyuz8xPZl2GlZgxeIpr68gsPHIzzWizibwFfbRYKCZO4wD0r7JCDOYs9KjcIPpCG6d3HUqijClgdQSBnLwHTdE04ZtdzO8akvy0hMzRCooI5TSc8IAHos53Gp9aaW92sPA8za+WRP6OSH6UsOW4N+iQc4jyl7/fckMSgIZlJouNqqV+P8iqIFJGs70Tj5L8G/m+P2lc3kcE4Vjmj+Fc0xG5+I/PsSOpcc6DfDfZdVDRe8yklYd/qC1jI89OCeqjxu3XcUGHj9s= terraform@gitlab.com"
var updatedRSAPubKeyWithoutComment string = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDStVqW58VZ5afXFphIvu2JahndXslJZMkgWsNiYCNdk/NvrEbc4i7yZVoDPFQsbS9I6Ty1RMW7qy3KxJalMsVHcw8arCQFDxs/ka1NHGCUPl68t5ZxUOl900KRQ0lOzGnDQMqG/UUZdPw4CCmigTr6Z9ZBcD1fXAiUwbXR4tWrr5z9KWXC2HgF4WkIJUTIct7ilY1m9W0y79dI/+K8bZrurn3q2QK83pxqqWkLwvUsCxtlhMpwuyflyzyuz8xPZl2GlZgxeIpr68gsPHIzzWizibwFfbRYKCZO4wD0r7JCDOYs9KjcIPpCG6d3HUqijClgdQSBnLwHTdE04ZtdzO8akvy0hMzRCooI5TSc8IAHos53Gp9aaW92sPA8za+WRP6OSH6UsOW4N+iQc4jyl7/fckMSgIZlJouNqqV+P8iqIFJGs70Tj5L8G/m+P2lc3kcE4Vjmj+Fc0xG5+I/PsSOpcc6DfDfZdVDRe8yklYd/qC1jI89OCeqjxu3XcUGHj9s="

func TestAccGitlabUserSSHKey_basic(t *testing.T) {
	var key gitlab.SSHKey
	testUser := testAccCreateUsers(t, 1)[0]

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabUserSSHKeyDestroy,
		Steps: []resource.TestStep{
			// Create a user + sshkey
			{
				Config: testAccGitlabUserSSHKeyConfig(testUser.ID, testRSAPubKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabUserSSHKeyExists("gitlab_user_sshkey.foo_key", &key),
					testAccCheckGitlabUserSSHKeyAttributes(&key, &testAccGitlabUserSSHKeyExpectedAttributes{
						Title: "foo-key",
						Key:   testRSAPubKey,
					}),
				),
			},
			// Only update key comment (which is a no-op plan)
			{
				Config:   testAccGitlabUserSSHKeyConfig(testUser.ID, testRSAPubKeyUpdatedComment),
				PlanOnly: true,
			},
			// Update the key and title
			{
				Config: testAccGitlabUserSSHKeyUpdateConfig(testUser.ID, updatedRSAPubKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabUserSSHKeyExists("gitlab_user_sshkey.foo_key", &key),
					testAccCheckGitlabUserSSHKeyAttributes(&key, &testAccGitlabUserSSHKeyExpectedAttributes{
						Title:     "key",
						Key:       updatedRSAPubKey,
						ExpiresAt: "3016-01-21T00:00:00Z",
					}),
				),
			},
			{
				ResourceName:      "gitlab_user_sshkey.foo_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Change pub key to one without a comment
			{
				Config: testAccGitlabUserSSHKeyConfig(testUser.ID, updatedRSAPubKeyWithoutComment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabUserSSHKeyExists("gitlab_user_sshkey.foo_key", &key),
					testAccCheckGitlabUserSSHKeyAttributes(&key, &testAccGitlabUserSSHKeyExpectedAttributes{
						Title: "foo-key",
						Key:   updatedRSAPubKeyWithoutComment,
					}),
				),
			},
			{
				ResourceName:      "gitlab_user_sshkey.foo_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGitlabUserSSHKeyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_user_sshkey" {
			continue
		}

		userID, keyID, err := resourceGitlabUserSSHKeyParseID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to parse user ssh key resource ID: %s", err)
		}

		keys, _, err := testGitlabClient.Users.ListSSHKeysForUser(userID, &gitlab.ListSSHKeysForUserOptions{})
		if err != nil {
			return err
		}

		var gotKey *gitlab.SSHKey

		for _, k := range keys {
			if k.ID == keyID {
				gotKey = k
				break
			}
		}
		if gotKey != nil {
			return fmt.Errorf("SSH Key still exists")
		}

		return nil
	}
	return nil
}

func testAccCheckGitlabUserSSHKeyExists(n string, key *gitlab.SSHKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		userID, keyID, err := resourceGitlabUserSSHKeyParseID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to parse user ssh key resource ID: %s", err)
		}

		keys, _, err := testGitlabClient.Users.ListSSHKeysForUser(userID, &gitlab.ListSSHKeysForUserOptions{})
		if err != nil {
			return err
		}

		var gotKey *gitlab.SSHKey

		for _, k := range keys {
			if k.ID == keyID {
				gotKey = k
				break
			}
		}
		if gotKey == nil {
			return fmt.Errorf("Could not find sshkey %d for user %d", keyID, userID)
		}

		*key = *gotKey
		return nil
	}
}

type testAccGitlabUserSSHKeyExpectedAttributes struct {
	Title     string
	Key       string
	CreatedAt string
	ExpiresAt string
}

func testAccCheckGitlabUserSSHKeyAttributes(key *gitlab.SSHKey, want *testAccGitlabUserSSHKeyExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if key.Title != want.Title {
			return fmt.Errorf("got title %q; want %q", key.Title, want.Title)
		}

		k := strings.Join(strings.Split(key.Key, " ")[:2], " ")
		wk := strings.Join(strings.Split(want.Key, " ")[:2], " ")

		if k != wk {
			return fmt.Errorf("got key %q; want %q", k, wk)
		}

		return nil
	}
}

func testAccGitlabUserSSHKeyConfig(userID int, pubKey string) string {
	return fmt.Sprintf(`
resource "gitlab_user_sshkey" "foo_key" {
  title = "foo-key"
  key = "%s"
  user_id = %d
}
  `, pubKey, userID)
}

func testAccGitlabUserSSHKeyUpdateConfig(userID int, pubKey string) string {
	return fmt.Sprintf(`
resource "gitlab_user_sshkey" "foo_key" {
  title = "key"
  key = "%s"
  user_id = %d
  expires_at = "3016-01-21T00:00:00Z"
}
  `, pubKey, userID)
}
