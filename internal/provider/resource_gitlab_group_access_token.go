package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	gitlab "github.com/xanzy/go-gitlab"
)

var validGroupAccessTokenScopes = []string{
	"api",
	"read_api",
	"read_registry",
	"write_registry",
	"read_repository",
	"write_repository",
}
var validAccessLevels = []string{
	"guest",
	"reporter",
	"developer",
	"maintainer",
	"owner",
}

var _ = registerResource("gitlab_group_access_token", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_group_access`" + `token resource allows to manage the lifecycle of a group access token.

-> Group Access Token were introduced in GitLab 14.7

**Upstream API**: [GitLab REST API](https://docs.gitlab.com/ee/api/group_access_tokens.html)`,

		CreateContext: resourceGitlabGroupAccessTokenCreate,
		ReadContext:   resourceGitlabGroupAccessTokenRead,
		DeleteContext: resourceGitlabGroupAccessTokenDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"group": {
				Description: "The ID or path of the group to add the group access token to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "The name of the group access token.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"scopes": {
				Description: fmt.Sprintf("The scope for the group access token. It determines the actions which can be performed when authenticating with this token. Valid values are: %s.", renderValueListForDocs(validGroupAccessTokenScopes)),
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(validGroupAccessTokenScopes, false),
				},
			},
			"access_level": {
				Description:      fmt.Sprintf("The access level for the group access token. Valid values are: %s.", renderValueListForDocs(validAccessLevels)),
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          accessLevelValueToName[gitlab.MaintainerPermissions],
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(validAccessLevels, false)),
			},
			"expires_at": {
				Description:      "The token expires at midnight UTC on that date. The date must be in the format YYYY-MM-DD. Default is never.",
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateDiagFunc: isISO6801Date,
			},
			"token": {
				Description: "The group access token. This is only populated when creating a new group access token. This attribute is not available for imported resources.",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
			"active": {
				Description: "True if the token is active.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"created_at": {
				Description: "Time the token has been created, RFC3339 format.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"revoked": {
				Description: "True if the token is revoked.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"user_id": {
				Description: "The user id associated to the token.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
})

func resourceGitlabGroupAccessTokenCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	group := d.Get("group").(string)
	options := &gitlab.CreateGroupAccessTokenOptions{
		Name:   gitlab.String(d.Get("name").(string)),
		Scopes: stringSetToStringSlice(d.Get("scopes").(*schema.Set)),
	}
	if v, ok := d.GetOk("access_level"); ok {
		accessLevel := accessLevelNameToValue[v.(string)]
		options.AccessLevel = &accessLevel
	}

	log.Printf("[DEBUG] create gitlab GroupAccessToken %s (scopes: %s, access_level: %v) for group ID %s", *options.Name, options.Scopes, options.AccessLevel, group)

	if v, ok := d.GetOk("expires_at"); ok {
		parsedExpiresAt, err := time.Parse("2006-01-02", v.(string))
		if err != nil {
			return diag.Errorf("Invalid expires_at date: %v", err)
		}
		parsedExpiresAtISOTime := gitlab.ISOTime(parsedExpiresAt)
		options.ExpiresAt = &parsedExpiresAtISOTime
		log.Printf("[DEBUG] create gitlab GroupAccessToken %s with expires_at %s for group ID %s", *options.Name, *options.ExpiresAt, group)
	}

	groupAccessToken, _, err := client.GroupAccessTokens.CreateGroupAccessToken(group, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] created gitlab GroupAccessToken %d - %s for group ID %s", groupAccessToken.ID, *options.Name, group)

	tokenId := strconv.Itoa(groupAccessToken.ID)
	d.SetId(buildTwoPartID(&group, &tokenId))
	// NOTE: the token can only be read once after creating it
	d.Set("token", groupAccessToken.Token)

	return resourceGitlabGroupAccessTokenRead(ctx, d, meta)
}

func resourceGitlabGroupAccessTokenRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	group, tokenId, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing ID: %s", d.Id())
	}

	client := meta.(*gitlab.Client)

	groupAccessTokenId, err := strconv.Atoi(tokenId)
	if err != nil {
		return diag.Errorf("%s cannot be converted to int", tokenId)
	}

	log.Printf("[DEBUG] read gitlab GroupAccessToken %d, group ID %s", groupAccessTokenId, group)
	groupAccessToken, _, err := client.GroupAccessTokens.GetGroupAccessToken(group, groupAccessTokenId, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] GitLab GroupAccessToken %d, group ID %s not found, removing from state", groupAccessTokenId, group)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("group", group)
	d.Set("name", groupAccessToken.Name)
	if groupAccessToken.ExpiresAt != nil {
		d.Set("expires_at", groupAccessToken.ExpiresAt.String())
	}
	d.Set("active", groupAccessToken.Active)
	d.Set("created_at", groupAccessToken.CreatedAt.Format(time.RFC3339))
	d.Set("access_level", accessLevelValueToName[groupAccessToken.AccessLevel])
	d.Set("revoked", groupAccessToken.Revoked)
	d.Set("user_id", groupAccessToken.UserID)

	if err = d.Set("scopes", groupAccessToken.Scopes); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGitlabGroupAccessTokenDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	group, tokenId, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing ID: %s", d.Id())
	}

	client := meta.(*gitlab.Client)

	groupAccessTokenId, err := strconv.Atoi(tokenId)
	if err != nil {
		return diag.Errorf("%s cannot be converted to int", tokenId)
	}

	log.Printf("[DEBUG] Delete gitlab GroupAccessToken %s", d.Id())
	_, err = client.GroupAccessTokens.RevokeGroupAccessToken(group, groupAccessTokenId, gitlab.WithContext(ctx))
	return diag.FromErr(err)
}
