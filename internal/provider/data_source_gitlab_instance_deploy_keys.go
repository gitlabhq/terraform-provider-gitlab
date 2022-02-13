package provider

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_instance_deploy_keys", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_instance_deploy_keys`" + ` data source allows to retrieve a list of deploy keys for a GitLab instance.

-> This data source requires administration privileges.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/deploy_keys.html#list-all-deploy-keys)`,

		ReadContext: dataSourceGitlabInstanceDeployKeysRead,
		Schema: map[string]*schema.Schema{
			"public": {
				Description: "Only return deploy keys that are public.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"deploy_keys": {
				Description: "The list of all deploy keys across all projects of the GitLab instance.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The ID of the deploy key.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"title": {
							Description: "The title of the deploy key.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"created_at": {
							Description: "The creation date of the deploy key. In RFC3339 format.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"key": {
							Description: "The deploy key.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"fingerprint": {
							Description: "The fingerprint of the deploy key.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"projects_with_write_access": {
							Description: "The list of projects that the deploy key has write access to.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Description: "The ID of the project.",
										Type:        schema.TypeInt,
										Computed:    true,
									},
									"description": {
										Description: "The description of the project.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"name": {
										Description: "The name of the project.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"name_with_namespace": {
										Description: "The name of the project with namespace.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"path": {
										Description: "The path of the project.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"path_with_namespace": {
										Description: "The path of the project with namespace.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"created_at": {
										Description: "The creation date of the project. In RFC3339 format.",
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
})

func dataSourceGitlabInstanceDeployKeysRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	// Get group memberships
	options := &gitlab.ListInstanceDeployKeysOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 20,
			Page:    1,
		},
		Public: gitlab.Bool(d.Get("public").(bool)),
	}

	log.Printf("[INFO] Reading Instance Deploy Keys, with: %v", options)

	var instanceDeployKeys []*gitlab.InstanceDeployKey
	for options.Page != 0 {
		paginatedInstancedeployKeys, resp, err := client.DeployKeys.ListAllDeployKeys(options, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
		instanceDeployKeys = append(instanceDeployKeys, paginatedInstancedeployKeys...)

		options.Page = resp.NextPage
	}

	// NOTE: this data source doesn't have a "real" id, but the query to the API
	//       should actually return the same response for the same options,
	//       therefore the `Public` field option is used as id.
	d.SetId(fmt.Sprintf("%b", options.Public))
	if err := d.Set("deploy_keys", flattenGitlabInstanceDeployKeys(instanceDeployKeys)); err != nil {
		return diag.Errorf("error setting deploy_keys: %s", err)
	}
	return nil
}

func flattenGitlabInstanceDeployKeys(keys []*gitlab.InstanceDeployKey) []interface{} {
	result := []interface{}{}
	for _, instanceDeployKey := range keys {
		values := map[string]interface{}{
			"id":                         instanceDeployKey.ID,
			"title":                      instanceDeployKey.Title,
			"created_at":                 instanceDeployKey.CreatedAt.Format(time.RFC3339),
			"key":                        instanceDeployKey.Key,
			"fingerprint":                instanceDeployKey.Fingerprint,
			"projects_with_write_access": flattenGitlabInstanceDeployKeysProjectsWithWriteAccess(instanceDeployKey.ProjectsWithWriteAccess),
		}
		result = append(result, values)
	}
	return result
}

func flattenGitlabInstanceDeployKeysProjectsWithWriteAccess(projects []*gitlab.DeployKeyProject) []interface{} {
	result := []interface{}{}
	for _, project := range projects {
		values := map[string]interface{}{
			"id":                  project.ID,
			"description":         project.Description,
			"name":                project.Name,
			"name_with_namespace": project.NameWithNamespace,
			"path":                project.Path,
			"path_with_namespace": project.PathWithNamespace,
			"created_at":          project.CreatedAt.Format(time.RFC3339),
		}
		result = append(result, values)
	}
	return result
}
