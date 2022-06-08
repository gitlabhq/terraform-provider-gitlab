//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/xanzy/go-gitlab"
)

func TestAccDataSourceGitlabGroupVariables_basic(t *testing.T) {
	testGroup := testAccCreateGroups(t, 1)[0]
	testVariables := make([]*gitlab.GroupVariable, 0)
	for i := 0; i < 25; i++ {
		testVariables = append(testVariables, testAccCreateGroupVariable(t, testGroup.ID))
	}

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "gitlab_group_variables" "this" {
						group = %d
					}
				`, testGroup.ID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_group_variables.this", "variables.#", fmt.Sprintf("%d", len(testVariables))),
					resource.TestCheckResourceAttr("data.gitlab_group_variables.this", "variables.0.key", testVariables[0].Key),
					resource.TestCheckResourceAttr("data.gitlab_group_variables.this", "variables.0.value", testVariables[0].Value),
					resource.TestCheckResourceAttr("data.gitlab_group_variables.this", "variables.24.key", testVariables[24].Key),
					resource.TestCheckResourceAttr("data.gitlab_group_variables.this", "variables.24.value", testVariables[24].Value),
				),
			},
		},
	})
}
