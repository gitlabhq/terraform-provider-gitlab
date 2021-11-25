package gitlab

import (
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/xanzy/go-gitlab"
	"net/http"
	"sort"
	"strings"
	"testing"
)

func TestAccGitlabRepositoryFiles_create(t *testing.T) {
	expectedFiles := []*gitlab.File{
		{
			FilePath: "my-file-1.txt",
			Content:  "some content 1",
		},
		{
			FilePath: "my-file-2.txt",
			Content:  "some content 2",
		}}
	actualFiles := make([]*gitlab.File, len(expectedFiles))
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabRepositoryFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabRepositoryFilesConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabRepositoryFilesExists("gitlab_repository_files.this", actualFiles, expectedFiles),
					testAccCheckGitlabRepositoryFilesAttributes(actualFiles, expectedFiles),
				),
			},
		},
	})
}

func TestAccGitlabRepositoryFiles_update(t *testing.T) {
	expectedFiles := []*gitlab.File{
		{
			FilePath: "my-file-1.txt",
			Content:  "some content 1",
		},
		{
			FilePath: "my-file-2.txt",
			Content:  "some content 2",
		}}
	expectedUpdatedFiles := []*gitlab.File{
		{
			FilePath: "my-file-2.txt",
			Content:  "some content new content",
		}, {
			FilePath: "my-file-3.txt",
			Content:  "some content 1",
		}}
	actualFiles := make([]*gitlab.File, len(expectedFiles))
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabRepositoryFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabRepositoryFilesConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabRepositoryFilesExists("gitlab_repository_files.this", actualFiles, expectedFiles),
					testAccCheckGitlabRepositoryFilesAttributes(actualFiles, expectedFiles),
				),
			},
			{
				Config: testAccGitlabRepositoryFilesUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabRepositoryFilesExists("gitlab_repository_files.this", actualFiles, expectedUpdatedFiles),
					testAccCheckGitlabRepositoryFilesAttributes(actualFiles, expectedUpdatedFiles),
				),
			},
		},
	})
}

func TestAccGitlabRepositoryFiles_remoteChangeIsDetected(t *testing.T) {
	expectedFiles := []*gitlab.File{
		{
			FilePath: "my-file-1.txt",
			Content:  "some content 1",
		},
		{
			FilePath: "my-file-2.txt",
			Content:  "some content 2",
		}}
	actualFiles := make([]*gitlab.File, len(expectedFiles))
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabRepositoryFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabRepositoryFilesConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabRepositoryFilesExists("gitlab_repository_files.this", actualFiles, expectedFiles),
					testAccCheckGitlabRepositoryFilesAttributes(actualFiles, expectedFiles),
				),
			},
			{
				Config: testAccGitlabRepositoryFilesConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabRepositoryFilesRemoveOneFile("gitlab_repository_files.this", expectedFiles[0]),
					testAccCheckGitlabRepositoryFilesChangeContentOfOneFile("gitlab_repository_files.this", expectedFiles[1]),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccGitlabRepositoryFilesConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabRepositoryFilesExists("gitlab_repository_files.this", actualFiles, expectedFiles),
					testAccCheckGitlabRepositoryFilesAttributes(actualFiles, expectedFiles),
				),
			},
		},
	})
}

func TestAccGitlabRepositoryFiles_createOnNewBranch(t *testing.T) {
	expectedFiles := []*gitlab.File{
		{
			FilePath: "my-file-1.txt",
			Content:  "some content 1",
		},
		{
			FilePath: "my-file-2.txt",
			Content:  "some content 2",
		}}
	actualFiles := make([]*gitlab.File, len(expectedFiles))
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabRepositoryFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabRepositoryFilesStartBranchConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabRepositoryFilesExists("gitlab_repository_files.this", actualFiles, expectedFiles),
					testAccCheckGitlabRepositoryFilesAttributes(actualFiles, expectedFiles),
				),
				// see https://gitlab.com/gitlab-org/gitlab/-/issues/342200
				SkipFunc: func() (bool, error) {
					return true, nil
				},
			},
		},
	})
}

