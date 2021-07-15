package gitlab

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabServicePipelinesEmail() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabServicePipelinesEmailCreate,
		Read:   resourceGitlabServicePipelinesEmailRead,
		Update: resourceGitlabServicePipelinesEmailCreate,
		Delete: resourceGitlabServicePipelinesEmailDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"recipients": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"notify_only_broken_pipelines": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"branches_to_be_notified": {
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

func resourceGitlabServicePipelinesEmailCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	d.SetId(project)
	options := &gitlab.SetPipelinesEmailServiceOptions{
		Recipients:                gitlab.String(strings.Join(*stringSetToStringSlice(d.Get("recipients").(*schema.Set)), ",")),
		NotifyOnlyBrokenPipelines: gitlab.Bool(d.Get("notify_only_broken_pipelines").(bool)),
		BranchesToBeNotified:      gitlab.String(d.Get("branches_to_be_notified").(string)),
	}

	log.Printf("[DEBUG] create gitlab pipelines emails service for project %s", project)

	_, err := client.Services.SetPipelinesEmailService(project, options)
	if err != nil {
		return err
	}

	return resourceGitlabServicePipelinesEmailRead(d, meta)
}

func resourceGitlabServicePipelinesEmailRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Id()

	log.Printf("[DEBUG] read gitlab pipelines emails service for project %s", project)

	service, _, err := client.Services.GetPipelinesEmailService(project)
	if err != nil {
		return err
	}

	d.Set("project", project)
	resourceGitlabServicePipelinesEmailSetToState(d, service)
	return nil
}

func resourceGitlabServicePipelinesEmailDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Id()

	log.Printf("[DEBUG] delete gitlab pipelines email service for project %s", project)

	_, err := client.Services.DeletePipelinesEmailService(project)
	return err
}
