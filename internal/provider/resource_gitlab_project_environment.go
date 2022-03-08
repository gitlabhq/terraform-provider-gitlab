package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	gitlab "github.com/xanzy/go-gitlab"
)

// https://docs.gitlab.com/ee/ci/environments/protected_environments.html
var _ = registerResource("gitlab_project_environment", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_project_environment`" + ` resource you to create and manage an environment in your GitLab project`,

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
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"external_url": {
				Description:  "Place to link to for this environment",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"state": {
				Description: "State the environment is in. Accepted values: available or stopped.",
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

	d.SetId(buildTwoPartID(&project, gitlab.String(fmt.Sprintf("%v", environment.ID))))

	return resourceGitlabProjectEnvironmentRead(ctx, d, meta)
}

func resourceGitlabProjectEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] read gitlab environment %s", d.Id())

	project, environmentID, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	environmentIDInt, err := strconv.Atoi(environmentID)

	if err == nil {
		fmt.Printf("%T \n %v", environmentIDInt, environmentIDInt)
	}

	log.Printf("[DEBUG] Project %s read gitlab environment %d", project, environmentIDInt)

	client := meta.(*gitlab.Client)

	environment, _, err := client.Environments.GetEnvironment(project, environmentIDInt, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] Project %s gitlab environment %q not found", project, environmentID)
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting gitlab project %q environment %q: %v", project, environmentID, err)
	}

	d.SetId(buildTwoPartID(&project, gitlab.String(fmt.Sprintf("%v", environment.ID))))
	d.Set("project", project)
	d.Set("name", environment.Name)
	d.Set("state", environment.State)

	return nil
}

func resourceGitlabProjectEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] update gitlab environment %s", d.Id())

	project, environmentID, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	environmentIDInt, err := strconv.Atoi(environmentID)
	if err != nil {
		return diag.Errorf("error converting environment ID to int: %v", err)
	}

	name := d.Get("name").(string)
	options := gitlab.EditEnvironmentOptions{
		Name: &name,
	}
	if d.HasChange("external_url") {
		options.ExternalURL = gitlab.String(d.Get("external_url").(string))
	}

	log.Printf("[DEBUG] Project %s update gitlab environment %d", project, environmentIDInt)

	client := meta.(*gitlab.Client)

	if _, _, err := client.Environments.EditEnvironment(project, environmentIDInt, &options, gitlab.WithContext(ctx)); err != nil {
		return diag.Errorf("error editing gitlab project %q environment %q: %v", project, environmentID, err)
	}

	return resourceGitlabProjectEnvironmentRead(ctx, d, meta)
}

func resourceGitlabProjectEnvironmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	project, environmentID, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	environmentIDInt, err := strconv.Atoi(environmentID)
	if err != nil {
		return diag.Errorf("error converting environment ID to int: %v", err)
	}

	log.Printf("[DEBUG] Project %s delete gitlab project-level environment %v", project, environmentIDInt)

	client := meta.(*gitlab.Client)

	_, err = client.Environments.StopEnvironment(project, environmentIDInt, gitlab.WithContext(ctx))
	if err != nil {
		return diag.Errorf("error stopping gitlab project %q environment %q: %v", project, environmentID, err)
	}
	_, err = client.Environments.DeleteEnvironment(project, environmentIDInt, gitlab.WithContext(ctx))
	if err != nil {
		return diag.Errorf("error deleting gitlab project %q environment %q: %v", project, environmentID, err)
	}

	return nil
}
