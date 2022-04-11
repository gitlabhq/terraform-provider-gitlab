package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_deploy_key", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_deploy_key`" + ` resource allows to manage the lifecycle of a deploy key.

-> To enable an already existing deploy key for another project use the ` + "`gitlab_project_deploy_key`" + ` resource.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/deploy_keys.html)`,

		CreateContext: resourceGitlabDeployKeyCreate,
		ReadContext:   resourceGitlabDeployKeyRead,
		DeleteContext: resourceGitlabDeployKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGitlabDeployKeyStateImporter,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The name or id of the project to add the deploy key to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"title": {
				Description: "A title to describe the deploy key with.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"key": {
				Description: "The public ssh key body.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old == strings.TrimSpace(new)
				},
			},
			"can_push": {
				Description: "Allow this deploy key to be used to push changes to the project. Defaults to `false`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
			},
		},
	}
})

func resourceGitlabDeployKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	options := &gitlab.AddDeployKeyOptions{
		Title:   gitlab.String(d.Get("title").(string)),
		Key:     gitlab.String(strings.TrimSpace(d.Get("key").(string))),
		CanPush: gitlab.Bool(d.Get("can_push").(bool)),
	}

	log.Printf("[DEBUG] create gitlab deployment key %s", *options.Title)

	deployKey, _, err := client.DeployKeys.AddDeployKey(project, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", deployKey.ID))

	return resourceGitlabDeployKeyRead(ctx, d, meta)
}

func resourceGitlabDeployKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	deployKeyID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] read gitlab deploy key %s/%d", project, deployKeyID)

	deployKey, _, err := client.DeployKeys.GetDeployKey(project, deployKeyID, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] gitlab deploy key not found %s/%d", project, deployKeyID)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("title", deployKey.Title)
	d.Set("key", deployKey.Key)
	d.Set("can_push", deployKey.CanPush)

	return nil
}

func resourceGitlabDeployKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	deployKeyID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Delete gitlab deploy key %s", d.Id())

	_, err = client.DeployKeys.DeleteDeployKey(project, deployKeyID, gitlab.WithContext(ctx))

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGitlabDeployKeyStateImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := strings.Split(d.Id(), ":")
	if len(s) != 2 {
		d.SetId("")
		return nil, fmt.Errorf("invalid Deploy Key import format; expected '{project_id}:{deploy_key_id}' was %v", s)
	}
	project, id := s[0], s[1]

	d.SetId(id)
	d.Set("project", project)

	return []*schema.ResourceData{d}, nil
}
