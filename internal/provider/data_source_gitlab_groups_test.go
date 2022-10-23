//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGitlabGroups_basic(t *testing.T) {
	prefixFoo := "acctest-group-foo"
	groupsFoo := testAccCreateGroupsWithPrefix(t, 2, prefixFoo)

	prefixLotsOf := "acctest-group-lotsof"
	testAccCreateGroupsWithPrefix(t, 42, prefixLotsOf)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGitlabGroupsConfigSearchSort(prefixFoo),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_groups.foos", "groups.#", "2"),
					resource.TestCheckResourceAttr("data.gitlab_groups.foos", "groups.0.group_id", fmt.Sprint(groupsFoo[0].ID)),
					resource.TestCheckResourceAttr("data.gitlab_groups.foos", "groups.0.full_path", groupsFoo[0].FullPath),
					resource.TestCheckResourceAttr("data.gitlab_groups.foos", "groups.0.name", groupsFoo[0].Name),
					resource.TestCheckResourceAttr("data.gitlab_groups.foos", "groups.0.full_name", groupsFoo[0].FullName),
					resource.TestCheckResourceAttr("data.gitlab_groups.foos", "groups.0.web_url", groupsFoo[0].WebURL),
					resource.TestCheckResourceAttr("data.gitlab_groups.foos", "groups.0.path", groupsFoo[0].Path),
					resource.TestCheckResourceAttr("data.gitlab_groups.foos", "groups.0.description", groupsFoo[0].Description),
					resource.TestCheckResourceAttr("data.gitlab_groups.foos", "groups.0.lfs_enabled", strconv.FormatBool(groupsFoo[0].LFSEnabled)),
					resource.TestCheckResourceAttr("data.gitlab_groups.foos", "groups.0.request_access_enabled", strconv.FormatBool(groupsFoo[0].RequestAccessEnabled)),
					resource.TestCheckResourceAttr("data.gitlab_groups.foos", "groups.0.visibility_level", fmt.Sprint(groupsFoo[0].Visibility)),
					resource.TestCheckResourceAttr("data.gitlab_groups.foos", "groups.0.parent_id", fmt.Sprint(groupsFoo[0].ParentID)),
					resource.TestCheckResourceAttr("data.gitlab_groups.foos", "groups.0.runners_token", groupsFoo[0].RunnersToken),
					resource.TestCheckResourceAttr("data.gitlab_groups.foos", "groups.0.default_branch_protection", fmt.Sprint(groupsFoo[0].DefaultBranchProtection)),
					resource.TestCheckResourceAttr("data.gitlab_groups.foos", "groups.0.prevent_forking_outside_group", strconv.FormatBool(groupsFoo[0].PreventForkingOutsideGroup)),
				),
			},
			{
				Config: testAccDataSourceGitlabLotsOfGroupsSearch(prefixLotsOf),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_groups.lotsof", "groups.#", "42"),
				),
			},
		},
	})
}

func testAccDataSourceGitlabGroupsConfigSearchSort(prefix string) string {
	return fmt.Sprintf(`
data "gitlab_groups" "foos" {
  sort = "asc"
  search = "%s"
  order_by = "id"
}
	`, prefix)
}

func testAccDataSourceGitlabLotsOfGroupsSearch(prefix string) string {
	return fmt.Sprintf(`
data "gitlab_groups" "lotsof" {
	search = "%s"
}
	`, prefix)
}
