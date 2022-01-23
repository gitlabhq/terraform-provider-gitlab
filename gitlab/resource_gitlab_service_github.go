package gitlab

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabServiceGithub() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGitlabServiceGithubCreate,
		ReadContext:   resourceGitlabServiceGithubRead,
		UpdateContext: resourceGitlabServiceGithubUpdate,
		DeleteContext: resourceGitlabServiceGithubDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGitlabServiceGithubImportState,
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

func resourceGitlabServiceGithubCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	log.Printf("[DEBUG] create gitlab github service for project %s", project)

	opts := &gitlab.SetGithubServiceOptions{
		Token:         gitlab.String(d.Get("token").(string)),
		RepositoryURL: gitlab.String(d.Get("repository_url").(string)),
		StaticContext: gitlab.Bool(d.Get("static_context").(bool)),
	}

	_, err := client.Services.SetGithubService(project, opts, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGitlabServiceGithubRead(ctx, d, meta)
}

func resourceGitlabServiceGithubRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	log.Printf("[DEBUG] read gitlab github service for project %s", project)

	service, _, err := client.Services.GetGithubService(project, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] gitlab service github not found %s / %s / %s",
				project,
				service.Title,
				service.Properties.RepositoryURL)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	resourceGitlabServiceGithubSetToState(d, service)

	return nil
}

func resourceGitlabServiceGithubUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceGitlabServiceGithubCreate(ctx, d, meta)
}

func resourceGitlabServiceGithubDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	log.Printf("[DEBUG] delete gitlab github service for project %s", project)

	_, err := client.Services.DeleteGithubService(project, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGitlabServiceGithubImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.Set("project", d.Id())

	return []*schema.ResourceData{d}, nil
}
