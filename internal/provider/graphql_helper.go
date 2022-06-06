package provider

import (
	"context"

	"github.com/xanzy/go-gitlab"
)

// Helper method for modifying client requests appropriately for sending a GraphQL call instead of a REST call.
func SendGraphQLRequest(ctx context.Context, client *gitlab.Client, query GraphQLQuery, response interface{}) (interface{}, error) {
	request, err := client.NewRequest("POST", "", query, nil)
	if err != nil {
		return nil, err
	}
	// Overwrite the path of the existing request, as otherwise the go-gitlab client appends /api/v4 instead.
	request.URL.Path = "/api/graphql"
	if _, err = client.Do(request, response); err != nil {
		return nil, err
	}
	return response, nil
}

// Represents a GraphQL call to the API. All GraphQL calls are a string passed to the "query" parameter, so they should be included here.
type GraphQLQuery struct {
	Query string `json:"query"`
}
