package provider

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

func testGitlabDeployKeyStateDataV0() map[string]interface{} {
	return map[string]interface{}{
		"id":      "5",
		"project": "1",
	}
}

func testGitlabDeployKeyStateDataV1() map[string]interface{} {
	v0 := testGitlabDeployKeyStateDataV0()
	return map[string]interface{}{
		"id":      fmt.Sprintf("%s:%s", v0["project"], v0["id"]),
		"project": "1",
	}
}

func TestGitlabDeployKeyStateUpgradeV0(t *testing.T) {
	expected := testGitlabDeployKeyStateDataV1()
	actual, err := resourceGitlabDeployKeyStateUpgradeV0(context.Background(), testGitlabDeployKeyStateDataV0(), nil)
	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\n\nexpected:\n\n%#v\n\ngot:\n\n%#v\n\n", expected, actual)
	}
}

func TestAccGitlabDeployKey_basic(t *testing.T) {
	var deployKey gitlab.ProjectDeployKey
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabDeployKeyDestroy,
		Steps: []resource.TestStep{
			// Create a project and deployKey with default options
			{
				Config: testAccGitlabDeployKeyConfig(rInt, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabDeployKeyExists("gitlab_deploy_key.foo", &deployKey),
					testAccCheckGitlabDeployKeyAttributes(&deployKey, &testAccGitlabDeployKeyExpectedAttributes{
						Title: fmt.Sprintf("deployKey-%d", rInt),
						Key:   "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCj13ozEBZ0s4el4k6mYqoyIKKKMh9hHY0sAYqSPXs2zGuVFZss1P8TPuwmdXVjHR7TiRXwC49zDrkyWJgiufggYJ1VilOohcMOODwZEJz+E5q4GCfHuh90UEh0nl8B2R0Uoy0LPeg93uZzy0hlHApsxRf/XZJz/1ytkZvCtxdllxfImCVxJReMeRVEqFCTCvy3YuJn0bce7ulcTFRvtgWOpQsr6GDK8YkcCCv2eZthVlrEwy6DEpAKTRiRLGgUj4dPO0MmO4cE2qD4ualY01PhNORJ8Q++I+EtkGt/VALkecwFuBkl18/gy+yxNJHpKc/8WVVinDeFrd/HhiY9yU0d richardc@tamborine.example.1",
					}),
				),
			},
			// Update the project deployKey to toggle all the values to their inverse
			{
				Config: testAccGitlabDeployKeyUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabDeployKeyExists("gitlab_deploy_key.foo", &deployKey),
					testAccCheckGitlabDeployKeyAttributes(&deployKey, &testAccGitlabDeployKeyExpectedAttributes{
						Title: fmt.Sprintf("modifiedDeployKey-%d", rInt),
						Key:   "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC6pSke2kb7YBjo65xDKegbOQsAtnMupRcFxXji7L1iXivGwORq0qpC2xzbhez5jk1WgPckEaNv2/Bz0uEW6oSIXw1KT1VN2WzEUfQCbpNyZPtn4iV3nyl6VQW/Nd1SrxiFJtH1H4vu+eCo4McMXTjuBBD06fiJNrHaSw734LjQgqtXWJuVym9qS5MqraZB7wDwTQwSM6kslL7KTgmo3ONsTLdb2zZhv6CS+dcFKinQo7/ttTmeMuXGbPOVuNfT/bePVIN1MF1TislHa2L2dZdGeoynNJT4fVPjA2Xl6eHWh4ySbvnfPznASsjBhP0n/QKprYJ/5fQShdBYBcuQiIMd richardc@tamborine.example.2",
					}),
				),
			},
			// Update the project deployKey to toggle the options back
			{
				Config: testAccGitlabDeployKeyConfig(rInt, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabDeployKeyExists("gitlab_deploy_key.foo", &deployKey),
					testAccCheckGitlabDeployKeyAttributes(&deployKey, &testAccGitlabDeployKeyExpectedAttributes{
						Title: fmt.Sprintf("deployKey-%d", rInt),
						Key:   "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCj13ozEBZ0s4el4k6mYqoyIKKKMh9hHY0sAYqSPXs2zGuVFZss1P8TPuwmdXVjHR7TiRXwC49zDrkyWJgiufggYJ1VilOohcMOODwZEJz+E5q4GCfHuh90UEh0nl8B2R0Uoy0LPeg93uZzy0hlHApsxRf/XZJz/1ytkZvCtxdllxfImCVxJReMeRVEqFCTCvy3YuJn0bce7ulcTFRvtgWOpQsr6GDK8YkcCCv2eZthVlrEwy6DEpAKTRiRLGgUj4dPO0MmO4cE2qD4ualY01PhNORJ8Q++I+EtkGt/VALkecwFuBkl18/gy+yxNJHpKc/8WVVinDeFrd/HhiY9yU0d richardc@tamborine.example.1",
					}),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_deploy_key.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGitlabDeployKey_suppressfunc(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabDeployKeyDestroy,
		Steps: []resource.TestStep{
			// Create a project and deployKey with newline as suffix
			{
				Config: testAccGitlabDeployKeyConfig(rInt, ""),
			},
		},
	})
}

func testAccCheckGitlabDeployKeyExists(n string, deployKey *gitlab.ProjectDeployKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		project, rawDeployKeyID, err := parseTwoPartID(rs.Primary.ID)
		if err != nil {
			return err
		}

		deployKeyID, err := strconv.Atoi(rawDeployKeyID)
		if err != nil {
			return err
		}

		gotDeployKey, _, err := testGitlabClient.DeployKeys.GetDeployKey(project, deployKeyID)
		if err != nil {
			return err
		}
		*deployKey = *gotDeployKey
		return nil
	}
}

