package provider

/*
import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabTopic_basic(t *testing.T) {
	var topic gitlab.Topic
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabTopicDestroy,
		Steps: []resource.TestStep{
			// Create a topic with default options
			{
				Config: testAccGitlabTopicCreateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabTopicExists("gitlab_group_label.fixme", &topic),
					testAccCheckGitlabTopicAttributes(&topic, &testAccGitlabTopicExpectedAttributes{
						Name:        fmt.Sprintf("FIXME-%d", rInt),
						Color:       "#ffcc00",
						Description: "fix this test",
					}),
				),
			},
			// Update the topics values
			{
				Config: testAccGitlabTopicUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabTopicExists("gitlab_group_label.fixme", &topic),
					testAccCheckGitlabTopicAttributes(&topic, &testAccGitlabTopicExpectedAttributes{
						Name:        fmt.Sprintf("FIXME-%d", rInt),
						Color:       "#ff0000",
						Description: "red label",
					}),
				),
			},
			// Update the topic back to its initial state
			{
				Config: testAccGitlabTopicCreateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabTopicExists("gitlab_group_label.fixme", &topic),
					testAccCheckGitlabTopicAttributes(&topic, &testAccGitlabTopicExpectedAttributes{
						Name:        fmt.Sprintf("FIXME-%d", rInt),
						Color:       "#ff0000",
						Description: "red label",
					}),
				),
			},
		},
	})
}

// lintignore: AT002 // TODO: Resolve this tfproviderlint issue
func TestAccGitlabTopic_import(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "gitlab_group_label.fixme"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabTopicConfig(rInt),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: getTopicImportID(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func getTopicImportID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("Not Found: %s", n)
		}

		labelID := rs.Primary.ID
		if labelID == "" {
			return "", fmt.Errorf("No deploy key ID is set")
		}
		groupID := rs.Primary.Attributes["group"]
		if groupID == "" {
			return "", fmt.Errorf("No group ID is set")
		}

		return fmt.Sprintf("%s:%s", groupID, labelID), nil
	}
}

func testAccCheckGitlabTopicExists(n string, label *gitlab.Topic) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		labelName := rs.Primary.ID
		groupName := rs.Primary.Attributes["group"]
		if groupName == "" {
			return fmt.Errorf("No group ID is set")
		}
		conn := testAccProvider.Meta().(*gitlab.Client)

		labels, _, err := conn.Topics.ListTopics(groupName, &gitlab.ListTopicsOptions{PerPage: 1000})
		if err != nil {
			return err
		}
		for _, gotLabel := range labels {
			if gotLabel.Name == labelName {
				*label = *gotLabel
				return nil
			}
		}
		return fmt.Errorf("Label does not exist")
	}
}

type testAccGitlabTopicExpectedAttributes struct {
	Name        string
	Color       string
	Description string
}

func testAccCheckGitlabTopicAttributes(label *gitlab.Topic, want *testAccGitlabTopicExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if label.Name != want.Name {
			return fmt.Errorf("got name %q; want %q", label.Name, want.Name)
		}

		if label.Description != want.Description {
			return fmt.Errorf("got description %q; want %q", label.Description, want.Description)
		}

		if label.Color != want.Color {
			return fmt.Errorf("got color %q; want %q", label.Color, want.Color)
		}

		return nil
	}
}

func testAccCheckGitlabTopicDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*gitlab.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_group" {
			continue
		}

		group, resp, err := conn.Groups.GetGroup(rs.Primary.ID, nil)
		if err == nil {
			if group != nil && fmt.Sprintf("%d", group.ID) == rs.Primary.ID {
				if group.MarkedForDeletionOn == nil {
					return fmt.Errorf("Group still exists")
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

func testAccGitlabTopicCreateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_topic" "foo" {
  name             = "foo-%d"
  path             = "foo-%d"
  description      = "Terraform acceptance tests"
  visibility_level = "public"
}
	`, rInt)
}

func testAccGitlabTopicUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_topic" "foo" {
  name             = "foo-%d"
  path             = "foo-%d"
  description      = "Terraform acceptance tests"
  visibility_level = "public"
}
	`, rInt)
}
*/
