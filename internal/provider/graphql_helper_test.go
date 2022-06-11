//go:build acceptance
// +build acceptance

package provider

import (
	"context"
	"log"
	"testing"
)

func TestAcc_GraphQL_basic(t *testing.T) {

	query := GraphQLQuery{
		`query {currentUser {name, bot, gitpodEnabled, groupCount, id, namespace{id}, publicEmail, username}}`,
	}

	var response CurrentUserResponse
	_, err := SendGraphQLRequest(context.Background(), testGitlabClient, query, &response)
	if err != nil {
		log.Println(err)
		t.Fail()
	}

	if response.Data.CurrentUser.Name != "Administrator" {
		t.Fail()
	}
}
