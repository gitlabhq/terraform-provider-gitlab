package gitlab

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabTopic(t *testing.T) {
	var topic gitlab.Topic
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabTopicDestroy,
		Steps: []resource.TestStep{
			// Create a topic with default options
			{
				Config: testAccGitlabTopicRequiredConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabTopicExists("gitlab_topic.foo", &topic),
					testAccCheckGitlabTopicAttributes(&topic, &testAccGitlabTopicExpectedAttributes{
						Name: fmt.Sprintf("foo-req-%d", rInt),
					}),
				),
			},
			// Update the topics values
			{
				Config: testAccGitlabTopicFullConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabTopicExists("gitlab_topic.foo", &topic),
					testAccCheckGitlabTopicAttributes(&topic, &testAccGitlabTopicExpectedAttributes{
						Name:        fmt.Sprintf("foo-full-%d", rInt),
						Description: "Terraform acceptance tests",
					}),
				),
			},
			// Update the topics values back to their initial state
			{
				Config: testAccGitlabTopicRequiredConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabTopicExists("gitlab_topic.foo", &topic),
					testAccCheckGitlabTopicAttributes(&topic, &testAccGitlabTopicExpectedAttributes{
						Name: fmt.Sprintf("foo-req-%d", rInt),
					}),
				),
			},
			// Updating the topic to have a description before it is deleted
			{
				Config: testAccGitlabTopicFullConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabTopicExists("gitlab_topic.foo", &topic),
					testAccCheckGitlabTopicAttributes(&topic, &testAccGitlabTopicExpectedAttributes{
						Name:        fmt.Sprintf("foo-full-%d", rInt),
						Description: "Terraform acceptance tests",
					}),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_topic.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGitlabTopicExists(n string, assign *gitlab.Topic) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not Found: %s", n)
		}

		conn := testAccProvider.Meta().(*gitlab.Client)

		topicID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		topic, _, err := conn.Topics.GetTopic(topicID)
		*assign = *topic

		return err
	}
}

type testAccGitlabTopicExpectedAttributes struct {
	Name        string
	Description string
}

func testAccCheckGitlabTopicAttributes(topic *gitlab.Topic, want *testAccGitlabTopicExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if topic.Name != want.Name {
			return fmt.Errorf("got name %q; want %q", topic.Name, want.Name)
		}

		if topic.Description != want.Description {
			return fmt.Errorf("got description %q; want %q", topic.Description, want.Description)
		}

		return nil
	}
}

func testAccCheckGitlabTopicDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*gitlab.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_topic" {
			continue
		}

		topicID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		topic, resp, err := conn.Topics.GetTopic(topicID)
		if err == nil {
			if topic != nil && fmt.Sprintf("%d", topic.ID) == rs.Primary.ID {

				if topic.Description != "" {
					return fmt.Errorf("topic still has a description")
				}
				// TODO: Return error as soon as deleting a topic is supported
				return nil
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func testAccGitlabTopicRequiredConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_topic" "foo" {
  name             = "foo-req-%d"
}`, rInt)
}

func testAccGitlabTopicFullConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_topic" "foo" {
  name             = "foo-full-%d"
  description      = "Terraform acceptance tests"
}`, rInt)
}
