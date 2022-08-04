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

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "gitlab_repository_file" "foo" {
						project        = "%[1]d"
						file_path      = "testfile-meow"
						branch         = "%[2]s"
						content        = base64encode("Meow goes the cat")
						commit_message = "feat: Meow"
					}

					data "gitlab_repository_tree" "this" {
						project = %[1]d
						ref     = gitlab_repository_file.foo.branch
					}
				`, testProject.ID, testProject.DefaultBranch),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_repository_tree.this", "tree.#", "2"),
					resource.TestCheckResourceAttr("data.gitlab_repository_tree.this", "tree.0.name", "README.md"),
					resource.TestCheckResourceAttr("data.gitlab_repository_tree.this", "tree.0.type", "blob"),
					resource.TestCheckResourceAttr("data.gitlab_repository_tree.this", "tree.0.path", "README.md"),
					resource.TestCheckResourceAttr("data.gitlab_repository_tree.this", "tree.0.mode", "100644"),

					resource.TestCheckResourceAttr("data.gitlab_repository_tree.this", "tree.1.name", "testfile-meow"),
					resource.TestCheckResourceAttr("data.gitlab_repository_tree.this", "tree.1.type", "blob"),
					resource.TestCheckResourceAttr("data.gitlab_repository_tree.this", "tree.1.path", "testfile-meow"),
					resource.TestCheckResourceAttr("data.gitlab_repository_tree.this", "tree.1.mode", "100644"),
				),
			},
		},
	})
}
