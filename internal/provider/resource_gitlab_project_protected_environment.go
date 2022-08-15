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

~> In order to use a user or group in the ` + "`deploy_access_levels`" + ` configuration,
   you need to make sure that users have access to the project and groups must have this project shared.
   You may use the ` + "`gitlab_project_membership`" + ` and ` + "`gitlab_project_shared_group`" + ` resources to achieve this.
   Unfortunately, the GitLab API does not complain about users and groups without access to the project and just ignores those.
   In case this happens you will get perpetual state diffs.

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
			"required_approval_count": {
				Description: "The number of approvals required to deploy to this environment.",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"deploy_access_levels": {
				Description: "Array of access levels allowed to deploy, with each described by a hash.",
				Type:        schema.TypeList,
				ForceNew:    true,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_level": {
							Description:  fmt.Sprintf("Levels of access required to deploy to this protected environment. Valid values are %s.", renderValueListForDocs(validProtectedEnvironmentDeploymentLevelNames)),
							Type:         schema.TypeString,
							ForceNew:     true,
							Optional:     true,
							Computed:     true, // When user_id or group_id is specified, the GitLab API still returns an access_level in the response.
							ValidateFunc: validation.StringInSlice(validProtectedEnvironmentDeploymentLevelNames, false),
						},
						"access_level_description": {
							Description: "Readable description of level of access.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"user_id": {
							Description:  "The ID of the user allowed to deploy to this protected environment. The user must be a member of the project.",
							Type:         schema.TypeInt,
							ForceNew:     true,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
						"group_id": {
							Description:  "The ID of the group allowed to deploy to this protected environment. The project must be shared with the group.",
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
		return diag.FromErr(err)
	}

	options := &gitlab.ProtectRepositoryEnvironmentsOptions{
		Name:               gitlab.String(d.Get("environment").(string)),
		DeployAccessLevels: &deployAccessLevels,
	}

	if v, ok := d.GetOk("required_approval_count"); ok {
		options.RequiredApprovalCount = gitlab.Int(v.(int))
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
	d.Set("required_approval_count", protectedEnvironment.RequiredApprovalCount)

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
	result := make([]*gitlab.EnvironmentAccessOptions, len(vs))

	for i, v := range vs {
		opts := v.(map[string]interface{})
		option := &gitlab.EnvironmentAccessOptions{}
		count := 0

		if accessLevel, ok := opts["access_level"]; ok && accessLevel != "" {
			option.AccessLevel = gitlab.AccessLevel(accessLevelNameToValue[accessLevel.(string)])
			count++
		}

		if userID, ok := opts["user_id"]; ok && userID != 0 {
			option.UserID = gitlab.Int(userID.(int))
			count++
		}

		if groupID, ok := opts["group_id"]; ok && groupID != 0 {
			option.GroupID = gitlab.Int(groupID.(int))
			count++
		}

		// This is a manual "ExactlyOneOf" schema check, since this cannot be validated at the
		// schema-level inside of a list.
		// See: https://github.com/hashicorp/terraform-plugin-sdk/blob/0f834ffb1619ce1ef8d3f5255911108ede086ef9/helper/schema/schema.go#L278
		if count != 1 {
			return nil, fmt.Errorf(`illegal deploy_access_levels.%d: exactly one of "access_level", "user_id", or "group_id" must be specified (got %d)`, i, count)
		}

		result[i] = option
	}

	return result, nil
}

func flattenDeployAccessLevels(accessDescriptions []*gitlab.EnvironmentAccessDescription) []map[string]interface{} {
	result := make([]map[string]interface{}, len(accessDescriptions))

	for i, accessDescription := range accessDescriptions {
		v := make(map[string]interface{})
		v["access_level_description"] = accessDescription.AccessLevelDescription
		if accessDescription.AccessLevel != 0 {
			v["access_level"] = accessLevelValueToName[accessDescription.AccessLevel]
		}
		if accessDescription.UserID != 0 {
			v["user_id"] = accessDescription.UserID
		}
		if accessDescription.GroupID != 0 {
			v["group_id"] = accessDescription.GroupID
		}
		result[i] = v
	}

	return result
}
