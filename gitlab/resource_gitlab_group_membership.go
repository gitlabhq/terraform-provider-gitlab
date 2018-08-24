package gitlab

import (
	"fmt"
	"log"
	"strconv"
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
			"expires_at": {
				Type:     schema.TypeString, // Format YYYY-MM-DD
				Optional: true,
			},
		},
	}
}

func resourceGitlabGroupMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	userId := d.Get("user_id").(int)
	groupId := d.Get("group_id").(string)
	expiresAt := d.Get("expires_at").(string)
	accessLevel := strings.ToLower(d.Get("access_level").(string))
	accessLevelId, ok := accessLevelID[accessLevel]

	if !ok {
		return fmt.Errorf("Invalid access level '%s'", accessLevel)
	}
	options := &gitlab.AddGroupMemberOptions{
		UserID:      &userId,
		AccessLevel: &accessLevelId,
		ExpiresAt:   &expiresAt,
	}
	log.Printf("[DEBUG] create gitlab group groupMember for %d in %s", options.UserID, groupId)

	groupMember, _, err := client.GroupMembers.AddGroupMember(groupId, options)
	if err != nil {
		return err
	}
	userIdString := strconv.Itoa(groupMember.ID)
	d.SetId(buildTwoPartID(&groupId, &userIdString))
	return resourceGitlabGroupMembershipRead(d, meta)
}

func resourceGitlabGroupMembershipRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	id := d.Id()
	log.Printf("[DEBUG] read gitlab group groupMember %s", id)

	groupId, userId, e := groupIdAndUserIdFromId(id)
	if e != nil {
		return e
	}

	groupMember, resp, err := client.GroupMembers.GetGroupMember(groupId, userId)
	if err != nil {
		if resp.StatusCode == 404 {
			log.Printf("[WARN] removing group groupMember %v for %s from state because it no longer exists in gitlab", userId, groupId)
			d.SetId("")
			return nil
		}
		return err
	}

	resourceGitlabGroupMembershipSetToState(d, groupMember, &groupId)
	return nil
}

func groupIdAndUserIdFromId(id string) (string, int, error) {
	groupId, userIdString, err := parseTwoPartID(id)
	userId, e := strconv.Atoi(userIdString)
	if err != nil {
		e = err
	}
	if e != nil {
		log.Printf("[WARN] cannot get group member id from input: %v", id)
	}
	return groupId, userId, e
}

func resourceGitlabGroupMembershipUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	userId := d.Get("user_id").(int)
	groupId := d.Get("group_id").(string)
	expiresAt := d.Get("expires_at").(string)
	accessLevel := strings.ToLower(d.Get("access_level").(string))
	accessLevelId, ok := accessLevelID[accessLevel]
	if !ok {
		return fmt.Errorf("Invalid access level '%s'", accessLevel)
	}
	options := gitlab.EditGroupMemberOptions{
		AccessLevel: &accessLevelId,
		ExpiresAt:   &expiresAt,
	}
	log.Printf("[DEBUG] update gitlab group membership %v for %s", userId, groupId)

	_, _, err := client.GroupMembers.EditGroupMember(groupId, userId, &options)
	if err != nil {
		return err
	}

	return resourceGitlabGroupMembershipRead(d, meta)
}

func resourceGitlabGroupMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	id := d.Id()
	groupId, userId, e := groupIdAndUserIdFromId(id)
	if e != nil {
		return e
	}

	log.Printf("[DEBUG] Delete gitlab group membership %v for %s", userId, groupId)

	_, err := client.GroupMembers.RemoveGroupMember(groupId, userId)
	return err
}

func resourceGitlabGroupMembershipSetToState(d *schema.ResourceData, groupMember *gitlab.GroupMember, group_id *string) {
	d.Set("username", groupMember.Username)
	d.Set("email", groupMember.Email)
	d.Set("Name", groupMember.Name)
	d.Set("State", groupMember.State)
	d.Set("AccessLevel", groupMember.AccessLevel)
	d.Set("ExpiresAt", groupMember.ExpiresAt)

	userId := strconv.Itoa(groupMember.ID)
	d.SetId(buildTwoPartID(group_id, &userId))
}
