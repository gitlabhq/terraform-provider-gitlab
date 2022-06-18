//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

func TestAccGitlabRepositoryFile_basic(t *testing.T) {
	var file gitlab.File
	testProject := testAccCreateProject(t)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabRepositoryFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabRepositoryFileConfig(testProject.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabRepositoryFileExists("gitlab_repository_file.this", &file),
					testAccCheckGitlabRepositoryFileAttributes(&file, &testAccGitlabRepositoryFileAttributes{
						FilePath: "meow.txt",
						Content:  "bWVvdyBtZW93IG1lb3c=",
					}),
					resource.TestCheckResourceAttr("gitlab_repository_file.this", "content", "bWVvdyBtZW93IG1lb3c="),
				),
			},
			{
				ResourceName:            "gitlab_repository_file.this",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"author_email", "author_name", "commit_message"},
			},
			{
				Config: testAccGitlabRepositoryFileUpdateConfig(testProject.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabRepositoryFileExists("gitlab_repository_file.this", &file),
					testAccCheckGitlabRepositoryFileAttributes(&file, &testAccGitlabRepositoryFileAttributes{
						FilePath: "meow.txt",
						Content:  "bWVvdyBtZW93IG1lb3cgbWVvdyBtZW93Cg==",
					}),
					resource.TestCheckResourceAttr("gitlab_repository_file.this", "content", "bWVvdyBtZW93IG1lb3cgbWVvdyBtZW93Cg=="),
				),
			},
			{
				ResourceName:            "gitlab_repository_file.this",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"author_email", "author_name", "commit_message"},
			},
		},
	})
}

func TestAccGitlabRepositoryFile_createSameFileDifferentRepository(t *testing.T) {
	var fooFile gitlab.File
	var barFile gitlab.File
	firstTestProject := testAccCreateProject(t)
	secondTestProject := testAccCreateProject(t)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabRepositoryFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabRepositoryFileSameFileDifferentRepositoryConfig(firstTestProject.ID, secondTestProject.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabRepositoryFileExists("gitlab_repository_file.foo_file", &fooFile),
					testAccCheckGitlabRepositoryFileAttributes(&fooFile, &testAccGitlabRepositoryFileAttributes{
						FilePath: "meow.txt",
						Content:  "bWVvdyBtZW93IG1lb3c=",
					}),
					testAccCheckGitlabRepositoryFileExists("gitlab_repository_file.bar_file", &barFile),
					testAccCheckGitlabRepositoryFileAttributes(&barFile, &testAccGitlabRepositoryFileAttributes{
						FilePath: "meow.txt",
						Content:  "bWVvdyBtZW93IG1lb3c=",
					}),
				),
			},
		},
	})
}

func TestAccGitlabRepositoryFile_concurrentResources(t *testing.T) {
	testProject := testAccCreateProject(t)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabRepositoryFileDestroy,
		Steps: []resource.TestStep{
			// NOTE: we don't need to check anything here, just make sure no terraform errors are being raised,
			//       the other test cases will do the actual testing :)
			{
				Config: testAccGitlabRepositoryFileConcurrentResourcesConfig(testProject.ID),
			},
			{
				Config: testAccGitlabRepositoryFileConcurrentResourcesConfigUpdate(testProject.ID),
			},
			{
				Config:  testAccGitlabRepositoryFileConcurrentResourcesConfigUpdate(testProject.ID),
				Destroy: true,
			},
		},
	})
}

func TestAccGitlabRepositoryFile_createOnNewBranch(t *testing.T) {
	var file gitlab.File
	testProject := testAccCreateProject(t)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabRepositoryFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabRepositoryFileStartBranchConfig(testProject.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabRepositoryFileExists("gitlab_repository_file.this", &file),
					testAccCheckGitlabRepositoryFileAttributes(&file, &testAccGitlabRepositoryFileAttributes{
						FilePath: "meow.txt",
						Content:  "bWVvdyBtZW93IG1lb3c=",
					}),
				),
			},
		},
	})
}

