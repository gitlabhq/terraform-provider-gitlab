package gitlab

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/xanzy/go-gitlab"
)

func dataSourceGitlabProjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGitlabProjectsRead,

		Schema: map[string]*schema.Schema{
			"order_by": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "id",
				ValidateFunc: validation.StringInSlice([]string{"id", "name",
					"username", "created_at", "updated_at"}, true),
			},
			"search": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"projects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGitlabProjectsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	log.Printf("[DEBUG] Reading Gitlab projects")

	listProjectsOptions, id, err := expandGitlabProjectsOptions(d)
	if err != nil {
		return err
	}

	projectList := []*gitlab.Project{}

	for {
		projects, response, err := client.Projects.ListProjects(listProjectsOptions, nil)
		projectList = append(projectList, projects...)

		if err != nil {
			return err
		}
		listProjectsOptions.ListOptions.Page++

		log.Printf("[INFO] Currentpage: %d, Total: %d", response.CurrentPage, response.TotalPages)
		if response.CurrentPage == response.TotalPages {
			break
		}
	}

	d.SetId(fmt.Sprintf("%d", id))
	d.Set("projects", flattenProjects(projectList))
	return nil
}

func flattenProjects(projects []*gitlab.Project) []interface{} {
	projectsList := []interface{}{}

	for _, project := range projects {
		values := map[string]interface{}{
			"id":   project.ID,
			"name": project.Name,
		}

		projectsList = append(projectsList, values)
	}

	return projectsList
}

func expandGitlabProjectsOptions(d *schema.ResourceData) (*gitlab.ListProjectsOptions, int, error) {
	listProjectsOptions := &gitlab.ListProjectsOptions{}
	var optionsHash strings.Builder

	if data, ok := d.GetOk("order_by"); ok {
		orderBy := data.(string)
		listProjectsOptions.OrderBy = &orderBy
		optionsHash.WriteString(orderBy)
	}

	if data, ok := d.GetOk("search"); ok {
		search := data.(string)
		listProjectsOptions.Search = &search
		optionsHash.WriteString(search)
	}

	listProjectsOptions.ListOptions.PerPage = 100
	listProjectsOptions.ListOptions.Page = 1

	id := schema.HashString(optionsHash.String())

	return listProjectsOptions, id, nil
}
