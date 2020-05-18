package gitlab

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabGroupCreate,
		Read:   resourceGitlabGroupRead,
		Update: resourceGitlabGroupUpdate,
		Delete: resourceGitlabGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"full_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"full_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"web_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lfs_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"request_access_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"visibility_level": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "internal", "public"}, true),
			},
			"subgroup_creation_level": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"owner", "maintainer"}, true),
			},
			"project_creation_level": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"noone", "developer", "maintainer"}, true),
			},
			"parent_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  0,
			},
			"runners_token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceGitlabGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	options := &gitlab.CreateGroupOptions{
		Name:                 gitlab.String(d.Get("name").(string)),
		LFSEnabled:           gitlab.Bool(d.Get("lfs_enabled").(bool)),
		RequestAccessEnabled: gitlab.Bool(d.Get("request_access_enabled").(bool)),
	}

	if v, ok := d.GetOk("path"); ok {
		options.Path = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		options.Description = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("visibility_level"); ok {
		options.Visibility = stringToVisibilityLevel(v.(string))
	}

	if v, ok := d.GetOk("subgroup_creation_level"); ok {
		options.SubGroupCreationLevel = stringToSubGroupCreationLevel(v.(string))
	}

	if v, ok := d.GetOk("project_creation_level"); ok {
		options.ProjectCreationLevel = stringToProjectCreationLevel(v.(string))
	}

	if v, ok := d.GetOk("parent_id"); ok {
		options.ParentID = gitlab.Int(v.(int))
	}

	log.Printf("[DEBUG] create gitlab group %q", *options.Name)

	group, _, err := client.Groups.CreateGroup(options)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", group.ID))

	return resourceGitlabGroupRead(d, meta)
}

func resourceGitlabGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	log.Printf("[DEBUG] read gitlab group %s", d.Id())

	group, resp, err := client.Groups.GetGroup(d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			log.Printf("[DEBUG] gitlab group %s not found so removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}
	if group.MarkedForDeletionOn != nil {
		log.Printf("[DEBUG] gitlab group %s is marked for deletion", d.Id())
		d.SetId("")
		return nil
	}

	d.SetId(fmt.Sprintf("%d", group.ID))
	d.Set("name", group.Name)
	d.Set("path", group.Path)
	d.Set("full_path", group.FullPath)
	d.Set("full_name", group.FullName)
	d.Set("web_url", group.WebURL)
	d.Set("description", group.Description)
	d.Set("lfs_enabled", group.LFSEnabled)
	d.Set("request_access_enabled", group.RequestAccessEnabled)
	d.Set("visibility_level", group.Visibility)
	d.Set("subgroup_creation_level", group.SubGroupCreationLevel)
	d.Set("project_creation_level", group.ProjectCreationLevel)
	d.Set("parent_id", group.ParentID)
	d.Set("runners_token", group.RunnersToken)

	return nil
}

func resourceGitlabGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	options := &gitlab.UpdateGroupOptions{}

	if d.HasChange("name") {
		options.Name = gitlab.String(d.Get("name").(string))
	}

	if d.HasChange("path") {
		options.Path = gitlab.String(d.Get("path").(string))
	}

	if d.HasChange("description") {
		options.Description = gitlab.String(d.Get("description").(string))
	}

	if d.HasChange("lfs_enabled") {
		options.LFSEnabled = gitlab.Bool(d.Get("lfs_enabled").(bool))
	}

	if d.HasChange("request_access_enabled") {
		options.RequestAccessEnabled = gitlab.Bool(d.Get("request_access_enabled").(bool))
	}

	// Always set visibility ; workaround for
	// https://gitlab.com/gitlab-org/gitlab-ce/issues/38459
	if v, ok := d.GetOk("visibility_level"); ok {
		options.Visibility = stringToVisibilityLevel(v.(string))
	}

	if d.HasChange("subgroup_creation_level") {
		options.SubGroupCreationLevel = stringToSubGroupCreationLevel(d.Get("subgroup_creation_level").(string))
	}

	if d.HasChange("project_creation_level") {
		options.ProjectCreationLevel = stringToProjectCreationLevel(d.Get("project_creation_level").(string))
	}

	log.Printf("[DEBUG] update gitlab group %s", d.Id())

	_, _, err := client.Groups.UpdateGroup(d.Id(), options)
	if err != nil {
		return err
	}

	return resourceGitlabGroupRead(d, meta)
}

func resourceGitlabGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	log.Printf("[DEBUG] Delete gitlab group %s", d.Id())

	_, err := client.Groups.DeleteGroup(d.Id())
	if err != nil && !strings.Contains(err.Error(), "Group has been already marked for deletion") {
		return fmt.Errorf("error deleting group %s: %s", d.Id(), err)
	}

	// Wait for the group to be deleted.
	// Deleting a group in gitlab is async.
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Deleting"},
		Target:  []string{"Deleted"},
		Refresh: func() (interface{}, string, error) {
			out, response, err := client.Groups.GetGroup(d.Id())
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					return out, "Deleted", nil
				}
				log.Printf("[ERROR] Received error: %#v", err)
				return out, "Error", err
			}
			if out.MarkedForDeletionOn != nil {
				// Represents a Gitlab EE soft-delete
				return out, "Deleted", nil
			}
			return out, "Deleting", nil
		},

		Timeout:    10 * time.Minute,
		MinTimeout: 3 * time.Second,
		Delay:      5 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for group (%s) to become deleted: %s", d.Id(), err)
	}
	return err
}
