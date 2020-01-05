package tls

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	expectedPublic = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDPLaq43D9C596ko9yQipWUf2Fb
RhFs18D3wBDBqXLIoP7W3rm5S292/JiNPa+mX76IYFF416zTBGG9J5w4d4VFrROn
8IuMWqHgdXsCUf2szN7EnJcVBsBzTxxWqz4DjX315vbm/PFOLlKzC0Ngs4h1iDiC
D9Hk2MajZuFnJiqj1QIDAQAB
-----END PUBLIC KEY-----`
	expectedPublicSSH            = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDPLaq43D9C596ko9yQipWUf2FbRhFs18D3wBDBqXLIoP7W3rm5S292/JiNPa+mX76IYFF416zTBGG9J5w4d4VFrROn8IuMWqHgdXsCUf2szN7EnJcVBsBzTxxWqz4DjX315vbm/PFOLlKzC0Ngs4h1iDiCD9Hk2MajZuFnJiqj1Q==`
	expectedPublicFingerprintMD5 = `62:c2:c6:7a:d0:27:72:e7:0d:bc:4e:97:42:0e:9e:e6`
)

func TestAccPublicKey_dataSource(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourcePublicKeyConfig, testPrivateKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.tls_public_key.test", "public_key_pem", strings.TrimSpace(expectedPublic)+"\n"),
					resource.TestCheckResourceAttr("data.tls_public_key.test", "public_key_openssh", strings.TrimSpace(expectedPublicSSH)+"\n"),
					resource.TestCheckResourceAttr("data.tls_public_key.test", "public_key_fingerprint_md5", strings.TrimSpace(expectedPublicFingerprintMD5)),
				),
			},
			{
				Config: `
					resource "tls_private_key" "test" {
						algorithm = "RSA"
					}
					data "tls_public_key" "test" {
						private_key_pem = "${tls_private_key.test.private_key_pem}"
					}
				`,
				Check: resource.TestCheckResourceAttrPair(
					"data.tls_public_key.test", "public_key_pem",
					"tls_private_key.test", "public_key_pem"),
			},
			{
				Config: `
					resource "tls_private_key" "key" {
						algorithm   = "ECDSA"
						ecdsa_curve = "P384"
					}
					data "tls_public_key" "pub" {
						private_key_pem = "${tls_private_key.key.private_key_pem}"
					}
				`,
				Check: resource.TestCheckResourceAttrPair(
					"data.tls_public_key.pub", "public_key_pem",
					"tls_private_key.key", "public_key_pem"),
			},
			{
				Config:      fmt.Sprintf(testAccDataSourcePublicKeyConfig, "corrupt"),
				ExpectError: regexp.MustCompile("failed to decode PEM block containing private key of type \"unknown\""),
			},
		},
	})
}

const testAccDataSourcePublicKeyConfig = `
data "tls_public_key" "test" {
  private_key_pem = <<EOF
	%s
	EOF
}
`
