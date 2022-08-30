package provider

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mitchellh/hashstructure"
	"github.com/xanzy/go-gitlab"
)

// Schemas

// WARN: go-gitlab may not be up-to-date with Gitlab exposed options
// https://docs.gitlab.com/ee/api/groups.html#list-a-groups-projects
// https://docs.gitlab.com/ee/api/projects.html#list-all-projects

// Helper functions
func flattenProjectPermissions(permissions *gitlab.Permissions) []map[string]interface{} {
	m := make(map[string]interface{}, 2)
	if permissions != nil {
		if permissions.ProjectAccess != nil {
			m["project_access"] = map[string]int{
				"access_level":       int(permissions.ProjectAccess.AccessLevel),
				"notification_level": int(permissions.ProjectAccess.NotificationLevel),
			}
		}
		if permissions.GroupAccess != nil {
			m["group_access"] = map[string]int{
				"access_level":       int(permissions.GroupAccess.AccessLevel),
				"notification_level": int(permissions.GroupAccess.NotificationLevel),
			}
		}
	}
	return []map[string]interface{}{m}
}

func flattenProjectNamespace(namespace *gitlab.ProjectNamespace) (values []map[string]interface{}) {
	if namespace != nil {
		values = []map[string]interface{}{
			{
				"id":        namespace.ID,
				"name":      namespace.Name,
				"path":      namespace.Path,
				"kind":      namespace.Kind,
				"full_path": namespace.FullPath,
			},
		}
	}
	return values
}

func flattenProjectLinks(links *gitlab.Links) (values map[string]string) {
	if links != nil {
		values = map[string]string{
			"self":           links.Self,
			"issues":         links.Issues,
			"merge_requests": links.MergeRequests,
			"repo_branches":  links.RepoBranches,
			"labels":         links.Labels,
			"events":         links.Events,
			"members":        links.Members,
		}
	}
	return values
}

func flattenForkedFromProject(forked *gitlab.ForkParent) (values []map[string]interface{}) {
	if forked != nil {
		values = []map[string]interface{}{
			{
				"http_url_to_repo":    forked.HTTPURLToRepo,
				"id":                  forked.ID,
				"name":                forked.Name,
				"name_with_namespace": forked.NameWithNamespace,
				"path":                forked.Path,
				"path_with_namespace": forked.PathWithNamespace,
				"web_url":             forked.WebURL,
			},
		}
	}
	return values
}

func flattenGitlabBasicUser(user *gitlab.User) (values []map[string]interface{}) {
	if user != nil {
		values = []map[string]interface{}{
			{
				"id":          user.ID,
				"username":    user.Username,
				"name":        user.Name,
				"state":       user.State,
				"avatar_url":  user.AvatarURL,
				"website_url": user.WebsiteURL,
			},
		}
	}
	return values
}

