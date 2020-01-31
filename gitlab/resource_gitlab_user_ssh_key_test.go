package gitlab

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabUserSSHKey_basic(t *testing.T) {
	var key *gitlab.SSHKey
	rn := "gitlab_user_ssh_key.test"
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	title := fmt.Sprintf("tf-acc-test-%s", randString)
	keyRe := regexp.MustCompile("^ecdsa-sha2-nistp384 ")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabUserSSHKeyDestroy,
		Steps: []resource.TestStep{
			// create Key
			{
				Config: testAccGitlabUserSSHKeyConfig(title),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabUserSSHKeyExists(rn, key),
					resource.TestCheckResourceAttr(rn, "title", title),
					resource.TestMatchResourceAttr(rn, "key", keyRe),
				),
			},
		},
	})
}

func testAccCheckGitlabUserSSHKeyExists(n string, key *gitlab.SSHKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		KeyID := rs.Primary.ID
		if KeyID == "" {
			return fmt.Errorf("No Key ID is set")
		}

		conn := testAccProvider.Meta().(*gitlab.Client)

		id, _ := strconv.Atoi(KeyID)

		gotKey, _, err := conn.Users.GetSSHKey(id)
		if err != nil {
			return err
		}
		key = gotKey
		return nil
	}
}

func testAccCheckGitlabUserSSHKeyDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*gitlab.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_user_ssh_key" {
			continue
		}

		id, _ := strconv.Atoi(rs.Primary.ID)

		userKey, resp, err := conn.Users.GetSSHKey(id)
		if err == nil {
			if userKey != nil && fmt.Sprintf("%d", userKey.ID) == rs.Primary.ID {
				return fmt.Errorf("User SSH Key still exists")
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func testAccGitlabUserSSHKeyConfig(title string) string {
	return fmt.Sprintf(`
	resource "gitlab_user_ssh_key" "test"{
		title	=	"%s"
		key   = "${tls_private_key.test.public_key_openssh}"

	}

	resource "tls_private_key" "test" {
		algorithm   = "ECDSA"
		ecdsa_curve = "P384"
	  }
	`, title)
}
