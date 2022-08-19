//go:build acceptance
// +build acceptance

package provider

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

type testAccGitlabProjectExpectedAttributes struct {
	DefaultBranch string
}

func TestAccGitlabProject_minimal(t *testing.T) {
	var received gitlab.Project
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project" "this" {
						name = "foo-%d"
						visibility_level = "public"
					}`, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.this", &received),
				),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_project.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGitlabProject_basic(t *testing.T) {
	var received, defaults, defaultsMainBranch gitlab.Project
	rInt := acctest.RandInt()

	defaults = gitlab.Project{
		Namespace:                        &gitlab.ProjectNamespace{ID: 0},
		Name:                             fmt.Sprintf("foo-%d", rInt),
		Path:                             fmt.Sprintf("foo.%d", rInt),
		Description:                      "Terraform acceptance tests",
		TagList:                          []string{"foo", "bar"},
		RequestAccessEnabled:             true,
		IssuesEnabled:                    true,
		MergeRequestsEnabled:             true,
		JobsEnabled:                      true,
		ApprovalsBeforeMerge:             0,
		WikiEnabled:                      true,
		SnippetsEnabled:                  true,
		ContainerRegistryEnabled:         true,
		LFSEnabled:                       true,
		SharedRunnersEnabled:             true,
		Visibility:                       gitlab.PublicVisibility,
		MergeMethod:                      gitlab.FastForwardMerge,
		OnlyAllowMergeIfPipelineSucceeds: true,
		OnlyAllowMergeIfAllDiscussionsAreResolved: true,
		SquashOption:                    gitlab.SquashOptionDefaultOff,
		AllowMergeOnSkippedPipeline:     false,
		Archived:                        false, // needless, but let's make this explicit
		PackagesEnabled:                 true,
		PrintingMergeRequestLinkEnabled: true,
		PagesAccessLevel:                gitlab.PublicAccessControl,
		IssuesTemplate:                  "",
		MergeRequestsTemplate:           "",
		CIConfigPath:                    ".gitlab-ci.yml@mynamespace/myproject",
		CIForwardDeploymentEnabled:      true,
		ResolveOutdatedDiffDiscussions:  true,
		AnalyticsAccessLevel:            gitlab.EnabledAccessControl,
		AutoCancelPendingPipelines:      "enabled",
		AutoDevopsDeployStrategy:        "continuous",
		AutoDevopsEnabled:               true,
		AutocloseReferencedIssues:       true,
		BuildGitStrategy:                "fetch",
		BuildTimeout:                    42 * 60,
		BuildsAccessLevel:               gitlab.EnabledAccessControl,
		ContainerExpirationPolicy: &gitlab.ContainerExpirationPolicy{
			Enabled:   true,
			Cadence:   "1month",
			KeepN:     10,
			OlderThan: "10d",
		},
		ContainerRegistryAccessLevel:     gitlab.EnabledAccessControl,
		EmailsDisabled:                   true,
		ForkingAccessLevel:               gitlab.EnabledAccessControl,
		IssuesAccessLevel:                gitlab.EnabledAccessControl,
		MergeRequestsAccessLevel:         gitlab.EnabledAccessControl,
		OperationsAccessLevel:            gitlab.EnabledAccessControl,
		PublicBuilds:                     false,
		RepositoryAccessLevel:            gitlab.EnabledAccessControl,
		RepositoryStorage:                "default",
		SecurityAndComplianceAccessLevel: gitlab.EnabledAccessControl,
		SnippetsAccessLevel:              gitlab.EnabledAccessControl,
		Topics:                           []string{"foo", "bar"},
		WikiAccessLevel:                  gitlab.EnabledAccessControl,
		SquashCommitTemplate:             "hello squash",
		MergeCommitTemplate:              "hello merge",
	}

	defaultsMainBranch = defaults
	defaultsMainBranch.DefaultBranch = "main"

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			// Create a project with all the features on (note: "archived" is "false")
			{
				Config: testAccGitlabProjectConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.foo", &received),
					testAccCheckAggregateGitlabProject(&defaults, &received),
				),
			},
			// Update the project to turn the features off (note: "archived" is "true")
			{
				Config: testAccGitlabProjectUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.foo", &received),
					testAccCheckAggregateGitlabProject(&gitlab.Project{
						Namespace:                        &gitlab.ProjectNamespace{ID: 0},
						Name:                             fmt.Sprintf("foo-%d", rInt),
						Path:                             fmt.Sprintf("foo.%d", rInt),
						Description:                      "Terraform acceptance tests!",
						TagList:                          []string{"foo", "bar"},
						JobsEnabled:                      false,
						ApprovalsBeforeMerge:             0,
						RequestAccessEnabled:             false,
						ContainerRegistryEnabled:         false,
						LFSEnabled:                       false,
						SharedRunnersEnabled:             false,
						Visibility:                       gitlab.PublicVisibility,
						MergeMethod:                      gitlab.FastForwardMerge,
						PrintingMergeRequestLinkEnabled:  true,
						OnlyAllowMergeIfPipelineSucceeds: true,
						OnlyAllowMergeIfAllDiscussionsAreResolved: true,
						SquashOption:                   gitlab.SquashOptionDefaultOn,
						AllowMergeOnSkippedPipeline:    true,
						Archived:                       true,
						PackagesEnabled:                false,
						PagesAccessLevel:               gitlab.DisabledAccessControl,
						CIForwardDeploymentEnabled:     false,
						ResolveOutdatedDiffDiscussions: false,
						AnalyticsAccessLevel:           gitlab.DisabledAccessControl,
						AutoCancelPendingPipelines:     "disabled",
						AutoDevopsDeployStrategy:       "manual",
						AutoDevopsEnabled:              false,
						AutocloseReferencedIssues:      false,
						BuildGitStrategy:               "fetch",
						BuildTimeout:                   10 * 60,
						BuildsAccessLevel:              gitlab.DisabledAccessControl,
						ContainerExpirationPolicy: &gitlab.ContainerExpirationPolicy{
							Enabled:   true,
							Cadence:   "4h",
							KeepN:     10,
							OlderThan: "10d",
						},
						ContainerRegistryAccessLevel:     gitlab.DisabledAccessControl,
						EmailsDisabled:                   false,
						ForkingAccessLevel:               gitlab.DisabledAccessControl,
						IssuesAccessLevel:                gitlab.DisabledAccessControl,
						MergeRequestsAccessLevel:         gitlab.DisabledAccessControl,
						OperationsAccessLevel:            gitlab.DisabledAccessControl,
						PublicBuilds:                     false,
						RepositoryAccessLevel:            gitlab.DisabledAccessControl,
						RepositoryStorage:                "default",
						SecurityAndComplianceAccessLevel: gitlab.DisabledAccessControl,
						SnippetsAccessLevel:              gitlab.DisabledAccessControl,
						Topics:                           []string{},
						WikiAccessLevel:                  gitlab.DisabledAccessControl,
						SquashCommitTemplate:             "goodby squash",
						MergeCommitTemplate:              "goodby merge",
					}, &received),
				),
			},
			// Update the project to turn the features on again (note: "archived" is "false")
			{
				Config: testAccGitlabProjectConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.foo", &received),
					testAccCheckAggregateGitlabProject(&defaults, &received),
				),
			},
			// Update the project creating the default branch
			{
				// Get the ID from the project data at the previous step
				SkipFunc: testAccGitlabProjectConfigDefaultBranchSkipFunc(&received, "main"),
				Config:   testAccGitlabProjectConfigDefaultBranch(rInt, "main"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.foo", &received),
					testAccCheckAggregateGitlabProject(&defaultsMainBranch, &received),
				),
			},
			// Test import without push rules (checks read function)
			{
				ResourceName:      "gitlab_project.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Add all push rules to an existing project
			{
				SkipFunc: isRunningInCE,
				Config: testAccGitlabProjectConfigPushRules(rInt, `
