//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabDeployKeyEnable_basic(t *testing.T) {
	testProjectParent := testAccCreateProject(t)
	testProjectKeyShared := testAccCreateProject(t)

	key := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDblguSWgpqiXIjHPSas4+N3Dten7MTLJMlGQXxGpaqN9nGPdNmuRB2YXyjT/nrryoY/qrtuVkPnis5WVo8N/s3hAnJbeJPUS2WKEGjpBlL34AQ+ANnlmGY8L6zr82Hp2Ommb7XGGtlq5D3yLCgTfcXLjC51tgcdwHsdH1U+RisgLwaTSrP/HF4G7IAr5ATsyYjtCwQRQ8ijdf5A34+XN6h8J6TLXKab5eZDuH38s9LxJuS7MRxx/P2UTOsqfjtrZWoQgE5adEGvnDxKyruex9PzNbCNVahzsma7tdikDbzxlHLIZ1aht6rKuai3iyLgcZfGIYtkq4xvg/bnNXxSsGf worker@kg.getwifi.com"

	canPushDeployKeyOptions := gitlab.AddDeployKeyOptions{
		Title:   gitlab.String("main"),
		Key:     gitlab.String(key),
		CanPush: gitlab.Bool(true),
	}

	parentProjectDeployKey := testAccCreateDeployKey(t, testProjectParent.ID, &canPushDeployKeyOptions)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabDeployKeyEnableDestroy,
		Steps: []resource.TestStep{
			// Enable a deployKey on project with default options
			{
				Config: testAccGitlabDeployKeyEnableConfig(testProjectKeyShared.ID, parentProjectDeployKey.ID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gitlab_deploy_key_enable.foo", "key"),
					resource.TestCheckResourceAttrSet("gitlab_deploy_key_enable.foo", "title"),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_deploy_key_enable.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Define canPush to true
			{
				Config: testAccGitlabDeployKeyEnableConfigCanPush(testProjectKeyShared.ID, parentProjectDeployKey.ID, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gitlab_deploy_key_enable.foo", "key"),
					resource.TestCheckResourceAttrSet("gitlab_deploy_key_enable.foo", "title"),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_deploy_key_enable.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Define canPush to false
			{
				Config: testAccGitlabDeployKeyEnableConfigCanPush(testProjectKeyShared.ID, parentProjectDeployKey.ID, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gitlab_deploy_key_enable.foo", "key"),
					resource.TestCheckResourceAttrSet("gitlab_deploy_key_enable.foo", "title"),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_deploy_key_enable.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Get back to default options
			{
				Config: testAccGitlabDeployKeyEnableConfig(testProjectKeyShared.ID, parentProjectDeployKey.ID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gitlab_deploy_key_enable.foo", "key"),
					resource.TestCheckResourceAttrSet("gitlab_deploy_key_enable.foo", "title"),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_deploy_key_enable.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGitlabDeployKeyEnableDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		project, deployKeyID, err := resourceGitLabDeployKeyEnableParseId(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("unable to parse resource ID into project and deployKeyID: %w", err)
		}

		gotDeployKey, _, err := testGitlabClient.DeployKeys.GetDeployKey(project, deployKeyID)
		if err == nil {
			if gotDeployKey != nil {
				return fmt.Errorf("Deploy key still exists: %d", deployKeyID)
			}
		}
		if !is404(err) {
			return err
		}
	}
	return nil
}

func testAccGitlabDeployKeyEnableConfig(shareProjectId int, keyId int) string {
	return fmt.Sprintf(`
resource "gitlab_deploy_key_enable" "foo" {
  project = %[1]d
  key_id  = %[2]d
}
  `, shareProjectId, keyId)
}

func testAccGitlabDeployKeyEnableConfigCanPush(shareProjectId int, keyId int, canPush bool) string {
	return fmt.Sprintf(`
resource "gitlab_deploy_key_enable" "foo" {
  project  = %[1]d
  key_id   = %[2]d
  can_push = %[3]t
}
  `, shareProjectId, keyId, canPush)
}
