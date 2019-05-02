package gitlab

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

func TestAccGitlabProject_basic(t *testing.T) {
	var project gitlab.Project
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			// Create a project with all the features on
			{
				Config: testAccGitlabProjectConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.foo", &project),
					testAccCheckGitlabProjectAttributes(&project, &testAccGitlabProjectExpectedAttributes{
						Name:                             fmt.Sprintf("foo-%d", rInt),
						Path:                             fmt.Sprintf("foo.%d", rInt),
						Description:                      "Terraform acceptance tests",
						IssuesEnabled:                    true,
						MergeRequestsEnabled:             true,
						ApprovalsBeforeMerge:             0,
						WikiEnabled:                      true,
						SnippetsEnabled:                  true,
						Visibility:                       gitlab.PublicVisibility,
						MergeMethod:                      gitlab.FastForwardMerge,
						OnlyAllowMergeIfPipelineSucceeds: true,
						OnlyAllowMergeIfAllDiscussionsAreResolved: true,
					}),
				),
			},
			// Update the project to turn the features off
			{
				Config: testAccGitlabProjectUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.foo", &project),
					testAccCheckGitlabProjectAttributes(&project, &testAccGitlabProjectExpectedAttributes{
						Name:                             fmt.Sprintf("foo-%d", rInt),
						Path:                             fmt.Sprintf("foo.%d", rInt),
						Description:                      "Terraform acceptance tests!",
						ApprovalsBeforeMerge:             0,
						Visibility:                       gitlab.PublicVisibility,
						MergeMethod:                      gitlab.FastForwardMerge,
						OnlyAllowMergeIfPipelineSucceeds: true,
						OnlyAllowMergeIfAllDiscussionsAreResolved: true,
					}),
				),
			},
			// Update the project to turn the features on again
			{
				Config: testAccGitlabProjectConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.foo", &project),
					testAccCheckGitlabProjectAttributes(&project, &testAccGitlabProjectExpectedAttributes{
						Name:                             fmt.Sprintf("foo-%d", rInt),
						Path:                             fmt.Sprintf("foo.%d", rInt),
						Description:                      "Terraform acceptance tests",
						IssuesEnabled:                    true,
						MergeRequestsEnabled:             true,
						ApprovalsBeforeMerge:             0,
						WikiEnabled:                      true,
						SnippetsEnabled:                  true,
						Visibility:                       gitlab.PublicVisibility,
						MergeMethod:                      gitlab.FastForwardMerge,
						OnlyAllowMergeIfPipelineSucceeds: true,
						OnlyAllowMergeIfAllDiscussionsAreResolved: true,
					}),
				),
			},
			//Update the project to share the project with a group
			{
				Config: testAccGitlabProjectSharedWithGroup(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.foo", &project),
					testAccCheckGitlabProjectAttributes(&project, &testAccGitlabProjectExpectedAttributes{
						Name:                             fmt.Sprintf("foo-%d", rInt),
						Path:                             fmt.Sprintf("foo.%d", rInt),
						Description:                      "Terraform acceptance tests",
						IssuesEnabled:                    true,
						MergeRequestsEnabled:             true,
						WikiEnabled:                      true,
						SnippetsEnabled:                  true,
						Visibility:                       gitlab.PublicVisibility,
						MergeMethod:                      gitlab.FastForwardMerge,
						OnlyAllowMergeIfPipelineSucceeds: false,
						OnlyAllowMergeIfAllDiscussionsAreResolved: false,
						SharedWithGroups: []struct {
							GroupID          int
							GroupName        string
							GroupAccessLevel int
						}{{0, fmt.Sprintf("foo-name-%d", rInt), 30}},
					}),
				),
			},
			//Update the project to share the project with more groups
			{
				Config: testAccGitlabProjectSharedWithGroup2(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.foo", &project),
					testAccCheckGitlabProjectAttributes(&project, &testAccGitlabProjectExpectedAttributes{
						Name:                             fmt.Sprintf("foo-%d", rInt),
						Path:                             fmt.Sprintf("foo.%d", rInt),
						Description:                      "Terraform acceptance tests",
						IssuesEnabled:                    true,
						MergeRequestsEnabled:             true,
						WikiEnabled:                      true,
						SnippetsEnabled:                  true,
						Visibility:                       gitlab.PublicVisibility,
						MergeMethod:                      gitlab.FastForwardMerge,
						OnlyAllowMergeIfPipelineSucceeds: false,
						OnlyAllowMergeIfAllDiscussionsAreResolved: false,
						SharedWithGroups: []struct {
							GroupID          int
							GroupName        string
							GroupAccessLevel int
						}{{0, fmt.Sprintf("foo-name-%d", rInt), 10}, {0, fmt.Sprintf("foo2-name-%d", rInt), 30}},
					}),
				),
			},
			//Update the project to unshare the project
			{
				Config: testAccGitlabProjectConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.foo", &project),
					testAccCheckGitlabProjectAttributes(&project, &testAccGitlabProjectExpectedAttributes{
						Name:                             fmt.Sprintf("foo-%d", rInt),
						Path:                             fmt.Sprintf("foo.%d", rInt),
						Description:                      "Terraform acceptance tests",
						IssuesEnabled:                    true,
						MergeRequestsEnabled:             true,
						WikiEnabled:                      true,
						SnippetsEnabled:                  true,
						Visibility:                       gitlab.PublicVisibility,
						MergeMethod:                      gitlab.FastForwardMerge,
						OnlyAllowMergeIfPipelineSucceeds: true,
						OnlyAllowMergeIfAllDiscussionsAreResolved: true,
						SharedWithGroups: []struct {
							GroupID          int
							GroupName        string
							GroupAccessLevel int
						}{},
					}),
				),
			},
		},
	})
}

