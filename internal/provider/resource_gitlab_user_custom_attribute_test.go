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

func TestAccGitlabUserCustomAttribute_basic(t *testing.T) {
	var user gitlab.User
	var customAttribute gitlab.CustomAttribute
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "gitlab_user" "user" {
  name        = "foo-%d"
  username    = "foo-%d"
  password    = "foofoofoo"
  email       = "foo@email.com"
}

resource "gitlab_user_custom_attribute" "attr" {
	user  = gitlab_user.user.id
	key   = "foo"
	value = "bar"
}`, rInt, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabUserExists("gitlab_user.user", &user),
					testAccCheckGitlabUserCustomAttributeExists("gitlab_user_custom_attribute.attr", &customAttribute),
					testAccCheckGitlabUserCustomAttributes(&customAttribute, &testAccGitlabUserExpectedCustomAttributes{
						Key:   "foo",
						Value: "bar",
					}),
				),
			},
			// Update the custom attribute
			{
				Config: fmt.Sprintf(`
resource "gitlab_user" "user" {
  name        = "foo-%d"
  username    = "foo-%d"
  password    = "foofoofoo"
  email       = "foo@email.com"
}

resource "gitlab_user_custom_attribute" "attr" {
	user  = gitlab_user.user.id
	key   = "foo"
	value = "updated"
}`, rInt, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabUserExists("gitlab_user.user", &user),
					testAccCheckGitlabUserCustomAttributeExists("gitlab_user_custom_attribute.attr", &customAttribute),
					testAccCheckGitlabUserCustomAttributes(&customAttribute, &testAccGitlabUserExpectedCustomAttributes{
						Key:   "foo",
						Value: "updated",
					}),
				),
			},
			{
				ResourceName:      "gitlab_user_custom_attribute.attr",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGitlabUserCustomAttributeExists(n string, customAttribute *gitlab.CustomAttribute) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		id, key, err := parseId(rs.Primary.ID)
		if err != nil {
			return err
		}

		gotCustomAttribute, _, err := testGitlabClient.CustomAttribute.GetCustomUserAttribute(id, key)
		if err != nil {
			return err
		}
		*customAttribute = *gotCustomAttribute
		return nil
	}
}

type testAccGitlabUserExpectedCustomAttributes struct {
	Key   string
	Value string
}

func testAccCheckGitlabUserCustomAttributes(got *gitlab.CustomAttribute, want *testAccGitlabUserExpectedCustomAttributes) resource.TestCheckFunc {
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
