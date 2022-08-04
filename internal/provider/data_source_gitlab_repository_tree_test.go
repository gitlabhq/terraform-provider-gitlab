//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGitlabRepositoryTree_basic(t *testing.T) {
	testProject := testAccCreateProject(t)
	testFile := testAccCreateProjectFile(t, testProject.ID, "content", "SomeFile", testProject.DefaultBranch)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "gitlab_repository_tree" "this" {
						project = %[1]d
						ref     = "%[2]s"
					}
				`, testProject.ID, testFile.Branch),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_repository_tree.this", "tree.#", "2"),
					resource.TestCheckResourceAttr("data.gitlab_repository_tree.this", "tree.0.name", "README.md"),
					resource.TestCheckResourceAttr("data.gitlab_repository_tree.this", "tree.0.type", "blob"),
					resource.TestCheckResourceAttr("data.gitlab_repository_tree.this", "tree.0.path", "README.md"),
					resource.TestCheckResourceAttr("data.gitlab_repository_tree.this", "tree.0.mode", "100644"),

					resource.TestCheckResourceAttr("data.gitlab_repository_tree.this", "tree.1.name", testFile.FilePath),
					resource.TestCheckResourceAttr("data.gitlab_repository_tree.this", "tree.1.type", "blob"),
					resource.TestCheckResourceAttr("data.gitlab_repository_tree.this", "tree.1.path", testFile.FilePath),
					resource.TestCheckResourceAttr("data.gitlab_repository_tree.this", "tree.1.mode", "100644"),
				),
			},
		},
	})
}
