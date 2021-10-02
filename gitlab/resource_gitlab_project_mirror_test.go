package gitlab

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabProjectMirror_basic(t *testing.T) {
	ctx := testAccGitlabProjectStart(t)
	var miror gitlab.ProjectMirror

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabProjectMirrorDestroy,
		Steps: []resource.TestStep{
			// Create with default options
			{
				Config: testAccGitlabProjectMirrorConfig(ctx.project.PathWithNamespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectMirrorExists("gitlab_project_mirror.foo", &miror),
					testAccCheckGitlabProjectMirrorAttributes(&miror, &testAccGitlabProjectMirrorExpectedAttributes{
						URL:                   "https://example.com/mirror",
						Enabled:               true,
						OnlyProtectedBranches: true,
						KeepDivergentRefs:     true,
					}),
				),
			},
			// Update to toggle all the values to their inverse
			{
				Config: testAccGitlabProjectMirrorUpdateConfig(ctx.project.PathWithNamespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectMirrorExists("gitlab_project_mirror.foo", &miror),
					testAccCheckGitlabProjectMirrorAttributes(&miror, &testAccGitlabProjectMirrorExpectedAttributes{
						URL:                   "https://example.com/mirror",
						Enabled:               false,
						OnlyProtectedBranches: false,
						KeepDivergentRefs:     false,
					}),
				),
			},
			// Update to toggle the options back
			{
				Config: testAccGitlabProjectMirrorConfig(ctx.project.PathWithNamespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectMirrorExists("gitlab_project_mirror.foo", &miror),
					testAccCheckGitlabProjectMirrorAttributes(&miror, &testAccGitlabProjectMirrorExpectedAttributes{
						URL:                   "https://example.com/mirror",
						Enabled:               true,
						OnlyProtectedBranches: true,
						KeepDivergentRefs:     true,
					}),
				),
			},
			// Import
			{
				ResourceName:      "gitlab_project_mirror.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGitlabProjectMirror_withPassword(t *testing.T) {
	//var mirror gitlab.ProjectMirror
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabProjectMirrorDestroy,
		Steps: []resource.TestStep{
			// Create a project and mirror with a username / password.
			{
				Config: testAccGitlabProjectMirrorConfigWithPassword(rInt),
			},
		},
	})
}

func TestAccGitlabProjectMirror_withCount(t *testing.T) {
	//var mirror gitlab.ProjectMirror
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabProjectMirrorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabProjectMirrorConfigWithCount(rInt),
			},
		},
	})
}

// lintignore: AT002 // TODO: Resolve this tfproviderlint issue
func TestAccGitlabProjectMirror_import(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabProjectMirrorDestroy,
		Steps: []resource.TestStep{
			// Create a project and mirror with a username / password.
			{
				Config: testAccGitlabProjectMirrorConfigWithPassword(ctx.project.PathWithNamespace),
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

func testAccGitlabProjectMirrorConfig(project string) string {
	return fmt.Sprintf(`
resource "gitlab_project_mirror" "foo" {
  project = %q
  url = "https://example.com/mirror"
}
	`, project)
}

func testAccGitlabProjectMirrorConfigWithCount(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_project_mirror" "foo" {
  project = "${gitlab_project.foo.id}"
  url = "https://foo:%d@example.com/mirror-%d"
  count = 40
}
	`, rInt, rInt, rInt)
}

func testAccGitlabProjectMirrorConfigWithPassword(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_project_mirror" "foo" {
  project = "${gitlab_project.foo.id}"
  url = "https://foo:%d@example.com/mirror-%d"
}
	`, rInt, rInt, rInt)
}

func testAccGitlabProjectMirrorUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project_mirror" "foo" {
  project = %q
  url = "https://foo:bar@example.com/mirror"
}
	`, project)
}

func testAccGitlabProjectMirrorUpdateConfig(project string) string {
	return fmt.Sprintf(`
resource "gitlab_project_mirror" "foo" {
  project = %q
  url = "https://example.com/mirror"
  enabled = false
  only_protected_branches = false
  keep_divergent_refs = false
}
	`, project)
}
