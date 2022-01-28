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

func resourceGitlabDeployEnableKey() *schema.Resource {
	return &schema.Resource{
		Description: "This resource allows you to enable pre-existing deploy keys for your GitLab projects.\n\n" +
			"> **NOTE**: the GITLAB KEY_ID for the deploy key must be known",

		CreateContext: resourceGitlabDeployKeyEnableCreate,
		ReadContext:   resourceGitlabDeployKeyEnableRead,
		DeleteContext: resourceGitlabDeployKeyEnableDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGitlabDeployKeyEnableStateImporter,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The name or id of the project to add the deploy key to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"key_id": {
				Description: "The Gitlab key id for the pre-existing deploy key",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"title": {
				Description: "Deploy key's title.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"key": {
				Description: "Deploy key.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"can_push": {
				Description: "Can deploy key push to the project’s repository.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceGitlabDeployKeyEnableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	key_id, err := strconv.Atoi(d.Get("key_id").(string)) // nolint // TODO: Resolve this golangci-lint issue: ineffectual assignment to err (ineffassign)

	log.Printf("[DEBUG] enable gitlab deploy key %s/%d", project, key_id)

	deployKey, _, err := client.DeployKeys.EnableDeployKey(project, key_id, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s:%d", project, deployKey.ID))

	return resourceGitlabDeployKeyEnableRead(ctx, d, meta)
}

func resourceGitlabDeployKeyEnableRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	deployKeyID, err := strconv.Atoi(d.Get("key_id").(string))
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
	d.Set("key_id", strconv.Itoa(deployKey.ID))
	d.Set("key", deployKey.Key)
	d.Set("can_push", deployKey.CanPush)
	d.Set("project", project)
	return nil
}

func resourceGitlabDeployKeyEnableDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	deployKeyID, err := strconv.Atoi(d.Get("key_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Delete gitlab deploy key %s/%d", project, deployKeyID)

	response, err := client.DeployKeys.DeleteDeployKey(project, deployKeyID, gitlab.WithContext(ctx))

	if err != nil {
		return diag.FromErr(err)
	}

	// HTTP 2XX is success including 204 with no body
	if response != nil && response.StatusCode/100 == 2 {
		return nil
	}

	return nil
}

func resourceGitlabDeployKeyEnableStateImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := strings.Split(d.Id(), ":")
	if len(s) != 2 {
		d.SetId("")
		return nil, fmt.Errorf("Invalid Deploy Key import format; expected '{project_id}:{deploy_key_id}'")
	}
	project, id := s[0], s[1]

	d.SetId(fmt.Sprintf("%s:%s", project, id))
	d.Set("key_id", id)
	d.Set("project", project)

	return []*schema.ResourceData{d}, nil
}
