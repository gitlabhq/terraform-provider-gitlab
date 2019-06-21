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
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"branch": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"ref": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

	_, _, err := client.Branches.CreateBranch(project, options)

	if err != nil {
		return err
	}

	return resourceGitlabBranchRead(d, meta)
}

func resourceGitlabBranchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project, branch, err := projectAndBranchFromID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] read gitlab branch protection for project %s, branch %s", project, branch)

	pb, _, err := client.ProtectedBranches.GetProtectedBranch(project, branch)
	if err != nil {
		return err
	}

	d.Set("project", project)
	d.Set("branch", pb.Name)
	d.Set("merge_access_level", pb.MergeAccessLevels[0].AccessLevel)
	d.Set("push_access_level", pb.PushAccessLevels[0].AccessLevel)

	d.SetId(buildTwoPartID(&project, &pb.Name))

	return nil
}
