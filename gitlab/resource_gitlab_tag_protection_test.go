package gitlab

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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabTagProtectionDestroy,
		Steps: []resource.TestStep{
			// Create a project and Tag Protection with default options
			{
				Config: testAccGitlabTagProtectionConfig(rInt, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabTagProtectionExists("gitlab_tag_protection.TagProtect", &pt),
					testAccCheckGitlabTagProtectionAttributes(&pt, &testAccGitlabTagProtectionExpectedAttributes{
						Name:              fmt.Sprintf("TagProtect-%d", rInt),
						CreateAccessLevel: accessLevel[gitlab.DeveloperPermissions],
					}),
				),
			},
			// Update the Tag Protection
			{
				Config: testAccGitlabTagProtectionUpdateConfig(rInt, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabTagProtectionExists("gitlab_tag_protection.TagProtect", &pt),
					testAccCheckGitlabTagProtectionAttributes(&pt, &testAccGitlabTagProtectionExpectedAttributes{
						Name:              fmt.Sprintf("TagProtect-%d", rInt),
						CreateAccessLevel: accessLevel[gitlab.MasterPermissions],
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
						CreateAccessLevel: accessLevel[gitlab.DeveloperPermissions],
					}),
				),
			},
		},
	})
}

func TestAccGitlabTagProtection_wildcard(t *testing.T) {

	var pt gitlab.ProtectedTag
	rInt := acctest.RandInt()

	wildcard := "-*"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabTagProtectionDestroy,
		Steps: []resource.TestStep{
			// Create a project and Tag Protection with default options
			{
				Config: testAccGitlabTagProtectionConfig(rInt, wildcard),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabTagProtectionExists("gitlab_tag_protection.TagProtect", &pt),
					testAccCheckGitlabTagProtectionAttributes(&pt, &testAccGitlabTagProtectionExpectedAttributes{
						Name:              fmt.Sprintf("TagProtect-%d%s", rInt, wildcard),
						CreateAccessLevel: accessLevel[gitlab.DeveloperPermissions],
					}),
				),
			},
			// Update the Tag Protection
			{
				Config: testAccGitlabTagProtectionUpdateConfig(rInt, wildcard),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabTagProtectionExists("gitlab_tag_protection.TagProtect", &pt),
					testAccCheckGitlabTagProtectionAttributes(&pt, &testAccGitlabTagProtectionExpectedAttributes{
						Name:              fmt.Sprintf("TagProtect-%d%s", rInt, wildcard),
						CreateAccessLevel: accessLevel[gitlab.MasterPermissions],
					}),
				),
			},
			// Update the Tag Protection to get back to initial settings
			{
				Config: testAccGitlabTagProtectionConfig(rInt, wildcard),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabTagProtectionExists("gitlab_tag_protection.TagProtect", &pt),
					testAccCheckGitlabTagProtectionAttributes(&pt, &testAccGitlabTagProtectionExpectedAttributes{
						Name:              fmt.Sprintf("TagProtect-%d%s", rInt, wildcard),
						CreateAccessLevel: accessLevel[gitlab.DeveloperPermissions],
					}),
				),
			},
		},
	})
}

// lintignore: AT002 // TODO: Resolve this tfproviderlint issue
func TestAccGitlabTagProtection_import(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabTagProtectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabTagProtectionConfig(rInt, ""),
			},
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

		conn := testAccProvider.Meta().(*gitlab.Client)

		pts, _, err := conn.ProtectedTags.ListProtectedTags(project, nil)
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

		if pt.CreateAccessLevels[0].AccessLevel != accessLevelID[want.CreateAccessLevel] {
			return fmt.Errorf("got Create access levels %q; want %q", pt.CreateAccessLevels[0].AccessLevel, accessLevelID[want.CreateAccessLevel])
		}

		return nil
	}
}

func testAccCheckGitlabTagProtectionDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*gitlab.Client)
	var project string
	var tag string
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "gitlab_project" {
			project = rs.Primary.ID
		} else if rs.Type == "gitlab_tag_protection" {
			tag = rs.Primary.ID
		}
	}

	pt, _, err := conn.ProtectedTags.GetProtectedTag(project, tag)
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
