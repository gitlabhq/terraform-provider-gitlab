//go:build acceptance
// +build acceptance

package provider

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

func TestAccGitlabProjectTag_basic(t *testing.T) {
	var tag gitlab.Tag
	var tag2 gitlab.Tag
	rInt, rInt2, rInt3 := acctest.RandInt(), acctest.RandInt(), acctest.RandInt()
	project := testAccCreateProject(t)
	branches := testAccCreateBranches(t, project, 1)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabProjectTagConfig(rInt, rInt2, project.PathWithNamespace, branches[0].Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectTagExists("foo", &tag, rInt),
					testAccCheckGitlabProjectTagExists("foo2", &tag2, rInt2),
					testAccCheckGitlabProjectTagAttributes("foo", &tag, &testAccGitlabProjectTagExpectedAttributes{
						Name:    fmt.Sprintf("tag-%d", rInt),
						Message: "",
						Ref:     "main",
					}),
					testAccCheckGitlabProjectTagAttributes("foo2", &tag2, &testAccGitlabProjectTagExpectedAttributes{
						Name:    fmt.Sprintf("tag-%d", rInt2),
						Message: fmt.Sprintf("tag-%d", rInt2),
						Ref:     branches[0].Name,
					}),
				),
			},
			// Test ImportState
			{
				ResourceName:            "gitlab_project_tag.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ref"},
			},
			// update properties in resource
			{
				Config: testAccGitlabProjectTagConfig(rInt, rInt3, project.PathWithNamespace, branches[0].Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectTagExists("foo2", &tag2, rInt3),
					testAccCheckGitlabProjectTagAttributes("foo2", &tag2, &testAccGitlabProjectTagExpectedAttributes{
						Name:    fmt.Sprintf("tag-%d", rInt3),
						Message: fmt.Sprintf("tag-%d", rInt3),
						Ref:     branches[0].Name,
					}),
				),
			},
		},
	})
}

func testAccCheckGitlabProjectTagDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_project_tag" {
			continue
		}
		name := rs.Primary.Attributes["name"]
		project := rs.Primary.Attributes["project"]
		_, _, err := testGitlabClient.Tags.GetTag(project, name)
		if err != nil {
			if is404(err) {
				return nil
			}
			return err
		}
		return errors.New("Tag still exists")
	}
	return nil
}

func testAccCheckGitlabProjectTagAttributes(n string, tag *gitlab.Tag, want *testAccGitlabProjectTagExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[fmt.Sprintf("gitlab_project_tag.%s", n)]
		ref := rs.Primary.Attributes["ref"]
		if ref != want.Ref {
			return fmt.Errorf("Got ref %s; want %s", ref, want.Ref)
		}
		if tag.Name != want.Name {
			return fmt.Errorf("Got tag name %s; want %s", tag.Name, want.Name)
		}
		if tag.Message != want.Message {
			return fmt.Errorf("Got message %s; want %s", tag.Message, want.Message)
		}
		if tag.Commit == nil {
			return errors.New("The tag commit is nil but expected to be populated")
		}
		if tag.Commit.ID == "" {
			return errors.New("The commit has an empty ID")
		}
		if tag.Release != nil {
			if tag.Release.TagName != want.Name {
				return fmt.Errorf("Got release note tag name %s; want %s", tag.Release.TagName, want.Name)
			}
		}
		if tag.Protected != false {
			return fmt.Errorf("Got tag field protected %v; want false", tag.Protected)
		}
		return nil
	}
}

func testAccCheckGitlabProjectTagExists(n string, tag *gitlab.Tag, rInt int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[fmt.Sprintf("gitlab_project_tag.%s", n)]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}
		project, name, err := parseTwoPartID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error in splitting project and tag")
		}
		gotTag, _, err := testGitlabClient.Tags.GetTag(project, name)
		if err != nil {
			return err
		}
		*tag = *gotTag
		return err
	}
}

func testAccGitlabProjectTagConfig(rInt int, rInt2 int, project string, branch string) string {
	return fmt.Sprintf(`
    resource "gitlab_project_tag" "foo" {
        name    = "tag-%[1]d"
        ref     = "main"
        project = "%[3]s"
    }
    resource "gitlab_project_tag" "foo2" {
        name    = "tag-%[2]d"
        ref     = "%[4]s"
        project = "%[3]s"
        message = "tag-%[2]d"
    }
  `, rInt, rInt2, project, branch)
}

type testAccGitlabProjectTagExpectedAttributes struct {
	Name    string
	Message string
	Ref     string
}
