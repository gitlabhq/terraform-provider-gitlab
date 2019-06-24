package gitlab

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabEnvironmentCreation() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabEnvironmentCreate,
		Read:   resourceGitlabEnvironmentRead,
		Delete: resourceGitlabEnvironmentDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"external_url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceGitlabEnvironmentCreate(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(*gitlab.Client)
	if !ok {
		return errors.New("meta not of expected type *gitlab.Client")
	}

	options := &gitlab.CreateEnvironmentOptions{
		Name:        gitlab.String(d.Get("name").(string)),
		ExternalURL: gitlab.String(d.Get("external_url").(string)),
	}

	project := d.Get("project").(string)

	env, _, err := client.Environments.CreateEnvironment(project, options)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] created gitlab environment for project %s, environment %s", project, env)

	d.SetId(fmt.Sprintf("%d", env.ID))

	return resourceGitlabEnvironmentRead(d, meta)
}

func resourceGitlabEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(*gitlab.Client)
	if !ok {
		return errors.New("meta not of expected type *gitlab.Client")
	}

	project := d.Get("project").(string)

	log.Printf("[DEBUG] try to read environment for project %s, id %s", project, d.Id())

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return errors.New("Can not convert id to int id")
	}

	env, _, err := client.Environments.GetEnvironment(project, id)
	if err != nil {
		return err
	}

	d.Set("name", env.Name)
	d.Set("external_url", env.ExternalURL)
	d.Set("slug", env.Slug)

	log.Printf("[DEBUG] read gitlab environment for project %s, branch %s", project, env.Name)

	return nil
}

func resourceGitlabEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(*gitlab.Client)
	if !ok {
		return errors.New("meta not of expected type *gitlab.Cclient")
	}

	project := d.Get("project").(string)
	name := d.Get("name").(string)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return errors.New("Can not convert id to int id")
	}

	log.Printf("[DEBUG] Delete gitlab environment %s for project %s", name, project)

	_, err = client.Environments.DeleteEnvironment(project, id)
	return err
}