author_email_regex = "foo_author"
branch_name_regex = "foo_branch"
commit_message_regex = "foo_commit"
commit_message_negative_regex = "foo_not_commit"
file_name_regex = "foo_file"
commit_committer_check = true
deny_delete_tag = true
member_check = true
prevent_secrets = true
reject_unsigned_commits = true
max_file_size = 123
`),
				Check: testAccCheckGitlabProjectPushRules("gitlab_project.foo", &gitlab.ProjectPushRules{
					AuthorEmailRegex:           "foo_author",
					BranchNameRegex:            "foo_branch",
					CommitMessageRegex:         "foo_commit",
					CommitMessageNegativeRegex: "foo_not_commit",
					FileNameRegex:              "foo_file",
					CommitCommitterCheck:       true,
					DenyDeleteTag:              true,
					MemberCheck:                true,
					PreventSecrets:             true,
					RejectUnsignedCommits:      true,
					MaxFileSize:                123,
				}),
			},
			// Test import with a all push rules defined (checks read function)
			{
				SkipFunc:          isRunningInCE,
				ResourceName:      "gitlab_project.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update some push rules but not others
			{
				SkipFunc: isRunningInCE,
				Config: testAccGitlabProjectConfigPushRules(rInt, `
author_email_regex = "foo_author"
branch_name_regex = "foo_branch"
commit_message_regex = "foo_commit"
commit_message_negative_regex = "foo_not_commit"
file_name_regex = "foo_file_2"
commit_committer_check = true
deny_delete_tag = true
member_check = false
prevent_secrets = true
reject_unsigned_commits = true
max_file_size = 1234
`),
				Check: testAccCheckGitlabProjectPushRules("gitlab_project.foo", &gitlab.ProjectPushRules{
					AuthorEmailRegex:           "foo_author",
					BranchNameRegex:            "foo_branch",
					CommitMessageRegex:         "foo_commit",
					CommitMessageNegativeRegex: "foo_not_commit",
					FileNameRegex:              "foo_file_2",
					CommitCommitterCheck:       true,
					DenyDeleteTag:              true,
					MemberCheck:                false,
					PreventSecrets:             true,
					RejectUnsignedCommits:      true,
					MaxFileSize:                1234,
				}),
			},
			// Try to add push rules to an existing project in CE
			{
				SkipFunc:    isRunningInEE,
				Config:      testAccGitlabProjectConfigPushRules(rInt, `author_email_regex = "foo_author"`),
				ExpectError: regexp.MustCompile(regexp.QuoteMeta("Project push rules are not supported in your version of GitLab")),
			},
			// Update push rules
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabProjectConfigPushRules(rInt, `author_email_regex = "foo_author"`),
				Check: testAccCheckGitlabProjectPushRules("gitlab_project.foo", &gitlab.ProjectPushRules{
					AuthorEmailRegex: "foo_author",
				}),
			},
			// Remove the push_rules block entirely.
			// NOTE: The push rules will still exist upstream because the push_rules block is computed.
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabProjectConfigDefaultBranch(rInt, "main"),
				Check: testAccCheckGitlabProjectPushRules("gitlab_project.foo", &gitlab.ProjectPushRules{
					AuthorEmailRegex: "foo_author",
				}),
			},
			// Add different push rules after the block was removed previously
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabProjectConfigPushRules(rInt, `branch_name_regex = "(feature|hotfix)\\/*"`),
				Check: testAccCheckGitlabProjectPushRules("gitlab_project.foo", &gitlab.ProjectPushRules{
					BranchNameRegex: `(feature|hotfix)\/*`,
				}),
			},
			// Destroy the project so we can next test creating a project with push rules simultaneously
			{
				Config:  testAccGitlabProjectConfigDefaultBranch(rInt, "main"),
				Destroy: true,
				Check:   testAccCheckGitlabProjectDestroy,
			},
			// Create a new project with push rules
			{
				SkipFunc: isRunningInCE,
				Config: testAccGitlabProjectConfigPushRules(rInt, `
author_email_regex = "foo_author"
max_file_size = 123
`),
				Check: testAccCheckGitlabProjectPushRules("gitlab_project.foo", &gitlab.ProjectPushRules{
					AuthorEmailRegex: "foo_author",
					MaxFileSize:      123,
				}),
			},
			// Try to create a new project with all push rules in CE
			{
				SkipFunc:    isRunningInEE,
				Config:      testAccGitlabProjectConfigPushRules(rInt, `author_email_regex = "foo_author"`),
				ExpectError: regexp.MustCompile(regexp.QuoteMeta("Project push rules are not supported in your version of GitLab")),
			},
			// Create a project using template name
			{
				Config: testAccGitlabProjectConfigTemplateName(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.template-name", &received),
					testAccCheckGitlabProjectDefaultBranch(&received, &testAccGitlabProjectExpectedAttributes{
						DefaultBranch: "master",
					}),
					func(state *terraform.State) error {
						projectID := state.RootModule().Resources["gitlab_project.template-name"].Primary.ID

						_, _, err := testGitlabClient.RepositoryFiles.GetFile(projectID, ".ruby-version", &gitlab.GetFileOptions{Ref: gitlab.String("master")}, nil)
						if err != nil {
							return fmt.Errorf("failed to get '.ruby-version' file from template project: %w", err)
						}

						return nil
					},
				),
			},
			// Create a project using custom template name
			{
				Config:   testAccGitlabProjectConfigTemplateNameCustom(rInt),
				SkipFunc: isRunningInCE,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.template-name-custom", &received),
					testAccCheckGitlabProjectDefaultBranch(&received, &testAccGitlabProjectExpectedAttributes{
						DefaultBranch: "master",
					}),
					func(state *terraform.State) error {
						projectID := state.RootModule().Resources["gitlab_project.template-name-custom"].Primary.ID

						_, _, err := testGitlabClient.RepositoryFiles.GetFile(projectID, "Gemfile", &gitlab.GetFileOptions{Ref: gitlab.String("master")}, nil)
						if err != nil {
							return fmt.Errorf("failed to get 'Gemfile' file from template project: %w", err)
						}

						return nil
					},
				),
			},
			// Create a project using custom template project id
			{
				Config:   testAccGitlabProjectConfigTemplateProjectID(rInt),
				SkipFunc: isRunningInCE,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.template-id", &received),
					testAccCheckGitlabProjectDefaultBranch(&received, &testAccGitlabProjectExpectedAttributes{
						DefaultBranch: "master",
					}),
					func(state *terraform.State) error {
						projectID := state.RootModule().Resources["gitlab_project.template-id"].Primary.ID

						_, _, err := testGitlabClient.RepositoryFiles.GetFile(projectID, "Rakefile", &gitlab.GetFileOptions{Ref: gitlab.String("master")}, nil)
						if err != nil {
							return fmt.Errorf("failed to get 'Rakefile' file from template project: %w", err)
						}

						return nil
					},
				),
			},
			// Update to original project config
			{
				Config: testAccGitlabProjectConfig(rInt),
			},
		},
	})
}

func TestAccGitlabProject_initializeWithReadme(t *testing.T) {
	var project gitlab.Project
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabProjectConfigInitializeWithReadme(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.foo", &project),
					testAccCheckGitlabProjectDefaultBranch(&project, nil),
					func(state *terraform.State) error {
						_, _, err := testGitlabClient.RepositoryFiles.GetFile(project.ID, "README.md", &gitlab.GetFileOptions{Ref: gitlab.String("main")}, nil)
						if err != nil {
							return fmt.Errorf("failed to get 'README.md' file from project: %w", err)
						}

						return nil
					},
				),
			},
		},
	})
}

func TestAccGitlabProject_initializeWithoutReadme(t *testing.T) {
	var project gitlab.Project
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabProjectConfigInitializeWithoutReadme(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.foo", &project),
					func(s *terraform.State) error {
						branches, _, err := testGitlabClient.Branches.ListBranches(project.ID, nil)
						if err != nil {
							return fmt.Errorf("failed to list branches: %w", err)
						}

						if len(branches) != 0 {
							return fmt.Errorf("expected no branch for new project when initialized without README; found %d", len(branches))
						}
						return nil
					},
				),
			},
		},
	})
}

func TestAccGitlabProject_archiveOnDestroy(t *testing.T) {
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectArchivedOnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabProjectConfigArchiveOnDestroy(rInt),
			},
		},
	})
}

func TestAccGitlabProject_setSinglePushRuleToDefault(t *testing.T) {
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			{
				SkipFunc: isRunningInCE,
				Config: testAccGitlabProjectConfigPushRules(rInt, `
member_check = false
`),
				Check: testAccCheckGitlabProjectPushRules("gitlab_project.foo", &gitlab.ProjectPushRules{
					MemberCheck: false,
				}),
			},
		},
	})
}

func TestAccGitlabProject_groupWithoutDefaultBranchProtection(t *testing.T) {
	var project gitlab.Project
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabProjectConfigWithoutDefaultBranchProtection(rInt),
				Check:  testAccCheckGitlabProjectExists("gitlab_project.foo", &project),
			},
			{
				Config:  testAccGitlabProjectConfigWithoutDefaultBranchProtection(rInt),
				Destroy: true,
			},
			{
				Config: testAccGitlabProjectConfigWithoutDefaultBranchProtectionInitializeReadme(rInt),
				Check:  testAccCheckGitlabProjectExists("gitlab_project.foo", &project),
			},
		},
	})
}

func TestAccGitlabProject_IssueMergeRequestTemplates(t *testing.T) {
	var project gitlab.Project
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabProjectConfigIssueMergeRequestTemplates(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.foo", &project),
					func(s *terraform.State) error {
						if project.IssuesTemplate != "foo" {
							return fmt.Errorf("expected issues template to be 'foo'; got '%s'", project.IssuesTemplate)
						}

						if project.MergeRequestsTemplate != "bar" {
							return fmt.Errorf("expected merge requests template to be 'bar'; got '%s'", project.MergeRequestsTemplate)
						}

						return nil
					},
				),
			},
		},
	})
}

func TestAccGitlabProject_MergePipelines(t *testing.T) {
	var project gitlab.Project
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitLabProjectMergePipelinesEnabled(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.foo", &project),
					func(s *terraform.State) error {
						if project.MergePipelinesEnabled != true {
							return fmt.Errorf("expected merge pipelines to be enabled")
						}

						return nil
					},
				),
			},
		},
	})
}

func TestAccGitlabProject_MergeTrains(t *testing.T) {
	var project gitlab.Project
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitLabProjectMergeTrainsEnabled(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.foo", &project),
					func(s *terraform.State) error {
						if project.MergeTrainsEnabled != true {
							return fmt.Errorf("expected merge trains to be enabled")
						}

						return nil
					},
				),
			},
		},
	})
}

func TestAccGitlabProject_willErrorOnAPIFailure(t *testing.T) {
	var received gitlab.Project
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			// Step0 Create a project
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project" "this" {
						name = "foo-%d"

						visibility_level = "public"
					}`, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.this", &received),
				),
			},
			// Step1 Verify that passing bad values will fail.
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project" "this" {
						name = "foo-%d"
						repository_storage = "non-existing"

						visibility_level = "public"
					}`, rInt),
				// This will fail because the repository_storage is not valid.
				ExpectError: regexp.MustCompile(`\[is invalid\]`),
			},
			// Step2 Reset
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project" "this" {
						name = "foo-%d"

						visibility_level = "public"
					}`, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.this", &received),
				),
			},
		},
	})
}

// lintignore: AT002 // specialized import test
func TestAccGitlabProject_import(t *testing.T) {
	rInt := acctest.RandInt()
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			{
				SkipFunc: isRunningInEE,
				Config:   testAccGitlabProjectConfig(rInt),
			},
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabProjectConfigEE(rInt),
			},
			{
				ResourceName:      "gitlab_project.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// lintignore: AT002 // specialized import test
func TestAccGitlabProject_nestedImport(t *testing.T) {
	rInt := acctest.RandInt()
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectDestroy,
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

func TestAccGitlabProject_transfer(t *testing.T) {
	var transferred, received gitlab.Project
	rInt := acctest.RandInt()

	transferred = gitlab.Project{
		Namespace:                        &gitlab.ProjectNamespace{Name: fmt.Sprintf("foo2group-%d", rInt)},
		Name:                             fmt.Sprintf("foo-%d", rInt),
		Path:                             fmt.Sprintf("foo-%d", rInt),
		Description:                      "Terraform acceptance tests",
		TagList:                          []string{},
		RequestAccessEnabled:             true,
		IssuesEnabled:                    true,
		MergeRequestsEnabled:             true,
		JobsEnabled:                      true,
		ApprovalsBeforeMerge:             0,
		WikiEnabled:                      true,
		SnippetsEnabled:                  true,
		ContainerRegistryEnabled:         true,
		LFSEnabled:                       true,
		SharedRunnersEnabled:             true,
		Visibility:                       gitlab.PublicVisibility,
		MergeMethod:                      gitlab.NoFastForwardMerge,
		OnlyAllowMergeIfPipelineSucceeds: false,
		OnlyAllowMergeIfAllDiscussionsAreResolved: false,
		SquashOption:                    gitlab.SquashOptionDefaultOff,
		PackagesEnabled:                 true,
		PrintingMergeRequestLinkEnabled: true,
		PagesAccessLevel:                gitlab.PrivateAccessControl,
		CIForwardDeploymentEnabled:      true,
	}

	pathBeforeTransfer := fmt.Sprintf("foogroup-%d/foo-%d", rInt, rInt)
	pathAfterTransfer := fmt.Sprintf("foo2group-%d/foo-%d", rInt, rInt)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			// Create a project in a group
			{
				Config: testAccGitlabProjectTransferBetweenGroupsBefore(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.foo", &received),
					resource.TestCheckResourceAttrPtr("gitlab_project_variable.foo", "value", &pathBeforeTransfer),
				),
			},
			// Create a second group and set the transfer the project to this group
			{
				Config: testAccGitlabProjectTransferBetweenGroupsAfter(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.foo", &received),
					testAccCheckAggregateGitlabProject(&transferred, &received),
					resource.TestCheckResourceAttrPtr("gitlab_project_variable.foo", "value", &pathAfterTransfer),
				),
			},
		},
	})
}

// lintignore: AT002 // not a Terraform import test
func TestAccGitlabProject_importURL(t *testing.T) {
	rInt := acctest.RandInt()

	// Create a base project for importing.
	baseProject, _, err := testGitlabClient.Projects.CreateProject(&gitlab.CreateProjectOptions{
		Name:       gitlab.String(fmt.Sprintf("base-%d", rInt)),
		Visibility: gitlab.Visibility(gitlab.PublicVisibility),
	})
	if err != nil {
		t.Fatalf("failed to create base project: %v", err)
	}

	defer testGitlabClient.Projects.DeleteProject(baseProject.ID) // nolint // TODO: Resolve this golangci-lint issue: Error return value of `testGitlabClient.Projects.DeleteProject` is not checked (errcheck)

	// Add a file to the base project, for later verifying the import.
	_, _, err = testGitlabClient.RepositoryFiles.CreateFile(baseProject.ID, "foo.txt", &gitlab.CreateFileOptions{
		Branch:        gitlab.String("main"),
		CommitMessage: gitlab.String("add file"),
		Content:       gitlab.String(""),
	})
	if err != nil {
		t.Fatalf("failed to commit file to base project: %v", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabProjectConfigImportURL(rInt, baseProject.HTTPURLToRepo),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gitlab_project.imported", "import_url", baseProject.HTTPURLToRepo),
					func(state *terraform.State) error {
						projectID := state.RootModule().Resources["gitlab_project.imported"].Primary.ID

						_, _, err := testGitlabClient.RepositoryFiles.GetFile(projectID, "foo.txt", &gitlab.GetFileOptions{Ref: gitlab.String("main")}, nil)
						if err != nil {
							return fmt.Errorf("failed to get file from imported project: %w", err)
						}

						return nil
					},
				),
			},
		},
	})
}

func TestAccGitlabProject_initializeWithReadmeAndCustomDefaultBranch(t *testing.T) {
	var project gitlab.Project
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name        = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"

  initialize_with_readme = true
  default_branch         = "foo"
}`, rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gitlab_project.foo", "initialize_with_readme", "true"),
					resource.TestCheckResourceAttr("gitlab_project.foo", "default_branch", "foo"),
					testAccCheckGitlabProjectExists("gitlab_project.foo", &project),
					testAccCheckGitlabProjectDefaultBranch(&project, &testAccGitlabProjectExpectedAttributes{
						DefaultBranch: "foo",
					}),
					func(state *terraform.State) error {
						projectID := state.RootModule().Resources["gitlab_project.foo"].Primary.ID

						_, _, err := testGitlabClient.RepositoryFiles.GetFile(projectID, "README.md", &gitlab.GetFileOptions{Ref: gitlab.String("foo")}, nil)
						if err != nil {
							return fmt.Errorf("failed to get 'README.md' file from project: %w", err)
						}

						return nil
					},
				),
			},
			// Verify Import
			{
				ResourceName:            "gitlab_project.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initialize_with_readme"},
			},
		},
	})
}

func TestAccGitlabProject_CreateProjectInUserNamespace(t *testing.T) {
	var project gitlab.Project
	rInt := acctest.RandInt()

	user := testAccCreateUsers(t, 1)[0]

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccRequiresAtLeast(t, "14.10") },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project" "foo" {
						name              = "foo-%d"
						description       = "Terraform acceptance tests"
						visibility_level  = "public"

						namespace_id = %d
					}
				`, rInt, user.NamespaceID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.foo", &project),
					func(s *terraform.State) error {
						if project.Namespace.ID != user.NamespaceID {
							return fmt.Errorf("project was created in namespace %d but expected %d", project.Namespace.ID, user.NamespaceID)
						}
						return nil
					},
				),
			},
		},
	})
}

func TestAccGitlabProject_InstanceBranchProtectionDisabled(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					settings, _, err := testGitlabClient.Settings.GetSettings()
					if err != nil {
						t.Fatalf("failed to get settings: %v", err)
					}
					t.Cleanup(func() {
						if _, _, err := testGitlabClient.Settings.UpdateSettings(&gitlab.UpdateSettingsOptions{DefaultBranchProtection: gitlab.Int(settings.DefaultBranchProtection)}); err != nil {
							t.Fatalf("failed to update instance-wide default branch protection setting to default: %v", err)
						}
					})

					if _, _, err := testGitlabClient.Settings.UpdateSettings(&gitlab.UpdateSettingsOptions{DefaultBranchProtection: gitlab.Int(0)}); err != nil {
						t.Fatalf("failed to update instance-wide default branch protection setting: %v", err)
					}
				},
				Config: ` `, // requires a space for empty config
			},
			// Without explicit default branch
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project" "foo" {
						name                   = "foo-%d"
						description            = "Terraform acceptance tests"
						visibility_level       = "public"
						initialize_with_readme = true
					}
				`, rInt),
			},
			// Verify Import
			{
				ResourceName:            "gitlab_project.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initialize_with_readme"},
			},
			// Force a destroy for the project so that it can be recreated as the same resource
			{
				Config: ` `, // requires a space for empty config
			},
			// With explicit default branch set to instance-wide default
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project" "foo" {
						name                   = "foo-%d"
						description            = "Terraform acceptance tests"
						visibility_level       = "public"
						default_branch         = "main"
						initialize_with_readme = true
					}
				`, rInt),
			},
			// Verify Import
			{
				ResourceName:            "gitlab_project.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initialize_with_readme"},
			},
			// Force a destroy for the project so that it can be recreated as the same resource
			{
				Config: ` `, // requires a space for empty config
			},
			// With custom default branch
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project" "foo" {
						name                   = "foo-%d-custom-default-branch"
						description            = "Terraform acceptance tests"
						visibility_level       = "public"
						default_branch         = "foobar-non-default-branch"
						initialize_with_readme = true
					}
				`, rInt),
			},
			// Verify Import
			{
				ResourceName:            "gitlab_project.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initialize_with_readme"},
			},
			// Force a destroy for the project so that it can be recreated as the same resource
			{
				Config: ` `, // requires a space for empty config
			},
			// With `skip_wait_for_default_branch_protection` enabled
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project" "foo" {
						name                   = "foo-%d-custom-default-branch"
						description            = "Terraform acceptance tests"
						visibility_level       = "public"
						initialize_with_readme = true

						skip_wait_for_default_branch_protection = true
					}
				`, rInt),
			},
			// Verify Import
			{
				ResourceName:            "gitlab_project.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initialize_with_readme", "skip_wait_for_default_branch_protection"},
			},
			// Force a destroy for the project so that it can be recreated as the same resource
			{
				Config: ` `, // requires a space for empty config
			},
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project" "foo" {
						name                   = "foo-%d-custom-default-branch"
						description            = "Terraform acceptance tests"
						visibility_level       = "public"
						initialize_with_readme = true

						skip_wait_for_default_branch_protection = false
					}
				`, rInt),
			},
			// Check if plan is empty after changing `skip_wait_for_default_branch_protection` attribute
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project" "foo" {
						name                   = "foo-%d-custom-default-branch"
						description            = "Terraform acceptance tests"
						visibility_level       = "public"
						initialize_with_readme = true

						skip_wait_for_default_branch_protection = true
					}
				`, rInt),
				PlanOnly: true,
			},
			// Verify Import
			{
				ResourceName:            "gitlab_project.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initialize_with_readme", "skip_wait_for_default_branch_protection"},
			},
		},
	})
}

type testAccGitlabProjectMirroredExpectedAttributes struct {
	Mirror                           bool
	MirrorTriggerBuilds              bool
	MirrorOverwritesDivergedBranches bool
	OnlyMirrorProtectedBranches      bool
}

func testAccCheckGitlabProjectMirroredAttributes(project *gitlab.Project, want *testAccGitlabProjectMirroredExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if project.Mirror != want.Mirror {
			return fmt.Errorf("got mirror %t; want %t", project.Mirror, want.Mirror)
		}

		if project.MirrorTriggerBuilds != want.MirrorTriggerBuilds {
			return fmt.Errorf("got mirror_trigger_builds %t; want %t", project.MirrorTriggerBuilds, want.MirrorTriggerBuilds)
		}

		if project.MirrorOverwritesDivergedBranches != want.MirrorOverwritesDivergedBranches {
			return fmt.Errorf("got mirror_overwrites_diverged_branches %t; want %t", project.MirrorOverwritesDivergedBranches, want.MirrorOverwritesDivergedBranches)
		}

		if project.OnlyMirrorProtectedBranches != want.OnlyMirrorProtectedBranches {
			return fmt.Errorf("got only_mirror_protected_branches %t; want %t", project.OnlyMirrorProtectedBranches, want.OnlyMirrorProtectedBranches)
		}
		return nil
	}
}

// lintignore: AT002 // not a Terraform import test
func TestAccGitlabProject_ImportURLMirrored(t *testing.T) {
	var mirror gitlab.Project
	rInt := acctest.RandInt()

	// Create a base project for importing.
	baseProject, _, err := testGitlabClient.Projects.CreateProject(&gitlab.CreateProjectOptions{
		Name:       gitlab.String(fmt.Sprintf("base-%d", rInt)),
		Visibility: gitlab.Visibility(gitlab.PublicVisibility),
	})
	if err != nil {
		t.Fatalf("failed to create base project: %v", err)
	}

	defer testGitlabClient.Projects.DeleteProject(baseProject.ID) // nolint // TODO: Resolve this golangci-lint issue: Error return value of `testGitlabClient.Projects.DeleteProject` is not checked (errcheck)

	// Add a file to the base project, for later verifying the import.
	_, _, err = testGitlabClient.RepositoryFiles.CreateFile(baseProject.ID, "foo.txt", &gitlab.CreateFileOptions{
		Branch:        gitlab.String("main"),
		CommitMessage: gitlab.String("add file"),
		Content:       gitlab.String(""),
	})
	if err != nil {
		t.Fatalf("failed to commit file to base project: %v", err)
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			{
				// First, import, as mirrored
				Config:   testAccGitlabProjectConfigImportURLMirror(rInt, baseProject.HTTPURLToRepo),
				SkipFunc: isRunningInCE,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.imported", &mirror),
					resource.TestCheckResourceAttr("gitlab_project.imported", "import_url", baseProject.HTTPURLToRepo),
					testAccCheckGitlabProjectMirroredAttributes(&mirror, &testAccGitlabProjectMirroredExpectedAttributes{
						Mirror:                           true,
						MirrorTriggerBuilds:              true,
						MirrorOverwritesDivergedBranches: true,
						OnlyMirrorProtectedBranches:      true,
					}),

					func(state *terraform.State) error {
						projectID := state.RootModule().Resources["gitlab_project.imported"].Primary.ID

						_, _, err := testGitlabClient.RepositoryFiles.GetFile(projectID, "foo.txt", &gitlab.GetFileOptions{Ref: gitlab.String("main")}, nil)
						if err != nil {
							return fmt.Errorf("failed to get file from imported project: %w", err)
						}

						return nil
					},
				),
			},
			{
				// Second, disable all optional mirroring options
				Config:   testAccGitlabProjectConfigImportURLMirrorDisabledOptionals(rInt, baseProject.HTTPURLToRepo),
				SkipFunc: isRunningInCE,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.imported", &mirror),
					resource.TestCheckResourceAttr("gitlab_project.imported", "import_url", baseProject.HTTPURLToRepo),
					testAccCheckGitlabProjectMirroredAttributes(&mirror, &testAccGitlabProjectMirroredExpectedAttributes{
						Mirror:                           true,
						MirrorTriggerBuilds:              false,
						MirrorOverwritesDivergedBranches: false,
						OnlyMirrorProtectedBranches:      false,
					}),

					// Ensure the test file still is as expected
					func(state *terraform.State) error {
						projectID := state.RootModule().Resources["gitlab_project.imported"].Primary.ID

						_, _, err := testGitlabClient.RepositoryFiles.GetFile(projectID, "foo.txt", &gitlab.GetFileOptions{Ref: gitlab.String("main")}, nil)
						if err != nil {
							return fmt.Errorf("failed to get file from imported project: %w", err)
						}

						return nil
					},
				),
			},
			{
				// Third, disable mirroring, using the original ImportURL acceptance test
				Config:   testAccGitlabProjectConfigImportURLMirrorDisabled(rInt, baseProject.HTTPURLToRepo),
				SkipFunc: isRunningInCE,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.imported", &mirror),
					resource.TestCheckResourceAttr("gitlab_project.imported", "import_url", baseProject.HTTPURLToRepo),
					testAccCheckGitlabProjectMirroredAttributes(&mirror, &testAccGitlabProjectMirroredExpectedAttributes{
						Mirror:                           false,
						MirrorTriggerBuilds:              false,
						MirrorOverwritesDivergedBranches: false,
						OnlyMirrorProtectedBranches:      false,
					}),

					// Ensure the test file still is as expected
					func(state *terraform.State) error {
						projectID := state.RootModule().Resources["gitlab_project.imported"].Primary.ID

						_, _, err := testGitlabClient.RepositoryFiles.GetFile(projectID, "foo.txt", &gitlab.GetFileOptions{Ref: gitlab.String("main")}, nil)
						if err != nil {
							return fmt.Errorf("failed to get file from imported project: %w", err)
						}

						return nil
					},
				),
			},
		},
	})
}

func TestAccGitlabProject_templateMutualExclusiveNameAndID(t *testing.T) {
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckMutualExclusiveNameAndID(rInt),
				SkipFunc:    isRunningInCE,
				ExpectError: regexp.MustCompile(regexp.QuoteMeta(`"template_project_id": conflicts with template_name`)),
			},
		},
	})
}

// Gitlab update project API call requires one from a subset of project fields to be set (See #1157)
// If only a non-blessed field is changed, this test checks that the provider ensures the code won't return an error.
func TestAccGitlabProject_UpdateAnalyticsAccessLevel(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			// Create minimal test project
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project" "this" {
						name = "foo-%d"
						visibility_level                = "public"
						analytics_access_level = "private"
					}`, rInt),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_project.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update `analytics_access_level`
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project" "this" {
						name = "foo-%d"
						visibility_level = "public"
						analytics_access_level = "disabled"
					}`, rInt),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_project.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGitlabProject_containerExpirationPolicy(t *testing.T) {
	var received gitlab.Project
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project" "this" {
						name = "foo-%d"

						container_expiration_policy {
							enabled = true
							cadence = "1d"
						}

						visibility_level = "public"
					}`, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.this", &received),
					resource.TestCheckResourceAttr("gitlab_project.this", "container_expiration_policy.0.enabled", "true"),
					resource.TestCheckResourceAttr("gitlab_project.this", "container_expiration_policy.0.cadence", "1d"),
				),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_project.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Set more attributes
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project" "this" {
						name = "foo-%d"

						container_expiration_policy {
							enabled = true
							cadence = "1month"
							name_regex_keep = "bar"
						}

						visibility_level = "public"
					}`, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.this", &received),
					resource.TestCheckResourceAttr("gitlab_project.this", "container_expiration_policy.0.enabled", "true"),
					resource.TestCheckResourceAttr("gitlab_project.this", "container_expiration_policy.0.cadence", "1month"),
					resource.TestCheckResourceAttr("gitlab_project.this", "container_expiration_policy.0.name_regex_keep", "bar"),
					resource.TestCheckResourceAttrSet("gitlab_project.this", "container_expiration_policy.0.next_run_at"),
				),
			},
			// Clear attributes
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project" "this" {
						name = "foo-%d"

						container_expiration_policy {}

						visibility_level = "public"
					}`, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.this", &received),
					resource.TestCheckResourceAttr("gitlab_project.this", "container_expiration_policy.0.enabled", "true"),
					resource.TestCheckResourceAttr("gitlab_project.this", "container_expiration_policy.0.cadence", "1month"),
					resource.TestCheckResourceAttr("gitlab_project.this", "container_expiration_policy.0.name_regex_keep", "bar"),
					resource.TestCheckResourceAttrSet("gitlab_project.this", "container_expiration_policy.0.next_run_at"),
				),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_project.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGitlabProject_DeprecatedBuildCoverageRegex(t *testing.T) {
	var received gitlab.Project
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			{
				SkipFunc: isGitLabVersionAtLeast(context.Background(), testGitlabClient, "15.0"),
				Config: fmt.Sprintf(`
					resource "gitlab_project" "this" {
						name = "foo-%d"
						visibility_level = "public"

						build_coverage_regex = "helloWorld"
					}`, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectExists("gitlab_project.this", &received),
				),
			},
			{
				SkipFunc:          isGitLabVersionAtLeast(context.Background(), testGitlabClient, "15.0"),
				ResourceName:      "gitlab_project.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGitlabProject_SetDefaultFalseBooleansOnCreate(t *testing.T) {
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project" "this" {
						name             = "foo-%d"
						visibility_level = "public"

						initialize_with_readme              = false
						resolve_outdated_diff_discussions   = false
						auto_devops_enabled                 = false
						autoclose_referenced_issues         = false
						emails_disabled                     = false
						public_builds                       = false
						merge_pipelines_enabled             = false
						merge_trains_enabled                = false
						ci_forward_deployment_enabled       = false
					}`, rInt),
			},
			{
				ResourceName:            "gitlab_project.this",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initialize_with_readme"},
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
		if g, _, err := testGitlabClient.Projects.GetProject(repoName, nil); err == nil {
			*project = *g
		}
		return err
	}
}

func testAccCheckGitlabProjectDestroy(s *terraform.State) error {
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

func testAccCheckGitlabProjectArchivedOnDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_project" {
			continue
		}

		gotRepo, _, err := testGitlabClient.Projects.GetProject(rs.Primary.ID, nil)
		if err != nil {
			return fmt.Errorf("unable to get project %s, to check if it has been archived on the destroy", rs.Primary.ID)
		}

		if !gotRepo.Archived {
			return fmt.Errorf("expected project to be archived, but it isn't")
		}
		return nil
	}

	return fmt.Errorf("no project resources found in state, but expected a `gitlab_project` resource marked as archvied")
}

func testAccCheckAggregateGitlabProject(expected, received *gitlab.Project) resource.TestCheckFunc {
	var checks []resource.TestCheckFunc

	testResource := allResources["gitlab_project"]()
	expectedData := testResource.TestResourceData()
	receivedData := testResource.TestResourceData()
	for a, v := range testResource.Schema {
		attribute := a
		attrValue := v
		checks = append(checks, func(_ *terraform.State) error {
			if attrValue.Computed {
				if attrDefault, err := attrValue.DefaultValue(); err == nil {
					if attrDefault == nil {
						return nil // Skipping because we have no way of pre-computing computed vars
					}
				} else {
					return err
				}
			}

			if err := resourceGitlabProjectSetToState(context.Background(), testGitlabClient, expectedData, expected); err != nil {
				return err
			}

			if err := resourceGitlabProjectSetToState(context.Background(), testGitlabClient, receivedData, received); err != nil {
				return err
			}

			// ignored for now
			if attribute == "container_expiration_policy" {
				return nil
			}

			return testAccCompareGitLabAttribute(attribute, expectedData, receivedData)
		})
	}
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func testAccCheckGitlabProjectDefaultBranch(project *gitlab.Project, want *testAccGitlabProjectExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if want != nil && project.DefaultBranch != want.DefaultBranch {
			return fmt.Errorf("got default branch %q; want %q", project.DefaultBranch, want.DefaultBranch)
		}

		branches, _, err := testGitlabClient.Branches.ListBranches(project.ID, nil)
		if err != nil {
			return fmt.Errorf("failed to list branches: %w", err)
		}

		if len(branches) != 1 {
			return fmt.Errorf("expected 1 branch for new project; found %d", len(branches))
		}

		if !branches[0].Protected {
			return errors.New("expected default branch to be protected")
		}

		return nil
	}
}

func testAccCheckGitlabProjectPushRules(name string, wantPushRules *gitlab.ProjectPushRules) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		projectResource := state.RootModule().Resources[name].Primary

		gotPushRules, _, err := testGitlabClient.Projects.GetProjectPushRules(projectResource.ID, nil)
		if err != nil {
			return err
		}

		var messages []string

		if gotPushRules.AuthorEmailRegex != wantPushRules.AuthorEmailRegex {
			messages = append(messages, fmt.Sprintf("author_email_regex (got: %q, wanted: %q)",
				gotPushRules.AuthorEmailRegex, wantPushRules.AuthorEmailRegex))
		}

		if gotPushRules.BranchNameRegex != wantPushRules.BranchNameRegex {
			messages = append(messages, fmt.Sprintf("branch_name_regex (got: %q, wanted: %q)",
				gotPushRules.BranchNameRegex, wantPushRules.BranchNameRegex))
		}

		if gotPushRules.CommitMessageRegex != wantPushRules.CommitMessageRegex {
			messages = append(messages, fmt.Sprintf("commit_message_regex (got: %q, wanted: %q)",
				gotPushRules.CommitMessageRegex, wantPushRules.CommitMessageRegex))
		}

		if gotPushRules.CommitMessageNegativeRegex != wantPushRules.CommitMessageNegativeRegex {
			messages = append(messages, fmt.Sprintf("commit_message_negative_regex (got: %q, wanted: %q)",
				gotPushRules.CommitMessageNegativeRegex, wantPushRules.CommitMessageNegativeRegex))
		}

		if gotPushRules.FileNameRegex != wantPushRules.FileNameRegex {
			messages = append(messages, fmt.Sprintf("file_name_regex (got: %q, wanted: %q)",
				gotPushRules.FileNameRegex, wantPushRules.FileNameRegex))
		}

		if gotPushRules.CommitCommitterCheck != wantPushRules.CommitCommitterCheck {
			messages = append(messages, fmt.Sprintf("commit_committer_check (got: %t, wanted: %t)",
				gotPushRules.CommitCommitterCheck, wantPushRules.CommitCommitterCheck))
		}

		if gotPushRules.DenyDeleteTag != wantPushRules.DenyDeleteTag {
			messages = append(messages, fmt.Sprintf("deny_delete_tag (got: %t, wanted: %t)",
				gotPushRules.DenyDeleteTag, wantPushRules.DenyDeleteTag))
		}

		if gotPushRules.MemberCheck != wantPushRules.MemberCheck {
			messages = append(messages, fmt.Sprintf("member_check (got: %t, wanted: %t)",
				gotPushRules.MemberCheck, wantPushRules.MemberCheck))
		}

		if gotPushRules.PreventSecrets != wantPushRules.PreventSecrets {
			messages = append(messages, fmt.Sprintf("prevent_secrets (got: %t, wanted: %t)",
				gotPushRules.PreventSecrets, wantPushRules.PreventSecrets))
		}

		if gotPushRules.RejectUnsignedCommits != wantPushRules.RejectUnsignedCommits {
			messages = append(messages, fmt.Sprintf("reject_unsigned_commits (got: %t, wanted: %t)",
				gotPushRules.RejectUnsignedCommits, wantPushRules.RejectUnsignedCommits))
		}

		if gotPushRules.MaxFileSize != wantPushRules.MaxFileSize {
			messages = append(messages, fmt.Sprintf("max_file_size (got: %d, wanted: %d)",
				gotPushRules.MaxFileSize, wantPushRules.MaxFileSize))
		}

		if len(messages) > 0 {
			return fmt.Errorf("unexpected push_rules:\n\t- %s", strings.Join(messages, "\n\t- "))
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

func testAccGitlabProjectConfigWithoutDefaultBranchProtection(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
  name = "foogroup-%d"
  path = "foogroup-%d"
  default_branch_protection = 0
  visibility_level = "public"
}

resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"
  namespace_id = "${gitlab_group.foo.id}"
}
	`, rInt, rInt, rInt)
}

func testAccGitlabProjectConfigWithoutDefaultBranchProtectionInitializeReadme(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
  name = "foogroup2-%d"
  path = "foogroup2-%d"
  default_branch_protection = 0
  visibility_level = "public"
}

resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"
  namespace_id = "${gitlab_group.foo.id}"
  initialize_with_readme = true
}
	`, rInt, rInt, rInt)
}

func testAccGitlabProjectTransferBetweenGroupsBefore(rInt int) string {
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

resource "gitlab_project_variable" "foo" {
  project = "${gitlab_project.foo.id}"

  key = "FOO"
  value = "${gitlab_project.foo.path_with_namespace}"
}
	`, rInt, rInt, rInt)
}

func testAccGitlabProjectTransferBetweenGroupsAfter(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
  name = "foogroup-%d"
  path = "foogroup-%d"
  visibility_level = "public"
}

resource "gitlab_group" "foo2" {
  name = "foo2group-%d"
  path = "foo2group-%d"
  visibility_level = "public"
}

resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"
  namespace_id = "${gitlab_group.foo2.id}"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_project_variable" "foo" {
  project = "${gitlab_project.foo.id}"

  key = "FOO"
  value = "${gitlab_project.foo.path_with_namespace}"
}
	`, rInt, rInt, rInt, rInt, rInt)
}

func testAccGitlabProjectConfigDefaultBranch(rInt int, defaultBranch string) string {
	defaultBranchStatement := ""

	if len(defaultBranch) > 0 {
		defaultBranchStatement = fmt.Sprintf("default_branch = \"%s\"", defaultBranch)
	}

	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  path = "foo.%d"
  description = "Terraform acceptance tests"

  %s

  # NOTE: replaces by topics
  # tags = [
  # "tag1",
  # ]

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
  merge_method = "ff"
  only_allow_merge_if_pipeline_succeeds = true
  only_allow_merge_if_all_discussions_are_resolved = true
  squash_option = "default_off"
  pages_access_level = "public"
  allow_merge_on_skipped_pipeline = false
  ci_config_path = ".gitlab-ci.yml@mynamespace/myproject"
  resolve_outdated_diff_discussions = true
  analytics_access_level = "enabled"
  auto_cancel_pending_pipelines = "enabled"
  auto_devops_deploy_strategy = "continuous"
  auto_devops_enabled = true
  autoclose_referenced_issues = true
  build_git_strategy = "fetch"
  build_timeout = 42 * 60
  builds_access_level = "enabled"
  container_expiration_policy {
	enabled = true
  	cadence = "1month"
  }
  container_registry_access_level = "enabled"
  emails_disabled = true
  forking_access_level = "enabled"
  issues_access_level = "enabled"
  merge_requests_access_level = "enabled"
  operations_access_level = "enabled"
  public_builds = false
  repository_access_level = "enabled"
  repository_storage = "default"
  security_and_compliance_access_level = "enabled"
  snippets_access_level = "enabled"
  topics = ["foo", "bar"]
  wiki_access_level = "enabled"
  squash_commit_template = "hello squash"
  merge_commit_template = "hello merge"
  ci_default_git_depth = 42
}
	`, rInt, rInt, defaultBranchStatement)
}

func testAccGitlabProjectConfigDefaultBranchSkipFunc(project *gitlab.Project, defaultBranch string) func() (bool, error) {
	return func() (bool, error) {
		// Commit data
		commitMessage := "Initial Commit"
		commitFile := "file.txt"
		commitFileAction := gitlab.FileCreate
		commitActions := []*gitlab.CommitActionOptions{
			{
				Action:   &commitFileAction,
				FilePath: &commitFile,
				Content:  &commitMessage,
			},
		}
		options := &gitlab.CreateCommitOptions{
			Branch:        &defaultBranch,
			CommitMessage: &commitMessage,
			Actions:       commitActions,
		}

		_, _, err := testGitlabClient.Commits.CreateCommit(project.ID, options)

		return false, err
	}
}

func testAccGitlabProjectConfig(rInt int) string {
	return testAccGitlabProjectConfigDefaultBranch(rInt, "")
}

func testAccGitlabProjectUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  path = "foo.%d"
  description = "Terraform acceptance tests!"

  # NOTE: replaces by topics
  # tags = [
  # "tag1",
  # "tag2"
  # ]

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
  merge_method = "ff"
  only_allow_merge_if_pipeline_succeeds = true
  only_allow_merge_if_all_discussions_are_resolved = true
  squash_option = "default_on"
  allow_merge_on_skipped_pipeline = true

  request_access_enabled = false
  issues_enabled = false
  merge_requests_enabled = false
  pipelines_enabled = false
  approvals_before_merge = 0
  wiki_enabled = false
  snippets_enabled = false
  container_registry_enabled = false
  lfs_enabled = false
  shared_runners_enabled = false
  archived = true
  packages_enabled = false
  pages_access_level = "disabled"
  ci_forward_deployment_enabled = false
  merge_pipelines_enabled = false
  merge_trains_enabled = false
  resolve_outdated_diff_discussions = false
  analytics_access_level = "disabled"
  auto_cancel_pending_pipelines = "disabled"
  auto_devops_deploy_strategy = "manual"
  auto_devops_enabled = false
  autoclose_referenced_issues = false
  build_git_strategy = "fetch"
  build_timeout = 10 * 60
  builds_access_level = "disabled"
  container_expiration_policy {
	enabled = true
  	cadence = "3month"
  }
  container_registry_access_level = "disabled"
  emails_disabled = false
  forking_access_level = "disabled"
  issues_access_level = "disabled"
  merge_requests_access_level = "disabled"
  operations_access_level = "disabled"
  public_builds = false
  repository_access_level = "disabled"
  repository_storage = "default"
  security_and_compliance_access_level = "disabled"
  snippets_access_level = "disabled"
  topics = []
  wiki_access_level = "disabled"
  squash_commit_template = "goodby squash"
  merge_commit_template = "goodby merge"
  ci_default_git_depth = 84
}
	`, rInt, rInt)
}

func testAccGitlabProjectConfigInitializeWithReadme(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  path = "foo.%d"
  description = "Terraform acceptance tests"
  initialize_with_readme = true

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
	`, rInt, rInt)
}

func testAccGitlabProjectConfigInitializeWithoutReadme(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  path = "foo.%d"
  description = "Terraform acceptance tests"
  initialize_with_readme = false

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
	`, rInt, rInt)
}

func testAccGitlabProjectConfigImportURL(rInt int, importURL string) string {
	return fmt.Sprintf(`
resource "gitlab_project" "imported" {
  name = "imported-%d"
  default_branch = "main"
  import_url = "%s"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
`, rInt, importURL)
}

func testAccGitlabProjectConfigImportURLMirror(rInt int, importURL string) string {
	return fmt.Sprintf(`
resource "gitlab_project" "imported" {
  name = "imported-%d"
  default_branch = "main"
  import_url = "%s"
  mirror = true
  mirror_trigger_builds = true
  mirror_overwrites_diverged_branches = true
  only_mirror_protected_branches = true

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
`, rInt, importURL)
}

func testAccGitlabProjectConfigImportURLMirrorDisabledOptionals(rInt int, importURL string) string {
	return fmt.Sprintf(`
resource "gitlab_project" "imported" {
  name = "imported-%d"
  default_branch = "main"
  import_url = "%s"
  mirror = true
  mirror_trigger_builds = false
  mirror_overwrites_diverged_branches = false
  only_mirror_protected_branches = false

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
`, rInt, importURL)
}

func testAccGitlabProjectConfigImportURLMirrorDisabled(rInt int, importURL string) string {
	return fmt.Sprintf(`
resource "gitlab_project" "imported" {
  name = "imported-%d"
  default_branch = "main"
  import_url = "%s"
  mirror = false
  mirror_trigger_builds = false
  mirror_overwrites_diverged_branches = false
  only_mirror_protected_branches = false

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
`, rInt, importURL)
}

func testAccGitlabProjectConfigPushRules(rInt int, pushRules string) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%[1]d"
  path = "foo.%[1]d"
  description = "Terraform acceptance tests"

  push_rules {
%[2]s
  }

  resolve_outdated_diff_discussions = true
  analytics_access_level = "enabled"
  auto_cancel_pending_pipelines = "enabled"
  auto_devops_deploy_strategy = "continuous"
  auto_devops_enabled = true
  autoclose_referenced_issues = true
  build_git_strategy = "fetch"
  build_timeout = 42 * 60
  builds_access_level = "enabled"
  container_expiration_policy {
	enabled = true
  	cadence = "1month"
  }
  container_registry_access_level = "enabled"
  emails_disabled = true
  forking_access_level = "enabled"
  issues_access_level = "enabled"
  merge_requests_access_level = "enabled"
  operations_access_level = "enabled"
  public_builds = false
  repository_access_level = "enabled"
  repository_storage = "default"
  security_and_compliance_access_level = "enabled"
  snippets_access_level = "enabled"
  topics = ["foo", "bar"]
  wiki_access_level = "enabled"
  squash_commit_template = "hello squash"
  merge_commit_template = "hello merge"

  # So that acceptance tests can be run in a gitlab organization with no billing.
  visibility_level = "public"
}
	`, rInt, pushRules)
}

func testAccGitlabProjectConfigTemplateName(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "template-name" {
  name = "template-name-%d"
  path = "template-name.%d"
  description = "Terraform acceptance tests"
  template_name = "rails"
  default_branch = "master"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
	`, rInt, rInt)
}

// 2020-09-07: Currently Gitlab (version 13.3.6 ) doesn't allow in admin API
// ability to set a group as instance level templates.
// To test resource_gitlab_project_test template features we add
// group, project myrails and admin settings directly in scripts/healthcheck-and-setup.sh
// Once Gitlab add admin template in API we could manage group/project/settings
// directly in tests like TestAccGitlabProject_basic.
func testAccGitlabProjectConfigTemplateNameCustom(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "template-name-custom" {
  name = "template-name-custom-%d"
  path = "template-name-custom.%d"
  description = "Terraform acceptance tests"
  template_name = "myrails"
  use_custom_template = true
  default_branch = "master"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
	`, rInt, rInt)
}

func testAccGitlabProjectConfigTemplateProjectID(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "template-id" {
  name = "template-id-%d"
  path = "template-id.%d"
  description = "Terraform acceptance tests"
  template_project_id = 999
  use_custom_template = true
  default_branch = "master"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
	`, rInt, rInt)
}

func testAccCheckMutualExclusiveNameAndID(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "template-mutual-exclusive" {
  name = "template-mutual-exclusive-%d"
  path = "template-mutual-exclusive.%d"
  description = "Terraform acceptance tests"
  template_name = "rails"
  template_project_id = 999
  use_custom_template = true
  default_branch = "master"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
	`, rInt, rInt)
}

func testAccGitlabProjectConfigIssueMergeRequestTemplates(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  path = "foo.%d"
  description = "Terraform acceptance tests"
  issues_template = "foo"
  merge_requests_template = "bar"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
	`, rInt, rInt)
}

func testAccGitlabProjectConfigArchiveOnDestroy(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  path = "foo.%d"
  description = "Terraform acceptance tests"
  archive_on_destroy = true
  archived = false

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
	`, rInt, rInt)
}

func testAccGitLabProjectMergePipelinesEnabled(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  path = "foo.%d"
  description = "Terraform acceptance tests"
  merge_pipelines_enabled = true

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
	`, rInt, rInt)
}

func testAccGitLabProjectMergeTrainsEnabled(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  path = "foo.%d"
  description = "Terraform acceptance tests"
  merge_pipelines_enabled = true
  merge_trains_enabled = true

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
	`, rInt, rInt)
}

func testAccGitlabProjectConfigEE(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  path = "foo.%d"
  description = "Terraform acceptance tests"
  default_branch = "main"

  # NOTE: replaces by topics
  # tags = [
  # "tag1",
  # ]

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
  merge_method = "ff"
  only_allow_merge_if_pipeline_succeeds = true
  only_allow_merge_if_all_discussions_are_resolved = true
  squash_option = "default_off"
  pages_access_level = "public"
  allow_merge_on_skipped_pipeline = false
  ci_config_path = ".gitlab-ci.yml@mynamespace/myproject"
  resolve_outdated_diff_discussions = true
  analytics_access_level = "enabled"
  auto_cancel_pending_pipelines = "enabled"
  auto_devops_deploy_strategy = "continuous"
  auto_devops_enabled = true
  autoclose_referenced_issues = true
  build_git_strategy = "fetch"
  build_timeout = 42 * 60
  builds_access_level = "enabled"
  container_expiration_policy {
	enabled = true
  	cadence = "1month"
  }
  container_registry_access_level = "enabled"
  emails_disabled = true
  forking_access_level = "enabled"
  issues_access_level = "enabled"
  merge_requests_access_level = "enabled"
  operations_access_level = "enabled"
  public_builds = false
  repository_access_level = "enabled"
  repository_storage = "default"
  security_and_compliance_access_level = "enabled"
  snippets_access_level = "enabled"
  topics = ["foo", "bar"]
  wiki_access_level = "enabled"
  squash_commit_template = "hello squash"
  merge_commit_template = "hello merge"
  ci_default_git_depth = 42

  # EE features
  approvals_before_merge = 2
  external_authorization_classification_label = "test"
  requirements_access_level = "enabled"
  # are tested in separate test case
  # mirror_trigger_builds = true
  # mirror = true
}
	`, rInt, rInt)
}
