package gitlab 

import (
	// "errors"
	// "fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabBranch() *schema.Resource {
	// removed guest TODO check acceptable access levels 
	// ref force new false --- TODO resolve if incorrect
	// acceptedAccessLevels := []string{ "reporter", "developer", "maintainer"}

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
				Type:         schema.TypeString,
				ForceNew:     false,
				Required:     true,
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
	// requestOptions := &gitlab.RequestOptionFunc()
	branch, resp, err := client.Branches.CreateBranch(project, branchOptions)
	if err != nil {
		log.Printf("[DEBUG] failed to create gitlab branch %v response %v", branch, resp)
		return err
	}
	return resourceGitlabBranchRead(d, meta)
}

// TODO investigate setting
// type Branch struct {
// 	Commit             *Commit `json:"commit"`
// 	Name               string  `json:"name"`
// 	Protected          bool    `json:"protected"`
// 	Merged             bool    `json:"merged"`
// 	Default            bool    `json:"default"`
// 	CanPush            bool    `json:"can_push"`
// 	DevelopersCanPush  bool    `json:"developers_can_push"`
// 	DevelopersCanMerge bool    `json:"developers_can_merge"`
// 	WebURL             string  `json:"web_url"`
// }

func resourceGitlabBranchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project, name, err := projectAndBranchFromID(d.Id())
	log.Printf("[DEBUG] read gitlab branch %s", d.Id())
	branch, resp, err := client.Branches.GetBranch(project, name)
	if err != nil {
		log.Printf("[DEBUG] failed to read gitlab branch %s response %v", branch, resp)
	}
	d.Set("name", branch.Name) 	
	d.Set("ref", d.Get("ref").(string))
	d.Set("project", project)
	return err
}

func resourceGitlabBranchDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project, name, err := projectAndBranchFromID(d.Id())
	log.Printf("[DEBUG] delete gitlab branch %s", name)
	_, err := client.Branches.DeleteBranch(project, name)

	return err
}