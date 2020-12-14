package gitlab

import (
	"fmt"
	"testing"

	"github.com/Fourcast/go-gitlab"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccGitlabLabel_basic(t *testing.T) {
	var label gitlab.Label
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabLabelDestroy,
		Steps: []resource.TestStep{
			// Create a project and label with default options
			{
				Config: testAccGitlabLabelConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabLabelExists("gitlab_label.fixme", &label),
					testAccCheckGitlabLabelAttributes(&label, &testAccGitlabLabelExpectedAttributes{
						Name:        fmt.Sprintf("FIXME-%d", rInt),
						Color:       "#ffcc00",
						Description: "fix this test",
					}),
				),
			},
			// Update the label to change the parameters
			{
				Config: testAccGitlabLabelUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabLabelExists("gitlab_label.fixme", &label),
					testAccCheckGitlabLabelAttributes(&label, &testAccGitlabLabelExpectedAttributes{
						Name:        fmt.Sprintf("FIXME-%d", rInt),
						Color:       "#ff0000",
						Description: "red label",
					}),
				),
			},
			// Update the label to get back to initial settings
			{
				Config: testAccGitlabLabelConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabLabelExists("gitlab_label.fixme", &label),
					testAccCheckGitlabLabelAttributes(&label, &testAccGitlabLabelExpectedAttributes{
						Name:        fmt.Sprintf("FIXME-%d", rInt),
						Color:       "#ffcc00",
						Description: "fix this test",
					}),
				),
			},
			// Create a project and lots of labels with default options
			{
				Config: testAccGitlabLabelLargeConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabLabelExists("gitlab_label.fixme.20", &label),
					testAccCheckGitlabLabelExists("gitlab_label.fixme.30", &label),
					testAccCheckGitlabLabelExists("gitlab_label.fixme.40", &label),
					testAccCheckGitlabLabelExists("gitlab_label.fixme.10", &label),
					testAccCheckGitlabLabelAttributes(&label, &testAccGitlabLabelExpectedAttributes{
						Name:        "FIXME11",
						Color:       "#ffcc00",
						Description: "fix this test",
					}),
				),
			},
			// Update the lots of labels to change the parameters
			{
				Config: testAccGitlabLabelUpdateLargeConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabLabelExists("gitlab_label.fixme.20", &label),
					testAccCheckGitlabLabelExists("gitlab_label.fixme.30", &label),
					testAccCheckGitlabLabelExists("gitlab_label.fixme.40", &label),
					testAccCheckGitlabLabelExists("gitlab_label.fixme.10", &label),
					testAccCheckGitlabLabelAttributes(&label, &testAccGitlabLabelExpectedAttributes{
						Name:        "FIXME11",
						Color:       "#ff0000",
						Description: "red label",
					}),
				),
			},
		},
	})
}

func testAccCheckGitlabLabelExists(n string, label *gitlab.Label) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		labelName := rs.Primary.ID
		repoName := rs.Primary.Attributes["project"]
		if repoName == "" {
			return fmt.Errorf("No project ID is set")
		}
		conn := testAccProvider.Meta().(*gitlab.Client)

		labels, _, err := conn.Labels.ListLabels(repoName, &gitlab.ListLabelsOptions{ListOptions: gitlab.ListOptions{PerPage: 1000}})
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

type testAccGitlabLabelExpectedAttributes struct {
	Name        string
	Color       string
	Description string
}

func testAccCheckGitlabLabelAttributes(label *gitlab.Label, want *testAccGitlabLabelExpectedAttributes) resource.TestCheckFunc {
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

func testAccCheckGitlabLabelDestroy(s *terraform.State) error {
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

func testAccGitlabLabelConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_label" "fixme" {
  project = "${gitlab_project.foo.id}"
  name = "FIXME-%d"
  color = "#ffcc00"
  description = "fix this test"
}
	`, rInt, rInt)
}

func testAccGitlabLabelUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_label" "fixme" {
  project = "${gitlab_project.foo.id}"
  name = "FIXME-%d"
  color = "#ff0000"
  description = "red label"
}
	`, rInt, rInt)
}

func testAccGitlabLabelLargeConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_label" "fixme" {
  project = "${gitlab_project.foo.id}"
  name = format("FIXME%%02d", count.index+1)
  count = 100
  color = "#ffcc00"
  description = "fix this test"
}
	`, rInt)
}

func testAccGitlabLabelUpdateLargeConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_label" "fixme" {
  project = "${gitlab_project.foo.id}"
  name = format("FIXME%%02d", count.index+1)
  count = 99
  color = "#ff0000"
  description = "red label"
}
	`, rInt)
}
