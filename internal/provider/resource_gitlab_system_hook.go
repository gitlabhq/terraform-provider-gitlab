package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_system_hook", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_system_hook`" + ` resource allows to manage the lifecycle of a system hook.

-> This resource requires GitLab 14.9

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/system_hooks.html)`,

		CreateContext: resourceGitlabSystemHookCreate,
		ReadContext:   resourceGitlabSystemHookRead,
		DeleteContext: resourceGitlabSystemHookDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"url": {
				Description: "The hook URL.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"token": {
				Description: "Secret token to validate received payloads; this isn’t returned in the response. This attribute is not available for imported resources.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ForceNew:    true,
			},
			"push_events": {
				Description: "When true, the hook fires on push events.",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
			},
			"tag_push_events": {
				Description: "When true, the hook fires on new tags being pushed.",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
			},
			"merge_requests_events": {
				Description: "Trigger hook on merge requests events.",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
			},
			"repository_update_events": {
				Description: "Trigger hook on repository update events.",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
			},
			"enable_ssl_verification": {
				Description: "Do SSL verification when triggering the hook.",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
			},
			"created_at": {
				Description: "The date and time the hook was created in ISO8601 format.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
})

func resourceGitlabSystemHookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	options := &gitlab.AddHookOptions{
		URL: gitlab.String(d.Get("url").(string)),
	}
	// NOTE: `GetOkExists()` is deprecated, but until there is a replacement we need to use it.
	//       see https://github.com/hashicorp/terraform-plugin-sdk/pull/350#issuecomment-597888969

	// nolint:staticcheck // SA1019 ignore deprecated GetOkExists
	// lintignore: XR001 // TODO: replace with alternative for GetOkExists
	if v, ok := d.GetOkExists("token"); ok {
		options.Token = gitlab.String(v.(string))
	}
	// nolint:staticcheck // SA1019 ignore deprecated GetOkExists
	// lintignore: XR001 // TODO: replace with alternative for GetOkExists
	if v, ok := d.GetOkExists("push_events"); ok {
		options.PushEvents = gitlab.Bool(v.(bool))
	}
	// nolint:staticcheck // SA1019 ignore deprecated GetOkExists
	// lintignore: XR001 // TODO: replace with alternative for GetOkExists
	if v, ok := d.GetOkExists("tag_push_events"); ok {
		options.TagPushEvents = gitlab.Bool(v.(bool))
	}
	// nolint:staticcheck // SA1019 ignore deprecated GetOkExists
	// lintignore: XR001 // TODO: replace with alternative for GetOkExists
	if v, ok := d.GetOkExists("merge_requests_events"); ok {
		options.MergeRequestsEvents = gitlab.Bool(v.(bool))
	}
	// nolint:staticcheck // SA1019 ignore deprecated GetOkExists
	// lintignore: XR001 // TODO: replace with alternative for GetOkExists
	if v, ok := d.GetOkExists("repository_update_events"); ok {
		options.RepositoryUpdateEvents = gitlab.Bool(v.(bool))
	}
	// nolint:staticcheck // SA1019 ignore deprecated GetOkExists
	// lintignore: XR001 // TODO: replace with alternative for GetOkExists
	if v, ok := d.GetOkExists("enable_ssl_verification"); ok {
		options.EnableSSLVerification = gitlab.Bool(v.(bool))
	}

	log.Printf("[DEBUG] create gitlab system hook %q", *options.URL)

	hook, _, err := client.SystemHooks.AddHook(options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", hook.ID))
	d.Set("token", options.Token)
	return resourceGitlabSystemHookRead(ctx, d, meta)
}

func resourceGitlabSystemHookRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	hookID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] read gitlab system hook %d", hookID)

	hook, _, err := client.SystemHooks.GetHook(hookID, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] gitlab system hook not found %d, removing from state", hookID)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("url", hook.URL)
	d.Set("push_events", hook.PushEvents)
	d.Set("tag_push_events", hook.TagPushEvents)
	d.Set("merge_requests_events", hook.MergeRequestsEvents)
	d.Set("repository_update_events", hook.RepositoryUpdateEvents)
	d.Set("enable_ssl_verification", hook.EnableSSLVerification)
	d.Set("created_at", hook.CreatedAt.Format(time.RFC3339))
	return nil
}

func resourceGitlabSystemHookDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	hookID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Delete gitlab system hook %s", d.Id())

	_, err = client.SystemHooks.DeleteHook(hookID, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
