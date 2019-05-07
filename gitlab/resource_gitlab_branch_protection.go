package gitlab

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabBranchProtection() *schema.Resource {
	acceptedAccessLevels := make([]string, 0, len(accessLevelID))

	for k := range accessLevelID {
		acceptedAccessLevels = append(acceptedAccessLevels, k)
	}
	return &schema.Resource{
		Create: resourceGitlabBranchProtectionCreate,
		Read:   resourceGitlabBranchProtectionRead,
		Delete: resourceGitlabBranchProtectionDelete,
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
			"merge_access_level": {
				Type:         schema.TypeString,
				ValidateFunc: validateValueFunc(acceptedAccessLevels),
				Required:     true,
				ForceNew:     true,
			},
			"push_access_level": {
				Type:         schema.TypeString,
				ValidateFunc: validateValueFunc(acceptedAccessLevels),
				Required:     true,
				ForceNew:     true,
			},
		},
	}
}

func resourceGitlabBranchProtectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	branch := gitlab.String(d.Get("branch").(string))
	mergeAccessLevel := accessLevelID[d.Get("merge_access_level").(string)]
	pushAccessLevel := accessLevelID[d.Get("push_access_level").(string)]

	options := &gitlab.ProtectRepositoryBranchesOptions{
		Name:             branch,
		MergeAccessLevel: &mergeAccessLevel,
		PushAccessLevel:  &pushAccessLevel,
	}

	log.Printf("[DEBUG] create gitlab branch protection on %v for project %s", options.Name, project)

	bp, _, err := client.ProtectedBranches.ProtectRepositoryBranches(project, options)
	if err != nil {
		// Remove existing branch protection
		_, err = client.ProtectedBranches.UnprotectRepositoryBranches(project, *branch)
		if err != nil {
			return err
		}
		// Reprotect branch with updated values
		bp, _, err = client.ProtectedBranches.ProtectRepositoryBranches(project, options)
		if err != nil {
			return err
		}
	}

	d.SetId(buildTwoPartID(&project, &bp.Name))

	return resourceGitlabBranchProtectionRead(d, meta)
}

func resourceGitlabBranchProtectionRead(d *schema.ResourceData, meta interface{}) error {
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

func resourceGitlabBranchProtectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	branch := d.Get("branch").(string)

	log.Printf("[DEBUG] Delete gitlab protected branch %s for project %s", branch, project)

	_, err := client.ProtectedBranches.UnprotectRepositoryBranches(project, branch)
	return err
}

func projectAndBranchFromID(id string) (string, string, error) {
	project, branch, err := parseTwoPartID(id)

	if err != nil {
		log.Printf("[WARN] cannot get group member id from input: %v", id)
	}
	return project, branch, err
}
