package gitlab

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	gitlab "github.com/xanzy/go-gitlab"
)

var resourceGitLabProjectSchema = map[string]*schema.Schema{
	"name": {
		Description: "The name of the project.",
		Type:        schema.TypeString,
		Required:    true,
	},
	"path": {
		Description: "The path of the repository.",
		Type:        schema.TypeString,
		Optional:    true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			if new == "" {
				return true
			}
			return old == new
		},
	},
	"path_with_namespace": {
		Description: "The path of the repository with namespace.",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"namespace_id": {
		Description: "The namespace (group or user) of the project. Defaults to your user.",
		Type:        schema.TypeInt,
		Optional:    true,
		Computed:    true,
	},
	"description": {
		Description: "A description of the project.",
		Type:        schema.TypeString,
		Optional:    true,
	},
	"default_branch": {
		Description: "The default branch for the project.",
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
	},
	"import_url": {
		Description: "Git URL to a repository to be imported.",
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
	},
	"request_access_enabled": {
		Description: "Allow users to request member access.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
	},
	"issues_enabled": {
		Description: "Enable issue tracking for the project.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
	},
	"merge_requests_enabled": {
		Description: "Enable merge requests for the project.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
	},
	"pipelines_enabled": {
		Description: "Enable pipelines for the project.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
	},
	"approvals_before_merge": {
		Description: "Number of merge request approvals required for merging. Default is 0.",
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     0,
	},
	"wiki_enabled": {
		Description: "Enable wiki for the project.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
	},
	"snippets_enabled": {
		Description: "Enable snippets for the project.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
	},
	"container_registry_enabled": {
		Description: "Enable container registry for the project.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
	},
	"lfs_enabled": {
		Description: "Enable LFS for the project.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
	},
	"visibility_level": {
		Description:  "Set to `public` to create a public project.",
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"private", "internal", "public"}, true),
		Default:      "private",
	},
	"merge_method": {
		Description:  "Set to `ff` to create fast-forward merges",
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"merge", "rebase_merge", "ff"}, true),
		Default:      "merge",
	},
	"only_allow_merge_if_pipeline_succeeds": {
		Description: "Set to true if you want allow merges only if a pipeline succeeds.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
	},
	"only_allow_merge_if_all_discussions_are_resolved": {
		Description: "Set to true if you want allow merges only if all discussions are resolved.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
	},
	"allow_merge_on_skipped_pipeline": {
		Description: "Set to true if you want to treat skipped pipelines as if they finished with success.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
	},
	"ssh_url_to_repo": {
		Description: "URL that can be provided to `git clone` to clone the",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"http_url_to_repo": {
		Description: "URL that can be provided to `git clone` to clone the",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"web_url": {
		Description: "URL that can be used to find the project in a browser.",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"runners_token": {
		Description: "Registration token to use during runner setup.",
		Type:        schema.TypeString,
		Computed:    true,
		Sensitive:   true,
	},
	"shared_runners_enabled": {
		Description: "Enable shared runners for this project.",
		Type:        schema.TypeBool,
		Optional:    true,
		Computed:    true,
	},
	"tags": {
		Description: "Tags (topics) of the project.",
		Type:        schema.TypeSet,
		Optional:    true,
		ForceNew:    false,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Set:         schema.HashString,
	},
	"archived": {
		Description: "Whether the project is in read-only mode (archived). Repositories can be archived/unarchived by toggling this parameter.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
	},
	"initialize_with_readme": {
		Description: "Create main branch with first commit containing a README.md file.",
		Type:        schema.TypeBool,
		Optional:    true,
	},
	"squash_option": {
		Description:  "Squash commits when merge request. Valid values are `never`, `always`, `default_on`, or `default_off`. The default value is `default_off`.",
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "default_off",
		ValidateFunc: validation.StringInSlice([]string{"never", "default_on", "always", "default_off"}, true),
	},
	"remove_source_branch_after_merge": {
		Description: "Enable `Delete source branch` option by default for all new merge requests.",
		Type:        schema.TypeBool,
		Optional:    true,
	},
	"packages_enabled": {
		Description: "Enable packages repository for the project.",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
	},
	"push_rules": {
		Description: "Push rules for the project.",
		Type:        schema.TypeList,
		MaxItems:    1,
		Optional:    true,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"author_email_regex": {
					Description: "All commit author emails must match this regex, e.g. `@my-company.com$`.",
					Type:        schema.TypeString,
					Optional:    true,
				},
				"branch_name_regex": {
					Description: "All branch names must match this regex, e.g. `(feature|hotfix)\\/*`.",
					Type:        schema.TypeString,
					Optional:    true,
				},
				"commit_message_regex": {
					Description: "All commit messages must match this regex, e.g. `Fixed \\d+\\..*`.",
					Type:        schema.TypeString,
					Optional:    true,
				},
				"commit_message_negative_regex": {
					Description: "No commit message is allowed to match this regex, for example `ssh\\:\\/\\/`.",
					Type:        schema.TypeString,
					Optional:    true,
				},
				"file_name_regex": {
					Description: "All commited filenames must not match this regex, e.g. `(jar|exe)$`.",
					Type:        schema.TypeString,
					Optional:    true,
				},
				"commit_committer_check": {
					Description: "Users can only push commits to this repository that were committed with one of their own verified emails.",
					Type:        schema.TypeBool,
					Optional:    true,
				},
				"deny_delete_tag": {
					Description: "Deny deleting a tag.",
					Type:        schema.TypeBool,
					Optional:    true,
				},
				"member_check": {
					Description: "Restrict commits by author (email) to existing GitLab users.",
					Type:        schema.TypeBool,
					Optional:    true,
				},
				"prevent_secrets": {
					Description: "GitLab will reject any files that are likely to contain secrets.",
					Type:        schema.TypeBool,
					Optional:    true,
				},
				"reject_unsigned_commits": {
					Description: "Reject commit when it’s not signed through GPG.",
					Type:        schema.TypeBool,
					Optional:    true,
				},
				"max_file_size": {
					Description:  "Maximum file size (MB).",
					Type:         schema.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntAtLeast(0),
				},
			},
		},
	},
	"template_name": {
		Description:   "When used without use_custom_template, name of a built-in project template. When used with use_custom_template, name of a custom project template. This option is mutually exclusive with `template_project_id`.",
		Type:          schema.TypeString,
		Optional:      true,
		ConflictsWith: []string{"template_project_id"},
		ForceNew:      true,
	},
	"template_project_id": {
		Description:   "When used with use_custom_template, project ID of a custom project template. This is preferable to using template_name since template_name may be ambiguous (enterprise edition). This option is mutually exclusive with `template_name`.",
		Type:          schema.TypeInt,
		Optional:      true,
		ConflictsWith: []string{"template_name"},
		ForceNew:      true,
	},
	"use_custom_template": {
		Description: "Use either custom instance or group (with group_with_project_templates_id) project template (enterprise edition).",
		Type:        schema.TypeBool,
		Optional:    true,
	},
	"group_with_project_templates_id": {
		Description: "For group-level custom templates, specifies ID of group from which all the custom project templates are sourced. Leave empty for instance-level templates. Requires use_custom_template to be true (enterprise edition).",
		Type:        schema.TypeInt,
		Optional:    true,
	},
	"pages_access_level": {
		Description:  "Enable pages access control",
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "private",
		ValidateFunc: validation.StringInSlice([]string{"public", "private", "enabled", "disabled"}, true),
	},
	// The GitLab API requires that import_url is also set when mirror options are used
	// Ref: https://github.com/gitlabhq/terraform-provider-gitlab/pull/449#discussion_r549729230
	"mirror": {
		Description:  "Enable project pull mirror.",
		Type:         schema.TypeBool,
		Optional:     true,
		Default:      false,
		RequiredWith: []string{"import_url"},
	},
	"mirror_trigger_builds": {
		Description:  "Enable trigger builds on pushes for a mirrored project.",
		Type:         schema.TypeBool,
		Optional:     true,
		Default:      false,
		RequiredWith: []string{"import_url"},
	},
	"mirror_overwrites_diverged_branches": {
		Description:  "Enable overwrite diverged branches for a mirrored project.",
		Type:         schema.TypeBool,
		Optional:     true,
		Default:      false,
		RequiredWith: []string{"import_url"},
	},
	"only_mirror_protected_branches": {
		Description:  "Enable only mirror protected branches for a mirrored project.",
		Type:         schema.TypeBool,
		Optional:     true,
		Default:      false,
		RequiredWith: []string{"import_url"},
	},
	"build_coverage_regex": {
		Description: "Test coverage parsing for the project.",
		Type:        schema.TypeString,
		Optional:    true,
	},
	"issues_template": {
		Description: "Sets the template for new issues in the project.",
		Type:        schema.TypeString,
		Optional:    true,
	},
	"merge_requests_template": {
		Description: "Sets the template for new merge requests in the project.",
		Type:        schema.TypeString,
		Optional:    true,
	},
	"ci_config_path": {
		Description: "Custom Path to CI config file.",
		Type:        schema.TypeString,
		Optional:    true,
	},
	"archive_on_destroy": {
		Description: "Set to `true` to archive the project instead of deleting on destroy. If set to `true` it will entire omit the `DELETE` operation.",
		Type:        schema.TypeBool,
		Optional:    true,
	},
}

