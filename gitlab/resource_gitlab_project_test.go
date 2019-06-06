package gitlab

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

func TestAccGitlabProject_basic(t *testing.T) {
	var received, defaults gitlab.Project
	rInt := acctest.RandInt()

	defaults = gitlab.Project{
		Namespace:                        &gitlab.ProjectNamespace{ID: 0},
		Name:                             fmt.Sprintf("foo-%d", rInt),
		Path:                             fmt.Sprintf("foo.%d", rInt),
		Description:                      "Terraform acceptance tests",
		TagList:                          []string{"tag1"},
		IssuesEnabled:                    true,
		MergeRequestsEnabled:             true,
		ApprovalsBeforeMerge:             0,
		WikiEnabled:                      true,
		SnippetsEnabled:                  true,
		ContainerRegistryEnabled:         true,
		SharedRunnersEnabled:             true,
		Visibility:                       gitlab.PublicVisibility,
		MergeMethod:                      gitlab.FastForwardMerge,
		OnlyAllowMergeIfPipelineSucceeds: true,
		OnlyAllowMergeIfAllDiscussionsAreResolved: true,
		Archived: false, // needless, but let's make this explicit
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			// Step0 Create a project with all the features on (note: "archived" is "false")
			{
				Config: testAccGitlabProjectConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.foo", &received),
					testAccCheckAggregateGitlabProject(&defaults, &received),
				),
			},
			// Step1 Update the project to turn the features off (note: "archived" is "true")
			{
				Config: testAccGitlabProjectUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.foo", &received),
					testAccCheckAggregateGitlabProject(&gitlab.Project{
						Namespace:                        &gitlab.ProjectNamespace{ID: 0},
						Name:                             fmt.Sprintf("foo-%d", rInt),
						Path:                             fmt.Sprintf("foo.%d", rInt),
						Description:                      "Terraform acceptance tests!",
						TagList:                          []string{"tag1", "tag2"},
						ApprovalsBeforeMerge:             0,
						ContainerRegistryEnabled:         false,
						SharedRunnersEnabled:             false,
						Visibility:                       gitlab.PublicVisibility,
						MergeMethod:                      gitlab.FastForwardMerge,
						OnlyAllowMergeIfPipelineSucceeds: true,
						OnlyAllowMergeIfAllDiscussionsAreResolved: true,
						Archived: true,
					}, &received),
				),
			},
			// Step2 Update the project to turn the features on again (note: "archived" is "false")
			{
				Config: testAccGitlabProjectConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.foo", &received),
					testAccCheckAggregateGitlabProject(&defaults, &received),
				),
			},
			// Step3 Update the project to share the project with a group
			{
				Config: testAccGitlabProjectSharedWithGroup(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.foo", &received),
					testAccCheckAggregateGitlabProject(
						&gitlab.Project{
							Namespace:                        &gitlab.ProjectNamespace{ID: 0},
							Name:                             fmt.Sprintf("foo-%d", rInt),
							Path:                             fmt.Sprintf("foo.%d", rInt),
							Description:                      "Terraform acceptance tests",
							IssuesEnabled:                    true,
							MergeRequestsEnabled:             true,
							WikiEnabled:                      true,
							SnippetsEnabled:                  true,
							ContainerRegistryEnabled:         true,
							SharedRunnersEnabled:             false,
							Visibility:                       gitlab.PublicVisibility,
							MergeMethod:                      gitlab.FastForwardMerge,
							OnlyAllowMergeIfPipelineSucceeds: false,
							OnlyAllowMergeIfAllDiscussionsAreResolved: false,
							TagList:              []string{},
							ApprovalsBeforeMerge: 0,
							SharedWithGroups: []struct {
								GroupID          int    `json:"group_id"`
								GroupName        string `json:"group_name"`
								GroupAccessLevel int    `json:"group_access_level"`
							}{
								{0, fmt.Sprintf("foo-name-%d", rInt), 30},
							},
						},
						&received),
				),
			},
			// Step4 Update the project to share the project with more groups
			{
				Config: testAccGitlabProjectSharedWithGroup2(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.foo", &received),
					testAccCheckAggregateGitlabProject(&gitlab.Project{
						Namespace:                        &gitlab.ProjectNamespace{ID: 0},
						Name:                             fmt.Sprintf("foo-%d", rInt),
						Path:                             fmt.Sprintf("foo.%d", rInt),
						Description:                      "Terraform acceptance tests",
						IssuesEnabled:                    true,
						MergeRequestsEnabled:             true,
						WikiEnabled:                      true,
						SnippetsEnabled:                  true,
						ContainerRegistryEnabled:         true,
						Visibility:                       gitlab.PublicVisibility,
						MergeMethod:                      gitlab.FastForwardMerge,
						OnlyAllowMergeIfPipelineSucceeds: false,
						OnlyAllowMergeIfAllDiscussionsAreResolved: false,
						SharedWithGroups: []struct {
							GroupID          int    `json:"group_id"`
							GroupName        string `json:"group_name"`
							GroupAccessLevel int    `json:"group_access_level"`
						}{
							{0, fmt.Sprintf("foo-name-%d", rInt), 10},
							{0, fmt.Sprintf("foo2-name-%d", rInt), 30},
						},
					}, &received),
				),
			},
			// Step5 Update the project to unshare the project
			{
				Config: testAccGitlabProjectConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.foo", &received),
					testAccCheckAggregateGitlabProject(&defaults, &received),
				),
			},
		},
	})
}

