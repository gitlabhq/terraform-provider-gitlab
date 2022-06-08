//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabTagProtection_basic(t *testing.T) {
	var pt gitlab.ProtectedTag
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabTagProtectionDestroy,
		Steps: []resource.TestStep{
			// Create a project and Tag Protection with default options
			{
				Config: testAccGitlabTagProtectionConfig(rInt, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabTagProtectionExists("gitlab_tag_protection.TagProtect", &pt),
					testAccCheckGitlabTagProtectionAttributes(&pt, &testAccGitlabTagProtectionExpectedAttributes{
						Name:              fmt.Sprintf("TagProtect-%d", rInt),
						CreateAccessLevel: accessLevelValueToName[gitlab.DeveloperPermissions],
					}),
				),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_tag_protection.TagProtect",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update the Tag Protection
			{
				Config: testAccGitlabTagProtectionUpdateConfig(rInt, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabTagProtectionExists("gitlab_tag_protection.TagProtect", &pt),
					testAccCheckGitlabTagProtectionAttributes(&pt, &testAccGitlabTagProtectionExpectedAttributes{
						Name:              fmt.Sprintf("TagProtect-%d", rInt),
						CreateAccessLevel: accessLevelValueToName[gitlab.MasterPermissions],
					}),
				),
			},
			// Update the Tag Protection to get back to initial settings
			{
				Config: testAccGitlabTagProtectionConfig(rInt, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabTagProtectionExists("gitlab_tag_protection.TagProtect", &pt),
					testAccCheckGitlabTagProtectionAttributes(&pt, &testAccGitlabTagProtectionExpectedAttributes{
						Name:              fmt.Sprintf("TagProtect-%d", rInt),
						CreateAccessLevel: accessLevelValueToName[gitlab.DeveloperPermissions],
					}),
				),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_tag_protection.TagProtect",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGitlabTagProtection_wildcard(t *testing.T) {
	var pt gitlab.ProtectedTag
	rInt := acctest.RandInt()

	wildcard := "-*"

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabTagProtectionDestroy,
		Steps: []resource.TestStep{
			// Create a project and Tag Protection with default options
			{
				Config: testAccGitlabTagProtectionConfig(rInt, wildcard),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabTagProtectionExists("gitlab_tag_protection.TagProtect", &pt),
					testAccCheckGitlabTagProtectionAttributes(&pt, &testAccGitlabTagProtectionExpectedAttributes{
						Name:              fmt.Sprintf("TagProtect-%d%s", rInt, wildcard),
						CreateAccessLevel: accessLevelValueToName[gitlab.DeveloperPermissions],
					}),
				),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_tag_protection.TagProtect",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update the Tag Protection
			{
				Config: testAccGitlabTagProtectionUpdateConfig(rInt, wildcard),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabTagProtectionExists("gitlab_tag_protection.TagProtect", &pt),
					testAccCheckGitlabTagProtectionAttributes(&pt, &testAccGitlabTagProtectionExpectedAttributes{
						Name:              fmt.Sprintf("TagProtect-%d%s", rInt, wildcard),
						CreateAccessLevel: accessLevelValueToName[gitlab.MasterPermissions],
					}),
				),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_tag_protection.TagProtect",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update the Tag Protection to get back to initial settings
			{
				Config: testAccGitlabTagProtectionConfig(rInt, wildcard),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabTagProtectionExists("gitlab_tag_protection.TagProtect", &pt),
					testAccCheckGitlabTagProtectionAttributes(&pt, &testAccGitlabTagProtectionExpectedAttributes{
						Name:              fmt.Sprintf("TagProtect-%d%s", rInt, wildcard),
						CreateAccessLevel: accessLevelValueToName[gitlab.DeveloperPermissions],
					}),
				),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_tag_protection.TagProtect",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGitlabTagProtectionExists(n string, pt *gitlab.ProtectedTag) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}
		project, tag, err := projectAndTagFromID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error in Splitting Project and Tag Ids")
		}

		pts, _, err := testGitlabClient.ProtectedTags.ListProtectedTags(project, nil)
		if err != nil {
			return err
		}
		for _, gotpt := range pts {
			if gotpt.Name == tag {
				*pt = *gotpt
				return nil
			}
		}
		return fmt.Errorf("Protected Tag does not exist")
	}
}

type testAccGitlabTagProtectionExpectedAttributes struct {
	Name              string
	CreateAccessLevel string
}

func testAccCheckGitlabTagProtectionAttributes(pt *gitlab.ProtectedTag, want *testAccGitlabTagProtectionExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if pt.Name != want.Name {
			return fmt.Errorf("got name %q; want %q", pt.Name, want.Name)
		}

		if pt.CreateAccessLevels[0].AccessLevel != accessLevelNameToValue[want.CreateAccessLevel] {
			return fmt.Errorf("got Create access levels %q; want %q", pt.CreateAccessLevels[0].AccessLevel, accessLevelNameToValue[want.CreateAccessLevel])
		}

		return nil
	}
}

func testAccCheckGitlabTagProtectionDestroy(s *terraform.State) error {
	var project string
	var tag string
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "gitlab_project" {
			project = rs.Primary.ID
		} else if rs.Type == "gitlab_tag_protection" {
			tag = rs.Primary.ID
		}
	}

	pt, _, err := testGitlabClient.ProtectedTags.GetProtectedTag(project, tag)
	if err == nil {
		if pt != nil {
			return fmt.Errorf("project tag protection %s still exists", tag)
		}
	}
	if !is404(err) {
		return err
	}
	return nil
}

func testAccGitlabTagProtectionConfig(rInt int, postfix string) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_tag_protection" "TagProtect" {
  project = "${gitlab_project.foo.id}"
  tag = "TagProtect-%d%s"
  create_access_level = "developer"
}
	`, rInt, rInt, postfix)
}

func testAccGitlabTagProtectionUpdateConfig(rInt int, postfix string) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_tag_protection" "TagProtect" {
	project = "${gitlab_project.foo.id}"
	tag = "TagProtect-%d%s"
	create_access_level = "maintainer"
}
	`, rInt, rInt, postfix)
}
