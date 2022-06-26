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
)

func TestAccGitlabReleaseLink_basic(t *testing.T) {

	rInt1, rInt2 := acctest.RandInt(), acctest.RandInt()
	project := testAccCreateProject(t)
	releases := testAccCreateReleases(t, project, 1)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabReleaseLinkDestroy,
		Steps: []resource.TestStep{
			{
				// create Release link with required values only
				Config: fmt.Sprintf(`
				resource "gitlab_release_link" "this" {
					project  = "%s"
					tag_name = "%s"
					name     = "test-%d"
					url      = "https://test/%d"
				}`, project.PathWithNamespace, releases[0].TagName, rInt1, rInt1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gitlab_release_link.this", "link_id"),
					resource.TestCheckResourceAttrSet("gitlab_release_link.this", "direct_asset_url"),
					resource.TestCheckResourceAttrSet("gitlab_release_link.this", "external"),
				),
			},
			{
				// verify import
				ResourceName:      "gitlab_release_link.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// update some Release link attributes
				Config: fmt.Sprintf(`
				resource "gitlab_release_link" "this" {
					project   = "%d"
					tag_name  = "%s"
					name      = "test-%d"
					url       = "https://test/%d"
					filepath  = "/test/%d"
					link_type = "runbook"
				}`, project.ID, releases[0].TagName, rInt2, rInt2, rInt2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gitlab_release_link.this", "link_id"),
					resource.TestCheckResourceAttrSet("gitlab_release_link.this", "direct_asset_url"),
					resource.TestCheckResourceAttrSet("gitlab_release_link.this", "external"),
				),
			},
			{
				// verify import
				ResourceName:      "gitlab_release_link.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGitlabReleaseLinkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_release_link" {
			continue
		}
		project, tagName, linkID, err := resourceGitLabReleaseLinkParseId(rs.Primary.ID)
		if err != nil {
			return err
		}

		releaseLink, _, err := testGitlabClient.ReleaseLinks.GetReleaseLink(project, tagName, linkID)
		if err == nil && releaseLink != nil {
			return errors.New("Release link still exists")
		}
		if !is404(err) {
			return err
		}
		return nil
	}
	return nil
}
