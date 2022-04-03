package provider

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	gitlab "github.com/xanzy/go-gitlab"
)

var validPersonalAccessTokenScopes = []string{
	"api",
	"read_user",
	"read_api",
	"read_repository",
	"write_repository",
	"read_registry",
	"write_registry",
	"sudo",
}

var _ = registerResource("gitlab_personal_access_token", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_personal_access_token`" + ` resource allows to manage the lifecycle of a personal access token for a specified user.

-> **Note** In order to utilize this resource without failure the executing user must be an admin. 

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/personal_access_tokens.html)`,

		CreateContext: resourceGitlabPersonalAccessTokenCreate,
		ReadContext:   resourceGitlabPersonalAccessTokenRead,
		DeleteContext: resourceGitlabPersonalAccessTokenDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"user_id": {
				Description: "The id of the user.",
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
			},
			"name": {
				Description: "The name of the personal access token.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"scopes": {
				Description: fmt.Sprintf("The scope for the personal access token. It determines the actions which can be performed when authenticating with this token. Valid values are: %s.", renderValueListForDocs(validPersonalAccessTokenScopes)),
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(validPersonalAccessTokenScopes, false),
				},
			},
			"active": {
				Description: "True if the token is active.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"revoked": {
				Description: "True if the token is revoked.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"created_at": {
				Description: "Time the token has been created, RFC3339 format.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"expires_at": {
				Description:      "The token expires at midnight UTC on that date. The date must be in the format YYYY-MM-DD. Default is never.",
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateDiagFunc: isISO6801Date,
			},
			"token": {
				Description: "The personal access token. This is only populated when creating a new personal access token. This attribute is not available for imported resources.",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
})

func resourceGitlabPersonalAccessTokenCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	userId := d.Get("user_id").(int)

	currentUserAdmin, err := isCurrentUserAdmin(client)
	if err != nil {
		return diag.Errorf("[ERROR] cannot query the user API for current user: %v", err)
	}

	if !currentUserAdmin {
		return diag.Errorf("[ERROR] Current user needs to be admin when creating a personal access token: %v", err)
	}

	options := &gitlab.CreatePersonalAccessTokenOptions{
		Name:   gitlab.String(d.Get("name").(string)),
		Scopes: stringSetToStringSlice(d.Get("scopes").(*schema.Set)),
	}

	log.Printf("[DEBUG] create gitlab PersonalAccessToken %s (scopes: %s) for user ID %d", *options.Name, options.Scopes, userId)

	if v, ok := d.GetOk("expires_at"); ok {
		parsedExpiresAt, err := time.Parse("2006-01-02", v.(string))
		if err != nil {
			return diag.Errorf("Invalid expires_at date: %v", err)
		}

		parsedExpiresAtISOTime := gitlab.ISOTime(parsedExpiresAt)
		options.ExpiresAt = &parsedExpiresAtISOTime
	}

	personalAccessToken, _, err := client.Users.CreatePersonalAccessToken(userId, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	tokenId := strconv.Itoa(personalAccessToken.ID)
	userIdString := strconv.Itoa(userId)

	d.SetId(buildTwoPartID(&userIdString, &tokenId))
	// NOTE: the token can only be read once after creating it
	d.Set("token", personalAccessToken.Token)

	return resourceGitlabPersonalAccessTokenRead(ctx, d, meta)
}

func resourceGitlabPersonalAccessTokenRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	userId, tokenId, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing ID: %s", d.Id())
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return diag.Errorf("Error converting user ID to string: %v", userId)
	}

	client := meta.(*gitlab.Client)

	personalAccessTokenId, err := strconv.Atoi(tokenId)
	if err != nil {
		return diag.Errorf("%s cannot be converted to int", tokenId)
	}

	log.Printf("[DEBUG] read gitlab PersonalAccessToken %d, user ID %s", personalAccessTokenId, userId)

	personalAccessToken, err := resourceGitlabPersonalAccessTokenFind(ctx, client, userIdInt, personalAccessTokenId)
	if errors.Is(err, errResourceGitlabPersonalAccessTokenNotFound) {
		log.Printf("[DEBUG] failed to read gitlab PersonalAccessToken %d, user ID %s", personalAccessTokenId, userId)
		d.SetId("")
	}

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("user_id", userIdInt)
	d.Set("name", personalAccessToken.Name)
	if personalAccessToken.ExpiresAt != nil {
		d.Set("expires_at", personalAccessToken.ExpiresAt.String())
	}
	d.Set("active", personalAccessToken.Active)
	d.Set("created_at", personalAccessToken.CreatedAt.Format(time.RFC3339))
	d.Set("revoked", personalAccessToken.Revoked)

	err = d.Set("scopes", personalAccessToken.Scopes)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGitlabPersonalAccessTokenDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, tokenId, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing ID: %s", d.Id())
	}

	client := meta.(*gitlab.Client)

	personalAccessTokenId, err := strconv.Atoi(tokenId)
	if err != nil {
		return diag.Errorf("%s cannot be converted to int", tokenId)
	}

	log.Printf("[DEBUG] Delete gitlab PersonalAccessToken %s", d.Id())
	_, err = client.PersonalAccessTokens.RevokePersonalAccessToken(personalAccessTokenId, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(err)
}

var errResourceGitlabPersonalAccessTokenNotFound = errors.New("personal access token not found")

// resourceGitlabPersonalAccessTokenFind finds the personal access token with the specified tokenID.
// It returns a errResourceGitlabPersonalAccessTokenNotFound error if the token is not found or in a revoked state.
func resourceGitlabPersonalAccessTokenFind(ctx context.Context, client *gitlab.Client, userId int, personalAccessTokenID int) (*gitlab.PersonalAccessToken, error) {
	//there is a slight possibility to not find an existing item, for example
	// 1. item is #101 (ie, in the 2nd page)
	// 2. I load first page (ie. I don't find my target item)
	// 3. A concurrent operation remove item 99 (ie, my target item shift to 1st page)
	// 4. a concurrent operation add an item
	// 5: I load 2nd page  (ie. I don't find my target item)
	// 6. Total pages and total items properties are unchanged (from the perspective of the reader)

	page := 1
	for page != 0 {
		personalAccessTokens, response, err := client.PersonalAccessTokens.ListPersonalAccessTokens(&gitlab.ListPersonalAccessTokensOptions{UserID: &userId, ListOptions: gitlab.ListOptions{Page: page, PerPage: 100}}, gitlab.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		for _, personalAccessToken := range personalAccessTokens {
			if personalAccessToken.ID == personalAccessTokenID && !personalAccessToken.Revoked {
				return personalAccessToken, nil
			}
		}

		page = response.NextPage
	}

	return nil, errResourceGitlabPersonalAccessTokenNotFound
}
