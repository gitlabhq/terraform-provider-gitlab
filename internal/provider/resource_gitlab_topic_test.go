package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabTopic_basic(t *testing.T) {
	var topic gitlab.Topic
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabTopicDestroy,
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
			// Verify import
			{
				ResourceName:      "gitlab_topic.foo",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"soft_destroy",
				},
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
			// Verify import
			{
				ResourceName:      "gitlab_topic.foo",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"soft_destroy",
				},
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
			// Verify import
			{
				ResourceName:      "gitlab_topic.foo",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"soft_destroy",
				},
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
				ImportStateVerifyIgnore: []string{
					"soft_destroy",
				},
			},
		},
	})
}

func TestAccGitlabTopic_softDestroy(t *testing.T) {
	var topic gitlab.Topic
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabTopicSoftDestroy,
		Steps: []resource.TestStep{
			// Create a topic with soft_destroy enabled
			{
				Config: testAccGitlabTopicSoftDestroyConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabTopicExists("gitlab_topic.foo", &topic),
				),
			},
		},
	})
}

func testAccCheckGitlabTopicExists(n string, assign *gitlab.Topic) resource.TestCheckFunc {
	return func(s *terraform.State) (err error) {

		defer func() {
			if err != nil {
				err = fmt.Errorf("checking for gitlab topic existence failed: %w", err)
			}
		}()

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not Found: %s", n)
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		topic, _, err := testGitlabClient.Topics.GetTopic(id)
		*assign = *topic

		return err
	}
}

type testAccGitlabTopicExpectedAttributes struct {
	Name        string
	Description string
	SoftDestroy bool
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

func testAccCheckGitlabTopicDestroy(s *terraform.State) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("destroying gitlab topic failed: %w", err)
		}
	}()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_topic" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		topic, _, err := testGitlabClient.Topics.GetTopic(id)
		if err == nil {
			if topic != nil && fmt.Sprintf("%d", topic.ID) == rs.Primary.ID {
				return fmt.Errorf("topic %s still exists", rs.Primary.ID)
			}
		}
		if !is404(err) {
			return err
		}
		return nil
	}
	return nil
}

func testAccCheckGitlabTopicSoftDestroy(s *terraform.State) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("destroying gitlab topic failed: %w", err)
		}
	}()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_topic" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		topic, _, err := testGitlabClient.Topics.GetTopic(id)
		if err == nil {
			if topic != nil && fmt.Sprintf("%d", topic.ID) == rs.Primary.ID {
				if topic.Description != "" {
					return fmt.Errorf("topic still has a description")
				}
				return nil
			}
		}
		if !is404(err) {
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

func testAccGitlabTopicSoftDestroyConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_topic" "foo" {
  name             = "foo-soft-destroy-%d"
  description      = "Terraform acceptance tests"

  soft_destroy     = true
}`, rInt)
}
