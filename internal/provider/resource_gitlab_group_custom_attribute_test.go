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

func TestAccGitlabGroupCustomAttribute_basic(t *testing.T) {
	var group gitlab.Group
	var customAttribute gitlab.CustomAttribute
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "gitlab_group" "group" {
  	name = "foo-name-%d"
    path = "foo-path-%d"
}

resource "gitlab_group_custom_attribute" "attr" {
	group = gitlab_group.group.id
	key   = "foo"
	value = "bar"
}`, rInt, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupExists("gitlab_group.group", &group),
					testAccCheckGitlabGroupCustomAttributeExists("gitlab_group_custom_attribute.attr", &customAttribute),
					testAccCheckGitlabGroupCustomAttributes(&customAttribute, &testAccGitlabGroupExpectedCustomAttributes{
						Key:   "foo",
						Value: "bar",
					}),
				),
			},
			// Update the custom attribute
			{
				Config: fmt.Sprintf(`
resource "gitlab_group" "group" {
	name = "foo-name-%d"
	path = "foo-path-%d"
}

resource "gitlab_group_custom_attribute" "attr" {
	group = gitlab_group.group.id
	key   = "foo"
	value = "updated"
}`, rInt, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupExists("gitlab_group.group", &group),
					testAccCheckGitlabGroupCustomAttributeExists("gitlab_group_custom_attribute.attr", &customAttribute),
					testAccCheckGitlabGroupCustomAttributes(&customAttribute, &testAccGitlabGroupExpectedCustomAttributes{
						Key:   "foo",
						Value: "updated",
					}),
				),
			},
			{
				ResourceName:      "gitlab_group_custom_attribute.attr",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGitlabGroupCustomAttributeExists(n string, customAttribute *gitlab.CustomAttribute) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		id, key, err := parseId(rs.Primary.ID)
		if err != nil {
			return err
		}

		gotCustomAttribute, _, err := testGitlabClient.CustomAttribute.GetCustomGroupAttribute(id, key)
		if err != nil {
			return err
		}
		*customAttribute = *gotCustomAttribute
		return nil
	}
}

type testAccGitlabGroupExpectedCustomAttributes struct {
	Key   string
	Value string
}

func testAccCheckGitlabGroupCustomAttributes(got *gitlab.CustomAttribute, want *testAccGitlabGroupExpectedCustomAttributes) resource.TestCheckFunc {
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