func TestAccGitlabProject_willError(t *testing.T) {
	var received, defaults gitlab.Project
	rInt := acctest.RandInt()
	defaults = gitlab.Project{
		Namespace:                        &gitlab.ProjectNamespace{ID: 0},
		Name:                             fmt.Sprintf("foo-%d", rInt),
		Path:                             fmt.Sprintf("foo.%d", rInt),
		Description:                      "Terraform acceptance tests",
		TagList:                          []string{"tag1"},
		IssuesEnabled:                    true,
		MergeRequestsEnabled:             true,
		ApprovalsBeforeMerge:             0,
		WikiEnabled:                      true,
		SnippetsEnabled:                  true,
		ContainerRegistryEnabled:         true,
		SharedRunnersEnabled:             true,
		Visibility:                       gitlab.PublicVisibility,
		MergeMethod:                      gitlab.FastForwardMerge,
		OnlyAllowMergeIfPipelineSucceeds: true,
		OnlyAllowMergeIfAllDiscussionsAreResolved: true,
	}
	willError := defaults
	willError.TagList = []string{"notatag"}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			// Step0 Create a project
			{
				Config: testAccGitlabProjectConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.foo", &received),
					testAccCheckAggregateGitlabProject(&defaults, &received),
				),
			},
			// Step1 Verify that passing bad values will fail.
			{
				Config:      testAccGitlabProjectConfig(rInt),
				ExpectError: regexp.MustCompile(`\stags\sexpected\s.+notatag.+\sreceived`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAggregateGitlabProject(&willError, &received),
				),
			},
			// Step2 Reset
			{
				Config: testAccGitlabProjectConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.foo", &received),
					testAccCheckAggregateGitlabProject(&defaults, &received),
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
		var err error
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}
		repoName := rs.Primary.ID
		if repoName == "" {
			return fmt.Errorf("No project ID is set")
		}
		conn := testAccProvider.Meta().(*gitlab.Client)
		if g, _, err := conn.Projects.GetProject(repoName, nil); err == nil {
			*project = *g
		}
		return err
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

// testAccSkipGitLabProjectAttributes are Resource attributes that should be
// skipped and handled another way, e.g. shared_with_groups
var testAccSkipGitLabProjectAttributes = []string{
	"shared_with_groups",
}

func testAccCheckAggregateGitlabProject(expected, received *gitlab.Project) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{
		testAccCheckGitLabProjectGroups(expected, received),
	}
	testResource := resourceGitlabProject()
	expectedData := testResource.TestResourceData()
	receivedData := testResource.TestResourceData()
	for a, v := range testResource.Schema {
		attribute := a
		attrValue := v
		checks = append(checks, func(_ *terraform.State) error {
			if testAccIsSkippedAttribute(attribute, testAccSkipGitLabProjectAttributes) {
				return nil // skipping because we said so.
			}
			if attrValue.Computed {
				if attrDefault, err := attrValue.DefaultValue(); err == nil {
					if attrDefault == nil {
						return nil // Skipping because we have no way of pre-computing computed vars
					}
				} else {
					return err
				}

			}
			resourceGitlabProjectSetToState(expectedData, expected)
			resourceGitlabProjectSetToState(receivedData, received)
			return testAccCompareGitLabAttribute(attribute, expectedData, receivedData)
		})
	}
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func testAccCheckGitLabProjectGroups(expected, received *gitlab.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		groupsToCheck := expected.SharedWithGroups
		for _, group := range received.SharedWithGroups {
			for i, groupToCheck := range groupsToCheck {
				if group.GroupName == groupToCheck.GroupName &&
					group.GroupAccessLevel == groupToCheck.GroupAccessLevel {
					groupsToCheck = append(groupsToCheck[:i], groupsToCheck[i+1:]...)
					break
				}
			}
		}
		if len(groupsToCheck) != 0 {
			return fmt.Errorf(
				`attribute shared_with_groups expected "%v" received "%v"`,
				received.SharedWithGroups,
				expected.SharedWithGroups)
		}
		return nil
	}
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

  tags = [
	"tag1",
  ]

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

  tags = [
	"tag1",
	"tag2",
  ]

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
  container_registry_enabled = false
	shared_runners_enabled = false
	archived = true
}
	`, rInt, rInt)
}

func testAccGitlabProjectSharedWithGroup(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name                                             = "foo-%d"
  path                                             = "foo.%d"
  description                                      = "Terraform acceptance tests"
  visibility_level                                 = "public"
  merge_method                                     = "ff"
  only_allow_merge_if_pipeline_succeeds            = false
  only_allow_merge_if_all_discussions_are_resolved = false

  shared_with_groups {
     group_id           = "${gitlab_group.foo.id}"
     group_access_level = "developer"
  }
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

  shared_with_groups {
      group_id           = "${gitlab_group.foo.id}"
      group_access_level = "guest"
  }
  shared_with_groups {
      group_id           = "${gitlab_group.foo2.id}"
      group_access_level = "developer"
  }
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
