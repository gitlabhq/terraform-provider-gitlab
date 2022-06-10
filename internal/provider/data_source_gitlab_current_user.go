package provider

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_current_user", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_current_user`" + ` data source allows details of the current user (determined by ` + "`token`" + ` provider attribute) to be retrieved.

**Upstream API**: [GitLab GraphQL API docs](https://docs.gitlab.com/ee/api/graphql/reference/index.html#querycurrentuser)`,

		ReadContext: dataSourceGitlabCurrentUserRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"global_id": {
				Description: "Global ID of the user. This is in the form of a GraphQL globally unique ID.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"username": {
				Description: "Username of the user. Unique within this instance of GitLab.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Human-readable name of the user. Returns **** if the user is a project bot and the requester does not have permission to view the project.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"bot": {
				Description: "Indicates if the user is a bot.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"group_count": {
				Description: "Group count for the user.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"namespace_id": {
				Description: "Personal namespace of the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"global_namespace_id": {
				Description: "Personal namespace of the user. This is in the form of a GraphQL globally unique ID.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"public_email": {
				Description: "User’s public email.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
})

func dataSourceGitlabCurrentUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	query := GraphQLQuery{
		`query {currentUser {name, bot, groupCount, id, namespace{id}, publicEmail, username}}`,
	}
	log.Printf("[DEBUG] executing GraphQL Query %s to retrieve current user", query.Query)

	var response CurrentUserResponse
	if _, err := SendGraphQLRequest(ctx, client, query, &response); err != nil {
		return diag.FromErr(err)
	}

	userID, err := extractIIDFromGlobalID(response.Data.CurrentUser.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	namespaceID, err := extractIIDFromGlobalID(response.Data.CurrentUser.Namespace.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", userID))
	d.Set("global_id", response.Data.CurrentUser.ID)
	d.Set("username", response.Data.CurrentUser.Username)
	d.Set("name", response.Data.CurrentUser.Name)
	d.Set("bot", response.Data.CurrentUser.Bot)
	d.Set("group_count", response.Data.CurrentUser.GroupCount)
	d.Set("namespace_id", fmt.Sprintf("%d", namespaceID))
	d.Set("global_namespace_id", response.Data.CurrentUser.Namespace.ID)
	d.Set("public_email", response.Data.CurrentUser.PublicEmail)

	return nil
}

// Struct representing current user based on the input API token
type CurrentUserResponse struct {
	Data struct {
		CurrentUser GraphQLUser `json:"currentUser"`
	} `json:"data"`
}

type GraphQLUser struct {
	Name       string `json:"name"`
	Bot        bool   `json:"bot"`
	GroupCount int    `json:"groupCount"`
	ID         string `json:"id"` // This is purposefully a string, as in some APIs it comes back as a globally unique ID
	Namespace  struct {
		ID string `json:"id"`
	} `json:"namespace"`
	PublicEmail string `json:"publicEmail"`
	Username    string `json:"username"`
}
