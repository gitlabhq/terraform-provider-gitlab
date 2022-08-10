package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_group_membership", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_group_membership`" + ` resource allows to manage the lifecycle of a users group membersip.

-> If a group should grant membership to another group use the ` + "`gitlab_group_share_group`" + ` resource instead.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/members.html)`,

		CreateContext: resourceGitlabGroupMembershipCreate,
		ReadContext:   resourceGitlabGroupMembershipRead,
		UpdateContext: resourceGitlabGroupMembershipUpdate,
		DeleteContext: resourceGitlabGroupMembershipDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"group_id": {
				Description: "The id of the group.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},
			"user_id": {
				Description: "The id of the user.",
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
			},
			"access_level": {
				Description:      fmt.Sprintf("Access level for the member. Valid values are: %s.", renderValueListForDocs(validGroupAccessLevelNames)),
				Type:             schema.TypeString,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(validGroupAccessLevelNames, false)),
				Required:         true,
			},
			"expires_at": {
				Description:  "Expiration date for the group membership. Format: `YYYY-MM-DD`",
				Type:         schema.TypeString,
				ValidateFunc: validateDateFunc,
				Optional:     true,
			},
			"skip_subresources_on_destroy": {
				Description: "Whether the deletion of direct memberships of the removed member in subgroups and projects should be skipped. Only used during a destroy.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"unassign_issuables_on_destroy": {
				Description: "Whether the removed member should be unassigned from any issues or merge requests inside a given group or project. Only used during a destroy.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
		},
	}
})

func resourceGitlabGroupMembershipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	userId := d.Get("user_id").(int)
	groupId := d.Get("group_id").(string)
	expiresAt := d.Get("expires_at").(string)
	accessLevelId := accessLevelNameToValue[d.Get("access_level").(string)]

	options := &gitlab.AddGroupMemberOptions{
		UserID:      &userId,
		AccessLevel: &accessLevelId,
		ExpiresAt:   &expiresAt,
	}
	log.Printf("[DEBUG] create gitlab group groupMember for %d in %s", options.UserID, groupId)

	groupMember, _, err := client.GroupMembers.AddGroupMember(groupId, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}
	userIdString := strconv.Itoa(groupMember.ID)
	d.SetId(buildTwoPartID(&groupId, &userIdString))
	return resourceGitlabGroupMembershipRead(ctx, d, meta)
}

func resourceGitlabGroupMembershipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	id := d.Id()
	log.Printf("[DEBUG] read gitlab group groupMember %s", id)

	groupId, userId, err := groupIdAndUserIdFromId(id)
	if err != nil {
		return diag.FromErr(err)
	}

	groupMember, _, err := client.GroupMembers.GetGroupMember(groupId, userId, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] gitlab group membership for %s not found so removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
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

func resourceGitlabGroupMembershipUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	userId := d.Get("user_id").(int)
	groupId := d.Get("group_id").(string)
	expiresAt := d.Get("expires_at").(string)
	accessLevelId := accessLevelNameToValue[strings.ToLower(d.Get("access_level").(string))]

	options := gitlab.EditGroupMemberOptions{
		AccessLevel: &accessLevelId,
		ExpiresAt:   &expiresAt,
	}
	log.Printf("[DEBUG] update gitlab group membership %v for %s", userId, groupId)

	_, _, err := client.GroupMembers.EditGroupMember(groupId, userId, &options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGitlabGroupMembershipRead(ctx, d, meta)
}

func resourceGitlabGroupMembershipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	id := d.Id()
	groupId, userId, err := groupIdAndUserIdFromId(id)
	if err != nil {
		return diag.FromErr(err)
	}

	options := gitlab.RemoveGroupMemberOptions{
		SkipSubresources:  gitlab.Bool(d.Get("skip_subresources_on_destroy").(bool)),
		UnassignIssuables: gitlab.Bool(d.Get("unassign_issuables_on_destroy").(bool)),
	}

	log.Printf("[DEBUG] Delete gitlab group membership %v for %s with options: %+v", userId, groupId, options)

	_, err = client.GroupMembers.RemoveGroupMember(groupId, userId, &options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGitlabGroupMembershipSetToState(d *schema.ResourceData, groupMember *gitlab.GroupMember, groupId *string) {

	d.Set("group_id", groupId)
	d.Set("user_id", groupMember.ID)
	d.Set("access_level", accessLevelValueToName[groupMember.AccessLevel])
	if groupMember.ExpiresAt != nil {
		d.Set("expires_at", groupMember.ExpiresAt.String())
	} else {
		d.Set("expires_at", "")
	}
	userId := strconv.Itoa(groupMember.ID)
	d.SetId(buildTwoPartID(groupId, &userId))
}
