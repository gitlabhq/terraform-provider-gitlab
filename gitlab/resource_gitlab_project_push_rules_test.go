package gitlab

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	gitlab "github.com/xanzy/go-gitlab"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccGitlabProjectPushRules_basic(t *testing.T) {
	var pushRules gitlab.ProjectPushRules
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabProjectPushRulesDestroy,
		Steps: []resource.TestStep{
			// Create project and push rules with basic options
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabProjectPushRulesConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectPushRulesExists("gitlab_project_push_rules.foo", &pushRules),
					testAccCheckGitlabProjectPushRulesAttributes(&pushRules, &testAccGitlabProjectPushRulesExpectedAttributes{
						CommitMessageRegex: "^(foo|bar).*",
						BranchNameRegex:    "^(foo|bar).*",
						AuthorEmailRegex:   "^(foo|bar).*",
						FileNameRegex:      "^(foo|bar).*",
						DenyDeleteTag:      true,
						MemberCheck:        true,
						PreventSecrets:     true,
						MaxFileSize:        10,
					}),
				),
			},
			// Update the project push rules
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabProjectPushRulesUpdate(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectPushRulesExists("gitlab_project_push_rules.foo", &pushRules),
					testAccCheckGitlabProjectPushRulesAttributes(&pushRules, &testAccGitlabProjectPushRulesExpectedAttributes{
						CommitMessageRegex: "^(fu|baz).*",
						BranchNameRegex:    "^(fu|baz).*",
						AuthorEmailRegex:   "^(fu|baz).*",
						FileNameRegex:      "^(fu|baz).*",
						DenyDeleteTag:      false,
						MemberCheck:        false,
						PreventSecrets:     false,
						MaxFileSize:        42,
					}),
				),
			},
			// Update the project push rules to original config
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabProjectPushRulesConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectPushRulesExists("gitlab_project_push_rules.foo", &pushRules),
					testAccCheckGitlabProjectPushRulesAttributes(&pushRules, &testAccGitlabProjectPushRulesExpectedAttributes{
						CommitMessageRegex: "^(foo|bar).*",
						BranchNameRegex:    "^(foo|bar).*",
						AuthorEmailRegex:   "^(foo|bar).*",
						FileNameRegex:      "^(foo|bar).*",
						DenyDeleteTag:      true,
						MemberCheck:        true,
						PreventSecrets:     true,
						MaxFileSize:        10,
					}),
				),
			},
		},
	})
}

func TestAccGitlabProjectPushRules_import(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabProjectPushRulesDestroy,
		Steps: []resource.TestStep{
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabProjectPushRulesConfig(rInt),
			},
			{
				SkipFunc:          isRunningInCE,
				ResourceName:      "gitlab_project_push_rules.foo",
				ImportStateId:     fmt.Sprintf("foo-%d", rInt),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGitlabProjectPushRulesExists(n string, pushRules *gitlab.ProjectPushRules) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		repoName := rs.Primary.Attributes["project"]
		if repoName == "" {
			return fmt.Errorf("No project ID is set")
		}
		conn := testAccProvider.Meta().(*gitlab.Client)
		gotPushRules, _, err := conn.Projects.GetProjectPushRules(repoName)
		if err != nil {
			return err
		}
		*pushRules = *gotPushRules
		return nil
	}
}

type testAccGitlabProjectPushRulesExpectedAttributes struct {
	CommitMessageRegex string
	BranchNameRegex    string
	AuthorEmailRegex   string
	FileNameRegex      string
	DenyDeleteTag      bool
	MemberCheck        bool
	PreventSecrets     bool
	MaxFileSize        int
}

func testAccCheckGitlabProjectPushRulesAttributes(pushRules *gitlab.ProjectPushRules, want *testAccGitlabProjectPushRulesExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if pushRules.CommitMessageRegex != want.CommitMessageRegex {
			return fmt.Errorf("got commit_message_regex %s; want %s", pushRules.CommitMessageRegex, want.CommitMessageRegex)
		}
		if pushRules.BranchNameRegex != want.BranchNameRegex {
			return fmt.Errorf("got branch_name_regex %s; want %s", pushRules.BranchNameRegex, want.BranchNameRegex)
		}
		if pushRules.AuthorEmailRegex != want.AuthorEmailRegex {
			return fmt.Errorf("got author_email_regex %s; want %s", pushRules.AuthorEmailRegex, want.AuthorEmailRegex)
		}
		if pushRules.FileNameRegex != want.FileNameRegex {
			return fmt.Errorf("got file_name_regex %s; want %s", pushRules.FileNameRegex, want.FileNameRegex)
		}
		if pushRules.DenyDeleteTag != want.DenyDeleteTag {
			return fmt.Errorf("got deny_delete_tag %t; want %t", pushRules.DenyDeleteTag, want.DenyDeleteTag)
		}
		if pushRules.MemberCheck != want.MemberCheck {
			return fmt.Errorf("got member_check %t; want %t", pushRules.MemberCheck, want.MemberCheck)
		}
		if pushRules.PreventSecrets != want.PreventSecrets {
			return fmt.Errorf("got prevent_secrets %t; want %t", pushRules.PreventSecrets, want.PreventSecrets)
		}
		if pushRules.MaxFileSize != want.MaxFileSize {
			return fmt.Errorf("got max_file_size %d; want %d", pushRules.MaxFileSize, want.MaxFileSize)
		}
		return nil
	}
}

func testAccCheckGitlabProjectPushRulesDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*gitlab.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_project" {
			continue
		}

		gotRepo, resp, err := conn.Projects.GetProject(rs.Primary.ID, nil)
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

func testAccGitlabProjectPushRulesConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance test - Push Rule"
  visibility_level = "public"
}

resource "gitlab_project_push_rules" "foo" {
  project = "${gitlab_project.foo.id}"
  commit_message_regex = "^(foo|bar).*"
  branch_name_regex = "^(foo|bar).*"
  author_email_regex = "^(foo|bar).*"
  file_name_regex = "^(foo|bar).*"
  deny_delete_tag = true
  member_check = true
  prevent_secrets = true
  max_file_size = 10
}
`, rInt)
}

func testAccGitlabProjectPushRulesUpdate(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance test - Push Rule"
  visibility_level = "public"
}

resource "gitlab_project_push_rules" "foo" {
  project = "${gitlab_project.foo.id}"
  commit_message_regex = "^(fu|baz).*"
  branch_name_regex = "^(fu|baz).*"
  author_email_regex = "^(fu|baz).*"
  file_name_regex = "^(fu|baz).*"
  deny_delete_tag = false
  member_check = false
  prevent_secrets = false
  max_file_size = 42
}
`, rInt)
}
