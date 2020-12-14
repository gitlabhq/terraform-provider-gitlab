package gitlab

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/Fourcast/go-gitlab"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccGitlabDeployToken_basic(t *testing.T) {
	var deployToken gitlab.DeployToken
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabDeployTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabDeployTokenConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabDeployTokenExists("gitlab_deploy_token.foo", &deployToken),
					testAccCheckGitlabDeployTokenAttributes(&deployToken, &testAccCheckGitlabDeployTokenExpectedAttributes{
						Name:     fmt.Sprintf("deployToken-%d", rInt),
						Username: "my-username",
					}),
				),
			},
		},
	})
}

func testAccCheckGitlabDeployTokenExists(n string, deployToken *gitlab.DeployToken) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		deployTokenID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		conn := testAccProvider.Meta().(*gitlab.Client)

		projectName := rs.Primary.Attributes["project"]
		groupName := rs.Primary.Attributes["group"]

		var gotDeployTokens []*gitlab.DeployToken

		if projectName != "" {
			gotDeployTokens, _, err = conn.DeployTokens.ListProjectDeployTokens(projectName, nil)
		} else if groupName != "" {
			gotDeployTokens, _, err = conn.DeployTokens.ListGroupDeployTokens(groupName, nil)
		} else {
			return fmt.Errorf("No project or group ID is set")
		}

		if err != nil {
			return err
		}

		for _, token := range gotDeployTokens {
			if token.ID == deployTokenID {
				*deployToken = *token
				return nil
			}
		}

		return fmt.Errorf("Deploy Token doesn't exist")
	}
}

type testAccCheckGitlabDeployTokenExpectedAttributes struct {
	Name     string
	Username string
}

func testAccCheckGitlabDeployTokenAttributes(deployToken *gitlab.DeployToken, want *testAccCheckGitlabDeployTokenExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if deployToken.Name != want.Name {
			return fmt.Errorf("got name %q; want %q", deployToken.Name, want.Name)
		}

		if deployToken.Username != want.Username {
			return fmt.Errorf("got username %q; want %q", deployToken.Username, want.Username)
		}

		return nil
	}
}

func testAccCheckGitlabDeployTokenDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*gitlab.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_deploy_token" {
			continue
		}

		deployTokenID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		project := rs.Primary.Attributes["project"]
		group := rs.Primary.Attributes["group"]

		var gotDeployTokens []*gitlab.DeployToken
		var resp *gitlab.Response

		if project != "" {
			gotDeployTokens, resp, err = conn.DeployTokens.ListProjectDeployTokens(project, nil)
		} else if group != "" {
			gotDeployTokens, resp, err = conn.DeployTokens.ListGroupDeployTokens(group, nil)
		} else {
			return fmt.Errorf("somehow neither project nor group were set")
		}

		if err == nil {
			for _, token := range gotDeployTokens {
				if token.ID == deployTokenID {
					return fmt.Errorf("Deploy token still exists")
				}
			}
		}

		if resp.StatusCode != http.StatusNotFound {
			return err
		}
	}

	return nil
}

func testAccGitlabDeployTokenConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name        = "foo-%d"
  description = "Terraform acceptance test"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_deploy_token" "foo" {
  project  = "${gitlab_project.foo.id}"
  name     = "deployToken-%d"
  username = "my-username"

  expires_at = "2021-03-14T07:20:50Z"

  scopes = [
	"read_registry",
	"read_repository",
  ]
}
  `, rInt, rInt)
}
