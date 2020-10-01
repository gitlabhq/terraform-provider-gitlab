package gitlab

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func resourceGitlabDeployEnableKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabDeployKeyEnableCreate,
		Read:   resourceGitlabDeployKeyEnableRead,
		Delete: resourceGitlabDeployKeyEnableDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGitlabDeployKeyEnableStateImporter,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"title": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"can_push": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceGitlabDeployKeyEnableCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	keyId, err := strconv.Atoi(d.Get("key_id").(string))

	log.Printf("[DEBUG] enable gitlab deploy key %s/%d", project, keyId)

	deployKey, _, err := client.DeployKeys.EnableDeployKey(project, keyId)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s:%d", project, deployKey.ID))

	return resourceGitlabDeployKeyEnableRead(d, meta)
}

func resourceGitlabDeployKeyEnableRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	deployKeyID, err := strconv.Atoi(d.Get("key_id").(string))
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] read gitlab deploy key %s/%d", project, deployKeyID)

	deployKey, _, err := client.DeployKeys.GetDeployKey(project, deployKeyID)
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] gitlab deploy key not found %s/%d", project, deployKeyID)
			d.SetId("")
			return nil
		}
		return err
	}

	_ = d.Set("title", deployKey.Title)
	_ = d.Set("key_id", deployKey.ID)
	_ = d.Set("key", deployKey.Key)
	_ = d.Set("can_push", deployKey.CanPush)
	_ = d.Set("project", project)
	return nil
}

func resourceGitlabDeployKeyEnableDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	deployKeyID, err := strconv.Atoi(d.Get("key_id").(string))
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Delete gitlab deploy key %s/%d", project, deployKeyID)

	response, err := client.DeployKeys.DeleteDeployKey(project, deployKeyID)

	// HTTP 2XX is success including 204 with no body
	if response != nil && response.StatusCode/100 == 2 {
		return nil
	}
	return err
}

func resourceGitlabDeployKeyEnableStateImporter(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	s := strings.Split(d.Id(), ":")
	if len(s) != 2 {
		d.SetId("")
		return nil, fmt.Errorf("invalid deploy key import format; expected '{project_id}:{deploy_key_id}'")
	}
	project, id := s[0], s[1]

	d.SetId(fmt.Sprintf("%s:%s", project, id))
	_ = d.Set("key_id", id)
	_ = d.Set("project", project)

	return []*schema.ResourceData{d}, nil
}
