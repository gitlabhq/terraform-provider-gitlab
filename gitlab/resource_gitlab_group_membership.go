package gitlab

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func resourceGitlabGroupMembership() *schema.Resource {
	acceptedAccessLevels := make([]string, 0, len(accessLevelID))
	for k := range accessLevelID {
		acceptedAccessLevels = append(acceptedAccessLevels, k)
	}
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
				Optional: true,
				Computed: true,
				ConflictsWith: []string{
					"username",
				},
			},
			"username": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
				ConflictsWith: []string{
					"user_id",
				},
			},
			"access_level": {
				Type:         schema.TypeString,
				ValidateFunc: validateValueFunc(acceptedAccessLevels),
				Required:     true,
			},
			"expires_at": {
				Type:         schema.TypeString, // Format YYYY-MM-DD
				ValidateFunc: validateDateFunc,
				Optional:     true,
			},
		},
	}
}

func resourceGitlabGroupMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	var userId int

	userIdData, userIdOk := d.GetOk("user_id")
	usernameData, usernameOk := d.GetOk("username")
	groupId := d.Get("group_id").(string)
	expiresAt := d.Get("expires_at").(string)
	accessLevelId := accessLevelID[d.Get("access_level").(string)]

	if usernameOk {
		username := strings.ToLower(usernameData.(string))

		listUsersOptions := &gitlab.ListUsersOptions{
			Username: gitlab.String(username),
		}

		var users []*gitlab.User
		users, _, err := client.Users.ListUsers(listUsersOptions)
		if err != nil {
			return err
		}

		if len(users) == 0 {
			return fmt.Errorf("couldn't find a user matching: %s", username)
		} else if len(users) != 1 {
			return fmt.Errorf("more than one user found matching: %s", username)
		}

		userId = users[0].ID
	} else if userIdOk {
		userId = userIdData.(int)
	} else {
		return fmt.Errorf("one and only one of user_id or username must be set")
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

	groupMember, _, err := client.GroupMembers.GetGroupMember(groupId, userId)
	if err != nil {
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

	var userId int

	userIdData, userIdOk := d.GetOk("user_id")
	usernameData, usernameOk := d.GetOk("username")
	groupId := d.Get("group_id").(string)
	expiresAt := d.Get("expires_at").(string)
	accessLevelId := accessLevelID[strings.ToLower(d.Get("access_level").(string))]

	if usernameOk {
		username := strings.ToLower(usernameData.(string))

		listUsersOptions := &gitlab.ListUsersOptions{
			Username: gitlab.String(username),
		}

		var users []*gitlab.User
		users, _, err := client.Users.ListUsers(listUsersOptions)
		if err != nil {
			return err
		}

		if len(users) == 0 {
			return fmt.Errorf("couldn't find a user matching: %s", username)
		} else if len(users) != 1 {
			return fmt.Errorf("more than one user found matching: %s", username)
		}

		userId = users[0].ID
	} else if userIdOk {
		userId = userIdData.(int)
	} else {
		return fmt.Errorf("one and only one of user_id or username must be set")
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

func resourceGitlabGroupMembershipSetToState(d *schema.ResourceData, groupMember *gitlab.GroupMember, groupId *string) {

	d.Set("group_id", groupId)
	d.Set("user_id", groupMember.ID)
	d.Set("username", groupMember.Username)
	d.Set("access_level", accessLevel[groupMember.AccessLevel])
	d.Set("expires_at", groupMember.ExpiresAt)

	userId := strconv.Itoa(groupMember.ID)
	d.SetId(buildTwoPartID(groupId, &userId))
}
