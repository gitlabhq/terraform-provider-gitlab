//go:build acceptance
// +build acceptance

package provider

import (
	"context"
	"testing"
)

func TestAcc_GraphQL_basic(t *testing.T) {

	query := GraphQLQuery{
		`query {currentUser {name}}`,
	}

	var response CurrentUserResponse
	_, _ = SendGraphQLRequest(context.Background(), testGitlabClient, query, &response)

	if response.Data.CurrentUser.Name != "Administrator" {
		t.Fail()
	}
}

type CurrentUserResponse struct {
	Data struct {
		CurrentUser struct {
			Name string `json:"name"`
		} `json:"currentUser"`
	} `json:"data"`
}