func TestAccGitlabRepositoryFile_base64EncodingWithTextContent(t *testing.T) {
	var file gitlab.File
	testProject := testAccCreateProject(t)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabRepositoryFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "gitlab_repository_file" "this" {
						project = %d
						file_path = "meow.txt"
						branch = "main"
						content = "Hello World, meow"
						author_email = "meow@catnip.com"
						author_name = "Meow Meowington"
						commit_message = "feature: add launch codes"
					}
				`, testProject.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabRepositoryFileExists("gitlab_repository_file.this", &file),
					testAccCheckGitlabRepositoryFileAttributes(&file, &testAccGitlabRepositoryFileAttributes{
						FilePath: "meow.txt",
						Content:  "SGVsbG8gV29ybGQsIG1lb3c=",
					}),
					resource.TestCheckResourceAttr("gitlab_repository_file.this", "content", "Hello World, meow"),
				),
			},
			// Update content
			{
				Config: fmt.Sprintf(`
					resource "gitlab_repository_file" "this" {
						project = %d
						file_path = "meow.txt"
						branch = "main"
						content = "Hello World, meow UPDATED"
						author_email = "meow@catnip.com"
						author_name = "Meow Meowington"
						commit_message = "feature: add launch codes"
					}
				`, testProject.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabRepositoryFileExists("gitlab_repository_file.this", &file),
					testAccCheckGitlabRepositoryFileAttributes(&file, &testAccGitlabRepositoryFileAttributes{
						FilePath: "meow.txt",
						Content:  "SGVsbG8gV29ybGQsIG1lb3cgVVBEQVRFRA==",
					}),
					resource.TestCheckResourceAttr("gitlab_repository_file.this", "content", "Hello World, meow UPDATED"),
				),
			},
			// Update to already base64 encoded content, but same content
			{
				Config: fmt.Sprintf(`
					resource "gitlab_repository_file" "this" {
						project = %d
						file_path = "meow.txt"
						branch = "main"
						content = "SGVsbG8gV29ybGQsIG1lb3cgVVBEQVRFRA=="
						author_email = "meow@catnip.com"
						author_name = "Meow Meowington"
						commit_message = "feature: add launch codes"
					}
				`, testProject.ID),
				Check: resource.TestCheckResourceAttr("gitlab_repository_file.this", "content", "SGVsbG8gV29ybGQsIG1lb3cgVVBEQVRFRA=="),
			},
			// Update content and already base64 encode it
			{
				Config: fmt.Sprintf(`
					resource "gitlab_repository_file" "this" {
						project = %d
						file_path = "meow.txt"
						branch = "main"
						content = "bWVvdyBtZW93IG1lb3c="
						author_email = "meow@catnip.com"
						author_name = "Meow Meowington"
						commit_message = "feature: add launch codes"
					}
				`, testProject.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabRepositoryFileExists("gitlab_repository_file.this", &file),
					testAccCheckGitlabRepositoryFileAttributes(&file, &testAccGitlabRepositoryFileAttributes{
						FilePath: "meow.txt",
						Content:  "bWVvdyBtZW93IG1lb3c=",
					}),
					resource.TestCheckResourceAttr("gitlab_repository_file.this", "content", "bWVvdyBtZW93IG1lb3c="),
				),
			},
			// Update to not-yet base64 encoded content, but same content
			{
				Config: fmt.Sprintf(`
					resource "gitlab_repository_file" "this" {
						project = %d
						file_path = "meow.txt"
						branch = "main"
						content = "meow meow meow"
						author_email = "meow@catnip.com"
						author_name = "Meow Meowington"
						commit_message = "feature: add launch codes"
					}
				`, testProject.ID),
				Check: resource.TestCheckResourceAttr("gitlab_repository_file.this", "content", "meow meow meow"),
			},
		},
	})
}

func TestAccGitlabRepositoryFile_createWithExecuteFilemode(t *testing.T) {
	testProject := testAccCreateProject(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccRequiresAtLeast(t, "14.10") },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabRepositoryFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "gitlab_repository_file" "this" {
						project = %d
						file_path = "meow.txt"
						branch = "main"
						content = "bWVvdyBtZW93IG1lb3cgbWVvdyBtZW93Cg=="
						author_email = "meow@catnip.com"
						author_name = "Meow Meowington"
						commit_message = "feature: change launch codes"
						execute_filemode = false
					}
				`, testProject.ID),
			},
			// Verify Import
			{
				ResourceName:            "gitlab_repository_file.this",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"author_email", "author_name", "commit_message"},
			},
			{
				Config: fmt.Sprintf(`
					resource "gitlab_repository_file" "this" {
						project = %d
						file_path = "meow.txt"
						branch = "main"
						content = "bWVvdyBtZW93IG1lb3cgbWVvdyBtZW93Cg=="
						author_email = "meow@catnip.com"
						author_name = "Meow Meowington"
						commit_message = "feature: change launch codes"
						execute_filemode = true
					}
				`, testProject.ID),
			},
			// Verify Import
			{
				ResourceName:            "gitlab_repository_file.this",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"author_email", "author_name", "commit_message"},
			},
		},
	})
}

