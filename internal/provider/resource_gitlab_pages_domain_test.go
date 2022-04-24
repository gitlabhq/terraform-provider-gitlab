package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabPagesDomain_basic(t *testing.T) {
	var pagesDomain gitlab.PagesDomain
	rInt := acctest.RandInt()
	testAccCheck(t)
	project := testAccCreateProject(t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabPagesDestroy,
		Steps: []resource.TestStep{
			// Create a pages domain with all options
			{
				Config: testAccGitlabPagesDomainCreate(rInt, project.PathWithNamespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabPagesDomainExists("gitlab_pages_domain.this", &pagesDomain),
					resource.TestCheckResourceAttrSet("gitlab_pages_domain.this", "auto_ssl_enabled"),
				),
			},
			// Verify import
			{
				ResourceName:            "gitlab_pages_domain.this",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token"},
			},
			// Update the pages domain to toggle all the values to their inverse
			{
				Config: testAccGitlabPagesDomainUpdate(rInt, project.PathWithNamespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabPagesDomainExists("gitlab_pages_domain.this", &pagesDomain),
				),
			},
		},
	})
}

func testAccCheckGitlabPagesDomainExists(n string, pagesDomain *gitlab.PagesDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		projectID := rs.Primary.Attributes["project"]
		if projectID == "" {
			return fmt.Errorf("No project ID is set")
		}

		domain := rs.Primary.Attributes["domain"]
		if domain == "" {
			return fmt.Errorf("No domain is set")
		}

		gotPagesDomain, _, err := testGitlabClient.PagesDomains.GetPagesDomain(projectID, domain)
		if err != nil {
			return err
		}
		*pagesDomain = *gotPagesDomain
		return nil
	}
}

func testAccCheckGitlabPagesDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_pages_domain" {
			continue
		}

		projectID := rs.Primary.Attributes["project"]
		if projectID == "" {
			return fmt.Errorf("No project ID is set")
		}

		domain := rs.Primary.Attributes["domain"]
		if domain == "" {
			return fmt.Errorf("No domain is set")
		}

		gotPagesDomain, err := testGitlabClient.PagesDomains.DeletePagesDomain(projectID, domain)
		if err == nil {
			if gotPagesDomain != nil {
				return fmt.Errorf("Pages domain %s still exists after deletion", domain)
			}
		}
		if !is404(err) {
			return err
		}
		return nil
	}
	return nil
}

func testAccGitlabPagesDomainCreate(rInt int, project string) string {
	return fmt.Sprintf(`
resource "gitlab_pages_domain" "this" {
  project                     = "%[2]s"
  domain                      = "page-%[1]d.example.com"
  auto_ssl_enabled            = false
}
	`, rInt, project)
}

func testAccGitlabPagesDomainUpdate(rInt int, project string) string {
	return fmt.Sprintf(`
resource "gitlab_pages_domain" "this" {
  project                     = "%[2]s"
  domain                      = "page-%[1]d.example.com"
  auto_ssl_enabled            = true
}
	`, rInt+1, project)
}
