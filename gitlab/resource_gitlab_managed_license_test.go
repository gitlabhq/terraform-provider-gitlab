package gitlab

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
	"strconv"
	"testing"
)

func TestAccGitlabManagedLicense_basic(t *testing.T) {
	var managedLicense gitlab.ManagedLicense
	rInt := acctest.RandInt()
	client := testAccNewClient(t)
	testAccCheckEE(t, client)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() {},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckManagedLicenseDestroy,
		Steps: []resource.TestStep{
			{
				// Create a managed license with an "approved" state
				Config: testManagedLicenseConfig(rInt, "approved"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabManagedLicenseExists("gitlab_managed_license.fix_me", &managedLicense),
				),
			},
			{
				// Update the managed license to have a blacklisted state
				Config: testManagedLicenseConfig(rInt, "blacklisted"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabManagedLicenseStatus("gitlab_managed_license.fix_me", &managedLicense),
				),
			},
			{
				ResourceName:      "gitlab_managed_license.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGitlabManagedLicenseStatus(n string, license *gitlab.ManagedLicense) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		licenseName := rs.Primary.ID
		project := rs.Primary.Attributes["project"]
		if project == "" {
			return fmt.Errorf("no project ID is set")
		}

		conn := testAccProvider.Meta().(*gitlab.Client)

		licenses, _, err := conn.ManagedLicenses.ListManagedLicenses(project)
		if err != nil {
			return err
		}

		for _, gotLicense := range licenses {
			if gotLicense.Name == licenseName {
				*license = *gotLicense
				return nil
			}
		}
		return fmt.Errorf("managed license does not exist")
	}
}

func testAccCheckManagedLicenseDestroy(state *terraform.State) error {
	conn := testAccProvider.Meta().(*gitlab.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "gitlab_managed_licence" {
			continue
		}

		id, _ := strconv.Atoi(rs.Primary.ID)
		pid := rs.Primary.Attributes["project"]

		license, _, err := conn.ManagedLicenses.GetManagedLicense(pid, id)
		if err == nil {
			if license != nil && license.ID == id {
				return fmt.Errorf("license still exists")
			}
		}

		if !is404(err) {
			return err
		}
		return nil
	}
	return nil
}

func testAccCheckGitlabManagedLicenseExists(n string, license *gitlab.ManagedLicense) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		licenseName := rs.Primary.ID
		project := rs.Primary.Attributes["project"]
		if project == "" {
			return fmt.Errorf("no project ID is set")
		}

		conn := testAccProvider.Meta().(*gitlab.Client)

		licenses, _, err := conn.ManagedLicenses.ListManagedLicenses(project)
		if err != nil {
			return err
		}

		for _, gotLicense := range licenses {
			if gotLicense.Name == licenseName {
				*license = *gotLicense
				return nil
			}
		}
		return fmt.Errorf("managed license does not exist")
	}
}

func testManagedLicenseConfig(rInt int, status string) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_managed_license" "fixme" {
  project = "${gitlab_project.foo.id}"
  name = "FIXME-%d"
  approval_status = "%s"
}
	`, rInt, rInt, status)
}
