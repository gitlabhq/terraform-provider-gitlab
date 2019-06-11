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
	"archived": {
		Type:        schema.TypeBool,
		Description: "Whether the project is archived.",
		Optional:    true,
		Default:     false,
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
	d.Set("archived", project.Archived)
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

	// need to manage partial state since project creation may require
	// more than a single API call, and they may all fail independently;
	// the default set of attributes is prepopulated with those used above
	d.Partial(true)
	setProperties := []string{
		"name",
		"issues_enabled",
		"merge_requests_enabled",
		"approvals_before_merge",
		"wiki_enabled",
		"snippets_enabled",
		"container_registry_enabled",
		"visibility_level",
		"merge_method",
		"only_allow_merge_if_pipeline_succeeds",
		"only_allow_merge_if_all_discussions_are_resolved",
		"shared_runners_enabled",
	}

	if v, ok := d.GetOk("path"); ok {
		options.Path = gitlab.String(v.(string))
		setProperties = append(setProperties, "path")
	}

	if v, ok := d.GetOk("namespace_id"); ok {
		options.NamespaceID = gitlab.Int(v.(int))
		setProperties = append(setProperties, "namespace_id")
	}

	if v, ok := d.GetOk("description"); ok {
		options.Description = gitlab.String(v.(string))
		setProperties = append(setProperties, "description")
	}

	if v, ok := d.GetOk("tags"); ok {
		options.TagList = stringSetToStringSlice(v.(*schema.Set))
		setProperties = append(setProperties, "tags")
	}

	log.Printf("[DEBUG] create gitlab project %q", *options.Name)

	project, _, err := client.Projects.CreateProject(options)
	if err != nil {
		return err
	}

	for _, setProperty := range setProperties {
		log.Printf("[DEBUG] partial gitlab project %s creation of property %q", d.Id(), setProperty)
		d.SetPartial(setProperty)
	}

	// from this point onwards no matter how we return, resource creation
	// is committed to state since we set its ID
	d.SetId(fmt.Sprintf("%d", project.ID))

	if v, ok := d.GetOk("shared_with_groups"); ok {
		for _, option := range expandSharedWithGroupsOptions(v) {
			if _, err := client.Projects.ShareProjectWithGroup(project.ID, option); err != nil {
				return err
			}
		}
		d.SetPartial("shared_with_groups")
	}

	v := d.Get("archived")
	if v.(bool) {
		// strange as it may seem, this project is created in archived state...
		err := archiveProject(d, meta)
		if err != nil {
			log.Printf("[WARN] New project (%s) could not be created in archived state: error %#v", d.Id(), err)
			return err
		}
		d.SetPartial(("archived"))
	}

	// everything went OK, we can revert to ordinary state management
	// and let the Gitlab server fill in the resource state via a read
	d.Partial(false)
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

	// need to manage partial state since project archiving requires
	// a separate API call which could fail
	d.Partial(true)
	updatedProperties := []string{}

	if d.HasChange("name") {
		options.Name = gitlab.String(d.Get("name").(string))
		updatedProperties = append(updatedProperties, "name")
	}

	if d.HasChange("path") && (d.Get("path").(string) != "") {
		options.Path = gitlab.String(d.Get("path").(string))
		updatedProperties = append(updatedProperties, "path")
	}

	if d.HasChange("description") {
		options.Description = gitlab.String(d.Get("description").(string))
		updatedProperties = append(updatedProperties, "description")
	}

	if d.HasChange("default_branch") {
		options.DefaultBranch = gitlab.String(d.Get("default_branch").(string))
		updatedProperties = append(updatedProperties, "default_branch")
	}

	if d.HasChange("visibility_level") {
		options.Visibility = stringToVisibilityLevel(d.Get("visibility_level").(string))
		updatedProperties = append(updatedProperties, "visibility_level")
	}

	if d.HasChange("merge_method") {
		options.MergeMethod = stringToMergeMethod(d.Get("merge_method").(string))
		updatedProperties = append(updatedProperties, "merge_method")
	}

	if d.HasChange("only_allow_merge_if_pipeline_succeeds") {
		options.OnlyAllowMergeIfPipelineSucceeds = gitlab.Bool(d.Get("only_allow_merge_if_pipeline_succeeds").(bool))
		updatedProperties = append(updatedProperties, "only_allow_merge_if_pipeline_succeeds")
	}

	if d.HasChange("only_allow_merge_if_all_discussions_are_resolved") {
		options.OnlyAllowMergeIfAllDiscussionsAreResolved = gitlab.Bool(d.Get("only_allow_merge_if_all_discussions_are_resolved").(bool))
		updatedProperties = append(updatedProperties, "only_allow_merge_if_all_discussions_are_resolved")
	}

	if d.HasChange("issues_enabled") {
		options.IssuesEnabled = gitlab.Bool(d.Get("issues_enabled").(bool))
		updatedProperties = append(updatedProperties, "issues_enabled")
	}

	if d.HasChange("merge_requests_enabled") {
		options.MergeRequestsEnabled = gitlab.Bool(d.Get("merge_requests_enabled").(bool))
		updatedProperties = append(updatedProperties, "merge_requests_enabled")
	}

	if d.HasChange("approvals_before_merge") {
		options.ApprovalsBeforeMerge = gitlab.Int(d.Get("approvals_before_merge").(int))
		updatedProperties = append(updatedProperties, "approvals_before_merge")
	}

	if d.HasChange("wiki_enabled") {
		options.WikiEnabled = gitlab.Bool(d.Get("wiki_enabled").(bool))
		updatedProperties = append(updatedProperties, "wiki_enabled")
	}

	if d.HasChange("snippets_enabled") {
		options.SnippetsEnabled = gitlab.Bool(d.Get("snippets_enabled").(bool))
		updatedProperties = append(updatedProperties, "snippets_enabled")
	}

	if d.HasChange("shared_runners_enabled") {
		options.SharedRunnersEnabled = gitlab.Bool(d.Get("shared_runners_enabled").(bool))
		updatedProperties = append(updatedProperties, "shared_runners_enabled")
	}

	if d.HasChange("tags") {
		options.TagList = stringSetToStringSlice(d.Get("tags").(*schema.Set))
		updatedProperties = append(updatedProperties, "tags")
	}

	if d.HasChange("container_registry_enabled") {
		options.ContainerRegistryEnabled = gitlab.Bool(d.Get("container_registry_enabled").(bool))
		updatedProperties = append(updatedProperties, "container_registry_enabled")
	}

	if *options != (gitlab.EditProjectOptions{}) {
		log.Printf("[DEBUG] update gitlab project %s", d.Id())
		_, _, err := client.Projects.EditProject(d.Id(), options)
		if err != nil {
			return err
		}
		for _, updatedProperty := range updatedProperties {
			log.Printf("[DEBUG] partial gitlab project %s update of property %q", d.Id(), updatedProperty)
			d.SetPartial(updatedProperty)
		}
	}

	if d.HasChange("shared_with_groups") {
		err := updateSharedWithGroups(d, meta)
		// TODO: check if handling partial state update in this simplistic
		// way is ok when an error in the "shared groups" API calls occurs
		if err != nil {
			d.SetPartial("shared_with_groups")
		}
	}

	if d.HasChange("archived") {
		v := d.Get("archived")
		if v.(bool) {
			err := archiveProject(d, meta)
			if err != nil {
				log.Printf("[WARN] Project (%s) could not be archived: error %#v", d.Id(), err)
				return err
			}
		} else {
			err := unarchiveProject(d, meta)
			if err != nil {
				log.Printf("[WARN] Project (%s) could not be unarchived: error %#v", d.Id(), err)
				return err
			}
		}
		d.SetPartial("archived")
	}

	d.Partial(false)
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

// archiveProject calls the Gitlab server to archive a project; if the
// project is already archived, the call will do nothing (the API is
// idempotent).
func archiveProject(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] Project (%s) will be archived", d.Id())
	client := meta.(*gitlab.Client)
	out, _, err := client.Projects.ArchiveProject(d.Id())
	if err != nil {
		log.Printf("[ERROR] Error archiving project (%s), received %#v", d.Id(), err)
		return err
	}
	if !out.Archived {
		log.Printf("[ERROR] Project (%s) is still not archived", d.Id())
		return fmt.Errorf("error archiving project (%s): its status on the server is still unarchived", d.Id())
	}
	log.Printf("[TRACE] Project (%s) archived", d.Id())
	return nil
}

// unarchiveProject calls the Gitlab server to unarchive a project; if the
// project is already not archived, the call will do nothing (the API is
// idempotent).
func unarchiveProject(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Project (%s) will be unarchived", d.Id())
	client := meta.(*gitlab.Client)
	out, _, err := client.Projects.UnarchiveProject(d.Id())
	if err != nil {
		log.Printf("[ERROR] Error unarchiving project (%s), received %#v", d.Id(), err)
		return err
	}
	if out.Archived {
		log.Printf("[ERROR] Project (%s) is still archived", d.Id())
		return fmt.Errorf("error unarchiving project (%s): its status on the server is still archived", d.Id())
	}
	log.Printf("[TRACE] Project (%s) unarchived", d.Id())
	return nil
}