func resourceGitlabProject() *schema.Resource {
	return &schema.Resource{
		Description: "This resource allows you to create and manage projects within your GitLab group or within your user.",

		CreateContext: resourceGitlabProjectCreate,
		ReadContext:   resourceGitlabProjectRead,
		UpdateContext: resourceGitlabProjectUpdate,
		DeleteContext: resourceGitlabProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resourceGitLabProjectSchema,
	}
}

func resourceGitlabProjectSetToState(d *schema.ResourceData, project *gitlab.Project) error {
	d.SetId(fmt.Sprintf("%d", project.ID))
	d.Set("name", project.Name)
	d.Set("path", project.Path)
	d.Set("path_with_namespace", project.PathWithNamespace)
	d.Set("description", project.Description)
	d.Set("default_branch", project.DefaultBranch)
	d.Set("request_access_enabled", project.RequestAccessEnabled)
	d.Set("issues_enabled", project.IssuesEnabled)
	d.Set("merge_requests_enabled", project.MergeRequestsEnabled)
	d.Set("pipelines_enabled", project.JobsEnabled)
	d.Set("approvals_before_merge", project.ApprovalsBeforeMerge)
	d.Set("wiki_enabled", project.WikiEnabled)
	d.Set("snippets_enabled", project.SnippetsEnabled)
	d.Set("container_registry_enabled", project.ContainerRegistryEnabled)
	d.Set("lfs_enabled", project.LFSEnabled)
	d.Set("visibility_level", string(project.Visibility))
	d.Set("merge_method", string(project.MergeMethod))
	d.Set("only_allow_merge_if_pipeline_succeeds", project.OnlyAllowMergeIfPipelineSucceeds)
	d.Set("only_allow_merge_if_all_discussions_are_resolved", project.OnlyAllowMergeIfAllDiscussionsAreResolved)
	d.Set("allow_merge_on_skipped_pipeline", project.AllowMergeOnSkippedPipeline)
	d.Set("namespace_id", project.Namespace.ID)
	d.Set("ssh_url_to_repo", project.SSHURLToRepo)
	d.Set("http_url_to_repo", project.HTTPURLToRepo)
	d.Set("web_url", project.WebURL)
	d.Set("runners_token", project.RunnersToken)
	d.Set("shared_runners_enabled", project.SharedRunnersEnabled)
	if err := d.Set("tags", project.TagList); err != nil {
		return err
	}
	d.Set("archived", project.Archived)
	d.Set("squash_option", project.SquashOption)
	d.Set("remove_source_branch_after_merge", project.RemoveSourceBranchAfterMerge)
	d.Set("packages_enabled", project.PackagesEnabled)
	d.Set("pages_access_level", string(project.PagesAccessLevel))
	d.Set("mirror", project.Mirror)
	d.Set("mirror_trigger_builds", project.MirrorTriggerBuilds)
	d.Set("mirror_overwrites_diverged_branches", project.MirrorOverwritesDivergedBranches)
	d.Set("only_mirror_protected_branches", project.OnlyMirrorProtectedBranches)
	d.Set("build_coverage_regex", project.BuildCoverageRegex)
	d.Set("issues_template", project.IssuesTemplate)
	d.Set("merge_requests_template", project.MergeRequestsTemplate)
	d.Set("ci_config_path", project.CIConfigPath)
	return nil
}

func resourceGitlabProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	options := &gitlab.CreateProjectOptions{
		Name:                             gitlab.String(d.Get("name").(string)),
		RequestAccessEnabled:             gitlab.Bool(d.Get("request_access_enabled").(bool)),
		IssuesEnabled:                    gitlab.Bool(d.Get("issues_enabled").(bool)),
		MergeRequestsEnabled:             gitlab.Bool(d.Get("merge_requests_enabled").(bool)),
		JobsEnabled:                      gitlab.Bool(d.Get("pipelines_enabled").(bool)),
		ApprovalsBeforeMerge:             gitlab.Int(d.Get("approvals_before_merge").(int)),
		WikiEnabled:                      gitlab.Bool(d.Get("wiki_enabled").(bool)),
		SnippetsEnabled:                  gitlab.Bool(d.Get("snippets_enabled").(bool)),
		ContainerRegistryEnabled:         gitlab.Bool(d.Get("container_registry_enabled").(bool)),
		LFSEnabled:                       gitlab.Bool(d.Get("lfs_enabled").(bool)),
		Visibility:                       stringToVisibilityLevel(d.Get("visibility_level").(string)),
		MergeMethod:                      stringToMergeMethod(d.Get("merge_method").(string)),
		OnlyAllowMergeIfPipelineSucceeds: gitlab.Bool(d.Get("only_allow_merge_if_pipeline_succeeds").(bool)),
		OnlyAllowMergeIfAllDiscussionsAreResolved: gitlab.Bool(d.Get("only_allow_merge_if_all_discussions_are_resolved").(bool)),
		AllowMergeOnSkippedPipeline:               gitlab.Bool(d.Get("allow_merge_on_skipped_pipeline").(bool)),
		SharedRunnersEnabled:                      gitlab.Bool(d.Get("shared_runners_enabled").(bool)),
		SquashOption:                              stringToSquashOptionValue(d.Get("squash_option").(string)),
		RemoveSourceBranchAfterMerge:              gitlab.Bool(d.Get("remove_source_branch_after_merge").(bool)),
		PackagesEnabled:                           gitlab.Bool(d.Get("packages_enabled").(bool)),
		Mirror:                                    gitlab.Bool(d.Get("mirror").(bool)),
		MirrorTriggerBuilds:                       gitlab.Bool(d.Get("mirror_trigger_builds").(bool)),
		BuildCoverageRegex:                        gitlab.String(d.Get("build_coverage_regex").(string)),
		CIConfigPath:                              gitlab.String(d.Get("ci_config_path").(string)),
	}

	if v, ok := d.GetOk("path"); ok {
		options.Path = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("namespace_id"); ok {
		options.NamespaceID = gitlab.Int(v.(int))
	}

	if v, ok := d.GetOk("description"); ok {
		options.Description = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("default_branch"); ok {
		options.DefaultBranch = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		options.TagList = stringSetToStringSlice(v.(*schema.Set))
	}

	if v, ok := d.GetOk("initialize_with_readme"); ok {
		options.InitializeWithReadme = gitlab.Bool(v.(bool))
	}

	if v, ok := d.GetOk("import_url"); ok {
		options.ImportURL = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("template_name"); ok {
		options.TemplateName = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("template_project_id"); ok {
		options.TemplateProjectID = gitlab.Int(v.(int))
	}

	if v, ok := d.GetOk("use_custom_template"); ok {
		options.UseCustomTemplate = gitlab.Bool(v.(bool))
	}

	if v, ok := d.GetOk("group_with_project_templates_id"); ok {
		options.GroupWithProjectTemplatesID = gitlab.Int(v.(int))
	}

	if v, ok := d.GetOk("pages_access_level"); ok {
		options.PagesAccessLevel = stringToAccessControlValue(v.(string))
	}

	if v, ok := d.GetOk("ci_config_path"); ok {
		options.CIConfigPath = gitlab.String(v.(string))
	}

	log.Printf("[DEBUG] create gitlab project %q", *options.Name)

	project, _, err := client.Projects.CreateProject(options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	// from this point onwards no matter how we return, resource creation
	// is committed to state since we set its ID
	d.SetId(fmt.Sprintf("%d", project.ID))

	// An import can be triggered by import_url or by creating the project from a template.
	if project.ImportStatus != "none" {
		log.Printf("[DEBUG] waiting for project %q import to finish", *options.Name)

		stateConf := &resource.StateChangeConf{
			Pending: []string{"scheduled", "started"},
			Target:  []string{"finished"},
			Timeout: 10 * time.Minute,
			Refresh: func() (interface{}, string, error) {
				status, _, err := client.ProjectImportExport.ImportStatus(d.Id(), gitlab.WithContext(ctx))
				if err != nil {
					return nil, "", err
				}

				return status, status.ImportStatus, nil
			},
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return diag.Errorf("error while waiting for project %q import to finish: %s", *options.Name, err)
		}

		// Read the project again, so that we can detect the default branch.
		project, _, err = client.Projects.GetProject(project.ID, nil, gitlab.WithContext(ctx))
		if err != nil {
			return diag.Errorf("Failed to get project %q after completing import: %s", d.Id(), err)
		}
	}

	if d.Get("archived").(bool) {
		// strange as it may seem, this project is created in archived state...
		if _, _, err := client.Projects.ArchiveProject(d.Id(), gitlab.WithContext(ctx)); err != nil {
			return diag.Errorf("new project %q could not be archived: %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("push_rules"); ok {
		err := editOrAddPushRules(ctx, client, d.Id(), d)
		var httpError *gitlab.ErrorResponse
		if errors.As(err, &httpError) && httpError.Response.StatusCode == http.StatusNotFound {
			log.Printf("[DEBUG] Failed to edit push rules for project %q: %v", d.Id(), err)
			return diag.Errorf("Project push rules are not supported in your version of GitLab")
		}
		if err != nil {
			return diag.Errorf("Failed to edit push rules for project %q: %s", d.Id(), err)
		}
	}

	// default_branch cannot always be set during creation.
	// If the branch does not exist, the update will fail, so we also create it here.
	// See: https://gitlab.com/gitlab-org/gitlab/-/issues/333426
	// This logic may be removed when the above issue is resolved.
	if v, ok := d.GetOk("default_branch"); ok && project.DefaultBranch != "" && project.DefaultBranch != v.(string) {
		oldDefaultBranch := project.DefaultBranch
		newDefaultBranch := v.(string)

		log.Printf("[DEBUG] create branch %q for project %q", newDefaultBranch, d.Id())
		_, _, err := client.Branches.CreateBranch(project.ID, &gitlab.CreateBranchOptions{
			Branch: gitlab.String(newDefaultBranch),
			Ref:    gitlab.String(oldDefaultBranch),
		}, gitlab.WithContext(ctx))
		if err != nil {
			return diag.Errorf("Failed to create branch %q for project %q: %s", newDefaultBranch, d.Id(), err)
		}

		log.Printf("[DEBUG] set new default branch to %q for project %q", newDefaultBranch, d.Id())
		_, _, err = client.Projects.EditProject(project.ID, &gitlab.EditProjectOptions{
			DefaultBranch: gitlab.String(newDefaultBranch),
		}, gitlab.WithContext(ctx))
		if err != nil {
			return diag.Errorf("Failed to set default branch to %q for project %q: %s", newDefaultBranch, d.Id(), err)
		}

		log.Printf("[DEBUG] protect new default branch %q for project %q", newDefaultBranch, d.Id())
		_, _, err = client.ProtectedBranches.ProtectRepositoryBranches(project.ID, &gitlab.ProtectRepositoryBranchesOptions{
			Name: gitlab.String(newDefaultBranch),
		}, gitlab.WithContext(ctx))
		if err != nil {
			return diag.Errorf("Failed to protect default branch %q for project %q: %s", newDefaultBranch, d.Id(), err)
		}

		log.Printf("[DEBUG] unprotect old default branch %q for project %q", oldDefaultBranch, d.Id())
		_, err = client.ProtectedBranches.UnprotectRepositoryBranches(project.ID, oldDefaultBranch, gitlab.WithContext(ctx))
		if err != nil {
			return diag.Errorf("Failed to unprotect undesired default branch %q for project %q: %s", oldDefaultBranch, d.Id(), err)
		}

		log.Printf("[DEBUG] delete old default branch %q for project %q", oldDefaultBranch, d.Id())
		_, err = client.Branches.DeleteBranch(project.ID, oldDefaultBranch, gitlab.WithContext(ctx))
		if err != nil {
			return diag.Errorf("Failed to clean up undesired default branch %q for project %q: %s", oldDefaultBranch, d.Id(), err)
		}
	}

	var editProjectOptions gitlab.EditProjectOptions

	if v, ok := d.GetOk("mirror_overwrites_diverged_branches"); ok {
		editProjectOptions.MirrorOverwritesDivergedBranches = gitlab.Bool(v.(bool))
		editProjectOptions.ImportURL = gitlab.String(d.Get("import_url").(string))
	}

	if v, ok := d.GetOk("only_mirror_protected_branches"); ok {
		editProjectOptions.OnlyMirrorProtectedBranches = gitlab.Bool(v.(bool))
		editProjectOptions.ImportURL = gitlab.String(d.Get("import_url").(string))
	}

	if v, ok := d.GetOk("issues_template"); ok {
		editProjectOptions.IssuesTemplate = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("merge_requests_template"); ok {
		editProjectOptions.MergeRequestsTemplate = gitlab.String(v.(string))
	}

	if (editProjectOptions != gitlab.EditProjectOptions{}) {
		if _, _, err := client.Projects.EditProject(d.Id(), &editProjectOptions, gitlab.WithContext(ctx)); err != nil {
			return diag.Errorf("Could not update project %q: %s", d.Id(), err)
		}
	}

	return resourceGitlabProjectRead(ctx, d, meta)
}

func resourceGitlabProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	log.Printf("[DEBUG] read gitlab project %s", d.Id())

	project, _, err := client.Projects.GetProject(d.Id(), nil, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}
	if project.MarkedForDeletionAt != nil {
		log.Printf("[DEBUG] gitlab project %s is marked for deletion", d.Id())
		d.SetId("")
		return nil
	}

	if err := resourceGitlabProjectSetToState(d, project); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] read gitlab project %q push rules", d.Id())

	pushRules, _, err := client.Projects.GetProjectPushRules(d.Id(), gitlab.WithContext(ctx))
	var httpError *gitlab.ErrorResponse
	if errors.As(err, &httpError) && httpError.Response.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] Failed to get push rules for project %q: %v", d.Id(), err)
	} else if err != nil {
		return diag.Errorf("Failed to get push rules for project %q: %s", d.Id(), err)
	}

	if err := d.Set("push_rules", flattenProjectPushRules(pushRules)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceGitlabProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	options := &gitlab.EditProjectOptions{}
	transferOptions := &gitlab.TransferProjectOptions{}

	if d.HasChange("name") {
		options.Name = gitlab.String(d.Get("name").(string))
	}

	if d.HasChange("path") && (d.Get("path").(string) != "") {
		options.Path = gitlab.String(d.Get("path").(string))
	}

	if d.HasChange("namespace_id") {
		transferOptions.Namespace = gitlab.Int(d.Get("namespace_id").(int))
	}

	if d.HasChange("description") {
		options.Description = gitlab.String(d.Get("description").(string))
	}

	if d.HasChange("default_branch") {
		options.DefaultBranch = gitlab.String(d.Get("default_branch").(string))
	}

	if d.HasChange("visibility_level") {
		options.Visibility = stringToVisibilityLevel(d.Get("visibility_level").(string))
	}

	if d.HasChange("merge_method") {
		options.MergeMethod = stringToMergeMethod(d.Get("merge_method").(string))
	}

	if d.HasChange("only_allow_merge_if_pipeline_succeeds") {
		options.OnlyAllowMergeIfPipelineSucceeds = gitlab.Bool(d.Get("only_allow_merge_if_pipeline_succeeds").(bool))
	}

	if d.HasChange("only_allow_merge_if_all_discussions_are_resolved") {
		options.OnlyAllowMergeIfAllDiscussionsAreResolved = gitlab.Bool(d.Get("only_allow_merge_if_all_discussions_are_resolved").(bool))
	}

	if d.HasChange("allow_merge_on_skipped_pipeline") {
		options.AllowMergeOnSkippedPipeline = gitlab.Bool(d.Get("allow_merge_on_skipped_pipeline").(bool))
	}

	if d.HasChange("request_access_enabled") {
		options.RequestAccessEnabled = gitlab.Bool(d.Get("request_access_enabled").(bool))
	}

	if d.HasChange("issues_enabled") {
		options.IssuesEnabled = gitlab.Bool(d.Get("issues_enabled").(bool))
	}

	if d.HasChange("merge_requests_enabled") {
		options.MergeRequestsEnabled = gitlab.Bool(d.Get("merge_requests_enabled").(bool))
	}

	if d.HasChange("pipelines_enabled") {
		options.JobsEnabled = gitlab.Bool(d.Get("pipelines_enabled").(bool))
	}

	if d.HasChange("approvals_before_merge") {
		options.ApprovalsBeforeMerge = gitlab.Int(d.Get("approvals_before_merge").(int))
	}

	if d.HasChange("wiki_enabled") {
		options.WikiEnabled = gitlab.Bool(d.Get("wiki_enabled").(bool))
	}

	if d.HasChange("snippets_enabled") {
		options.SnippetsEnabled = gitlab.Bool(d.Get("snippets_enabled").(bool))
	}

	if d.HasChange("shared_runners_enabled") {
		options.SharedRunnersEnabled = gitlab.Bool(d.Get("shared_runners_enabled").(bool))
	}

	if d.HasChange("tags") {
		options.TagList = stringSetToStringSlice(d.Get("tags").(*schema.Set))
	}

	if d.HasChange("container_registry_enabled") {
		options.ContainerRegistryEnabled = gitlab.Bool(d.Get("container_registry_enabled").(bool))
	}

	if d.HasChange("lfs_enabled") {
		options.LFSEnabled = gitlab.Bool(d.Get("lfs_enabled").(bool))
	}

	if d.HasChange("squash_option") {
		options.SquashOption = stringToSquashOptionValue(d.Get("squash_option").(string))
	}

	if d.HasChange("remove_source_branch_after_merge") {
		options.RemoveSourceBranchAfterMerge = gitlab.Bool(d.Get("remove_source_branch_after_merge").(bool))
	}

	if d.HasChange("packages_enabled") {
		options.PackagesEnabled = gitlab.Bool(d.Get("packages_enabled").(bool))
	}

	if d.HasChange("pages_access_level") {
		options.PagesAccessLevel = stringToAccessControlValue(d.Get("pages_access_level").(string))
	}

	if d.HasChange("mirror") {
		options.ImportURL = gitlab.String(d.Get("import_url").(string))
		options.Mirror = gitlab.Bool(d.Get("mirror").(bool))
	}

	if d.HasChange("mirror_trigger_builds") {
		options.ImportURL = gitlab.String(d.Get("import_url").(string))
		options.MirrorTriggerBuilds = gitlab.Bool(d.Get("mirror_trigger_builds").(bool))
	}

	if d.HasChange("only_mirror_protected_branches") {
		options.ImportURL = gitlab.String(d.Get("import_url").(string))
		options.OnlyMirrorProtectedBranches = gitlab.Bool(d.Get("only_mirror_protected_branches").(bool))
	}

	if d.HasChange("mirror_overwrites_diverged_branches") {
		options.ImportURL = gitlab.String(d.Get("import_url").(string))
		options.MirrorOverwritesDivergedBranches = gitlab.Bool(d.Get("mirror_overwrites_diverged_branches").(bool))
	}

	if d.HasChange("build_coverage_regex") {
		options.BuildCoverageRegex = gitlab.String(d.Get("build_coverage_regex").(string))
	}

	if d.HasChange("issues_template") {
		options.IssuesTemplate = gitlab.String(d.Get("issues_template").(string))
	}

	if d.HasChange("merge_requests_template") {
		options.MergeRequestsTemplate = gitlab.String(d.Get("merge_requests_template").(string))
	}

	if d.HasChange("ci_config_path") {
		options.CIConfigPath = gitlab.String(d.Get("ci_config_path").(string))
	}

	if *options != (gitlab.EditProjectOptions{}) {
		log.Printf("[DEBUG] update gitlab project %s", d.Id())
		_, _, err := client.Projects.EditProject(d.Id(), options, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if *transferOptions != (gitlab.TransferProjectOptions{}) {
		log.Printf("[DEBUG] transferring project %s to namespace %d", d.Id(), transferOptions.Namespace)
		_, _, err := client.Projects.TransferProject(d.Id(), transferOptions, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("archived") {
		if d.Get("archived").(bool) {
			if _, _, err := client.Projects.ArchiveProject(d.Id(), gitlab.WithContext(ctx)); err != nil {
				return diag.Errorf("project %q could not be archived: %s", d.Id(), err)
			}
		} else {
			if _, _, err := client.Projects.UnarchiveProject(d.Id(), gitlab.WithContext(ctx)); err != nil {
				return diag.Errorf("project %q could not be unarchived: %s", d.Id(), err)
			}
		}
	}

	if d.HasChange("push_rules") {
		err := editOrAddPushRules(ctx, client, d.Id(), d)
		var httpError *gitlab.ErrorResponse
		if errors.As(err, &httpError) && httpError.Response.StatusCode == http.StatusNotFound {
			log.Printf("[DEBUG] Failed to get push rules for project %q: %v", d.Id(), err)
			return diag.Errorf("Project push rules are not supported in your version of GitLab")
		}
		if err != nil {
			return diag.Errorf("Failed to edit push rules for project %q: %s", d.Id(), err)
		}
	}

	return resourceGitlabProjectRead(ctx, d, meta)
}

func resourceGitlabProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	if !d.Get("archive_on_destroy").(bool) {
		log.Printf("[DEBUG] Delete gitlab project %s", d.Id())
		_, err := client.Projects.DeleteProject(d.Id(), gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}

		// Wait for the project to be deleted.
		// Deleting a project in gitlab is async.
		stateConf := &resource.StateChangeConf{
			Pending: []string{"Deleting"},
			Target:  []string{"Deleted"},
			Refresh: func() (interface{}, string, error) {
				out, response, err := client.Projects.GetProject(d.Id(), nil, gitlab.WithContext(ctx))
				if err != nil {
					if response.StatusCode == 404 {
						return out, "Deleted", nil
					}
					log.Printf("[ERROR] Received error: %#v", err)
					return out, "Error", err
				}
				if out.MarkedForDeletionAt != nil {
					// Represents a Gitlab EE soft-delete
					return out, "Deleted", nil
				}
				return out, "Deleting", nil
			},

			Timeout:    10 * time.Minute,
			MinTimeout: 3 * time.Second,
			Delay:      5 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for project (%s) to become deleted: %s", d.Id(), err)
		}

	} else {
		log.Printf("[DEBUG] Archive gitlab project %s", d.Id())
		_, _, err := client.Projects.ArchiveProject(d.Id(), gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func editOrAddPushRules(ctx context.Context, client *gitlab.Client, projectID string, d *schema.ResourceData) error {
	log.Printf("[DEBUG] Editing push rules for project %q", projectID)

	editOptions := expandEditProjectPushRuleOptions(d)
	_, _, err := client.Projects.EditProjectPushRule(projectID, editOptions, gitlab.WithContext(ctx))
	if err == nil {
		return nil
	}

	var httpErr *gitlab.ErrorResponse
	if !errors.As(err, &httpErr) || httpErr.Response.StatusCode != http.StatusNotFound {
		return err
	}

	// A 404 could mean that the push rules need to be re-created.

	log.Printf("[DEBUG] Failed to edit push rules for project %q: %v", projectID, err)
	log.Printf("[DEBUG] Creating new push rules for project %q", projectID)

	addOptions := expandAddProjectPushRuleOptions(d)
	_, _, err = client.Projects.AddProjectPushRule(projectID, addOptions, gitlab.WithContext(ctx))
	if err != nil {
		return err
	}

	return nil
}

func expandEditProjectPushRuleOptions(d *schema.ResourceData) *gitlab.EditProjectPushRuleOptions {
	options := &gitlab.EditProjectPushRuleOptions{}

	if d.HasChange("push_rules.0.author_email_regex") {
		options.AuthorEmailRegex = gitlab.String(d.Get("push_rules.0.author_email_regex").(string))
	}

	if d.HasChange("push_rules.0.branch_name_regex") {
		options.BranchNameRegex = gitlab.String(d.Get("push_rules.0.branch_name_regex").(string))
	}

	if d.HasChange("push_rules.0.commit_message_regex") {
		options.CommitMessageRegex = gitlab.String(d.Get("push_rules.0.commit_message_regex").(string))
	}

	if d.HasChange("push_rules.0.commit_message_negative_regex") {
		options.CommitMessageNegativeRegex = gitlab.String(d.Get("push_rules.0.commit_message_negative_regex").(string))
	}

	if d.HasChange("push_rules.0.file_name_regex") {
		options.FileNameRegex = gitlab.String(d.Get("push_rules.0.file_name_regex").(string))
	}

	if d.HasChange("push_rules.0.commit_committer_check") {
		options.CommitCommitterCheck = gitlab.Bool(d.Get("push_rules.0.commit_committer_check").(bool))
	}

	if d.HasChange("push_rules.0.deny_delete_tag") {
		options.DenyDeleteTag = gitlab.Bool(d.Get("push_rules.0.deny_delete_tag").(bool))
	}

	if d.HasChange("push_rules.0.member_check") {
		options.MemberCheck = gitlab.Bool(d.Get("push_rules.0.member_check").(bool))
	}

	if d.HasChange("push_rules.0.prevent_secrets") {
		options.PreventSecrets = gitlab.Bool(d.Get("push_rules.0.prevent_secrets").(bool))
	}

	if d.HasChange("push_rules.0.reject_unsigned_commits") {
		options.RejectUnsignedCommits = gitlab.Bool(d.Get("push_rules.0.reject_unsigned_commits").(bool))
	}

	if d.HasChange("push_rules.0.max_file_size") {
		options.MaxFileSize = gitlab.Int(d.Get("push_rules.0.max_file_size").(int))
	}

	return options
}

func expandAddProjectPushRuleOptions(d *schema.ResourceData) *gitlab.AddProjectPushRuleOptions {
	options := &gitlab.AddProjectPushRuleOptions{}

	if v, ok := d.GetOk("push_rules.0.author_email_regex"); ok {
		options.AuthorEmailRegex = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("push_rules.0.branch_name_regex"); ok {
		options.BranchNameRegex = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("push_rules.0.commit_message_regex"); ok {
		options.CommitMessageRegex = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("push_rules.0.commit_message_negative_regex"); ok {
		options.CommitMessageNegativeRegex = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("push_rules.0.file_name_regex"); ok {
		options.FileNameRegex = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("push_rules.0.commit_committer_check"); ok {
		options.CommitCommitterCheck = gitlab.Bool(v.(bool))
	}

	if v, ok := d.GetOk("push_rules.0.deny_delete_tag"); ok {
		options.DenyDeleteTag = gitlab.Bool(v.(bool))
	}

	if v, ok := d.GetOk("push_rules.0.member_check"); ok {
		options.MemberCheck = gitlab.Bool(v.(bool))
	}

	if v, ok := d.GetOk("push_rules.0.prevent_secrets"); ok {
		options.PreventSecrets = gitlab.Bool(v.(bool))
	}

	if v, ok := d.GetOk("push_rules.0.reject_unsigned_commits"); ok {
		options.RejectUnsignedCommits = gitlab.Bool(v.(bool))
	}

	if v, ok := d.GetOk("push_rules.0.max_file_size"); ok {
		options.MaxFileSize = gitlab.Int(v.(int))
	}

	return options
}

func flattenProjectPushRules(pushRules *gitlab.ProjectPushRules) (values []map[string]interface{}) {
	if pushRules == nil {
		return []map[string]interface{}{}
	}

	return []map[string]interface{}{
		{
			"author_email_regex":            pushRules.AuthorEmailRegex,
			"branch_name_regex":             pushRules.BranchNameRegex,
			"commit_message_regex":          pushRules.CommitMessageRegex,
			"commit_message_negative_regex": pushRules.CommitMessageNegativeRegex,
			"file_name_regex":               pushRules.FileNameRegex,
			"commit_committer_check":        pushRules.CommitCommitterCheck,
			"deny_delete_tag":               pushRules.DenyDeleteTag,
			"member_check":                  pushRules.MemberCheck,
			"prevent_secrets":               pushRules.PreventSecrets,
			"reject_unsigned_commits":       pushRules.RejectUnsignedCommits,
			"max_file_size":                 pushRules.MaxFileSize,
		},
	}
}
