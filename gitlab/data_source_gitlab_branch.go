package gitlab

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func dataSourceGitlabBranch() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGitlabBranchRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"web_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"can_push": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"merged": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"commit": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				// Schema: map[string]*schema.Schema{
				// "id": {
				// 	Type:     schema.TypeString,
				// 	Computed: true,
				// },
				// "author_email": {
				// 	Type:     schema.TypeString,
				// 	Computed: true,
				// },
				// "author_name": {
				// 	Type:     schema.TypeString,
				// 	Computed: true,
				// },
				// "authored_date": {
				// 	Type:     schema.TypeString,
				// 	Computed: true,
				// },
				// "committed_date": {
				// 	Type:     schema.TypeString,
				// 	Computed: true,
				// },
				// "committer_email": {
				// 	Type:     schema.TypeString,
				// 	Computed: true,
				// },
				// "committer_name": {
				// 	Type:     schema.TypeString,
				// 	Computed: true,
				// },
				// "short_id": {
				// 	Type:     schema.TypeString,
				// 	Computed: true,
				// },
				// "title": {
				// 	Type:     schema.TypeString,
				// 	Computed: true,
				// },
				// "message": {
				// 	Type:     schema.TypeString,
				// 	Computed: true,
				// },
				// "parent_ids": {
				// 	Type:     schema.TypeSet,
				// 	Computed: true,
				// 	Elem:     &schema.Schema{Type: schema.TypeString},
				// 	Set:      schema.HashString,
				// },
				// },
				// },
			},
		},
	}
}

func dataSourceGitlabBranchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	name := d.Get("name").(string)
	project := d.Get("project").(string)
	log.Printf("[DEBUG] read gitlab branch %s", name)
	branch, resp, err := client.Branches.GetBranch(project, name)
	if err != nil {
		log.Printf("[DEBUG] failed to read gitlab branch %s response %v", name, resp)
		return err
	}

	d.SetId(buildTwoPartID(&project, &name))
	d.Set("name", branch.Name)
	d.Set("project", project)
	d.Set("web_url", branch.WebURL)
	d.Set("default", branch.Default)
	d.Set("can_push", branch.CanPush)
	d.Set("protected", branch.Protected)
	d.Set("merged", branch.Merged)
	d.Set("developer_can_merge", branch.DevelopersCanMerge)
	d.Set("developer_can_push", branch.DevelopersCanPush)
	commit := flattenCommit(branch.Commit)
	if err := d.Set("commit", commit); err != nil {
		return err
	}
	return nil
}
