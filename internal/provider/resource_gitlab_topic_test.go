//go:build acceptance
// +build acceptance

package provider

import (
	"context"
	"fmt"
	"os"
	"regexp"
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
					"avatar", "avatar_hash", "soft_destroy",
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
					resource.TestCheckResourceAttrSet("gitlab_topic.foo", "avatar_url"),
					resource.TestCheckResourceAttr("gitlab_topic.foo", "avatar_hash", "8d29d9c393facb9d86314eb347a03fde503f2c0422bf55af7df086deb126107e"),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_topic.foo",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"avatar", "avatar_hash", "soft_destroy",
				},
			},
			// Update back to the default topics avatar
			{
				Config: testAccGitlabTopicFullConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabTopicExists("gitlab_topic.foo", &topic),
					testAccCheckGitlabTopicAttributes(&topic, &testAccGitlabTopicExpectedAttributes{
						Name:        fmt.Sprintf("foo-full-%d", rInt),
						Description: "Terraform acceptance tests",
					}),
					resource.TestCheckResourceAttrSet("gitlab_topic.foo", "avatar_url"),
					resource.TestCheckResourceAttr("gitlab_topic.foo", "avatar_hash", "8d29d9c393facb9d86314eb347a03fde503f2c0422bf55af7df086deb126107e"),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_topic.foo",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"avatar", "avatar_hash", "soft_destroy",
				},
			},
			// Update the topics avatar
			{
				Config: testAccGitlabTopicFullUpdatedAvatarConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabTopicExists("gitlab_topic.foo", &topic),
					testAccCheckGitlabTopicAttributes(&topic, &testAccGitlabTopicExpectedAttributes{
						Name:        fmt.Sprintf("foo-full-%d", rInt),
						Description: "Terraform acceptance tests",
					}),
					resource.TestCheckResourceAttrSet("gitlab_topic.foo", "avatar_url"),
					resource.TestCheckResourceAttr("gitlab_topic.foo", "avatar_hash", "a58bd926fd3baabd41c56e810f62ade8705d18a4e280fb35764edb4b778444db"),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_topic.foo",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"avatar", "avatar_hash", "soft_destroy",
				},
			},
			// Update back to the default topics avatar
			{
				Config: testAccGitlabTopicFullConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabTopicExists("gitlab_topic.foo", &topic),
					testAccCheckGitlabTopicAttributes(&topic, &testAccGitlabTopicExpectedAttributes{
						Name:        fmt.Sprintf("foo-full-%d", rInt),
						Description: "Terraform acceptance tests",
					}),
					resource.TestCheckResourceAttrSet("gitlab_topic.foo", "avatar_url"),
					resource.TestCheckResourceAttr("gitlab_topic.foo", "avatar_hash", "8d29d9c393facb9d86314eb347a03fde503f2c0422bf55af7df086deb126107e"),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_topic.foo",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"avatar", "avatar_hash", "soft_destroy",
				},
			},
			// Update the avatar image, but keep the filename to test the `CustomizeDiff` function
			{
				Config: testAccGitlabTopicFullConfig(rInt),
				PreConfig: func() {
					// overwrite the avatar image file
					if err := copyFile("testdata/gitlab_topic/avatar.png", "testdata/gitlab_topic/avatar.png.bak"); err != nil {
						t.Fatalf("failed to backup the avatar image file: %v", err)
					}
					if err := copyFile("testdata/gitlab_topic/avatar-update.png", "testdata/gitlab_topic/avatar.png"); err != nil {
						t.Fatalf("failed to overwrite the avatar image file: %v", err)
					}
					t.Cleanup(func() {
						if err := os.Rename("testdata/gitlab_topic/avatar.png.bak", "testdata/gitlab_topic/avatar.png"); err != nil {
							t.Fatalf("failed to restore the avatar image file: %v", err)
						}
					})
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabTopicExists("gitlab_topic.foo", &topic),
					testAccCheckGitlabTopicAttributes(&topic, &testAccGitlabTopicExpectedAttributes{
						Name:        fmt.Sprintf("foo-full-%d", rInt),
						Description: "Terraform acceptance tests",
					}),
					resource.TestCheckResourceAttrSet("gitlab_topic.foo", "avatar_url"),
					resource.TestCheckResourceAttr("gitlab_topic.foo", "avatar_hash", "a58bd926fd3baabd41c56e810f62ade8705d18a4e280fb35764edb4b778444db"),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_topic.foo",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"avatar", "avatar_hash", "soft_destroy",
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
					resource.TestCheckResourceAttr("gitlab_topic.foo", "avatar_url", ""),
					resource.TestCheckResourceAttr("gitlab_topic.foo", "avatar_hash", ""),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_topic.foo",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"avatar", "avatar_hash", "soft_destroy",
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
					"avatar", "avatar_hash", "soft_destroy",
				},
			},
		},
	})
}

