//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataGitlabReleaseLinks_basic(t *testing.T) {

	project := testAccCreateProject(t)
	releases := testAccCreateReleases(t, project, 2)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				// get release link used Project ID
				Config: fmt.Sprintf(`
				data "gitlab_release_links" "this" {
					project = "%d"
					tag_name = "%s"
				}`, project.ID, releases[0].TagName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_release_links.this", "release_links.#", fmt.Sprintf("%v", 2)),
					resource.TestCheckResourceAttr("data.gitlab_release_links.this", "release_links.0.name", releases[0].Assets.Links[0].Name),
					resource.TestCheckResourceAttr("data.gitlab_release_links.this", "release_links.0.url", releases[0].Assets.Links[0].URL),
					resource.TestCheckResourceAttr("data.gitlab_release_links.this", "release_links.0.direct_asset_url", releases[0].Assets.Links[0].DirectAssetURL),
					resource.TestCheckResourceAttr("data.gitlab_release_links.this", "release_links.1.name", releases[0].Assets.Links[1].Name),
					resource.TestCheckResourceAttr("data.gitlab_release_links.this", "release_links.1.url", releases[0].Assets.Links[1].URL),
					resource.TestCheckResourceAttr("data.gitlab_release_links.this", "release_links.1.direct_asset_url", releases[0].Assets.Links[1].DirectAssetURL),
				),
			},
			{
				// get release link used full Project path
				Config: fmt.Sprintf(`
				data "gitlab_release_links" "this" {
					project = "%s"
					tag_name = "%s"
				}`, project.PathWithNamespace, releases[1].TagName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_release_links.this", "release_links.#", fmt.Sprintf("%v", 2)),
					resource.TestCheckResourceAttr("data.gitlab_release_links.this", "release_links.0.name", releases[1].Assets.Links[0].Name),
					resource.TestCheckResourceAttr("data.gitlab_release_links.this", "release_links.0.url", releases[1].Assets.Links[0].URL),
					resource.TestCheckResourceAttr("data.gitlab_release_links.this", "release_links.0.direct_asset_url", releases[1].Assets.Links[0].DirectAssetURL),
					resource.TestCheckResourceAttr("data.gitlab_release_links.this", "release_links.1.name", releases[1].Assets.Links[1].Name),
					resource.TestCheckResourceAttr("data.gitlab_release_links.this", "release_links.1.url", releases[1].Assets.Links[1].URL),
					resource.TestCheckResourceAttr("data.gitlab_release_links.this", "release_links.1.direct_asset_url", releases[1].Assets.Links[1].DirectAssetURL),
				),
			},
			{
				// get empty list
				Config: fmt.Sprintf(`
				data "gitlab_release_links" "this" {
					project = "%s"
					tag_name = "%s"
				}`, project.PathWithNamespace, "error_tag"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_release_links.this", "release_links.#", fmt.Sprintf("%v", 0)),
				),
			},
		},
	})
}
