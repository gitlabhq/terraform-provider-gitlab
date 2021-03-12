package gitlab

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
	"log"
)

func resourceGitlabBranch() *schema.Resource {
	// removed guest TODO check acceptable access levels
	// ref force new false --- TODO resolve if incorrect
	// TODO project -> project_name
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
				Type:     schema.TypeString,
				ForceNew: true,
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
	name := d.Get("name").(string)
	project := d.Get("project").(string)
	ref := d.Get("ref").(string)
	log.Printf("[DEBUG] read gitlab branch %s", name)
	branch, resp, err := client.Branches.GetBranch(project, name)
	if err != nil {
		log.Printf("[DEBUG] failed to read gitlab branch %s response %v", name, resp)
		return err
	}
	d.SetId(fmt.Sprintf("%s-%s", project, name))
	d.Set("name", branch.Name)
	d.Set("project", project)
	d.Set("ref", ref)
	d.Set("web_url", branch.WebURL)
	d.Set("default", branch.Default)
	d.Set("can_push", branch.CanPush)
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
