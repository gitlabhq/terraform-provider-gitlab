package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/hashstructure"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_repository_tree", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_repository_tree`" + ` data source allows details of directories and files in a repository to be retrieved.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/repositories.html#list-repository-tree)`,

		ReadContext: dataSourceGitlabRepositoryTreeRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The ID or full path of the project owned by the authenticated user.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"ref": {
				Description: "The name of a repository branch or tag.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"path": {
				Description: "The path inside repository. Used to get content of subdirectories.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"recursive": {
				Description: "Boolean value used to get a recursive tree (false by default).",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"tree": {
				Description: "The list of files/directories returned by the search",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The SHA-1 hash of the tree or blob in the repository.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "Name of the blob or tree in the repository",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"type": {
							Description: "Type of object in the repository. Can be either type tree or of type blob",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"path": {
							Description: "Path of the object inside of the repository.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"mode": {
							Description: "Unix access mode of the file in the repository.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
})

func dataSourceGitlabRepositoryTreeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	options := &gitlab.ListTreeOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 20,
			Page:    1,
		},
		Path:      gitlab.String(d.Get("path").(string)),
		Ref:       gitlab.String(d.Get("ref").(string)),
		Recursive: gitlab.Bool(d.Get("recursive").(bool)),
	}

	var nodes []*gitlab.TreeNode
	for options.Page != 0 {

		paginatedNodes, resp, err := client.Repositories.ListTree(project, options, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}

		nodes = append(nodes, paginatedNodes...)

		options.Page = resp.NextPage
	}

	optionsHash, err := hashstructure.Hash(&options, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s:%d", project, optionsHash))
	if err := d.Set("tree", flattenGitlabRepositoryTree(project, nodes)); err != nil {
		return diag.Errorf("failed to set repository tree nodes to state: %v", err)
	}

	return nil
}

func flattenGitlabRepositoryTree(project string, treeNodes []*gitlab.TreeNode) []interface{} {
	treeNodeList := []interface{}{}

	for _, node := range treeNodes {

		values := map[string]interface{}{
			"id":   project,
			"name": node.Name,
			"type": node.Type,
			"path": node.Path,
			"mode": node.Mode,
		}

		treeNodeList = append(treeNodeList, values)
	}
	return treeNodeList
}
