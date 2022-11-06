package provider

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	gitlab "github.com/xanzy/go-gitlab"
)

var validBranchesToBeNotified = []string{
	"all", "default", "protected", "default_and_protected",
}

var _ = registerResource("gitlab_service_emails_on_push", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_service_emails_on_push`" + ` resource allows to manage the lifecycle of a project integration with Emails on Push Service.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/integrations.html#emails-on-push)`,

		CreateContext: resourceGitlabServiceEmailsOnPushCreate,
		ReadContext:   resourceGitlabServiceEmailsOnPushRead,
		UpdateContext: resourceGitlabServiceEmailsOnPushCreate,
		DeleteContext: resourceGitlabServiceEmailsOnPushDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Description:  "ID or full-path of the project you want to activate integration on.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"recipients": {
				Description:  "Emails separated by whitespace.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"disable_diffs": {
				Description: "Disable code diffs.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"send_from_committer_email": {
				Description: "Send from committer.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"push_events": {
				Description: "Enable notifications for push events.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"tag_push_events": {
				Description: "Enable notifications for tag push events.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"branches_to_be_notified": {
				Description:      fmt.Sprintf("Branches to send notifications for. Valid options are %s. Notifications are always fired for tag pushes.", renderValueListForDocs(validBranchesToBeNotified)),
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "all",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(validBranchesToBeNotified, false)),
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

func resourceGitlabServiceEmailsOnPushCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	options := &gitlab.SetEmailsOnPushServiceOptions{
		Recipients: gitlab.String(d.Get("recipients").(string)),
	}
	if v, ok := d.GetOk("disable_diffs"); ok {
		options.DisableDiffs = gitlab.Bool(v.(bool))
	}
	if v, ok := d.GetOk("send_from_committer_email"); ok {
		options.SendFromCommitterEmail = gitlab.Bool(v.(bool))
	}
	if v, ok := d.GetOk("push_events"); ok {
		options.PushEvents = gitlab.Bool(v.(bool))
	}
	if v, ok := d.GetOk("tag_push_events"); ok {
		options.TagPushEvents = gitlab.Bool(v.(bool))
	}
	if v, ok := d.GetOk("branches_to_be_notified"); ok {
		options.BranchesToBeNotified = gitlab.String(v.(string))
	}

	project := d.Get("project").(string)
	log.Printf("[DEBUG] create gitlab emails on push service for project %s", project)

	_, err := client.Services.SetEmailsOnPushService(project, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(project)

	return resourceGitlabServiceEmailsOnPushRead(ctx, d, meta)
}

func resourceGitlabServiceEmailsOnPushRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Id()

	log.Printf("[DEBUG] read gitlab emails on push service for project %s", project)

	service, _, err := client.Services.GetEmailsOnPushService(project, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] gitlab emails on push service not found for project %s, removing from state", project)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("project", project)
	d.Set("recipients", service.Properties.Recipients)
	d.Set("branches_to_be_notified", service.Properties.BranchesToBeNotified)
	d.Set("disable_diffs", service.Properties.DisableDiffs)
	d.Set("push_events", service.Properties.PushEvents)
	d.Set("send_from_committer_email", service.Properties.SendFromCommitterEmail)
	d.Set("tag_push_events", service.Properties.TagPushEvents)
	d.Set("active", service.Active)
	d.Set("slug", service.Slug)
	d.Set("title", service.Title)
	d.Set("created_at", service.CreatedAt.Format(time.RFC3339))
	if service.UpdatedAt != nil {
		d.Set("updated_at", service.UpdatedAt.Format(time.RFC3339))
	}

	return nil
}

func resourceGitlabServiceEmailsOnPushDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Id()

	log.Printf("[DEBUG] delete gitlab emails on push service for project %s", project)

	_, err := client.Services.DeleteEmailsOnPushService(project, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
