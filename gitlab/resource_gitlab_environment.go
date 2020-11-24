package gitlab

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func resourceGitlabEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabEnvironmentCreate,
		Read:   resourceGitlabEnvironmentRead,
		Update: resourceGitlabEnvironmentUpdate,
		Delete: resourceGitlabEnvironmentDelete,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"external_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"environment_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceGitlabEnvironmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	options := &gitlab.CreateEnvironmentOptions{
		Name:        gitlab.String(d.Get("name").(string)),
		ExternalURL: gitlab.String(d.Get("external_url").(string)),
	}

	environment, _, err := client.Environments.CreateEnvironment(project, options)
	if err != nil {
		return err
	}

	d.Set("environment_id", environment.ID)
	d.SetId(fmt.Sprintf("%s:%d", project, environment.ID))

	return resourceGitlabEnvironmentRead(d, meta)
}

func resourceGitlabEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	environmentID := d.Get("environment_id").(int)

	environment, _, err := client.Environments.GetEnvironment(project, environmentID)
	if err != nil {
		return err
	}

	d.Set("project", project)
	d.Set("name", environment.Name)
	d.Set("external_url", environment.ExternalURL)
	d.Set("slug", environment.Slug)
	d.Set("state", environment.State)

	return nil
}

func resourceGitlabEnvironmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	options := &gitlab.EditEnvironmentOptions{
		Name:        gitlab.String(d.Get("name").(string)),
		ExternalURL: gitlab.String(d.Get("external_url").(string)),
	}
	environmentID := d.Get("environment_id").(int)

	if d.HasChange("name") {
		options.Name = gitlab.String(d.Get("name").(string))
	}

	if d.HasChange("external_url") {
		options.ExternalURL = gitlab.String(d.Get("external_url").(string))
	}

	_, _, err := client.Environments.EditEnvironment(project, environmentID, options)
	if err != nil {
		return err
	}

	return resourceGitlabEnvironmentRead(d, meta)
}

func resourceGitlabEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	environmentID := d.Get("environment_id").(int)

	// Environment must be stopped prior to deletion or a 403 will be received
	stopResp, err := client.Environments.StopEnvironment(project, environmentID)
	if err != nil {
		return fmt.Errorf("%s failed to stop environment prior to delete: (%s) %v", d.Id(), stopResp.Status, err)
	}

	resp, err := client.Environments.DeleteEnvironment(project, environmentID)
	if err != nil {
		return fmt.Errorf("%s failed to delete environment: (%s) %v", d.Id(), resp.Status, err)
	}
	return err
}