func TestAccGitlabRepositoryFile_import(t *testing.T) {
	resourceName := "gitlab_repository_files.this"
	expectedFiles := []*gitlab.File{
		{
			FilePath: "my-file-1.txt",
			Content:  "some content 1",
		},
		{
			FilePath: "my-file-2.txt",
			Content:  "some content 2",
		}}
	actualFiles := make([]*gitlab.File, len(expectedFiles))
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabPipelineTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabRepositoryFilesConfig(rInt),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: getRepositoryFilesImportID(resourceName, expectedFiles),
				ImportState:       true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabRepositoryFilesExists("gitlab_repository_files.this", actualFiles, expectedFiles),
					testAccCheckGitlabRepositoryFilesAttributes(actualFiles, expectedFiles),
				),
			},
		},
	})
}
func getRepositoryFilesImportID(n string, expectedFiles []*gitlab.File) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("not Found: %s", n)
		}

		projectID := rs.Primary.Attributes["project"]
		if projectID == "" {
			return "", fmt.Errorf("no project ID is set")
		}
		branch := rs.Primary.Attributes["branch"]
		if branch == "" {
			return "", fmt.Errorf("no branch is set")
		}

		var filePaths []string
		for _, expectedFile := range expectedFiles {
			hash := fileHash(map[string]interface{}{
				"file_path": expectedFile.FilePath,
				"content":   expectedFile.Content,
			})

			filePath := rs.Primary.Attributes[fmt.Sprintf("file.%d.file_path", hash)]
			if filePath == "" {
				return "", fmt.Errorf("unable to find filepath %s in attributes", expectedFile.FilePath)
			}

			filePaths = append(filePaths, filePath)
		}

		return fmt.Sprintf("%s:%s:%s", projectID, branch, strings.Join(filePaths, ",")), nil
	}
}

func testAccCheckGitlabRepositoryFilesChangeContentOfOneFile(n string, file *gitlab.File) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		client := testAccProvider.Meta().(*gitlab.Client)
		projectId := rs.Primary.Attributes["project"]
		if projectId == "" {
			return fmt.Errorf("no project ID set")
		}
		branch := rs.Primary.Attributes["branch"]
		if branch == "" {
			return fmt.Errorf("no branch set")
		}
		_, _, err := client.RepositoryFiles.UpdateFile(projectId, file.FilePath, &gitlab.UpdateFileOptions{
			Branch:        gitlab.String(branch),
			CommitMessage: gitlab.String("deleting file"),
			Content:       gitlab.String("someone updated the content outside of state"),
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckGitlabRepositoryFilesRemoveOneFile(n string, file *gitlab.File) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		client := testAccProvider.Meta().(*gitlab.Client)
		projectId := rs.Primary.Attributes["project"]
		if projectId == "" {
			return fmt.Errorf("no project ID set")
		}
		branch := rs.Primary.Attributes["branch"]
		if branch == "" {
			return fmt.Errorf("no branch set")
		}

		_, err := client.RepositoryFiles.DeleteFile(projectId, file.FilePath, &gitlab.DeleteFileOptions{
			Branch:        gitlab.String(branch),
			CommitMessage: gitlab.String("deleting file"),
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckGitlabRepositoryFilesExists(n string, files []*gitlab.File, expectedFiles []*gitlab.File) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		// remove possible files from previous test step
		for i := 0; i < len(files); i++ {
			files[i] = nil
		}
		branch := rs.Primary.Attributes["branch"]
		if branch == "" {
			return fmt.Errorf("no branch set")
		}
		options := &gitlab.GetFileOptions{
			Ref: gitlab.String(branch),
		}
		projectId := rs.Primary.Attributes["project"]
		if projectId == "" {
			return fmt.Errorf("no project ID set")
		}

		conn := testAccProvider.Meta().(*gitlab.Client)

		for i, expectedFile := range expectedFiles {
			hash := fileHash(map[string]interface{}{
				"file_path": expectedFile.FilePath,
				"content":   expectedFile.Content,
			})
			filePath := rs.Primary.Attributes[fmt.Sprintf("file.%d.file_path", hash)]
			if filePath == "" {
				return fmt.Errorf("unable to find filepath %s in attributes", expectedFile.FilePath)
			}
			file, resp, err := conn.RepositoryFiles.GetFile(projectId, filePath, options)
			if err != nil || resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("cannot get file: %v", err)
			}

			files[i] = file
		}

		return nil
	}
}

func testAccCheckGitlabRepositoryFileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*gitlab.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_project" {
			continue
		}

		gotRepo, resp, err := client.Projects.GetProject(rs.Primary.ID, nil)
		if err == nil {
			if gotRepo != nil && fmt.Sprintf("%d", gotRepo.ID) == rs.Primary.ID {
				if gotRepo.MarkedForDeletionAt == nil {
					return fmt.Errorf("repository still exists")
				}
			}
		}
		if resp.StatusCode != http.StatusNotFound {
			return err
		}
		return nil
	}
	return nil
}

