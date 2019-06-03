package gitlab

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	gitlab "github.com/xanzy/go-gitlab"
)

var resourceGitLabProjectSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"path": {
		Type:     schema.TypeString,
		Optional: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			if new == "" {
				return true
			}
			return old == new
		},
	},
	"namespace_id": {
		Type:     schema.TypeInt,
		Optional: true,
		ForceNew: true,
		Computed: true,
	},
	"description": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"default_branch": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"issues_enabled": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
	"merge_requests_enabled": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
	"approvals_before_merge": {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  0,
	},
	"wiki_enabled": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
	"snippets_enabled": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
	"container_registry_enabled": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
	"visibility_level": {
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"private", "internal", "public"}, true),
		Default:      "private",
	},
	"import_url": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"merge_method": {
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"merge", "rebase_merge", "ff"}, true),
		Default:      "merge",
	},
	"only_allow_merge_if_pipeline_succeeds": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"only_allow_merge_if_all_discussions_are_resolved": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"ssh_url_to_repo": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"http_url_to_repo": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"web_url": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"runners_token": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"shared_runners_enabled": {
		Type:     schema.TypeBool,
		Optional: true,
		Computed: true,
	},
	"tags": {
		Type:     schema.TypeSet,
		Optional: true,
		ForceNew: false,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Set:      schema.HashString,
	},
	"shared_with_groups": {
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"group_id": {
					Type:     schema.TypeInt,
					Required: true,
				},
				"group_access_level": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						"no one", "guest", "reporter", "developer", "maintainer"}, false),
				},
				"group_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	},
}

func resourceGitlabProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabProjectCreate,
		Read:   resourceGitlabProjectRead,
		Update: resourceGitlabProjectUpdate,
		Delete: resourceGitlabProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: resourceGitLabProjectSchema,
	}
}

func resourceGitlabProjectSetToState(d *schema.ResourceData, project *gitlab.Project) {
	d.SetId(fmt.Sprintf("%d", project.ID))
	d.Set("name", project.Name)
	d.Set("path", project.Path)
	d.Set("description", project.Description)
	d.Set("default_branch", project.DefaultBranch)
	d.Set("issues_enabled", project.IssuesEnabled)
	d.Set("merge_requests_enabled", project.MergeRequestsEnabled)
	d.Set("approvals_before_merge", project.ApprovalsBeforeMerge)
	d.Set("wiki_enabled", project.WikiEnabled)
	d.Set("snippets_enabled", project.SnippetsEnabled)
	d.Set("container_registry_enabled", project.ContainerRegistryEnabled)
	d.Set("visibility_level", string(project.Visibility))
	d.Set("merge_method", string(project.MergeMethod))
	d.Set("only_allow_merge_if_pipeline_succeeds", project.OnlyAllowMergeIfPipelineSucceeds)
	d.Set("only_allow_merge_if_all_discussions_are_resolved", project.OnlyAllowMergeIfAllDiscussionsAreResolved)
	d.Set("namespace_id", project.Namespace.ID)
	d.Set("ssh_url_to_repo", project.SSHURLToRepo)
	d.Set("http_url_to_repo", project.HTTPURLToRepo)
	d.Set("web_url", project.WebURL)
	d.Set("runners_token", project.RunnersToken)
	d.Set("shared_runners_enabled", project.SharedRunnersEnabled)
	d.Set("shared_with_groups", flattenSharedWithGroupsOptions(project))
	d.Set("tags", project.TagList)
}

func resourceGitlabProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	options := &gitlab.CreateProjectOptions{
		Name:                             gitlab.String(d.Get("name").(string)),
		IssuesEnabled:                    gitlab.Bool(d.Get("issues_enabled").(bool)),
		MergeRequestsEnabled:             gitlab.Bool(d.Get("merge_requests_enabled").(bool)),
		ApprovalsBeforeMerge:             gitlab.Int(d.Get("approvals_before_merge").(int)),
		WikiEnabled:                      gitlab.Bool(d.Get("wiki_enabled").(bool)),
		SnippetsEnabled:                  gitlab.Bool(d.Get("snippets_enabled").(bool)),
		ContainerRegistryEnabled:         gitlab.Bool(d.Get("container_registry_enabled").(bool)),
		Visibility:                       stringToVisibilityLevel(d.Get("visibility_level").(string)),
		MergeMethod:                      stringToMergeMethod(d.Get("merge_method").(string)),
		OnlyAllowMergeIfPipelineSucceeds: gitlab.Bool(d.Get("only_allow_merge_if_pipeline_succeeds").(bool)),
		OnlyAllowMergeIfAllDiscussionsAreResolved: gitlab.Bool(d.Get("only_allow_merge_if_all_discussions_are_resolved").(bool)),
		SharedRunnersEnabled:                      gitlab.Bool(d.Get("shared_runners_enabled").(bool)),
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

	if v, ok := d.GetOk("tags"); ok {
		options.TagList = stringSetToStringSlice(v.(*schema.Set))
	}

	if v, ok := d.GetOk("import_url"); ok {
		options.ImportURL = gitlab.String(v.(string))
	}

	log.Printf("[DEBUG] create gitlab project %q", *options.Name)

	project, _, err := client.Projects.CreateProject(options)
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("shared_with_groups"); ok {
		for _, option := range expandSharedWithGroupsOptions(v) {
			if _, err := client.Projects.ShareProjectWithGroup(project.ID, option); err != nil {
				return err
			}
		}
	}

	d.SetId(fmt.Sprintf("%d", project.ID))

	return resourceGitlabProjectRead(d, meta)
}

func resourceGitlabProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	log.Printf("[DEBUG] read gitlab project %s", d.Id())

	project, _, err := client.Projects.GetProject(d.Id(), nil)
	if err != nil {
		return err
	}

	resourceGitlabProjectSetToState(d, project)
	return nil
}

func resourceGitlabProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	options := &gitlab.EditProjectOptions{}

	if d.HasChange("name") {
		options.Name = gitlab.String(d.Get("name").(string))
	}

	if d.HasChange("path") && (d.Get("path").(string) != "") {
		options.Path = gitlab.String(d.Get("path").(string))
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

	if d.HasChange("issues_enabled") {
		options.IssuesEnabled = gitlab.Bool(d.Get("issues_enabled").(bool))
	}

	if d.HasChange("merge_requests_enabled") {
		options.MergeRequestsEnabled = gitlab.Bool(d.Get("merge_requests_enabled").(bool))
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

	if *options != (gitlab.EditProjectOptions{}) {
		log.Printf("[DEBUG] update gitlab project %s", d.Id())
		_, _, err := client.Projects.EditProject(d.Id(), options)
		if err != nil {
			return err
		}
	}

	if d.HasChange("shared_with_groups") {
		updateSharedWithGroups(d, meta)
	}

	return resourceGitlabProjectRead(d, meta)
}

func resourceGitlabProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	log.Printf("[DEBUG] Delete gitlab project %s", d.Id())

	_, err := client.Projects.DeleteProject(d.Id())
	if err != nil {
		return err
	}

	// Wait for the project to be deleted.
	// Deleting a project in gitlab is async.
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Deleting"},
		Target:  []string{"Deleted"},
		Refresh: func() (interface{}, string, error) {
			out, response, err := client.Projects.GetProject(d.Id(), nil)
			if err != nil {
				if response.StatusCode == 404 {
					return out, "Deleted", nil
				}
				log.Printf("[ERROR] Received error: %#v", err)
				return out, "Error", err
			}
			return out, "Deleting", nil
		},

		Timeout:    10 * time.Minute,
		MinTimeout: 3 * time.Second,
		Delay:      5 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for project (%s) to become deleted: %s", d.Id(), err)
	}
	return nil
}

