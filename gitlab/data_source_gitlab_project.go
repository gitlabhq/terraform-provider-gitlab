package gitlab

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func dataSourceGitlabProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGitlabProjectRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func dataSourceGitlabProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	log.Printf("[INFO] Reading Gitlab project")

	v, _ := d.GetOk("id")

	found, _, err := client.Projects.GetProject(v)
	if err != nil {
		return err
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
