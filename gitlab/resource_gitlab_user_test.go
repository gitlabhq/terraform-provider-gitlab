package gitlab

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabUser_basic(t *testing.T) {
	var user gitlab.User
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabGroupDestroy,
		Steps: []resource.TestStep{
			// Create a user
			{
				Config: testAccGitlabUserConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabUserExists("gitlab_user.foo", &user),
					testAccCheckGitlabUserAttributes(&user, &testAccGitlabUserExpectedAttributes{
						Email:            fmt.Sprintf("listest%d@ssss.com", rInt),
						Password:         fmt.Sprintf("test%dtt", rInt),
						Username:         fmt.Sprintf("listest%d", rInt),
						Name:             fmt.Sprintf("foo %d", rInt),
						ProjectsLimit:    0,
						Admin:            false,
						CanCreateGroup:   false,
						SkipConfirmation: true,
						External:         false,
					}),
				),
			},
			// Update the user to change the name, email, projects_limit and more
			{
				Config: testAccGitlabUserUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabUserExists("gitlab_user.foo", &user),
					testAccCheckGitlabUserAttributes(&user, &testAccGitlabUserExpectedAttributes{
						Email:            fmt.Sprintf("listest%d@tttt.com", rInt),
						Password:         fmt.Sprintf("test%dtt", rInt),
						Username:         fmt.Sprintf("listest%d", rInt),
						Name:             fmt.Sprintf("bar %d", rInt),
						ProjectsLimit:    10,
						Admin:            true,
						CanCreateGroup:   true,
						SkipConfirmation: false,
						External:         false,
					}),
				),
			},
			// Update the user to put the name back
			{
				Config: testAccGitlabUserConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabUserExists("gitlab_user.foo", &user),
					testAccCheckGitlabUserAttributes(&user, &testAccGitlabUserExpectedAttributes{
						Email:            fmt.Sprintf("listest%d@ssss.com", rInt),
						Password:         fmt.Sprintf("test%dtt", rInt),
						Username:         fmt.Sprintf("listest%d", rInt),
						Name:             fmt.Sprintf("foo %d", rInt),
						ProjectsLimit:    0,
						Admin:            false,
						CanCreateGroup:   false,
						SkipConfirmation: false,
						External:         false,
					}),
				),
			},
			{
				ResourceName:      "gitlab_user.foo",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
				},
			},
		},
	})
}

func TestAccGitlabUser_password_reset(t *testing.T) {
	var user gitlab.User
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabGroupDestroy,
		Steps: []resource.TestStep{
			// Test that either password or reset_password is needed
			{
				Config:      testAccGitlabUserConfigWrong(rInt),
				ExpectError: regexp.MustCompile("At least one of either password or reset_password must be defined"),
			},
			// Create a user without a password
			{
				Config: testAccGitlabUserConfigPasswordReset(rInt),
				Check:  testAccCheckGitlabUserExists("gitlab_user.foo", &user),
			},
		},
	})
}

func testAccCheckGitlabUserExists(n string, user *gitlab.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		userID := rs.Primary.ID
		if userID == "" {
			return fmt.Errorf("No user ID is set")
		}
		conn := testAccProvider.Meta().(*gitlab.Client)

		id, _ := strconv.Atoi(userID)

		gotUser, _, err := conn.Users.GetUser(id)
		if err != nil {
			return err
		}
		*user = *gotUser
		return nil
	}
}

type testAccGitlabUserExpectedAttributes struct {
	Email            string
	Password         string
	Username         string
	Name             string
	ProjectsLimit    int
	Admin            bool
	CanCreateGroup   bool
	SkipConfirmation bool
	External         bool
}

func testAccCheckGitlabUserAttributes(user *gitlab.User, want *testAccGitlabUserExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if user.Name != want.Name {
			return fmt.Errorf("got name %q; want %q", user.Name, want.Name)
		}

		if user.Username != want.Username {
			return fmt.Errorf("got username %q; want %q", user.Username, want.Username)
		}

		if user.Email != want.Email {
			return fmt.Errorf("got email %q; want %q", user.Email, want.Email)
		}

		if user.CanCreateGroup != want.CanCreateGroup {
			return fmt.Errorf("got can_create_group %t; want %t", user.CanCreateGroup, want.CanCreateGroup)
		}

		if user.External != want.External {
			return fmt.Errorf("got is_external %t; want %t", user.External, want.External)
		}

		if user.IsAdmin != want.Admin {
			return fmt.Errorf("got is_admin %t; want %t", user.IsAdmin, want.Admin)
		}

		if user.ProjectsLimit != want.ProjectsLimit {
			return fmt.Errorf("got projects_limit %d; want %d", user.ProjectsLimit, want.ProjectsLimit)
		}

		return nil
	}
}

func testAccCheckGitlabUserDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*gitlab.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_user" {
			continue
		}

		id, _ := strconv.Atoi(rs.Primary.ID)

		user, resp, err := conn.Users.GetUser(id)
		if err == nil {
			if user != nil && fmt.Sprintf("%d", user.ID) == rs.Primary.ID {
				return fmt.Errorf("User still exists")
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func testAccGitlabUserConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_user" "foo" {
  name             = "foo %d"
  username         = "listest%d"
  password         = "test%dtt"
  email            = "listest%d@ssss.com"
  is_admin         = false
  projects_limit   = 0
  can_create_group = false
  is_external      = false
}
  `, rInt, rInt, rInt, rInt)
}

func testAccGitlabUserUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_user" "foo" {
  name             = "bar %d"
  username         = "listest%d"
  password         = "test%dtt"
  email            = "listest%d@tttt.com"
  is_admin         = true
  projects_limit   = 10
  can_create_group = true
  is_external      = false
}
  `, rInt, rInt, rInt, rInt)
}

func testAccGitlabUserConfigPasswordReset(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_user" "foo" {
  name             = "foo %d"
  username         = "listest%d"
  email            = "listest%d@ssss.com"
  reset_password   = true
}
  `, rInt, rInt, rInt)
}

func testAccGitlabUserConfigWrong(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_user" "foo" {
  name             = "foo %d"
  username         = "listest%d"
  email            = "listest%d@ssss.com"
}
  `, rInt, rInt, rInt)
}
