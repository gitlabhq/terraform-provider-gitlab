package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_deploy_key_enable", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_deploy_key_enable`" + ` resource allows to enable an already existing deploy key (see ` + "`gitlab_deploy_key resource`" + `) for a specific project.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/deploy_keys.html#enable-a-deploy-key)`,

		CreateContext: resourceGitlabDeployKeyEnableCreate,
		ReadContext:   resourceGitlabDeployKeyEnableRead,
		DeleteContext: resourceGitlabDeployKeyEnableDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Description: "Can deploy key push to the project's repository.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
			},
		},
	}
})

func resourceGitlabDeployKeyEnableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	key_id, err := strconv.Atoi(d.Get("key_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] enable gitlab deploy key %s/%d", project, key_id)

	deployKey, _, err := client.DeployKeys.EnableDeployKey(project, key_id, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	options := &gitlab.UpdateDeployKeyOptions{
		CanPush: gitlab.Bool(d.Get("can_push").(bool)),
	}
	_, _, err = client.DeployKeys.UpdateDeployKey(project, key_id, options)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s:%d", project, deployKey.ID))

	return resourceGitlabDeployKeyEnableRead(ctx, d, meta)
}

func resourceGitlabDeployKeyEnableRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	project, deployKeyID, err := resourceGitLabDeployKeyEnableParseId(d.Id())
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

	project, deployKeyID, err := resourceGitLabDeployKeyEnableParseId(d.Id())
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

func resourceGitLabDeployKeyEnableParseId(id string) (string, int, error) {
	projectID, deployTokenID, err := parseTwoPartID(id)
	if err != nil {
		return "", 0, err
	}

	deployTokenIID, err := strconv.Atoi(deployTokenID)
	if err != nil {
		return "", 0, err
	}

	return projectID, deployTokenIID, nil
}
