package gitlab

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabApplication_basic(t *testing.T) {
	var application gitlab.Application
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabApplicationConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabApplicationExists("gitlab_application.foo", &application),
					testAccCheckGitlabApplicationAttributes(&application, &testAccCheckGitlabApplicationExpectedAttributes{
						Name:        fmt.Sprintf("application-%d", rInt),
						RedirectURI: "http://redirect.uri",
					}),
				),
			},
		},
	})
}

func testAccCheckGitlabApplicationExists(n string, application *gitlab.Application) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		applicationID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		conn := testAccProvider.Meta().(*gitlab.Client)

		var gotApplications []*gitlab.Application

		gotApplications, _, err = conn.Applications.ListApplications(nil)

		if err != nil {
			return err
		}

		for _, app := range gotApplications {
			if app.ID == applicationID {
				*application = *app
				return nil
			}
		}

		return fmt.Errorf("Application doesn't exist")
	}
}

type testAccCheckGitlabApplicationExpectedAttributes struct {
	Name        string
	RedirectURI string
}

func testAccCheckGitlabApplicationAttributes(application *gitlab.Application, want *testAccCheckGitlabApplicationExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if application.ApplicationName != want.Name {
			return fmt.Errorf("got name %q; want %q", application.ApplicationName, want.Name)
		}

		if application.CallbackURL != want.RedirectURI {
			return fmt.Errorf("got username %q; want %q", application.CallbackURL, want.RedirectURI)
		}

		return nil
	}
}

func testAccCheckGitlabApplicationDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*gitlab.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_application" {
			continue
		}

		applicationID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		var gotApplications []*gitlab.Application
		var resp *gitlab.Response

		gotApplications, resp, err = conn.Applications.ListApplications(nil)

		if err == nil {
			for _, app := range gotApplications {
				if app.ID == applicationID {
					return fmt.Errorf("Application still exists")
				}
			}
		}

		if resp.StatusCode != http.StatusNotFound {
			return err
		}
	}

	return nil
}

func testAccGitlabApplicationConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_application" "foo" {
  name     = "application-%d"

  redirect_uri = "http://redirect.uri"
}
  `, rInt)
}
