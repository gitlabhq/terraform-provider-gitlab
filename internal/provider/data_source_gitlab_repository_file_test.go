//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataGitlabRepositoryFile_basic(t *testing.T) {
	project := testAccCreateProject(t)
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataGitlabRepositoryFile(project.PathWithNamespace),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGitlabRepositoryFile("gitlab_repository_file.foo", "data.gitlab_repository_file.foo"),
				),
			},
		},
	})
}

func testAccDataSourceGitlabRepositoryFile(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		file := s.RootModule().Resources[src]
		fileAttr := file.Primary.Attributes

		search := s.RootModule().Resources[n]
		searchAttr := search.Primary.Attributes

		testAttributes := []string{
			"project",
			"file_path",
			"size",
			"encoding",
			"content",
			"execute_filemode",
			"ref",
			"blob_id",
			"commit_id",
			"content_sha256",
			"last_commit_id",
		}

		for _, attribute := range testAttributes {
			if searchAttr[attribute] != fileAttr[attribute] {
				return fmt.Errorf("expected file's parameter `%s` to be: %s, but got: `%s`", attribute, fileAttr[attribute], searchAttr[attribute])
			}
		}
		return nil
	}
}

func testAccDataGitlabRepositoryFile(project string) string {
	return fmt.Sprintf(`
resource "gitlab_repository_file" "foo" {
	project = "%s"
	file_path = "testfile-meow"
	branch = "main"
	content = base64encode("Meow goes the cat")
	commit_message = "feat: Meow"
}

data "gitlab_repository_file" "foo" {
  project = gitlab_repository_file.foo.project
  file_path = gitlab_repository_file.foo.file_path
  ref = "main"
}
`, project)
}