func testAccCheckGitlabRepositoryFileExists(n string, file *gitlab.File) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		_, branch, fileID, err := resourceGitLabRepositoryFileParseId(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing repository file ID: %s", err)
		}
		// branch := rs.Primary.Attributes["branch"]
		if branch == "" {
			return fmt.Errorf("No branch set")
		}
		options := &gitlab.GetFileOptions{
			Ref: gitlab.String(branch),
		}
		repoName := rs.Primary.Attributes["project"]
		if repoName == "" {
			return fmt.Errorf("No project ID set")
		}

		gotFile, _, err := testGitlabClient.RepositoryFiles.GetFile(repoName, fileID, options)
		if err != nil {
			return fmt.Errorf("Cannot get file: %v", err)
		}

		if gotFile.FilePath == fileID {
			*file = *gotFile
			return nil
		}
		return fmt.Errorf("File does not exist")
	}
}

type testAccGitlabRepositoryFileAttributes struct {
	FilePath string
	Content  string
}

func testAccCheckGitlabRepositoryFileAttributes(got *gitlab.File, want *testAccGitlabRepositoryFileAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if got.FileName != want.FilePath {
			return fmt.Errorf("got name %q; want %q", got.FileName, want.FilePath)
		}

		if got.Content != want.Content {
			return fmt.Errorf("got content %q; want %q", got.Content, want.Content)
		}
		return nil
	}
}

func testAccCheckGitlabRepositoryFileDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_project" {
			continue
		}

		gotRepo, resp, err := testGitlabClient.Projects.GetProject(rs.Primary.ID, nil)
		if err == nil {
			if gotRepo != nil && fmt.Sprintf("%d", gotRepo.ID) == rs.Primary.ID {
				if gotRepo.MarkedForDeletionAt == nil {
					return fmt.Errorf("Repository still exists")
				}
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func testAccGitlabRepositoryFileConfig(projectID int) string {
	return fmt.Sprintf(`
resource "gitlab_repository_file" "this" {
  project = %d
  file_path = "meow.txt"
  branch = "main"
  content = "bWVvdyBtZW93IG1lb3c="
  author_email = "meow@catnip.com"
  author_name = "Meow Meowington"
  commit_message = "feature: add launch codes"
}
	`, projectID)
}

func testAccGitlabRepositoryFileStartBranchConfig(projectID int) string {
	return fmt.Sprintf(`
resource "gitlab_repository_file" "this" {
  project = %d
  file_path = "meow.txt"
  branch = "meow-branch"
  start_branch = "main"
  content = "bWVvdyBtZW93IG1lb3c="
  author_email = "meow@catnip.com"
  author_name = "Meow Meowington"
  commit_message = "feature: add launch codes"
}
	`, projectID)
}

func testAccGitlabRepositoryFileUpdateConfig(projectID int) string {
	return fmt.Sprintf(`
resource "gitlab_repository_file" "this" {
  project = %d
  file_path = "meow.txt"
  branch = "main"
  content = "bWVvdyBtZW93IG1lb3cgbWVvdyBtZW93Cg=="
  author_email = "meow@catnip.com"
  author_name = "Meow Meowington"
  commit_message = "feature: change launch codes"
}
	`, projectID)
}

func testAccGitlabRepositoryFileSameFileDifferentRepositoryConfig(firstProjectID, secondProjectID int) string {
	return fmt.Sprintf(`
resource "gitlab_repository_file" "foo_file" {
  project = %d
  file_path = "meow.txt"
  branch = "main"
  content = "bWVvdyBtZW93IG1lb3c="
  author_email = "meow@catnip.com"
  author_name = "Meow Meowington"
  commit_message = "feature: add launch codes"
}

resource "gitlab_repository_file" "bar_file" {
  project = %d
  file_path = "meow.txt"
  branch = "main"
  content = "bWVvdyBtZW93IG1lb3c="
  author_email = "meow@catnip.com"
  author_name = "Meow Meowington"
  commit_message = "feature: add launch codes"
}
	`, firstProjectID, secondProjectID)
}

func testAccGitlabRepositoryFileConcurrentResourcesConfig(projectID int) string {
	return fmt.Sprintf(`
resource "gitlab_repository_file" "this" {
  project = "%d"
  file_path = "file-${count.index}.txt"
  branch = "main"
  content = base64encode("content-${count.index}")
  commit_message = "Add file ${count.index}"

  count = 50
}
	`, projectID)
}

func testAccGitlabRepositoryFileConcurrentResourcesConfigUpdate(projectID int) string {
	return fmt.Sprintf(`
resource "gitlab_repository_file" "this" {
  project = "%d"
  file_path = "file-${count.index}.txt"
  branch = "main"
  content = base64encode("updated-content-${count.index}")
  commit_message = "Add file ${count.index}"

  count = 50
}
	`, projectID)
}
