//go:build acceptance
// +build acceptance

package provider

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabManagedLicense_basic(t *testing.T) {
	var managedLicense gitlab.ManagedLicense
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccCheckEE(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckManagedLicenseDestroy,
		Steps: []resource.TestStep{
			{
				// Create a managed license with an "allowed" state
				Config: testManagedLicenseConfig(rInt, "allowed"),
				Check:  testAccCheckGitlabManagedLicenseExists("gitlab_managed_license.fixme", &managedLicense),
			},
			{
				// Update the managed license to have a denied state
				Config: testManagedLicenseConfig(rInt, "denied"),
				Check:  testAccCheckGitlabManagedLicenseStatus("gitlab_managed_license.fixme", "denied", &managedLicense),
			},
			{
				ResourceName:      "gitlab_managed_license.fixme",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGitlabManagedLicense_deprecatedConfigValues(t *testing.T) {
	var managedLicense gitlab.ManagedLicense
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccRequiresLessThan(t, "15.0")
			testAccCheckEE(t)
		},
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckManagedLicenseDestroy,
		Steps: []resource.TestStep{
			{
				// Create a managed license with an "approved" state
				Config: testManagedLicenseConfig(rInt, "approved"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabManagedLicenseExists("gitlab_managed_license.fixme", &managedLicense),
				),
			},
			{
				// Update the managed license to have a "blacklisted" state
				Config: testManagedLicenseConfig(rInt, "blacklisted"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabManagedLicenseStatus("gitlab_managed_license.fixme", "blacklisted", &managedLicense),
				),
			},
			{
				ResourceName:      "gitlab_managed_license.fixme",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestCheckManagedLicenseStatusDiffFunc(t *testing.T) {
	if !checkDeprecatedValuesForDiff("", "approved", "allowed", nil) {
		t.Log("approved and allowed should be suppressed")
		t.Fail()
	}

	if !checkDeprecatedValuesForDiff("", "denied", "blacklisted", nil) {
		t.Log("denied and blacklisted should be suppressed")
		t.Fail()
	}

	if checkDeprecatedValuesForDiff("", "denied", "approved", nil) {
		t.Log("denied and approved should not be suppressed")
		t.Fail()
	}

	if checkDeprecatedValuesForDiff("", "allowed", "blacklisted", nil) {
		t.Log("allowed and blacklisted should not be suppressed")
		t.Fail()
	}
}

func testAccCheckGitlabManagedLicenseStatus(resource string, status string, license *gitlab.ManagedLicense) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("not Found: %s", resource)
		}

		project := rs.Primary.Attributes["project"]
		if project == "" {
			return fmt.Errorf("no project ID is set")
		}

		licenses, _, err := testGitlabClient.ManagedLicenses.ListManagedLicenses(project)
		if err != nil {
			return err
		}

		for _, gotLicense := range licenses {
			approvalStatus, err := stringToApprovalStatus(context.TODO(), testGitlabClient, status)
			if err != nil {
				return err
			}
			if gotLicense.ApprovalStatus == *approvalStatus {
				*license = *gotLicense
				return nil
			}
		}
		return fmt.Errorf("managed license does not exist")
	}
}

func testAccCheckManagedLicenseDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "gitlab_managed_licence" {
			continue
		}

		id, _ := strconv.Atoi(rs.Primary.ID)
		pid := rs.Primary.Attributes["project"]

		license, _, err := testGitlabClient.ManagedLicenses.GetManagedLicense(pid, id)
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

		licenseName := rs.Primary.Attributes["name"]
		project := rs.Primary.Attributes["project"]
		if project == "" {
			return fmt.Errorf("no project ID is set")
		}

		licenses, _, err := testGitlabClient.ManagedLicenses.ListManagedLicenses(project)
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
  name = "MIT"
  approval_status = "%s"
}
	`, rInt, status)
}
