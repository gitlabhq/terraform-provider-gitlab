package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_project_error_tracking_settings", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`" + `gitlab_project_error_tracking_settings` + "`" + ` resource allows to manage the lifecycle of the error tracking configuration for a project.

-> This resource requires maintainer privileges on the GitLab project you are configuring.

-> This resource can only be used after first enabling Error Tracking in the GitLab UI first on the project you wish to configure.  Failure to do so will result in an error.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/error_tracking.html)`,

		CreateContext: resourceGitlabProjectErrorTrackingSettingsCreate,
		ReadContext:   resourceGitlabProjectErrorTrackingSettingsRead,
		DeleteContext: resourceGitlabProjectErrorTrackingSettingsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Description: "The id of the project to configure error tracking for.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"integrated": {
				Description: "Enable or disable the integrated error tracking backend of a GitLab project.",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
			},
			"active": {
				Description: "Enable or disable the error tracking configuration of a project.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"project_name": {
				Description: "Project name in Sentry which is only utilized if an external Sentry instance is used.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"sentry_external_url": {
				Description: "Web URL of the project in Sentry which is only utilized if an external Sentry instance is used.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"api_url": {
				Description: "API URL of the project in Sentry which is only utilized if an external Sentry instance is used.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
})

func resourceGitlabProjectErrorTrackingSettingsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectID := d.Get("project_id").(string)
	d.SetId(projectID)

	options := gitlab.EnableDisableErrorTrackingOptions{
		Active:     gitlab.Bool(true),
		Integrated: gitlab.Bool(d.Get("integrated").(bool)),
	}

	log.Printf("[DEBUG] Project %s configure gitlab project-level error tracking %+v", projectID, options)

	client := meta.(*gitlab.Client)
	_, _, err := client.ErrorTracking.EnableDisableErrorTracking(projectID, &options, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			return diag.Errorf("Error Tracking must be enabled in the GitLab UI first before using this resource.")
		}
		return diag.FromErr(err)
	}

	return resourceGitlabProjectErrorTrackingSettingsRead(ctx, d, meta)
}

func resourceGitlabProjectErrorTrackingSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Id()

	log.Printf("[DEBUG] Read Gitlab Error Tracking Settings %s", project)

	ErrorTrackingSettings, _, err := client.ErrorTracking.GetErrorTrackingSettings(project, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] Gitlab Error tracking configuration not found %s, removing from state", project)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("project_id", project)
	d.Set("active", ErrorTrackingSettings.Active)
	d.Set("integrated", ErrorTrackingSettings.Integrated)
	d.Set("project_name", ErrorTrackingSettings.ProjectName)
	d.Set("sentry_external_url", ErrorTrackingSettings.SentryExternalURL)
	d.Set("api_url", ErrorTrackingSettings.APIURL)

	return nil
}

func resourceGitlabProjectErrorTrackingSettingsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	project := d.Get("project_id").(string)

	log.Printf("[DEBUG] Delete Gitlab Project Error Tracking settings %s", d.Id())

	options := gitlab.EnableDisableErrorTrackingOptions{
		Active: gitlab.Bool(false),
	}

	_, _, err := client.ErrorTracking.EnableDisableErrorTracking(project, &options, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
