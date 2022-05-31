//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/xanzy/go-gitlab"
)

func TestAcc_GraphQL_basic(t *testing.T) {
	project := testAccCreateProject(t)
	projectToParse := &gitlab.Project{}

	_, err := SendGraphQLRequest(nil, testGitlabClient, fmt.Sprintf(`
	query {
		project(fullPath: "%s") {
		  fullPath
		}
	  }
	`, project.Namespace.FullPath), *projectToParse)

	fmt.Println(err)

}
