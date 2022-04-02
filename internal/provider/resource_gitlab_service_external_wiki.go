package provider

import (
	"context"
	"log"

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
