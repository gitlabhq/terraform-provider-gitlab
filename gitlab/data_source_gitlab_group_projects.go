package gitlab

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/xanzy/go-gitlab"
)

func dataSourceGitlabGroupProjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGitlabGroupProjectsRead,

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
			"group_id": {
				Type:     schema.TypeInt,
				Required: true,
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

func dataSourceGitlabGroupProjectsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	log.Printf("[INFO] Reading Gitlab projects")

	listGroupProjectsOptions, id, err := expandGitlabGroupProjectsOptions(d)
	if err != nil {
		return err
	}

	groupId, _ := d.GetOk("group_id")

	projectList := []*gitlab.Project{}
	for {
		projects, response, err := client.Groups.ListGroupProjects(groupId, listGroupProjectsOptions, nil)
		projectList = append(projectList, projects...)

		if err != nil {
			return err
		}
		listGroupProjectsOptions.ListOptions.Page++

		log.Printf("[INFO] Currentpage: %d, Total: %d", response.CurrentPage, response.TotalPages)
		if response.CurrentPage == response.TotalPages {
			break
		}
	}

	d.SetId(fmt.Sprintf("%d", id))
	d.Set("projects", flattenGroupProjects(projectList))
	return nil
}

func flattenGroupProjects(projects []*gitlab.Project) []interface{} {
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

func expandGitlabGroupProjectsOptions(d *schema.ResourceData) (*gitlab.ListGroupProjectsOptions, int, error) {
	listGroupProjectsOptions := &gitlab.ListGroupProjectsOptions{}
	var optionsHash strings.Builder

	if data, ok := d.GetOk("order_by"); ok {
		orderBy := data.(string)
		listGroupProjectsOptions.OrderBy = &orderBy
		optionsHash.WriteString(orderBy)
	}

	if data, ok := d.GetOk("search"); ok {
		search := data.(string)
		listGroupProjectsOptions.Search = &search
		optionsHash.WriteString(search)
	}

	listGroupProjectsOptions.ListOptions.PerPage = 100
	listGroupProjectsOptions.ListOptions.Page = 1

	id := schema.HashString(optionsHash.String())

	return listGroupProjectsOptions, id, nil
}
