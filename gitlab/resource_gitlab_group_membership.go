package gitlab

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabGroupMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabGroupMembershipCreate,
		Read:   resourceGitlabGroupMembershipRead,
		Update: resourceGitlabGroupMembershipUpdate,
		Delete: resourceGitlabGroupMembershipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"group_id": {
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

func resourceGitlabGroupMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	group_id := d.Get("group_id").(string)
	access_level := strings.ToLower(d.Get("access_level").(string))
	access_level_id, ok := accessLevelID[access_level]
	if !ok {
		return fmt.Errorf("Invalid access level '%s'", access_level)
	}
	user_id := d.Get("user_id").(int)
	options := &gitlab.AddGroupMemberOptions{
		UserID:      &user_id,
		AccessLevel: &access_level_id,
	}
	log.Printf("[DEBUG] create gitlab group membership for %d in %s", options.UserID, group_id)

	membership, _, err := client.GroupMembers.AddGroupMember(group_id, options)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%d", membership.ID))

	return resourceGitlabGroupMembershipRead(d, meta)
}

func resourceGitlabGroupMembershipRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	log.Printf("[DEBUG] read gitlab group membership %s", d.Id())

	group_id := d.Get("group_id").(string)
	user_id := d.Get("user_id").(int)

	membership, resp, err := client.GroupMembers.GetGroupMember(group_id, user_id)
	if err != nil {
		if resp.StatusCode == 404 {
			log.Printf("[WARN] removing group membership %s for %s from state because it no longer exists in gitlab", d.Id(), group_id)
			d.SetId("")
			return nil
		}
		return err
	}

	resourceGitlabGroupMembershipSetToState(d, membership)
	return nil
}

func resourceGitlabGroupMembershipUpdate(d *schema.ResourceData, meta interface{}) error {
	if !d.HasChange("access_level") {
		return nil
	}

	client := meta.(*gitlab.Client)

	group_id := d.Get("group_id").(string)
	user_id := d.Get("user_id").(int)
	access_level := strings.ToLower(d.Get("access_level").(string))
	access_level_id, ok := accessLevelID[access_level]
	if !ok {
		return fmt.Errorf("Invalid access level '%s'", access_level)
	}
	options := gitlab.EditGroupMemberOptions{
		AccessLevel: &access_level_id,
	}
	log.Printf("[DEBUG] update gitlab group membership %s for %s", d.Id(), group_id)

	_, _, err := client.GroupMembers.EditGroupMember(group_id, user_id, &options)
	if err != nil {
		return err
	}

	return resourceGitlabGroupMembershipRead(d, meta)
}

func resourceGitlabGroupMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	group_id := d.Get("group_id").(string)
	user_id := d.Get("user_id").(int)
	log.Printf("[DEBUG] Delete gitlab group membership %s for %s", d.Id(), group_id)

	_, err := client.GroupMembers.RemoveGroupMember(group_id, user_id)
	return err
}

func resourceGitlabGroupMembershipSetToState(d *schema.ResourceData, membership *gitlab.GroupMember) {
	d.SetId(fmt.Sprintf("%d", membership.ID))
	d.Set("username", membership.Username)
	d.Set("email", membership.Email)
	d.Set("Name", membership.Name)
	d.Set("State", membership.State)
	d.Set("AccessLevel", membership.AccessLevel)
}
