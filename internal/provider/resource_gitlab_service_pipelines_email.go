package provider

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabServicePipelinesEmail() *schema.Resource {
	return &schema.Resource{
		Description: "This resource manages a [Pipelines email integration](https://docs.gitlab.com/ee/user/project/integrations/overview.html#integrations-listing) that emails the pipeline status to a list of recipients.",

		CreateContext: resourceGitlabServicePipelinesEmailCreate,
		ReadContext:   resourceGitlabServicePipelinesEmailRead,
		UpdateContext: resourceGitlabServicePipelinesEmailCreate,
		DeleteContext: resourceGitlabServicePipelinesEmailDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Description: "ID of the project you want to activate integration on.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"recipients": {
				Description: ") email addresses where notifications are sent.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"notify_only_broken_pipelines": {
				Description: "Notify only broken pipelines. Default is true.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"branches_to_be_notified": {
				Description:  "Branches to send notifications for. Valid options are `all`, `default`, `protected`, and `default_and_protected`. Default is `default`",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"all", "default", "protected", "default_and_protected"}, true),
				Default:      "default",
			},
		},
	}
}

func resourceGitlabServicePipelinesEmailSetToState(d *schema.ResourceData, service *gitlab.PipelinesEmailService) {
	d.Set("recipients", strings.Split(service.Properties.Recipients, ",")) // lintignore: XR004 // TODO: Resolve this tfproviderlint issue
	d.Set("notify_only_broken_pipelines", service.Properties.NotifyOnlyBrokenPipelines)
	d.Set("branches_to_be_notified", service.Properties.BranchesToBeNotified)
}

func resourceGitlabServicePipelinesEmailCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	d.SetId(project)
	options := &gitlab.SetPipelinesEmailServiceOptions{
		Recipients:                gitlab.String(strings.Join(*stringSetToStringSlice(d.Get("recipients").(*schema.Set)), ",")),
		NotifyOnlyBrokenPipelines: gitlab.Bool(d.Get("notify_only_broken_pipelines").(bool)),
		BranchesToBeNotified:      gitlab.String(d.Get("branches_to_be_notified").(string)),
	}

	log.Printf("[DEBUG] create gitlab pipelines emails service for project %s", project)

	_, err := client.Services.SetPipelinesEmailService(project, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGitlabServicePipelinesEmailRead(ctx, d, meta)
}

func resourceGitlabServicePipelinesEmailRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Id()

	log.Printf("[DEBUG] read gitlab pipelines emails service for project %s", project)

	service, _, err := client.Services.GetPipelinesEmailService(project, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] gitlab pipelines emails service not found for project %s", project)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("project", project)
	resourceGitlabServicePipelinesEmailSetToState(d, service)
	return nil
}

func resourceGitlabServicePipelinesEmailDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Id()

	log.Printf("[DEBUG] delete gitlab pipelines email service for project %s", project)

	_, err := client.Services.DeletePipelinesEmailService(project, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
