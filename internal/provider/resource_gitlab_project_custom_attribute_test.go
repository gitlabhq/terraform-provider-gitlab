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

func TestAccGitlabProjectCustomAttribute_basic(t *testing.T) {
	var project gitlab.Project
	var customAttribute gitlab.CustomAttribute
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "gitlab_project" "project" {
	name = "foo-%d"
}

resource "gitlab_project_custom_attribute" "attr" {
	project = gitlab_project.project.id
	key     = "foo"
	value   = "bar"
}`, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.project", &project),
					testAccCheckGitlabProjectCustomAttributeExists("gitlab_project_custom_attribute.attr", &customAttribute),
					testAccCheckGitlabProjectCustomAttributes(&customAttribute, &testAccGitlabProjectExpectedCustomAttributes{
						Key:   "foo",
						Value: "bar",
					}),
				),
			},
			// Update the custom attribute
			{
				Config: fmt.Sprintf(`
resource "gitlab_project" "project" {
	name = "foo-%d"
}

resource "gitlab_project_custom_attribute" "attr" {
	project = gitlab_project.project.id
	key     = "foo"
	value   = "updated"
}`, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.project", &project),
					testAccCheckGitlabProjectCustomAttributeExists("gitlab_project_custom_attribute.attr", &customAttribute),
					testAccCheckGitlabProjectCustomAttributes(&customAttribute, &testAccGitlabProjectExpectedCustomAttributes{
						Key:   "foo",
						Value: "updated",
					}),
				),
			},
			{
				ResourceName:      "gitlab_project_custom_attribute.attr",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGitlabProjectCustomAttributeExists(n string, customAttribute *gitlab.CustomAttribute) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		id, key, err := parseId(rs.Primary.ID)
		if err != nil {
			return err
		}

		gotCustomAttribute, _, err := testGitlabClient.CustomAttribute.GetCustomProjectAttribute(id, key)
		if err != nil {
			return err
		}
		*customAttribute = *gotCustomAttribute
		return nil
	}
}

type testAccGitlabProjectExpectedCustomAttributes struct {
	Key   string
	Value string
}

func testAccCheckGitlabProjectCustomAttributes(got *gitlab.CustomAttribute, want *testAccGitlabProjectExpectedCustomAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if got.Key != want.Key {
			return fmt.Errorf("got key %q; want %q", got.Key, want.Key)
		}

		if got.Value != want.Value {
			return fmt.Errorf("got value %q; want %q", got.Value, want.Value)
		}

		return nil
	}
}
