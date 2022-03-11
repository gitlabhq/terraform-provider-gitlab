package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_project_environment", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_project_environment`" + ` resource allows to manage the lifecycle of an environment in a project.

-> During a terraform destroy this resource by default will not attempt to stop the environment first. 
An environment is required to be in a stopped state before a deletetion of the environment can occur. 
Set the ` + "`stop_before_destroy`" + ` flag to attempt to automatically stop the environment before deletion.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/environments.html`,

		CreateContext: resourceGitlabProjectEnvironmentCreate,
		ReadContext:   resourceGitlabProjectEnvironmentRead,
		UpdateContext: resourceGitlabProjectEnvironmentUpdate,
		DeleteContext: resourceGitlabProjectEnvironmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"project": {
				Description:  "The ID or full path of the project to environment is created for.",
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"name": {
				Description:  "The name of the environment",
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"external_url": {
				Description:  "Place to link to for this environment",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"slug": {
				Description: "The name of the environment in lowercase, shortened to 63 bytes, and with everything except 0-9 and a-z replaced with -. No leading / trailing -. Use in URLs, host names and domain names.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"stop_before_destroy": {
				Description: "Determines whether the environment is attempted to be stopped before the environment is deleted.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"created_at": {
				Description: "The ISO8601 date/time that this environment was created at in UTC",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "The ISO8601 date/time that this environment was last updated at in UTC",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"state": {
				Description: fmt.Sprintf("State the environment is in. Valid values are %s.", renderValueListForDocs(validProjectEnvironmentStates)),
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
})

func resourceGitlabProjectEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	options := gitlab.CreateEnvironmentOptions{
		Name: &name,
	}
	if externalURL, ok := d.GetOk("external_url"); ok {
		options.ExternalURL = gitlab.String(externalURL.(string))
	}

	project := d.Get("project").(string)

	log.Printf("[DEBUG] Project %s create gitlab environment %q", project, *options.Name)

	client := meta.(*gitlab.Client)

	environment, _, err := client.Environments.CreateEnvironment(project, &options, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			return diag.Errorf("feature Environments is not available")
		}
		return diag.FromErr(err)
	}

	environmentID := fmt.Sprintf("%d", environment.ID)
	d.SetId(buildTwoPartID(&project, &environmentID))

	return resourceGitlabProjectEnvironmentRead(ctx, d, meta)
}

func resourceGitlabProjectEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] read gitlab environment %s", d.Id())

	project, environmentIDInt, err := environmentAndProjectFromResource(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Project %s read gitlab environment %d", project, environmentIDInt)

	client := meta.(*gitlab.Client)

	environment, _, err := client.Environments.GetEnvironment(project, environmentIDInt, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] Project %s gitlab environment %q not found, removing from state", project, environmentIDInt)
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting gitlab project %q environment %q: %v", project, environmentIDInt, err)
	}

	d.Set("project", project)
	d.Set("name", environment.Name)
	d.Set("state", environment.State)
	d.Set("external_url", environment.ExternalURL)
	d.Set("created_at", environment.CreatedAt.Format(time.RFC3339))

	if environment.UpdatedAt != nil {
		d.Set("updated_at", environment.UpdatedAt.Format(time.RFC3339))
	}

	return nil
}

func resourceGitlabProjectEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] update gitlab environment %s", d.Id())

	project, environmentIDInt, err := environmentAndProjectFromResource(d)
	if err != nil {
		return diag.FromErr(err)
	}

	options := &gitlab.EditEnvironmentOptions{
		Name: gitlab.String(d.Get("name").(string)),
	}

	if d.HasChange("external_url") {
		options.ExternalURL = gitlab.String(d.Get("external_url").(string))
	}

	log.Printf("[DEBUG] Project %s update gitlab environment %d", project, environmentIDInt)

	client := meta.(*gitlab.Client)

	if _, _, err := client.Environments.EditEnvironment(project, environmentIDInt, options, gitlab.WithContext(ctx)); err != nil {
		return diag.Errorf("error editing gitlab project %q environment %q: %v", project, environmentIDInt, err)
	}

	return resourceGitlabProjectEnvironmentRead(ctx, d, meta)
}

func resourceGitlabProjectEnvironmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, environmentIDInt, err := environmentAndProjectFromResource(d)
	if err != nil {
		return diag.FromErr(err)
	}

	stopBeforeDestroy := d.Get("stop_before_destroy").(bool)

	if stopBeforeDestroy {
		log.Printf("[DEBUG] Stopping environment %v for Project %s before destruction", environmentIDInt, project)
		_, err = client.Environments.StopEnvironment(project, environmentIDInt, gitlab.WithContext(ctx))
		if err != nil {
			return diag.Errorf("error stopping gitlab project %q environment %q: %v", project, environmentIDInt, err)
		}
	} else {
		environment, _, err := client.Environments.GetEnvironment(project, environmentIDInt, gitlab.WithContext(ctx))
		if err != nil {
			if is404(err) {
				log.Printf("[DEBUG] Project %s gitlab environment %q not found, removing from state", project, environmentIDInt)
				d.SetId("")
				return nil
			}
			return diag.Errorf("error getting gitlab project %q environment %q: %v", project, environmentIDInt, err)
		}

		if environment.State != "stopped" {
			return diag.Errorf("[ERROR] cannot destroy gitlab project %q environment %q: Environment must be in a stopped state before deletion. Set stop_before_destroy flag to attempt to auto stop the environment on destruction", project, environmentIDInt)
		}
	}

	_, err = client.Environments.DeleteEnvironment(project, environmentIDInt, gitlab.WithContext(ctx))
	if err != nil {
		return diag.Errorf("error deleting gitlab project %q environment %q: %v", project, environmentIDInt, err)
	}

	return nil
}

func environmentAndProjectFromResource(d *schema.ResourceData) (string, int, error) {
	project, environmentID, err := parseTwoPartID(d.Id())

	if err != nil {
		log.Printf("[ERROR] cannot get project and environmentID from input: %v", d.Id())
		return "", 0, err
	}

	environmentIDInt, e := strconv.Atoi(environmentID)

	if e != nil {
		log.Printf("[ERROR] cannot convert environment ID to int: %v", e)
		return "", 0, e
	}

	return project, environmentIDInt, nil
}
