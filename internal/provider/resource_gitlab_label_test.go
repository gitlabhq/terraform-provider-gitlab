//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabLabel_basic(t *testing.T) {
	var label gitlab.Label
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabLabelDestroy,
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
			// Verify Import
			{
				ResourceName:      "gitlab_label.fixme",
				ImportStateIdFunc: getLabelImportID("gitlab_label.fixme"),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func getLabelImportID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("Not Found: %s", n)
		}

		labelID := rs.Primary.ID
		if labelID == "" {
			return "", fmt.Errorf("No label key ID is set")
		}
		projectID := rs.Primary.Attributes["project"]
		if projectID == "" {
			return "", fmt.Errorf("No project ID is set")
		}

		return fmt.Sprintf("%s:%s", projectID, labelID), nil
	}
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

		labels, _, err := testGitlabClient.Labels.ListLabels(repoName, &gitlab.ListLabelsOptions{ListOptions: gitlab.ListOptions{PerPage: 1000}})
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
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_label" {
			continue
		}

		project := rs.Primary.Attributes["project"]
		labelName := rs.Primary.ID

		_, _, err := testGitlabClient.Labels.GetLabel(project, labelName)
		if err == nil {
			return fmt.Errorf("Label %s in project %s still exists", labelName, project)
		}
		if !is404(err) {
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
