package gitlab

import (
	"fmt"
	"log"

	gitlab "github.com/Fourcast/go-gitlab"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceGitlabServiceGithub() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabServiceGithubCreate,
		Read:   resourceGitlabServiceGithubRead,
		Update: resourceGitlabServiceGithubUpdate,
		Delete: resourceGitlabServiceGithubDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGitlabServiceGithubImportState,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"token": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"repository_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"static_context": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			// Computed from the GitLab API. Omitted event fields because they're always true in Github.
			"title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceGitlabServiceGithubSetToState(d *schema.ResourceData, service *gitlab.GithubService) {
	d.SetId(fmt.Sprintf("%d", service.ID))
	d.Set("repository_url", service.Properties.RepositoryURL)
	d.Set("static_context", service.Properties.StaticContext)

	d.Set("title", service.Title)
	d.Set("created_at", service.CreatedAt.String())
	d.Set("updated_at", service.UpdatedAt.String())
	d.Set("active", service.Active)
}

func resourceGitlabServiceGithubCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	log.Printf("[DEBUG] create gitlab github service for project %s", project)

	opts := &gitlab.SetGithubServiceOptions{
		Token:         gitlab.String(d.Get("token").(string)),
		RepositoryURL: gitlab.String(d.Get("repository_url").(string)),
		StaticContext: gitlab.Bool(d.Get("static_context").(bool)),
	}

	_, err := client.Services.SetGithubService(project, opts)
	if err != nil {
		return err
	}

	return resourceGitlabServiceGithubRead(d, meta)
}

func resourceGitlabServiceGithubRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	log.Printf("[DEBUG] read gitlab github service for project %s", project)

	service, _, err := client.Services.GetGithubService(project)
	if err != nil {
		return err
	}

	resourceGitlabServiceGithubSetToState(d, service)

	return nil
}

func resourceGitlabServiceGithubUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceGitlabServiceGithubCreate(d, meta)
}

func resourceGitlabServiceGithubDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	log.Printf("[DEBUG] delete gitlab github service for project %s", project)

	_, err := client.Services.DeleteGithubService(project)
	return err
}

func resourceGitlabServiceGithubImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.Set("project", d.Id())

	return []*schema.ResourceData{d}, nil
}
