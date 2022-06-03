package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_current_user", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_current_user`" + ` data source allows details of the current user (determined by the input provider token) to be retrieved.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/graphql/reference/index.html#querycurrentuser)`,

		ReadContext: dataSourceGitlabCurrentUserRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the branch.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
})

func dataSourceGitlabCurrentUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	query := GraphQLQuery{
		`query {currentUser {name, bot, gitpodEnabled, groupCount, id, namespace{id} publicEmail, username}}`,
	}
	var response CurrentUserResponse
	_, err := SendGraphQLRequest(context.Background(), client, query, &response)
	if err != nil {
		log.Println(err)

	}

	if response.Data.CurrentUser.Name != "Administrator" {
		log.Println("neat")
	}

	return nil
}

// Struct representing current user based on the input API token
type CurrentUserResponse struct {
	Data struct {
		CurrentUser GraphQLUser `json:"currentUser"`
	} `json:"data"`
}

type GraphQLUser struct {
	Name          string `json:"name"`
	Bot           bool   `json:"bot"`
	GitPodEnabled bool   `json:"gitpodEnabled"`
	GroupCount    int    `json:"groupCount"`
	ID            string `json:"id"` // This is purposefully a string, as in some APIs it comes back as a globally unique ID
	Namespace     struct {
		ID string `json:"id"`
	} `json:"namespace"`
	PublicEmail string `json:"publicEmail"`
	Username    string `json:"username"`
}
