//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGitlabReleaseLink_basic(t *testing.T) {

	project := testAccCreateProject(t)
	releases := testAccCreateReleases(t, project, 2)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				// get release link used Project ID
				Config: fmt.Sprintf(`
				data "gitlab_release_link" "this" {
					project  = "%d"
					tag_name = "%s"
					link_id  = "%d"
				}`, project.ID, releases[0].TagName, releases[0].Assets.Links[0].ID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_release_link.this", "name", releases[0].Assets.Links[0].Name),
					resource.TestCheckResourceAttr("data.gitlab_release_link.this", "url", releases[0].Assets.Links[0].URL),
					resource.TestCheckResourceAttr("data.gitlab_release_link.this", "direct_asset_url", releases[0].Assets.Links[0].DirectAssetURL),
				),
			},
			{
				// get release link used full Project path
				Config: fmt.Sprintf(`
				data "gitlab_release_link" "this" {
					project  = "%s"
					tag_name = "%s"
					link_id  = "%d"
				}`, project.PathWithNamespace, releases[1].TagName, releases[1].Assets.Links[0].ID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_release_link.this", "name", releases[1].Assets.Links[0].Name),
					resource.TestCheckResourceAttr("data.gitlab_release_link.this", "url", releases[1].Assets.Links[0].URL),
					resource.TestCheckResourceAttr("data.gitlab_release_link.this", "direct_asset_url", releases[1].Assets.Links[0].DirectAssetURL),
				),
			},
		},
	})
}