func expandSharedWithGroupsOptions(v interface{}) []*gitlab.ShareWithGroupOptions {
	shareWithGroupOptionsList := []*gitlab.ShareWithGroupOptions{}

	for _, config := range v.(*schema.Set).List() {
		data := config.(map[string]interface{})

		groupAccess := accessLevelNameToValue[data["group_access_level"].(string)]

		shareWithGroupOptions := &gitlab.ShareWithGroupOptions{
			GroupID:     gitlab.Int(data["group_id"].(int)),
			GroupAccess: &groupAccess,
		}

		shareWithGroupOptionsList = append(shareWithGroupOptionsList,
			shareWithGroupOptions)
	}

	return shareWithGroupOptionsList
}

func flattenSharedWithGroupsOptions(project *gitlab.Project) []interface{} {
	sharedWithGroups := project.SharedWithGroups
	sharedWithGroupsList := []interface{}{}

	for _, option := range sharedWithGroups {
		values := map[string]interface{}{
			"group_id": option.GroupID,
			"group_access_level": accessLevelValueToName[gitlab.AccessLevelValue(
				option.GroupAccessLevel)],
			"group_name": option.GroupName,
		}

		sharedWithGroupsList = append(sharedWithGroupsList, values)
	}

	return sharedWithGroupsList
}

func findGroupProjectSharedWith(target *gitlab.ShareWithGroupOptions,
	groups []*gitlab.ShareWithGroupOptions) (*gitlab.ShareWithGroupOptions, int, error) {
	for i, group := range groups {
		if *group.GroupID == *target.GroupID {
			return group, i, nil
		}
	}

	return nil, 0, fmt.Errorf("group not found")
}

func getGroupsProjectSharedWith(project *gitlab.Project) []*gitlab.ShareWithGroupOptions {
	sharedGroups := []*gitlab.ShareWithGroupOptions{}

	for _, group := range project.SharedWithGroups {
		sharedGroups = append(sharedGroups, &gitlab.ShareWithGroupOptions{
			GroupID: gitlab.Int(group.GroupID),
			GroupAccess: gitlab.AccessLevel(gitlab.AccessLevelValue(
				group.GroupAccessLevel)),
		})
	}

	return sharedGroups
}

func updateSharedWithGroups(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	var groupsToUnshare []*gitlab.ShareWithGroupOptions
	var groupsToShare []*gitlab.ShareWithGroupOptions

	// Get target groups from the TF config and current groups from Gitlab server
	targetGroups := expandSharedWithGroupsOptions(d.Get("shared_with_groups"))
	project, _, err := client.Projects.GetProject(d.Id(), nil)
	if err != nil {
		return err
	}
	currentGroups := getGroupsProjectSharedWith(project)

	for _, targetGroup := range targetGroups {
		currentGroup, index, err := findGroupProjectSharedWith(targetGroup, currentGroups)

		// If no corresponding group is found, it must be added
		if err != nil {
			groupsToShare = append(groupsToShare, targetGroup)
			continue
		}

		// If group is different it must be deleted and added again
		if *targetGroup.GroupAccess != *currentGroup.GroupAccess {
			groupsToShare = append(groupsToShare, targetGroup)
			groupsToUnshare = append(groupsToUnshare, targetGroup)
		}

		// Remove currentGroup from from list
		currentGroups = append(currentGroups[:index], currentGroups[index+1:]...)
	}

	// All groups still present in currentGroup must be deleted
	groupsToUnshare = append(groupsToUnshare, currentGroups...)

	// Unshare groups to delete and update
	for _, group := range groupsToUnshare {
		_, err := client.Projects.DeleteSharedProjectFromGroup(d.Id(), *group.GroupID)
		if err != nil {
			return err
		}
	}

	// Share groups to add and update
	for _, group := range groupsToShare {
		_, err := client.Projects.ShareProjectWithGroup(d.Id(), group)
		if err != nil {
			return err
		}
	}

	return nil
}