func flattenProjects(projects []*gitlab.Project) (values []map[string]interface{}) {
	if projects != nil { // nolint // TODO: Resolve this golangci-lint issue: S1031: unnecessary nil check around range (gosimple)
		for _, project := range projects {
			v := map[string]interface{}{
				"id":                                    project.ID,
				"description":                           project.Description,
				"default_branch":                        project.DefaultBranch,
				"public":                                project.Public,
				"visibility":                            string(project.Visibility),
				"ssh_url_to_repo":                       project.SSHURLToRepo,
				"http_url_to_repo":                      project.HTTPURLToRepo,
				"web_url":                               project.WebURL,
				"readme_url":                            project.ReadmeURL,
				"tag_list":                              project.TagList,
				"owner":                                 flattenGitlabBasicUser(project.Owner),
				"name":                                  project.Name,
				"name_with_namespace":                   project.NameWithNamespace,
				"path":                                  project.Path,
				"path_with_namespace":                   project.PathWithNamespace,
				"issues_enabled":                        project.IssuesEnabled,
				"open_issues_count":                     project.OpenIssuesCount,
				"merge_requests_enabled":                project.MergeRequestsEnabled,
				"approvals_before_merge":                project.ApprovalsBeforeMerge,
				"jobs_enabled":                          project.JobsEnabled,
				"wiki_enabled":                          project.WikiEnabled,
				"snippets_enabled":                      project.SnippetsEnabled,
				"resolve_outdated_diff_discussions":     project.ResolveOutdatedDiffDiscussions,
				"container_registry_enabled":            project.ContainerRegistryEnabled,
				"created_at":                            project.CreatedAt.String(),
				"last_activity_at":                      project.LastActivityAt.String(),
				"creator_id":                            project.CreatorID,
				"namespace":                             flattenProjectNamespace(project.Namespace),
				"import_status":                         project.ImportStatus,
				"import_error":                          project.ImportError,
				"permissions":                           flattenProjectPermissions(project.Permissions),
				"archived":                              project.Archived,
				"avatar_url":                            project.AvatarURL,
				"shared_runners_enabled":                project.SharedRunnersEnabled,
				"forks_count":                           project.ForksCount,
				"star_count":                            project.StarCount,
				"runners_token":                         project.RunnersToken,
				"public_builds":                         project.PublicBuilds,
				"only_allow_merge_if_pipeline_succeeds": project.OnlyAllowMergeIfPipelineSucceeds,
				"only_allow_merge_if_all_discussions_are_resolved": project.OnlyAllowMergeIfAllDiscussionsAreResolved,
				"allow_merge_on_skipped_pipeline":                  project.AllowMergeOnSkippedPipeline,
				"lfs_enabled":                                      project.LFSEnabled,
				"request_access_enabled":                           project.RequestAccessEnabled,
				"merge_method":                                     project.MergeMethod,
				"forked_from_project":                              flattenForkedFromProject(project.ForkedFromProject),
				"mirror":                                           project.Mirror,
				"mirror_user_id":                                   project.MirrorUserID,
				"mirror_trigger_builds":                            project.MirrorTriggerBuilds,
				"only_mirror_protected_branches":                   project.OnlyMirrorProtectedBranches,
				"mirror_overwrites_diverged_branches":              project.MirrorOverwritesDivergedBranches,
				"shared_with_groups":                               flattenSharedWithGroupsOptions(project),
				"statistics":                                       project.Statistics,
				"_links":                                           flattenProjectLinks(project.Links),
				"ci_config_path":                                   project.CIConfigPath,
				"custom_attributes":                                project.CustomAttributes,
				"packages_enabled":                                 project.PackagesEnabled,
				"build_coverage_regex":                             project.BuildCoverageRegex,
				"ci_forward_deployment_enabled":                    project.CIForwardDeploymentEnabled,
				"merge_pipelines_enabled":                          project.MergePipelinesEnabled,
				"merge_trains_enabled":                             project.MergeTrainsEnabled,
				"analytics_access_level":                           string(project.AnalyticsAccessLevel),
				"auto_cancel_pending_pipelines":                    project.AutoCancelPendingPipelines,
				"auto_devops_deploy_strategy":                      project.AutoDevopsDeployStrategy,
				"auto_devops_enabled":                              project.AutoDevopsEnabled,
				"autoclose_referenced_issues":                      project.AutocloseReferencedIssues,
				"build_git_strategy":                               project.BuildGitStrategy,
				"build_timeout":                                    project.BuildTimeout,
				"builds_access_level":                              string(project.BuildsAccessLevel),
				"container_expiration_policy":                      flattenContainerExpirationPolicy(project.ContainerExpirationPolicy),
				"container_registry_access_level":                  string(project.ContainerRegistryAccessLevel),
				"emails_disabled":                                  project.EmailsDisabled,
				"external_authorization_classification_label":      project.ExternalAuthorizationClassificationLabel,
				"forking_access_level":                             string(project.ForkingAccessLevel),
				"issues_access_level":                              string(project.IssuesAccessLevel),
				"merge_requests_access_level":                      string(project.MergeRequestsAccessLevel),
				"operations_access_level":                          string(project.OperationsAccessLevel),
				"repository_access_level":                          string(project.RepositoryAccessLevel),
				"repository_storage":                               project.RepositoryStorage,
				"requirements_access_level":                        string(project.RequirementsAccessLevel),
				"security_and_compliance_access_level":             string(project.SecurityAndComplianceAccessLevel),
				"snippets_access_level":                            string(project.SnippetsAccessLevel),
				"topics":                                           project.Topics,
				"wiki_access_level":                                string(project.WikiAccessLevel),
				"squash_commit_template":                           project.SquashCommitTemplate,
				"merge_commit_template":                            project.MergeCommitTemplate,
				"ci_default_git_depth":                             project.CIDefaultGitDepth,
			}
			values = append(values, v)
		}
	}
	return values
}