type testAccGitlabDeployKeyExpectedAttributes struct {
	Title   string
	Key     string
	CanPush bool
}

func testAccCheckGitlabDeployKeyAttributes(deployKey *gitlab.ProjectDeployKey, want *testAccGitlabDeployKeyExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if deployKey.Title != want.Title {
			return fmt.Errorf("got title %q; want %q", deployKey.Title, want.Title)
		}

		if deployKey.Key != want.Key {
			return fmt.Errorf("got key %q; want %q", deployKey.Key, want.Key)
		}

		if deployKey.CanPush != want.CanPush {
			return fmt.Errorf("got can_push %t; want %t", deployKey.CanPush, want.CanPush)
		}

		return nil
	}
}

func testAccCheckGitlabDeployKeyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_deploy_key" {
			continue
		}

		project, rawDeployKeyID, err := parseTwoPartID(rs.Primary.ID)
		if err != nil {
			return err
		}

		deployKeyID, err := strconv.Atoi(rawDeployKeyID)
		if err != nil {
			return err
		}

		gotDeployKey, _, err := testGitlabClient.DeployKeys.GetDeployKey(project, deployKeyID)
		if err == nil {
			if gotDeployKey != nil {
				return fmt.Errorf("Deploy key still exists")
			}
		}
		if !is404(err) {
			return err
		}
		return nil
	}
	return nil
}

func testAccGitlabDeployKeyConfig(rInt int, suffix string) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_deploy_key" "foo" {
  project = "${gitlab_project.foo.id}"
  title = "deployKey-%d"
  key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCj13ozEBZ0s4el4k6mYqoyIKKKMh9hHY0sAYqSPXs2zGuVFZss1P8TPuwmdXVjHR7TiRXwC49zDrkyWJgiufggYJ1VilOohcMOODwZEJz+E5q4GCfHuh90UEh0nl8B2R0Uoy0LPeg93uZzy0hlHApsxRf/XZJz/1ytkZvCtxdllxfImCVxJReMeRVEqFCTCvy3YuJn0bce7ulcTFRvtgWOpQsr6GDK8YkcCCv2eZthVlrEwy6DEpAKTRiRLGgUj4dPO0MmO4cE2qD4ualY01PhNORJ8Q++I+EtkGt/VALkecwFuBkl18/gy+yxNJHpKc/8WVVinDeFrd/HhiY9yU0d richardc@tamborine.example.1%s"
}
  `, rInt, rInt, suffix)
}

func testAccGitlabDeployKeyUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_deploy_key" "foo" {
  project = "${gitlab_project.foo.id}"
  title = "modifiedDeployKey-%d"
  key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC6pSke2kb7YBjo65xDKegbOQsAtnMupRcFxXji7L1iXivGwORq0qpC2xzbhez5jk1WgPckEaNv2/Bz0uEW6oSIXw1KT1VN2WzEUfQCbpNyZPtn4iV3nyl6VQW/Nd1SrxiFJtH1H4vu+eCo4McMXTjuBBD06fiJNrHaSw734LjQgqtXWJuVym9qS5MqraZB7wDwTQwSM6kslL7KTgmo3ONsTLdb2zZhv6CS+dcFKinQo7/ttTmeMuXGbPOVuNfT/bePVIN1MF1TislHa2L2dZdGeoynNJT4fVPjA2Xl6eHWh4ySbvnfPznASsjBhP0n/QKprYJ/5fQShdBYBcuQiIMd richardc@tamborine.example.2"
}
  `, rInt, rInt)
}
