package gitlab

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func dataSourceGitlabProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGitlabProjectRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

// Performs the lookup
func dataSourceGitlabProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	log.Printf("[INFO] Reading Gitlab project")

	searchName := strings.ToLower(d.Get("name").(string))

	o := &gitlab.ListProjectsOptions{
		Search: &searchName,
	}

	projects, _, err := client.Projects.ListProjects(o)
	if err != nil {
		return err
	}

	var found *gitlab.Project

	for _, project := range projects {
		if strings.ToLower(project.Name) == searchName {
			found = project
			break
		}
	}

	if found == nil {
		return fmt.Errorf("Unable to locate any project with the name: %s", searchName)
	}
	d.SetId(fmt.Sprintf("%d", found.ID))
	d.Set("name", found.Name)
	d.Set("path", found.Path)
	d.Set("description", found.Description)
	d.Set("default_branch", found.DefaultBranch)
	d.Set("issues_enabled", found.IssuesEnabled)
	d.Set("merge_requests_enabled", found.MergeRequestsEnabled)
	d.Set("wiki_enabled", found.WikiEnabled)
	d.Set("snippets_enabled", found.SnippetsEnabled)
	d.Set("visibility_level", string(found.Visibility))
	d.Set("namespace_id", found.Namespace.ID)
	d.Set("ssh_url_to_repo", found.SSHURLToRepo)
	d.Set("http_url_to_repo", found.HTTPURLToRepo)
	d.Set("web_url", found.WebURL)
	return nil
}
