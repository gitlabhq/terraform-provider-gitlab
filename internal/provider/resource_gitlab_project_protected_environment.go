package provider

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_project_protected_environment", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_project_protected_environment`" + ` resource allows to manage the lifecycle of a protected environment in a project.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/protected_environments.html)`,

		CreateContext: resourceGitlabProjectProtectedEnvironmentCreate,
		ReadContext:   resourceGitlabProjectProtectedEnvironmentRead,
		DeleteContext: resourceGitlabProjectProtectedEnvironmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"project": {
				Description:  "The ID or full path of the project which the protected environment is created against.",
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"environment": {
				Description:  "The name of the environment.",
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			// Uncomment and validate after 14.9 is released, update acceptance tests
			// 	"required_approval_count": {
			// 		Description:  "The number of approvals required to deploy to this environment. This is part of Deployment Approvals, which isn't yet available for use.",
			// 		Type:         schema.TypeString,
			// 		ForceNew:     true,
			// 		Required:     false,
			// 		ValidateFunc: validation.StringIsNotEmpty,
			// },
			"deploy_access_levels": {
				Description: "Array of access levels allowed to deploy, with each described by a hash.",
				Type:        schema.TypeList,
				MaxItems:    1,
				ForceNew:    true,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_level": {
							Description:  fmt.Sprintf("Levels of access required to deploy to this protected environment. Valid values are %s.", renderValueListForDocs(validProtectedEnvironmentDeploymentLevelNames)),
							Type:         schema.TypeString,
							ForceNew:     true,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(validProtectedEnvironmentDeploymentLevelNames, false),
						},
						"access_level_description": {
							Description: "Readable description of level of access.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"user_id": {
							Description:  "The ID of the user allowed to deploy to this protected environment.",
							Type:         schema.TypeInt,
							ForceNew:     true,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
						"group_id": {
							Description:  "The ID of the group allowed to deploy to this protected environment.",
							Type:         schema.TypeInt,
							ForceNew:     true,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
					},
				},
			},
		},
	}
})

func resourceGitlabProjectProtectedEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	deployAccessLevels, err := expandDeployAccessLevels(d.Get("deploy_access_levels").([]interface{}))
	if err != nil {
		return diag.Errorf("error expanding deploy_access_levels: %v", err)
	}

	options := &gitlab.ProtectRepositoryEnvironmentsOptions{
		Name: gitlab.String(d.Get("environment").(string)),
		DeployAccessLevels: &deployAccessLevels,
	}

	project := d.Get("project").(string)

	log.Printf("[DEBUG] Project %s create gitlab protected environment %q", project, *options.Name)

	client := meta.(*gitlab.Client)

	protectedEnvironment, _, err := client.ProtectedEnvironments.ProtectRepositoryEnvironments(project, options, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			return diag.Errorf("feature Protected Environments is not available")
		}
		return diag.FromErr(err)
	}

	d.SetId(buildTwoPartID(&project, &protectedEnvironment.Name))

	return resourceGitlabProjectProtectedEnvironmentRead(ctx, d, meta)
}

func resourceGitlabProjectProtectedEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] read gitlab protected environment %s", d.Id())

	project, environment, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("project", project)
	d.Set("environment", environment)

	log.Printf("[DEBUG] Project %s read gitlab protected environment %q", project, environment)

	client := meta.(*gitlab.Client)

	protectedEnvironment, _, err := client.ProtectedEnvironments.GetProtectedEnvironment(project, environment, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] Project %s gitlab protected environment %q not found, removing from state", project, environment)
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting gitlab project %q protected environment %q: %v", project, environment, err)
	}

	if err := d.Set("deploy_access_levels", flattenDeployAccessLevels(protectedEnvironment.DeployAccessLevels)); err != nil {
		return diag.Errorf("error setting deploy_access_levels: %v", err)
	}

	return nil
}

func resourceGitlabProjectProtectedEnvironmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	project, environmentName, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Project %s delete gitlab project-level protected environment %s", project, environmentName)

	client := meta.(*gitlab.Client)

	_, err = client.ProtectedEnvironments.UnprotectEnvironment(project, environmentName, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func expandDeployAccessLevels(vs []interface{}) ([]*gitlab.EnvironmentAccessOptions, error) {
	result := make([]*gitlab.EnvironmentAccessOptions, 0)

	for _, v := range vs {
		opts := v.(map[string]interface{})
		option := &gitlab.EnvironmentAccessOptions{}
		if accessLevel, exists := opts["access_level"]; exists {
			option.AccessLevel = gitlab.AccessLevel(accessLevelNameToValue[accessLevel.(string)])
		} else if userID, exists := opts["user_id"]; exists {
			option.UserID = gitlab.Int(userID.(int))
		} else if groupID, exists := opts["group_id"]; exists {
			option.GroupID = gitlab.Int(groupID.(int))
		}
		result = append(result, option)
	}

	return result, nil
}

func flattenDeployAccessLevels(vs []*gitlab.EnvironmentAccessDescription) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	for _, accessDescription := range vs {
		v := make(map[string]interface{})
		v["access_level"] = accessLevelValueToName[accessDescription.AccessLevel]
		v["access_level_description"] = accessDescription.AccessLevelDescription
		if accessDescription.UserID != 0 {
			v["user_id"] = accessDescription.UserID
		}
		if accessDescription.GroupID != 0 {
			v["group_id"] = accessDescription.GroupID
		}
		result = append(result, v)
	}

	return result
}