func TestAccGitlabTopic_withoutAvatarHash(t *testing.T) {
	var topic gitlab.Topic
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabTopicDestroy,
		Steps: []resource.TestStep{
			// Create a topic with avatar, but without giving a hash
			{
				Config: testAccGitlabTopicAvatarWithoutHashConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabTopicExists("gitlab_topic.foo", &topic),
					resource.TestCheckResourceAttrSet("gitlab_topic.foo", "avatar_url"),
				),
				ExpectNonEmptyPlan: true,
			},
			// Update the avatar image, but keep the filename to test the `CustomizeDiff` function
			{
				Config:             testAccGitlabTopicAvatarWithoutHashConfig(rInt),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccGitlabTopic_softDestroy(t *testing.T) {
	var topic gitlab.Topic
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
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

func TestAccGitlabTopic_titleSupport(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabTopicDestroy,
		Steps: []resource.TestStep{
			{
				SkipFunc: isGitLabVersionAtLeast(context.Background(), testGitlabClient, "15.0"),
				Config: fmt.Sprintf(`
					resource "gitlab_topic" "this" {
						name = "foo-%d"
						title = "Foo-%d"
					}
				`, rInt, rInt),
				ExpectError: regexp.MustCompile(`title is not supported by your version of GitLab. At least GitLab 15.0 is required`),
			},
			{
				SkipFunc: isGitLabVersionLessThan(context.Background(), testGitlabClient, "15.0"),
				Config: fmt.Sprintf(`
					resource "gitlab_topic" "this" {
						name = "foo-%d"
					}
				`, rInt),
				ExpectError: regexp.MustCompile(`title is a required attribute for GitLab 15.0 and newer. Please specify it in the configuration.`),
			},
			{
				SkipFunc: isGitLabVersionLessThan(context.Background(), testGitlabClient, "15.0"),
				Config: fmt.Sprintf(`
					resource "gitlab_topic" "this" {
						name = "foo-%d"
						title = "Foo-%d"
					}
				`, rInt, rInt),
				Check: resource.TestCheckResourceAttr("gitlab_topic.this", "title", fmt.Sprintf("Foo-%d", rInt)),
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
  name = "foo-req-%d"
  title = "Foo Req %d"
}`, rInt, rInt)
}

func testAccGitlabTopicFullConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_topic" "foo" {
  name        = "foo-full-%d"
  title       = "Foo Full %d"
  description = "Terraform acceptance tests"
  avatar      = "${path.module}/testdata/gitlab_topic/avatar.png"
  avatar_hash = filesha256("${path.module}/testdata/gitlab_topic/avatar.png")
}`, rInt, rInt)
}

func testAccGitlabTopicFullUpdatedAvatarConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_topic" "foo" {
  name        = "foo-full-%d"
  title       = "Foo Full %d"
  description = "Terraform acceptance tests"
  avatar 	  = "${path.module}/testdata/gitlab_topic/avatar-update.png"
  avatar_hash = filesha256("${path.module}/testdata/gitlab_topic/avatar-update.png")
}`, rInt, rInt)
}

func testAccGitlabTopicAvatarWithoutHashConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_topic" "foo" {
  name   = "foo-%d"
  title  = "Foo %d"
  avatar = "${path.module}/testdata/gitlab_topic/avatar.png"
}`, rInt, rInt)
}

func testAccGitlabTopicSoftDestroyConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_topic" "foo" {
  name        = "foo-soft-destroy-%d"
  title       = "Foo Soft Destroy %d"
  description = "Terraform acceptance tests"

  soft_destroy = true
}`, rInt, rInt)
}
