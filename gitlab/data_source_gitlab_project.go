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
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"default_branch": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"issues_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"merge_requests_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"wiki_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"snippets_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"visibility_level": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"namespace_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"ssh_url_to_repo": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"http_url_to_repo": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"web_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"runners_token": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"archived": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceGitlabProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	log.Printf("[INFO] Reading Gitlab project")

	v, _ := d.GetOk("id")

	found, _, err := client.Projects.GetProject(v, nil)
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
	d.Set("runners_token", found.RunnersToken)
	d.Set("archived", found.Archived)
	return nil
}
