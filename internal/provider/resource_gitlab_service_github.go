package provider

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
		Description: "**NOTE**: requires either EE (self-hosted) or Silver and above (GitLab.com).\n\n" +
			"This resource manages a [GitHub integration](https://docs.gitlab.com/ee/user/project/integrations/github.html) that updates pipeline statuses on a GitHub repo's pull requests.",

		CreateContext: resourceGitlabServiceGithubCreate,
		ReadContext:   resourceGitlabServiceGithubRead,
		UpdateContext: resourceGitlabServiceGithubUpdate,
		DeleteContext: resourceGitlabServiceGithubDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGitlabServiceGithubImportState,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Description: "ID of the project you want to activate integration on.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"token": {
				Description: "A GitHub personal access token with at least `repo:status` scope.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"repository_url": {
				Description: "The URL of the GitHub repo to integrate with, e,g, https://github.com/gitlabhq/terraform-provider-gitlab.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"static_context": {
				Description: "Append instance name instead of branch to the status. Must enable to set a GitLab status check as _required_ in GitHub. See [Static / dynamic status check names] to learn more.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},

			// Computed from the GitLab API. Omitted event fields because they're always true in Github.
			"title": {
				Description: "Title.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_at": {
				Description: "Create time.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "Update time.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"active": {
				Description: "Whether the integration is active.",
				Type:        schema.TypeBool,
				Computed:    true,
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
