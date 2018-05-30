package gitlab

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabProjectMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabProjectMembershipCreate,
		Read:   resourceGitlabProjectMembershipRead,
		Update: resourceGitlabProjectMembershipUpdate,
		Delete: resourceGitlabProjectMembershipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"user_id": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},
			"access_level": {
				Type:     schema.TypeString,
				Required: true,
			},
			// "expires_at": {
			// 	Type:     schema.TypeString, // Format YYYY-MM-DD
			// 	ForceNew: true,
			// 	Required: false,
			// 	Optional: true,
			// },
		},
	}
}

var accessLevelID = map[string]gitlab.AccessLevelValue{
	"guest":     gitlab.GuestPermissions,
	"reporter":  gitlab.ReporterPermissions,
	"developer": gitlab.DeveloperPermissions,
	"master":    gitlab.MasterPermissions,
	"owner":     gitlab.OwnerPermission,
}

func resourceGitlabProjectMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	project_id := d.Get("project_id").(string)
	access_level := strings.ToLower(d.Get("access_level").(string))
	access_level_id, ok := accessLevelID[access_level]
	if !ok {
		return fmt.Errorf("Invalid access level '%s'", access_level)
	}
	user_id := d.Get("user_id").(int)
	options := &gitlab.AddProjectMemberOptions{
		UserID:      &user_id,
		AccessLevel: &access_level_id,
	}
	log.Printf("[DEBUG] create gitlab project membership for %d in %s", options.UserID, project_id)

	membership, _, err := client.ProjectMembers.AddProjectMember(project_id, options)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%d", membership.ID))

	return resourceGitlabProjectMembershipRead(d, meta)
}

func resourceGitlabProjectMembershipRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	log.Printf("[DEBUG] read gitlab project membership %s", d.Id())

	project_id := d.Get("project_id").(string)
	user_id := d.Get("user_id").(int)

	membership, resp, err := client.ProjectMembers.GetProjectMember(project_id, user_id)
	if err != nil {
		if resp.StatusCode == 404 {
			log.Printf("[WARN] removing project membership %s for %s from state because it no longer exists in gitlab", d.Id(), project_id)
			d.SetId("")
			return nil
		}
		return err
	}

	resourceGitlabProjectMembershipSetToState(d, membership)
	return nil
}

func resourceGitlabProjectMembershipUpdate(d *schema.ResourceData, meta interface{}) error {
	if !d.HasChange("access_level") {
		return nil
	}

	client := meta.(*gitlab.Client)

	project_id := d.Get("project_id").(string)
	user_id := d.Get("user_id").(int)
	access_level := strings.ToLower(d.Get("access_level").(string))
	access_level_id, ok := accessLevelID[access_level]
	if !ok {
		return fmt.Errorf("Invalid access level '%s'", access_level)
	}
	options := gitlab.EditProjectMemberOptions{
		AccessLevel: &access_level_id,
	}
	log.Printf("[DEBUG] update gitlab project membership %s for %s", d.Id(), project_id)

	_, _, err := client.ProjectMembers.EditProjectMember(project_id, user_id, &options)
	if err != nil {
		return err
	}

	return resourceGitlabProjectMembershipRead(d, meta)
}

func resourceGitlabProjectMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project_id := d.Get("project_id").(string)
	user_id := d.Get("user_id").(int)
	log.Printf("[DEBUG] Delete gitlab project membership %s for %s", d.Id(), project_id)

	_, err := client.ProjectMembers.DeleteProjectMember(project_id, user_id)
	return err
}

func resourceGitlabProjectMembershipSetToState(d *schema.ResourceData, membership *gitlab.ProjectMember) {
	d.SetId(fmt.Sprintf("%d", membership.ID))
	d.Set("username", membership.Username)
	d.Set("email", membership.Email)
	d.Set("Name", membership.Name)
	d.Set("State", membership.State)
	d.Set("AccessLevel", membership.AccessLevel)
}
