package provider

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_project_access_token", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`" + `gitlab_project_access_token` + "`" + ` resource allows to manage the lifecycle of a project access token.

**Upstream API**: [GitLab API docs](https://docs.gitlab.com/ee/api/project_access_tokens.html)`,

		CreateContext: resourceGitlabProjectAccessTokenCreate,
		ReadContext:   resourceGitlabProjectAccessTokenRead,
		DeleteContext: resourceGitlabProjectAccessTokenDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The id of the project to add the project access token to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "A name to describe the project access token.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"scopes": {
				Description: "Valid values: `api`, `read_api`, `read_repository`, `write_repository`.",
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"api", "read_api", "read_repository", "write_repository"}, false),
				},
			},
			"expires_at": {
				Description:      "Time the token will expire it, YYYY-MM-DD format. Will not expire per default.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: isISO6801Date,
				ForceNew:         true,
			},
			"token": {
				Description: "The secret token. **Note**: the token is not available for imported resources.",
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
				Description: "The user_id associated to the token.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"access_level": {
				Description:      fmt.Sprintf("The access level for the project access token. Valid values are: %s. Default is `%s`.", renderValueListForDocs(validProjectAccessLevelNames), accessLevelValueToName[gitlab.MaintainerPermissions]),
				Type:             schema.TypeString,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(validProjectAccessLevelNames, false)),
				Optional:         true,
				Default:          accessLevelValueToName[gitlab.MaintainerPermissions],
				ForceNew:         true,
			},
		},
	}
})

func resourceGitlabProjectAccessTokenCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	accessLevelId := accessLevelNameToValue[d.Get("access_level").(string)]
	project := d.Get("project").(string)

	options := &gitlab.CreateProjectAccessTokenOptions{
		Name:        gitlab.String(d.Get("name").(string)),
		Scopes:      stringSetToStringSlice(d.Get("scopes").(*schema.Set)),
		AccessLevel: &accessLevelId,
	}

	log.Printf("[DEBUG] create gitlab ProjectAccessToken %s %s for project ID %s", *options.Name, options.Scopes, project)

	if v, ok := d.GetOk("expires_at"); ok {
		parsedExpiresAt, err := time.Parse("2006-01-02", v.(string))
		if err != nil {
			return diag.Errorf("Invalid expires_at date: %v", err)
		}
		parsedExpiresAtISOTime := gitlab.ISOTime(parsedExpiresAt)
		options.ExpiresAt = &parsedExpiresAtISOTime
	}

	projectAccessToken, _, err := client.ProjectAccessTokens.CreateProjectAccessToken(project, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	PATstring := strconv.Itoa(projectAccessToken.ID)
	d.SetId(buildTwoPartID(&project, &PATstring))
	d.Set("token", projectAccessToken.Token)

	return resourceGitlabProjectAccessTokenRead(ctx, d, meta)
}

func resourceGitlabProjectAccessTokenRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	project, PATstring, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing ID: %s", d.Id())
	}

	client := meta.(*gitlab.Client)

	projectAccessTokenID, err := strconv.Atoi(PATstring)
	if err != nil {
		return diag.Errorf("%s cannot be converted to int", PATstring)
	}

	log.Printf("[DEBUG] read gitlab ProjectAccessToken %d, project ID %s", projectAccessTokenID, project)

	projectAccessToken, err := resourceGitlabProjectAccessTokenFind(ctx, client, project, projectAccessTokenID)
	if errors.Is(err, errResourceGitlabProjectAccessTokenNotFound) {
		log.Printf("[DEBUG] failed to read gitlab ProjectAccessToken %d, project ID %s", projectAccessTokenID, project)
		d.SetId("")
	}
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("project", project)
	d.Set("name", projectAccessToken.Name)
	if projectAccessToken.ExpiresAt != nil {
		d.Set("expires_at", projectAccessToken.ExpiresAt.String())
	}
	d.Set("active", projectAccessToken.Active)
	d.Set("created_at", projectAccessToken.CreatedAt.String())
	d.Set("revoked", projectAccessToken.Revoked)
	d.Set("user_id", projectAccessToken.UserID)
	d.Set("access_level", accessLevelValueToName[projectAccessToken.AccessLevel])
	err = d.Set("scopes", projectAccessToken.Scopes)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGitlabProjectAccessTokenDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	project, patString, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing ID: %s", d.Id())
	}

	client := meta.(*gitlab.Client)

	projectAccessTokenID, err := strconv.Atoi(patString)
	if err != nil {
		return diag.Errorf("%s cannot be converted to int", patString)
	}

	log.Printf("[DEBUG] Delete gitlab ProjectAccessToken %s", d.Id())
	_, err = client.ProjectAccessTokens.RevokeProjectAccessToken(project, projectAccessTokenID, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Waiting for ProjectAccessToken %s to finish deleting", d.Id())

	err = resource.RetryContext(ctx, 5*time.Minute, func() *resource.RetryError {
		_, err := resourceGitlabProjectAccessTokenFind(ctx, client, project, projectAccessTokenID)
		if errors.Is(err, errResourceGitlabProjectAccessTokenNotFound) {
			return nil
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(errors.New("project access token was not deleted"))
	})

	return diag.FromErr(err)
}

var errResourceGitlabProjectAccessTokenNotFound = errors.New("project access token not found")

// resourceGitlabProjectAccessTokenFind finds the project access token with the specified tokenID.
// It returns a errResourceGitlabProjectAccessTokenNotFound error if the token is not found.
func resourceGitlabProjectAccessTokenFind(ctx context.Context, client *gitlab.Client, project interface{}, projectAccessTokenID int) (*gitlab.ProjectAccessToken, error) {
	//there is a slight possibility to not find an existing item, for example
	// 1. item is #101 (ie, in the 2nd page)
	// 2. I load first page (ie. I don't find my target item)
	// 3. A concurrent operation remove item 99 (ie, my target item shift to 1st page)
	// 4. a concurrent operation add an item
	// 5: I load 2nd page  (ie. I don't find my target item)
	// 6. Total pages and total items properties are unchanged (from the perspective of the reader)

	page := 1
	for page != 0 {
		projectAccessTokens, response, err := client.ProjectAccessTokens.ListProjectAccessTokens(project, &gitlab.ListProjectAccessTokensOptions{Page: page, PerPage: 100}, gitlab.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		for _, projectAccessToken := range projectAccessTokens {
			if projectAccessToken.ID == projectAccessTokenID {
				return projectAccessToken, nil
			}
		}

		page = response.NextPage
	}

	return nil, errResourceGitlabProjectAccessTokenNotFound
}
