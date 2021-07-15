package gitlab

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabProjectMirror_basic(t *testing.T) {
	ctx := testAccGitlabProjectStart(t)
	defer ctx.finish()

	var mirror gitlab.ProjectMirror

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabProjectMirrorDestroy,
		Steps: []resource.TestStep{
			// Create a project mirror with default options
			{
				Config: fmt.Sprintf(`
resource "gitlab_project_mirror" "foo" {
  project = %q
  url     = "https://example.com/mirror"
}`, ctx.project.PathWithNamespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectMirrorExists("gitlab_project_mirror.foo", &mirror),
					testAccCheckGitlabProjectMirrorAttributes(&mirror, &testAccGitlabProjectMirrorExpectedAttributes{
						URL:                   "https://example.com/mirror",
						Enabled:               true,
						OnlyProtectedBranches: true,
						KeepDivergentRefs:     true,
					}),
				),
			},
			// Mirror with an unprotected URL can be imported
			{
				ResourceName:      "gitlab_project_mirror.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update the project mirror to toggle all the values to their inverse
			{
				Config: fmt.Sprintf(`
resource "gitlab_project_mirror" "foo" {
  project                 = %q
  url                     = "https://example.com/mirror"
  enabled                 = false
  only_protected_branches = false
  keep_divergent_refs     = false
}`, ctx.project.PathWithNamespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectMirrorExists("gitlab_project_mirror.foo", &mirror),
					testAccCheckGitlabProjectMirrorAttributes(&mirror, &testAccGitlabProjectMirrorExpectedAttributes{
						URL:                   "https://example.com/mirror",
						Enabled:               false,
						OnlyProtectedBranches: false,
						KeepDivergentRefs:     false,
					}),
				),
			},
			// Update the project mirror URL to have basicauth
			{
				Config: fmt.Sprintf(`
resource "gitlab_project_mirror" "foo" {
  project = %q
  url     = "https://user:pass@example.com/mirror"
}`, ctx.project.PathWithNamespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectMirrorExists("gitlab_project_mirror.foo", &mirror),
					resource.TestCheckResourceAttr("gitlab_project_mirror.foo", "url", "https://user:pass@example.com/mirror"),
					testAccCheckGitlabProjectMirrorAttributes(&mirror, &testAccGitlabProjectMirrorExpectedAttributes{
						URL:                   "https://*****:*****@example.com/mirror",
						Enabled:               true,
						OnlyProtectedBranches: true,
						KeepDivergentRefs:     true,
					}),
				),
			},
			// Mirror with an authenticated URL cannot be fully imported
			{
				ResourceName:            "gitlab_project_mirror.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"url"},
			},
			// Updating the basicauth URL updates it in state
			{
				Config: fmt.Sprintf(`
resource "gitlab_project_mirror" "foo" {
  project = %q
  url     = "https://user:razz@example.com/mirror"
}`, ctx.project.PathWithNamespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectMirrorExists("gitlab_project_mirror.foo", &mirror),
					resource.TestCheckResourceAttr("gitlab_project_mirror.foo", "url", "https://user:razz@example.com/mirror"),
					testAccCheckGitlabProjectMirrorAttributes(&mirror, &testAccGitlabProjectMirrorExpectedAttributes{
						URL:                   "https://*****:*****@example.com/mirror",
						Enabled:               true,
						OnlyProtectedBranches: true,
						KeepDivergentRefs:     true,
					}),
				),
			},
		},
	})
}

func testAccCheckGitlabProjectMirrorExists(n string, mirror *gitlab.ProjectMirror) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		splitID := strings.Split(rs.Primary.ID, ":")

		mirrorID, err := strconv.Atoi(splitID[len(splitID)-1])
		if err != nil {
			return err
		}
		repoName := rs.Primary.Attributes["project"]
		if repoName == "" {
			return fmt.Errorf("No project ID is set")
		}
		conn := testAccProvider.Meta().(*gitlab.Client)

		mirrors, _, err := conn.ProjectMirrors.ListProjectMirror(repoName, nil)
		if err != nil {
			return err
		}

		for _, m := range mirrors {
			if m.ID == mirrorID {
				*mirror = *m
				return nil
			}
		}
		return errors.New("unable to find mirror")
	}
}

type testAccGitlabProjectMirrorExpectedAttributes struct {
	URL                   string
	Enabled               bool
	OnlyProtectedBranches bool
	KeepDivergentRefs     bool
}

func testAccCheckGitlabProjectMirrorAttributes(mirror *gitlab.ProjectMirror, want *testAccGitlabProjectMirrorExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if mirror.URL != want.URL {
			return fmt.Errorf("got url %q; want %q", mirror.URL, want.URL)
		}

		if mirror.Enabled != want.Enabled {
			return fmt.Errorf("got enabled %t; want %t", mirror.Enabled, want.Enabled)
		}
		if mirror.OnlyProtectedBranches != want.OnlyProtectedBranches {
			return fmt.Errorf("got only protected branches %t; want %t", mirror.OnlyProtectedBranches, want.OnlyProtectedBranches)
		}
		if mirror.KeepDivergentRefs != want.KeepDivergentRefs {
			return fmt.Errorf("got keep divergent refs %t; want %t", mirror.KeepDivergentRefs, want.KeepDivergentRefs)
		}
		return nil
	}
}

func testAccCheckGitlabProjectMirrorDestroy(s *terraform.State) error {
	var mirror gitlab.ProjectMirror
	if err := testAccCheckGitlabProjectMirrorExists("gitlab_project_mirror.foo", &mirror)(s); err != nil {
		return err
	}
	if mirror.Enabled {
		return errors.New("mirror is enabled")
	}
	return nil
}