var _ = registerDataSource("gitlab_projects", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_projects`" + ` data source allows details of multiple projects to be retrieved. Optionally filtered by the set attributes.

-> This data source supports all available filters exposed by the xanzy/go-gitlab package, which might not expose all available filters exposed by the Gitlab APIs.

-> The [owner sub-attributes](#nestedobjatt--projects--owner) are only populated if the Gitlab token used has an administrator scope.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/projects.html#list-all-projects)`,

		ReadContext: dataSourceGitlabProjectsRead,

		Schema: map[string]*schema.Schema{
			"max_queryable_pages": {
				Description: "The maximum number of project results pages that may be queried. Prevents overloading your Gitlab instance in case of a misconfiguration.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
			},
			"group_id": {
				Description: "The ID of the group owned by the authenticated user to look projects for within. Cannot be used with `min_access_level`, `with_programming_language` or `statistics`.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"page": {
				Description: "The first page to begin the query on.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
			},
			"per_page": {
				Description:  "The number of results to return per page.",
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      20,
				ValidateFunc: validation.IntAtMost(100),
			},
			"archived": {
				Description: "Limit by archived status.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"order_by": {
				Description: "Return projects ordered by `id`, `name`, `path`, `created_at`, `updated_at`, or `last_activity_at` fields. Default is `created_at`.",
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"id",
					"name",
					"username",
					"created_at",
					"updated_at"}, true),
			},
			"sort": {
				Description:  "Return projects sorted in `asc` or `desc` order. Default is `desc`.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"desc", "asc"}, true),
			},
			"search": {
				Description: "Return list of authorized projects matching the search criteria.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"simple": {
				Description: "Return only the ID, URL, name, and path of each project.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"owned": {
				Description: "Limit by projects owned by the current user.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"starred": {
				Description: "Limit by projects starred by the current user.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"visibility": {
				Description: "Limit by visibility `public`, `internal`, or `private`.",
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"public",
					"private",
					"internal"}, true),
			},
			"with_issues_enabled": {
				Description: "Limit by projects with issues feature enabled. Default is `false`.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"with_merge_requests_enabled": {
				Description: "Limit by projects with merge requests feature enabled. Default is `false`.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"with_custom_attributes": {
				Description: "Include custom attributes in response _(admins only)_.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"membership": {
				Description: "Limit by projects that the current user is a member of.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"min_access_level": {
				Description: "Limit to projects where current user has at least this access level, refer to the [official documentation](https://docs.gitlab.com/ee/api/members.html) for values. Cannot be used with `group_id`.",
				Type:        schema.TypeInt,
				Optional:    true,
				ValidateFunc: validation.IntInSlice([]int{
					int(gitlab.GuestPermissions),
					int(gitlab.DeveloperPermissions),
					int(gitlab.MaintainerPermissions),
					int(gitlab.MasterPermissions),
				}),
				ConflictsWith: []string{
					"group_id",
				},
			},
			"with_programming_language": {
				Description: "Limit by projects which use the given programming language. Cannot be used with `group_id`.",
				Type:        schema.TypeString,
				Optional:    true,
				ConflictsWith: []string{
					"group_id",
				},
			},
			"statistics": {
				Description: "Include project statistics. Cannot be used with `group_id`.",
				Type:        schema.TypeBool,
				Optional:    true,
				ConflictsWith: []string{
					"group_id",
				},
			},
			"with_shared": {
				Description: "Include projects shared to this group. Default is `true`. Needs `group_id`.",
				Type:        schema.TypeBool,
				Optional:    true,
				ConflictsWith: []string{
					"statistics",
					"with_programming_language",
					"min_access_level",
				},
			},
			"include_subgroups": {
				Description: "Include projects in subgroups of this group. Default is `false`. Needs `group_id`.",
				Type:        schema.TypeBool,
				Optional:    true,
				ConflictsWith: []string{
					"statistics",
					"with_programming_language",
					"min_access_level",
				},
			},
			"projects": {
				Description: "A list containing the projects matching the supplied arguments",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The ID of the project.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"description": {
							Description: "The description of the project.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"default_branch": {
							Description: "The default branch name of the project.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"public": {
							Description: "Whether the project is public.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"visibility": {
							Description: "The visibility of the project.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"ssh_url_to_repo": {
							Description: "The SSH clone URL of the project.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"http_url_to_repo": {
							Description: "The HTTP clone URL of the project.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"web_url": {
							Description: "The web url of the project.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"readme_url": {
							Description: "The remote url of the project.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"tag_list": {
							Description: "A set of the project topics (formerly called \"project tags\").",
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"owner": {
							Description: "The owner of the project, due to Terraform aggregate types limitations, this field's attributes are accessed with the `owner.0` prefix. Structure is documented below.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Description: "The ID of the owner.",
										Type:        schema.TypeInt,
										Computed:    true,
									},
									"username": {
										Description: "The username of the owner.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"name": {
										Description: "The name of the owner.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"state": {
										Description: "The state of the owner.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"avatar_url": {
										Description: "The avatar url of the owner.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"website_url": {
										Description: "The website url of the owner.",
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
						"name": {
							Description: "The name of the project.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name_with_namespace": {
							Description: "In `group / subgroup / project` or `user / project` format.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"path": {
							Description: "The path of the project.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"path_with_namespace": {
							Description: "In `group/subgroup/project` or `user/project` format.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"issues_enabled": {
							Description: "Whether issues are enabled for the project.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"open_issues_count": {
							Description: "The number of open issies for the project.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"merge_requests_enabled": {
							Description: "Whether merge requests are enabled for the project.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"approvals_before_merge": {
							Description: "The numbers of approvals needed in a merge requests.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"jobs_enabled": {
							Description: "Whether pipelines are enabled for the project.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"wiki_enabled": {
							Description: "Whether wiki is enabled for the project.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"snippets_enabled": {
							Description: "Whether snippets are enabled for the project.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"resolve_outdated_diff_discussions": {
							Description: "Whether resolve_outdated_diff_discussions is enabled for the project",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"container_registry_enabled": {
							Description: "Whether the container registry is enabled for the project.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"created_at": {
							Description: "Creation time for the project.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"last_activity_at": {
							Description: "Last activirty time for the project.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"creator_id": {
							Description: "Creator ID for the project.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"namespace": {
							Description: "Namespace of the project (parent group/s).",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Description: "The ID of the namespace.",
										Type:        schema.TypeInt,
										Computed:    true,
									},
									"name": {
										Description: "The name of the namespace.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"path": {
										Description: "The path of the namespace.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"kind": {
										Description: "The kind of the namespace.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"full_path": {
										Description: "The full path of the namespace.",
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
						"import_status": {
							Description: "The import status of the project.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"import_error": {
							Description: "The import error, if it exists, for the project.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"permissions": {
							Description: "Permissions for the project.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"project_access": {
										Description: "Project access level.",
										Type:        schema.TypeMap,
										Computed:    true,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"group_access": {
										Description: "Group access level.",
										Type:        schema.TypeMap,
										Computed:    true,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
								},
							},
						},
						"archived": {
							Description: "Whether the project is archived.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"avatar_url": {
							Description: "The avatar url of the project.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"shared_runners_enabled": {
							Description: "Whether shared runners are enabled for the project.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"forks_count": {
							Description: "The number of forks of the project.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"star_count": {
							Description: "The number of stars on the project.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"runners_token": {
							Description: "The runners token for the project.",
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
						},
						"public_builds": {
							Description: "Whether public builds are enabled for the project.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"only_allow_merge_if_pipeline_succeeds": {
							Description: "Whether only_allow_merge_if_pipeline_succeeds is enabled for the project.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"only_allow_merge_if_all_discussions_are_resolved": {
							Description: "Whether only_allow_merge_if_all_discussions_are_resolved is enabled for the project.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_merge_on_skipped_pipeline": {
							Description: "Whether allow_merge_on_skipped_pipeline is enabled for the project.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"lfs_enabled": {
							Description: "Whether LFS (large file storage) is enabled for the project.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"request_access_enabled": {
							Description: "Whether requesting access is enabled for the project.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"merge_method": {
							Description: "Merge method for the project.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"forked_from_project": {
							Description: "Present if the project is a fork. Contains information about the upstream project.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"http_url_to_repo": {
										Description: "The HTTP clone URL of the upstream project.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"id": {
										Description: "The ID of the upstream project.",
										Type:        schema.TypeInt,
										Computed:    true,
									},
									"name": {
										Description: "The name of the upstream project.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"name_with_namespace": {
										Description: "In `group / subgroup / project` or `user / project` format.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"path": {
										Description: "The path of the upstream project.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"path_with_namespace": {
										Description: "In `group/subgroup/project` or `user/project` format.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"web_url": {
										Description: "The web url of the upstream project.",
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
						"mirror": {
							Description: "Whether the pull mirroring is enabled for the project.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"mirror_user_id": {
							Description: "The mirror user ID for the project.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"mirror_trigger_builds": {
							Description: "Whether pull mirrororing triggers builds for the project.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"only_mirror_protected_branches": {
							Description: "Whether only_mirror_protected_branches is enabled for the project.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"mirror_overwrites_diverged_branches": {
							Description: "Wther mirror_overwrites_diverged_branches is enabled for the project.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"shared_with_groups": {
							Description: "Groups the the project is shared with.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_id": {
										Description: "The group ID.",
										Type:        schema.TypeInt,
										Computed:    true,
									},
									"group_access_level": {
										Description: "The group access level.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"group_name": {
										Description: "The group name.",
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
						"statistics": {
							Description: "Statistics for the project.",
							Type:        schema.TypeMap,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"_links": {
							Description: "Links for the project.",
							Type:        schema.TypeMap,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"ci_config_path": {
							Description: "CI config file path for the project.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"custom_attributes": {
							Description: "Custom attributes for the project.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeMap,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
						"packages_enabled": {
							Description: "Whether packages are enabled for the project.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"build_coverage_regex": {
							Description: "Build coverage regex for the project.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"ci_forward_deployment_enabled": {
							Description: "When a new deployment job starts, skip older deployment jobs that are still pending.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"merge_pipelines_enabled": {
							Description: "Enable or disable merge pipelines.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"merge_trains_enabled": {
							Description: "Enable or disable merge trains.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"analytics_access_level": {
							Description: fmt.Sprintf("Set the analytics access level. Valid values are %s.", renderValueListForDocs(validProjectAccessLevels)),
							Type:        schema.TypeString,
							Computed:    true,
						},
						"auto_cancel_pending_pipelines": {
							Description: "Auto-cancel pending pipelines. This isn’t a boolean, but enabled/disabled.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"auto_devops_deploy_strategy": {
							Description: fmt.Sprintf("Auto Deploy strategy. Valid values are %s.", renderValueListForDocs(validProjectAutoDevOpsDeployStrategyValues)),
							Type:        schema.TypeString,
							Computed:    true,
						},
						"auto_devops_enabled": {
							Description: "Enable Auto DevOps for this project.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"autoclose_referenced_issues": {
							Description: "Set whether auto-closing referenced issues on default branch.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"build_git_strategy": {
							Description: "The Git strategy. Defaults to fetch.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"build_timeout": {
							Description: "The maximum amount of time, in seconds, that a job can run.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"builds_access_level": {
							Description: fmt.Sprintf("Set the builds access level. Valid values are %s.", renderValueListForDocs(validProjectAccessLevels)),
							Type:        schema.TypeString,
							Computed:    true,
						},
						"container_expiration_policy": {
							Description: "Set the image cleanup policy for this project. **Note**: this field is sometimes named `container_expiration_policy_attributes` in the GitLab Upstream API.",
							Type:        schema.TypeList,
							Elem:        containerExpirationPolicyAttributesSchema,
							Computed:    true,
						},
						"container_registry_access_level": {
							Description: fmt.Sprintf("Set visibility of container registry, for this project. Valid values are %s.", renderValueListForDocs(validProjectAccessLevels)),
							Type:        schema.TypeString,
							Computed:    true,
						},
						"emails_disabled": {
							Description: "Disable email notifications.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"external_authorization_classification_label": {
							Description: "The classification label for the project.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"forking_access_level": {
							Description: fmt.Sprintf("Set the forking access level. Valid values are %s.", renderValueListForDocs(validProjectAccessLevels)),
							Type:        schema.TypeString,
							Computed:    true,
						},
						"issues_access_level": {
							Description: fmt.Sprintf("Set the issues access level. Valid values are %s.", renderValueListForDocs(validProjectAccessLevels)),
							Type:        schema.TypeString,
							Computed:    true,
						},
						"merge_requests_access_level": {
							Description: fmt.Sprintf("Set the merge requests access level. Valid values are %s.", renderValueListForDocs(validProjectAccessLevels)),
							Type:        schema.TypeString,
							Computed:    true,
						},
						"operations_access_level": {
							Description: fmt.Sprintf("Set the operations access level. Valid values are %s.", renderValueListForDocs(validProjectAccessLevels)),
							Type:        schema.TypeString,
							Computed:    true,
						},
						"repository_access_level": {
							Description: fmt.Sprintf("Set the repository access level. Valid values are %s.", renderValueListForDocs(validProjectAccessLevels)),
							Type:        schema.TypeString,
							Computed:    true,
						},
						"repository_storage": {
							Description: "	Which storage shard the repository is on. (administrator only)",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"requirements_access_level": {
							Description: fmt.Sprintf("Set the requirements access level. Valid values are %s.", renderValueListForDocs(validProjectAccessLevels)),
							Type:        schema.TypeString,
							Computed:    true,
						},
						"security_and_compliance_access_level": {
							Description: fmt.Sprintf("Set the security and compliance access level. Valid values are %s.", renderValueListForDocs(validProjectAccessLevels)),
							Type:        schema.TypeString,
							Computed:    true,
						},
						"snippets_access_level": {
							Description: fmt.Sprintf("Set the snippets access level. Valid values are %s.", renderValueListForDocs(validProjectAccessLevels)),
							Type:        schema.TypeString,
							Computed:    true,
						},
						"topics": {
							Description: "The list of topics for the project.",
							Type:        schema.TypeSet,
							Set:         schema.HashString,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
						},
						"wiki_access_level": {
							Description: fmt.Sprintf("Set the wiki access level. Valid values are %s.", renderValueListForDocs(validProjectAccessLevels)),
							Type:        schema.TypeString,
							Computed:    true,
						},
						"squash_commit_template": {
							Description: "Template used to create squash commit message in merge requests. (Introduced in GitLab 14.6.)",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"merge_commit_template": {
							Description: "Template used to create merge commit message in merge requests. (Introduced in GitLab 14.5.)",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"ci_default_git_depth": {
							Description: "Default number of revisions for shallow cloning.",
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
		},
	}
})

// CRUD methods

func dataSourceGitlabProjectsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	var projectList []*gitlab.Project

	// Permanent parameters

	page := d.Get("page").(int)
	perPage := d.Get("per_page").(int)
	maxQueryablePages := d.Get("max_queryable_pages").(int)

	// Conditional parameters
	// Only way I found to conditionally pass a search parameter to the List(Group/Project)Options
	// Json marshalling a complete List(Group/Project)Options JSON object had conversion issues with booleans
	var archivedPtr *bool
	var includeSubGroupsPtr *bool
	var membershipPtr *bool
	var minAccessLevelPtr *gitlab.AccessLevelValue
	var orderByPtr *string
	var ownedPtr *bool
	var searchPtr *string
	var simplePtr *bool
	var sortPtr *string
	var starredPtr *bool
	var statisticsPtr *bool
	var visibilityPtr *gitlab.VisibilityValue
	var withCustomAttributesPtr *bool
	var withIssuesEnabledPtr *bool
	var withMergeRequestsEnabledPtr *bool
	var withProgrammingLanguagePtr *string
	var withSharedPtr *bool

	// NOTE: `GetOkExists()` is deprecated, but until there is a replacement we need to use it.
	//       see https://github.com/hashicorp/terraform-plugin-sdk/pull/350#issuecomment-597888969

	// nolint:staticcheck // SA1019 ignore deprecated GetOkExists
	// lintignore: XR001 // TODO: replace with alternative for GetOkExists
	if data, ok := d.GetOkExists("archived"); ok {
		d := data.(bool)
		archivedPtr = &d
	}
	// nolint:staticcheck // SA1019 ignore deprecated GetOkExists
	// lintignore: XR001 // TODO: replace with alternative for GetOkExists
	if data, ok := d.GetOkExists("include_subgroups"); ok {
		d := data.(bool)
		includeSubGroupsPtr = &d
	}
	// nolint:staticcheck // SA1019 ignore deprecated GetOkExists
	// lintignore: XR001 // TODO: replace with alternative for GetOkExists
	if data, ok := d.GetOkExists("membership"); ok {
		d := data.(bool)
		membershipPtr = &d
	}
	if data, ok := d.GetOk("min_access_level"); ok {
		minAccessLevelPtr = gitlab.AccessLevel(gitlab.AccessLevelValue(data.(int)))
	}
	if data, ok := d.GetOk("order_by"); ok {
		d := data.(string)
		orderByPtr = &d
	}
	// nolint:staticcheck // SA1019 ignore deprecated GetOkExists
	// lintignore: XR001 // TODO: replace with alternative for GetOkExists
	if data, ok := d.GetOkExists("owned"); ok {
		d := data.(bool)
		ownedPtr = &d
	}
	if data, ok := d.GetOk("search"); ok {
		d := data.(string)
		searchPtr = &d
	}
	// nolint:staticcheck // SA1019 ignore deprecated GetOkExists
	// lintignore: XR001 // TODO: replace with alternative for GetOkExists
	if data, ok := d.GetOkExists("simple"); ok {
		d := data.(bool)
		simplePtr = &d
	}
	if data, ok := d.GetOk("sort"); ok {
		d := data.(string)
		sortPtr = &d
	}
	// nolint:staticcheck // SA1019 ignore deprecated GetOkExists
	// lintignore: XR001 // TODO: replace with alternative for GetOkExists
	if data, ok := d.GetOkExists("starred"); ok {
		d := data.(bool)
		starredPtr = &d
	}
	// nolint:staticcheck // SA1019 ignore deprecated GetOkExists
	// lintignore: XR001 // TODO: replace with alternative for GetOkExists
	if data, ok := d.GetOkExists("statistics"); ok {
		d := data.(bool)
		statisticsPtr = &d
	}
	if data, ok := d.GetOk("visibility"); ok {
		visibilityPtr = gitlab.Visibility(gitlab.VisibilityValue(data.(string)))
	}
	// nolint:staticcheck // SA1019 ignore deprecated GetOkExists
	// lintignore: XR001 // TODO: replace with alternative for GetOkExists
	if data, ok := d.GetOkExists("with_custom_attributes"); ok {
		d := data.(bool)
		withCustomAttributesPtr = &d
	}
	// nolint:staticcheck // SA1019 ignore deprecated GetOkExists
	// lintignore: XR001 // TODO: replace with alternative for GetOkExists
	if data, ok := d.GetOkExists("with_issues_enabled"); ok {
		d := data.(bool)
		withIssuesEnabledPtr = &d
	}
	// nolint:staticcheck // SA1019 ignore deprecated GetOkExists
	// lintignore: XR001 // TODO: replace with alternative for GetOkExists
	if data, ok := d.GetOkExists("with_merge_requests_enabled"); ok {
		d := data.(bool)
		withMergeRequestsEnabledPtr = &d
	}
	if data, ok := d.GetOk("with_programming_language"); ok {
		d := data.(string)
		withProgrammingLanguagePtr = &d
	}
	// nolint:staticcheck // SA1019 ignore deprecated GetOkExists
	// lintignore: XR001 // TODO: replace with alternative for GetOkExists
	if data, ok := d.GetOkExists("with_shared"); ok {
		d := data.(bool)
		withSharedPtr = &d
	}

	log.Printf("[DEBUG] Reading Gitlab projects")

	switch groupId, ok := d.GetOk("group_id"); ok {
	// GroupProject case
	case true:
		opts := &gitlab.ListGroupProjectsOptions{
			ListOptions: gitlab.ListOptions{
				Page:    page,
				PerPage: perPage,
			},
			Archived:                 archivedPtr,
			Visibility:               visibilityPtr,
			OrderBy:                  orderByPtr,
			Sort:                     sortPtr,
			Search:                   searchPtr,
			Simple:                   simplePtr,
			Owned:                    ownedPtr,
			Starred:                  starredPtr,
			WithIssuesEnabled:        withIssuesEnabledPtr,
			WithMergeRequestsEnabled: withMergeRequestsEnabledPtr,
			WithShared:               withSharedPtr,
			IncludeSubGroups:         includeSubGroupsPtr,
			WithCustomAttributes:     withCustomAttributesPtr,
		}

		for {
			projects, response, err := client.Groups.ListGroupProjects(groupId.(int), opts, gitlab.WithContext(ctx))
			if err != nil {
				return diag.FromErr(err)
			}
			projectList = append(projectList, projects...)
			opts.ListOptions.Page++

			log.Printf("[INFO] Currentpage: %d, Total: %d", response.CurrentPage, response.TotalPages)
			if response.CurrentPage == response.TotalPages || response.CurrentPage > maxQueryablePages {
				break
			}
		}
		h, err := hashstructure.Hash(*opts, nil)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(fmt.Sprintf("%d-%d", groupId.(int), h))
		if err := d.Set("projects", flattenProjects(projectList)); err != nil {
			return diag.FromErr(err)
		}

	// Project case
	default:
		opts := &gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{
				Page:    page,
				PerPage: perPage,
			},
			Archived:                 archivedPtr,
			OrderBy:                  orderByPtr,
			Sort:                     sortPtr,
			Search:                   searchPtr,
			Simple:                   simplePtr,
			Owned:                    ownedPtr,
			Membership:               membershipPtr,
			Starred:                  starredPtr,
			Statistics:               statisticsPtr,
			Visibility:               visibilityPtr,
			WithIssuesEnabled:        withIssuesEnabledPtr,
			WithMergeRequestsEnabled: withMergeRequestsEnabledPtr,
			MinAccessLevel:           minAccessLevelPtr,
			WithCustomAttributes:     withCustomAttributesPtr,
			WithProgrammingLanguage:  withProgrammingLanguagePtr,
		}

		for {
			projects, response, err := client.Projects.ListProjects(opts, nil, gitlab.WithContext(ctx))
			if err != nil {
				return diag.FromErr(err)
			}
			projectList = append(projectList, projects...)
			opts.ListOptions.Page++

			log.Printf("[INFO] Currentpage: %d, Total: %d", response.CurrentPage, response.TotalPages)
			if response.CurrentPage == response.TotalPages || response.CurrentPage > maxQueryablePages {
				break
			}
		}
		h, err := hashstructure.Hash(*opts, nil)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(fmt.Sprintf("%d", h))
		if err := d.Set("projects", flattenProjects(projectList)); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func flattenSharedWithGroupsOptions(project *gitlab.Project) []interface{} {
	var sharedWithGroupsList []interface{}

	for _, option := range project.SharedWithGroups {
		values := map[string]interface{}{
			"group_id":           option.GroupID,
			"group_access_level": accessLevelValueToName[gitlab.AccessLevelValue(option.GroupAccessLevel)],
			"group_name":         option.GroupName,
		}

		sharedWithGroupsList = append(sharedWithGroupsList, values)
	}

	return sharedWithGroupsList
}
