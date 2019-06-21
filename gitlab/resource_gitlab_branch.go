package gitlab

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabBranchCreation() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabBranchCreate,
		Read:   resourceGitlabBranchRead,
		Delete: resourceGitlabBranchDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ref": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceGitlabBranchCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	options := &gitlab.CreateBranchOptions{
		Branch: gitlab.String(d.Get("name").(string)),
		Ref:    gitlab.String(d.Get("ref").(string)),
	}

	project := d.Get("project").(string)

	branch, _, err := client.Branches.CreateBranch(project, options)

	log.Printf("[DEBUG] created gitlab branch for project %s, branch %s", project, branch)

	if err != nil {
		return err
	}

	return resourceGitlabBranchRead(d, meta)
}

func resourceGitlabBranchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	project := d.Get("project").(string)
	branch := d.Get("name").(string)

	pb, _, err := client.Branches.GetBranch(project, branch)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] read gitlab branch for project %s, branch %s", project, branch)

	d.SetId(buildTwoPartID(&project, &pb.Name))

	return nil
}

func resourceGitlabBranchDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	branch := d.Get("name").(string)

	log.Printf("[DEBUG] Delete gitlab branch %s for project %s", branch, project)

	_, err := client.Branches.DeleteBranch(project, branch)
	return err
}
