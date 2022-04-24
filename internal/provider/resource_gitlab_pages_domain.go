package provider

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_pages_domain", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_pages_domain`" + ` resource enables endpoints for connecting custom domain(s) and TLS certificates in GitLab Pages.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/pages_domains.html)`,

		CreateContext: resourceGitlabPagesDomainCreate,
		ReadContext:   resourceGitlabPagesDomainRead,
		UpdateContext: resourceGitlabPagesDomainUpdate,
		DeleteContext: resourceGitlabPagesDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Description: "The custom domain.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"project": {
				Description: "The ID or full path of the project which the branch is created against.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},
			"auto_ssl_enabled": {
				Description: "Enables automatic generation of SSL certificates issued by Let’s Encrypt for custom domains.",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				Default:     nil,
			},
			"certificate": {
				Description: "The certificate in PEM format with intermediates following in most specific to least specific order.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Default:     nil,
			},
			"key": {
				Description: "The certificate key in PEM format.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Default:     nil,
			},
			"url": {
				Description: "The URL for the given domain.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"certificate_data": {
				Description: "The certificate data.",
				Type:        schema.TypeMap,
				Computed:    true,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"verified": {
				Description: "The certificate data.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"verification_code": {
				Description: "The verification code for the domain.",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
})

func resourceGitlabPagesDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	projectId := d.Get("project").(string)
	domain := d.Get("domain").(string)
	auto_ssl_enabled := d.Get("auto_ssl_enabled").(bool)
	certificate := d.Get("certificate").(string)
	key := d.Get("key").(string)

	options := &gitlab.CreatePagesDomainOptions{
		Domain:         &domain,
		AutoSslEnabled: &auto_ssl_enabled,
		Certificate:    &certificate,
		Key:            &key,
	}

	log.Printf("[DEBUG] create gitlab pages domain %s", domain)

	_, _, err := client.PagesDomains.CreatePagesDomain(projectId, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(buildTwoPartID(&projectId, &domain))
	return resourceGitlabPagesDomainRead(ctx, d, meta)
}

func resourceGitlabPagesDomainUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	projectID := d.Get("project").(string)
	domain := d.Get("domain").(string)
	auto_ssl_enabled := d.Get("auto_ssl_enabled").(bool)
	certificate := d.Get("certificate").(string)
	key := d.Get("key").(string)

	options := &gitlab.UpdatePagesDomainOptions{
		AutoSslEnabled: &auto_ssl_enabled,
		Certificate:    &certificate,
		Key:            &key,
	}
	log.Printf("[DEBUG] update gitlab pages domain %s for %s", domain, projectID)

	_, _, err := client.PagesDomains.UpdatePagesDomain(projectID, domain, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceGitlabProjectMirrorRead(ctx, d, meta)
}

type Certificate struct {
	expired    bool
	expiration string
}

func resourceGitlabPagesDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	projectID, domain, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] read gitlab pages domain %s", domain)

	pagesDomain, _, err := client.PagesDomains.GetPagesDomain(projectID, domain, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] gitlab pages domain %s not found, removing from state", domain)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	var certificate_expiration string
	if pagesDomain.Certificate.Expiration == nil {
		certificate_expiration = ""
	} else {
		certificate_expiration = pagesDomain.Certificate.Expiration.UTC().Format(time.UnixDate)
	}
	certificate_data := map[string]string{
		"expired":    strconv.FormatBool(pagesDomain.Certificate.Expired),
		"expiration": certificate_expiration,
	}

	d.Set("project", projectID)
	d.Set("domain", pagesDomain.Domain)
	d.Set("url", pagesDomain.URL)
	d.Set("auto_ssl_enabled", pagesDomain.AutoSslEnabled)
	d.Set("certificate_data", certificate_data)
	d.Set("verified", pagesDomain.Verified)
	d.Set("verification_code", pagesDomain.VerificationCode)
	return nil
}

func resourceGitlabPagesDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	projectID, domain, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Delete gitlab pages domain %s", domain)

	_, err = client.PagesDomains.DeletePagesDomain(projectID, domain, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