func TestAccGitlabProject_import(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabProjectConfig(rInt),
			},
			{
				ResourceName:      "gitlab_project.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGitlabProject_nestedImport(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabProjectInGroupConfig(rInt),
			},
			{
				ResourceName:      "gitlab_project.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGitlabProjectExists(n string, project *gitlab.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		repoName := rs.Primary.ID
		if repoName == "" {
			return fmt.Errorf("No project ID is set")
		}
		conn := testAccProvider.Meta().(*gitlab.Client)

		gotProject, _, err := conn.Projects.GetProject(repoName, nil)
		if err != nil {
			return err
		}
		*project = *gotProject
		return nil
	}
}

type testAccGitlabProjectExpectedAttributes struct {
	Name                                      string
	Path                                      string
	Description                               string
	DefaultBranch                             string
	IssuesEnabled                             bool
	MergeRequestsEnabled                      bool
	ApprovalsBeforeMerge                      int
	WikiEnabled                               bool
	SnippetsEnabled                           bool
	Visibility                                gitlab.VisibilityValue
	MergeMethod                               gitlab.MergeMethodValue
	OnlyAllowMergeIfPipelineSucceeds          bool
	OnlyAllowMergeIfAllDiscussionsAreResolved bool
	SharedWithGroups                          []struct {
		GroupID          int
		GroupName        string
		GroupAccessLevel int
	}
}

func testAccCheckGitlabProjectAttributes(project *gitlab.Project, want *testAccGitlabProjectExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if project.Name != want.Name {
			return fmt.Errorf("got repo %q; want %q", project.Name, want.Name)
		}
		if project.Path != want.Path {
			return fmt.Errorf("got repo %q; want %q", project.Path, want.Path)
		}
		if project.Description != want.Description {
			return fmt.Errorf("got description %q; want %q", project.Description, want.Description)
		}

		if project.DefaultBranch != want.DefaultBranch {
			return fmt.Errorf("got default_branch %q; want %q", project.DefaultBranch, want.DefaultBranch)
		}

		if project.IssuesEnabled != want.IssuesEnabled {
			return fmt.Errorf("got issues_enabled %t; want %t", project.IssuesEnabled, want.IssuesEnabled)
		}

		if project.MergeRequestsEnabled != want.MergeRequestsEnabled {
			return fmt.Errorf("got merge_requests_enabled %t; want %t", project.MergeRequestsEnabled, want.MergeRequestsEnabled)
		}

		if project.ApprovalsBeforeMerge != want.ApprovalsBeforeMerge {
			return fmt.Errorf("got approvals_before_merge %d; want %d", project.ApprovalsBeforeMerge, want.ApprovalsBeforeMerge)
		}

		if project.WikiEnabled != want.WikiEnabled {
			return fmt.Errorf("got wiki_enabled %t; want %t", project.WikiEnabled, want.WikiEnabled)
		}

		if project.SnippetsEnabled != want.SnippetsEnabled {
			return fmt.Errorf("got snippets_enabled %t; want %t", project.SnippetsEnabled, want.SnippetsEnabled)
		}

		if project.Visibility != want.Visibility {
			return fmt.Errorf("got visibility %q; want %q", project.Visibility, want.Visibility)
		}

		groupsToCheck := want.SharedWithGroups
		for _, group := range project.SharedWithGroups {
			for i, groupToCheck := range groupsToCheck {
				if group.GroupName == groupToCheck.GroupName && group.GroupAccessLevel == groupToCheck.GroupAccessLevel {
					groupsToCheck = append(groupsToCheck[:i], groupsToCheck[i+1:]...)
					break
				}
			}
		}
		if len(groupsToCheck) != 0 {
			return fmt.Errorf("got shared with groups: %v; want %v", project.SharedWithGroups, want.SharedWithGroups)
		}

		if project.MergeMethod != want.MergeMethod {
			return fmt.Errorf("got merge_method %q; want %q", project.MergeMethod, want.MergeMethod)
		}

		if project.OnlyAllowMergeIfPipelineSucceeds != want.OnlyAllowMergeIfPipelineSucceeds {
			return fmt.Errorf("got only_allow_merge_if_pipeline_succeeds %t; want %t", project.OnlyAllowMergeIfPipelineSucceeds, want.OnlyAllowMergeIfPipelineSucceeds)
		}

		if project.OnlyAllowMergeIfAllDiscussionsAreResolved != want.OnlyAllowMergeIfAllDiscussionsAreResolved {
			return fmt.Errorf("got only_allow_merge_if_all_discussions_are_resolved %t; want %t", project.OnlyAllowMergeIfAllDiscussionsAreResolved, want.OnlyAllowMergeIfAllDiscussionsAreResolved)
		}

		return nil
	}
}

func testAccCheckGitlabProjectDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*gitlab.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_project" {
			continue
		}

		gotRepo, resp, err := conn.Projects.GetProject(rs.Primary.ID, nil)
		if err == nil {
			if gotRepo != nil && fmt.Sprintf("%d", gotRepo.ID) == rs.Primary.ID {
				return fmt.Errorf("Repository still exists")
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func testAccGitlabProjectInGroupConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
  name = "foogroup-%d"
  path = "foogroup-%d"
  visibility_level = "public"
}

resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"
  namespace_id = "${gitlab_group.foo.id}"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
	`, rInt, rInt, rInt)
}

func testAccGitlabProjectConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  path = "foo.%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
  merge_method = "ff"
  only_allow_merge_if_pipeline_succeeds = true
  only_allow_merge_if_all_discussions_are_resolved = true
}
	`, rInt, rInt)
}

func testAccGitlabProjectUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  path = "foo.%d"
  description = "Terraform acceptance tests!"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
  merge_method = "ff"
  only_allow_merge_if_pipeline_succeeds = true
  only_allow_merge_if_all_discussions_are_resolved = true

  issues_enabled = false
  merge_requests_enabled = false
  approvals_before_merge = 0
  wiki_enabled = false
  snippets_enabled = false
}
	`, rInt, rInt)
}

func testAccGitlabProjectSharedWithGroup(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name             = "foo-%d"
  path             = "foo.%d"
  description      = "Terraform acceptance tests"
  visibility_level = "public"
  merge_method = "ff"
  only_allow_merge_if_pipeline_succeeds = false
  only_allow_merge_if_all_discussions_are_resolved = false

  shared_with_groups = [
    {
      group_id           = "${gitlab_group.foo.id}"
      group_access_level = "developer"
    },
  ]
}

resource "gitlab_group" "foo" {
  name             = "foo-name-%d"
  path             = "foo-path-%d"
  description      = "Terraform acceptance tests!"
  visibility_level = "public"
}
	`, rInt, rInt, rInt, rInt)
}

func testAccGitlabProjectSharedWithGroup2(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name             = "foo-%d"
  path             = "foo.%d"
  description      = "Terraform acceptance tests"
  visibility_level = "public"
  merge_method = "ff"
  only_allow_merge_if_pipeline_succeeds = false
  only_allow_merge_if_all_discussions_are_resolved = false

  shared_with_groups = [
    {
      group_id           = "${gitlab_group.foo.id}"
      group_access_level = "guest"
    },
    {
      group_id           = "${gitlab_group.foo2.id}"
      group_access_level = "developer"
    },
  ]
}

resource "gitlab_group" "foo" {
  name             = "foo-name-%d"
  path             = "foo-path-%d"
  description      = "Terraform acceptance tests!"
  visibility_level = "public"
}

resource "gitlab_group" "foo2" {
  name             = "foo2-name-%d"
  path             = "foo2-path-%d"
  description      = "Terraform acceptance tests!"
  visibility_level = "public"
}
	`, rInt, rInt, rInt, rInt, rInt, rInt)
}
