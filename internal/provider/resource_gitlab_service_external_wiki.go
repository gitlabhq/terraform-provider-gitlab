package provider

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_service_external_wiki", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_service_external_wiki`" + ` resource allows to manage the lifecycle of a project integration with External Wiki Service.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/integrations.html#external-wiki)`,

		CreateContext: resourceGitlabServiceExternalWikiCreate,
		ReadContext:   resourceGitlabServiceExternalWikiRead,
		UpdateContext: resourceGitlabServiceExternalWikiCreate,
		DeleteContext: resourceGitlabServiceExternalWikiDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Description:  "ID of the project you want to activate integration on.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"external_wiki_url": {
				Description:  "The URL of the external wiki.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"title": {
				Description: "Title of the integration.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_at": {
				Description: "The ISO8601 date/time that this integration was activated at in UTC.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "The ISO8601 date/time that this integration was last updated at in UTC.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"slug": {
				Description: "The name of the integration in lowercase, shortened to 63 bytes, and with everything except 0-9 and a-z replaced with -. No leading / trailing -. Use in URLs, host names and domain names.",
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
})

func resourceGitlabServiceExternalWikiCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	d.SetId(project)

	options := &gitlab.SetExternalWikiServiceOptions{
		ExternalWikiURL: gitlab.String(d.Get("external_wiki_url").(string)),
	}

	log.Printf("[DEBUG] create gitlab external wiki service for project %s", project)

	_, err := client.Services.SetExternalWikiService(project, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGitlabServiceExternalWikiRead(ctx, d, meta)
}

func resourceGitlabServiceExternalWikiRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Id()

	log.Printf("[DEBUG] read gitlab external wiki service for project %s", project)

	service, _, err := client.Services.GetExternalWikiService(project, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] gitlab external wiki service not found for project %s", project)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("project", project)
	d.Set("external_wiki_url", service.Properties.ExternalWikiURL)
	d.Set("active", service.Active)
	d.Set("slug", service.Slug)
	d.Set("title", service.Title)
	d.Set("created_at", service.CreatedAt.Format(time.RFC3339))
	if service.UpdatedAt != nil {
		d.Set("updated_at", service.UpdatedAt.Format(time.RFC3339))
	}

	return nil
}

func resourceGitlabServiceExternalWikiDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Id()

	log.Printf("[DEBUG] delete gitlab external wiki service for project %s", project)

	_, err := client.Services.DeleteExternalWikiService(project, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
