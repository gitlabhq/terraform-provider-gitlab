package provider

import (
	"context"

	"github.com/xanzy/go-gitlab"
)

// Helper method for modifying client requests appropriately for sending a GraphQL call instead of a REST call.
func SendGraphQLRequest(ctx context.Context, client *gitlab.Client, graphQLCall GraphQLQuery, objectToParseForResponse interface{}) (interface{}, error) {

	request, err := client.NewRequest("POST", "", graphQLCall, nil)
	//Overwrite the path of the existing request, as otherwise the client appends /api/v4 instead.
	request.URL.Path = "/api/graphql"
	if err != nil {
		return nil, err
	}

	_, err = client.Do(request, objectToParseForResponse)
	if err != nil {
		return nil, err
	}

	return objectToParseForResponse, nil
}

// Represents a GraphQL call to the API. All graphQL calls are a string passed to the "query" parameter, so they should be included here.
type GraphQLQuery struct {
	Query string `json:"query"`
}
