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
			"parent_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  0,
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

	group, _, err := client.Groups.GetGroup(d.Id())
	if err != nil {
		return err
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
	d.Set("parent_id", group.ParentID)

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
	if err != nil {
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
		return fmt.Errorf("error waiting for group (%s) to become deleted: %s", d.Id(), err)
	}
	return err
}
