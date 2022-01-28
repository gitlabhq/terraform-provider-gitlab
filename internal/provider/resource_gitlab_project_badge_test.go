package provider

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabProjectBadge_basic(t *testing.T) {
	var badge gitlab.ProjectBadge
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectBadgeDestroy,
		Steps: []resource.TestStep{
			// Create a project and badge
			{
				Config: testAccGitlabProjectBadgeConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectBadgeExists("gitlab_project_badge.foo", &badge),
					testAccCheckGitlabProjectBadgeAttributes(&badge, &testAccGitlabProjectBadgeExpectedAttributes{
						LinkURL:  fmt.Sprintf("https://example.com/badge-%d", rInt),
						ImageURL: fmt.Sprintf("https://example.com/badge-%d.svg", rInt),
					}),
				),
			},
			// Test ImportState
			{
				ResourceName:      "gitlab_project_badge.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update the project badge
			{
				Config: testAccGitlabProjectBadgeUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectBadgeExists("gitlab_project_badge.foo", &badge),
					testAccCheckGitlabProjectBadgeAttributes(&badge, &testAccGitlabProjectBadgeExpectedAttributes{
						LinkURL:  fmt.Sprintf("https://example.com/badge-%d", rInt),
						ImageURL: fmt.Sprintf("https://example.com/badge-%d.svg", rInt),
					}),
				),
			},
		},
	})
}

func testAccCheckGitlabProjectBadgeExists(n string, badge *gitlab.ProjectBadge) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		splitID := strings.Split(rs.Primary.ID, ":")

		badgeID, err := strconv.Atoi(splitID[len(splitID)-1])
		if err != nil {
			return err
		}
		repoName := rs.Primary.Attributes["project"]
		if repoName == "" {
			return fmt.Errorf("No project ID is set")
		}

		gotBadge, _, err := testGitlabClient.ProjectBadges.GetProjectBadge(repoName, badgeID)
		if err != nil {
			return err
		}
		*badge = *gotBadge
		return nil
	}
}

type testAccGitlabProjectBadgeExpectedAttributes struct {
	LinkURL  string
	ImageURL string
}

func testAccCheckGitlabProjectBadgeAttributes(badge *gitlab.ProjectBadge, want *testAccGitlabProjectBadgeExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if badge.LinkURL != want.LinkURL {
			return fmt.Errorf("got link_url %q; want %q", badge.LinkURL, want.LinkURL)
		}

		if badge.ImageURL != want.ImageURL {
			return fmt.Errorf("got image_url %s; want %s", badge.ImageURL, want.ImageURL)
		}

		return nil
	}
}

func testAccCheckGitlabProjectBadgeDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_project" {
			continue
		}

		gotRepo, resp, err := testGitlabClient.Projects.GetProject(rs.Primary.ID, nil)
		if err == nil {
			if gotRepo != nil && fmt.Sprintf("%d", gotRepo.ID) == rs.Primary.ID {
				if gotRepo.MarkedForDeletionAt == nil {
					return fmt.Errorf("Repository still exists")
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

func testAccGitlabProjectBadgeConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_project_badge" "foo" {
  project   = "${gitlab_project.foo.id}"
  link_url  = "https://example.com/badge-%d"
  image_url = "https://example.com/badge-%d.svg"
}
	`, rInt, rInt, rInt)
}

func testAccGitlabProjectBadgeUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_project_badge" "foo" {
  project   = "${gitlab_project.foo.id}"
  link_url  = "https://example.com/badge-%d"
  image_url = "https://example.com/badge-%d.svg"
}
	`, rInt, rInt, rInt)
}
