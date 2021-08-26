package gitlab

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
	"log"
)

func resourceGitlabBranch() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabBranchCreate,
		Read:   resourceGitlabBranchRead,
		Delete: resourceGitlabBranchDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"project": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"ref": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  "master",
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
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"author_email": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"author_name": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"authored_date": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"committed_date": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"committer_email": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"committer_name": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"short_id": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"title": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceGitlabBranchCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	name := d.Get("name").(string)
	project := d.Get("project").(string)
	ref := d.Get("ref").(string)
	branchOptions := &gitlab.CreateBranchOptions{
		Branch: &name, Ref: &ref,
	}

	log.Printf("[DEBUG] create gitlab branch %s for project %s with ref %s", name, project, ref)
	branch, resp, err := client.Branches.CreateBranch(project, branchOptions)
	if err != nil {
		log.Printf("[DEBUG] failed to create gitlab branch %v response %v", branch, resp)
		return err
	}
	return resourceGitlabBranchRead(d, meta)
}

func resourceGitlabBranchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	name := d.Get("name").(string)
	project := d.Get("project").(string)
	ref := d.Get("ref").(string)
	log.Printf("[DEBUG] read gitlab branch %s", name)
	branch, resp, err := client.Branches.GetBranch(project, name)
	if err != nil {
		log.Printf("[DEBUG] failed to read gitlab branch %s response %v", name, resp)
		return err
	}
	d.SetId(buildTwoPartID(&project, &name))
	d.Set("name", branch.Name)
	d.Set("project", project)
	d.Set("ref", ref)
	d.Set("web_url", branch.WebURL)
	d.Set("default", branch.Default)
	d.Set("can_push", branch.CanPush)
	d.Set("merged", branch.Merged)
	commit := flattenCommit(branch.Commit)
	if err := d.Set("commit", commit); err != nil {
		return err
	}

	return nil
}

func resourceGitlabBranchDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	name := d.Get("name").(string)
	log.Printf("[DEBUG] delete gitlab branch %s", name)
	resp, err := client.Branches.DeleteBranch(project, name)
	if err != nil {
		log.Printf("[DEBUG] failed to delete gitlab branch %s response %v", name, resp)
	}
	return err
}

func flattenCommit(commit *gitlab.Commit) (values []map[string]interface{}) {
	if commit == nil {
		return []map[string]interface{}{}
	}

	return []map[string]interface{}{
		{
			"id":          commit.ID,
			"short_id":    commit.ShortID,
			"title":       commit.Title,
			"author_name": commit.AuthorName,
			"message":     commit.Message,
		},
	}
}
