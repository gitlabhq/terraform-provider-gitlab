//go:build acceptance
// +build acceptance

package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/xanzy/go-gitlab"
)

func TestAcc_GraphQL_basic(t *testing.T) {
	project := testAccCreateProject(t)
	projectToParse := &gitlab.Project{}

	_, err := SendGraphQLRequest(context.Background(), testGitlabClient, fmt.Sprintf(`
	query {
		project(fullPath: "%s") {
		  fullPath
		}
	  }
	`, project.NameWithNamespace), *projectToParse)

	fmt.Println(err)

}