func testAccCheckGitlabRepositoryFilesAttributes(actual []*gitlab.File, expect []*gitlab.File) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(actual) != len(expect) {
			return fmt.Errorf("files received are length of %d, expect length %d", len(actual), len(expect))
		}

		sort.SliceStable(actual, func(i, j int) bool {
			return actual[i].FilePath < actual[j].FilePath
		})

		for i, actualFile := range actual {
			if actualFile.FilePath != expect[i].FilePath {
				return fmt.Errorf("actual name %s; expect %s", actualFile.FilePath, expect[i].FilePath)
			}

			actualContent, err := base64.StdEncoding.DecodeString(actualFile.Content)
			if err != nil {
				return err
			}
			if string(actualContent) != expect[i].Content {
				return fmt.Errorf("actual content %s; expect %s", string(actualContent), expect[i].Content)
			}
		}

		return nil
	}
}

func testAccGitlabRepositoryFilesConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name        = "foo-%d"
  description = "Terraform acceptance tests"

  default_branch = "main"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level       = "public"
  initialize_with_readme = true
}

resource "gitlab_repository_files" "this" {
  project        = gitlab_project.foo.id
  branch         = "main"
  author_email   = "email@example.com"
  author_name    = "Ola Nordmann"
  commit_message = "feature: multiple files"

  file {
    file_path = "my-file-1.txt"
    content   = "some content 1"
  }

  file {
    file_path = "my-file-2.txt"
    content   = "some content 2"
  }
 }
 	`, rInt)
}

func testAccGitlabRepositoryFilesStartBranchConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name        = "foo-%d"
  description = "Terraform acceptance tests"

  default_branch = "main"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level       = "public"
  initialize_with_readme = true
}

resource "gitlab_repository_files" "this" {
  project        = gitlab_project.foo.id
  branch         = "main"
  author_email   = "email@example.com"
  author_name    = "Ola Nordmann"
  commit_message = "feature: multiple files"
  start_branch = "some-new-distant-branch"

  file {
    file_path = "my-file-1.txt"
    content   = "some content 1"
  }

  file {
    file_path = "my-file-2.txt"
    content   = "some content 2"
  }
 }
 	`, rInt)
}

func testAccGitlabRepositoryFilesUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name        = "foo-%d"
  description = "Terraform acceptance tests"

  default_branch = "main"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level       = "public"
  initialize_with_readme = true
}

resource "gitlab_repository_files" "this" {
  project        = gitlab_project.foo.id
  branch         = "main"
  author_email   = "email@example.com"
  author_name    = "Ola Nordmann"
  commit_message = "feature: multiple files"

  file {
    file_path = "my-file-3.txt"
    content   = "some content 1"
  }

  file {
    file_path = "my-file-2.txt"
    content   = "some content new content"
  }
 }
 	`, rInt)
}
