//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
)

/* Keys are generated via:

$ cat gen-key-script
Key-Type: eddsa
Key-Curve: ed25519
Key-Usage: sign
Subkey-Type: eddsa
Subkey-Curve: ed25519
Subkey-Usage: sign
Name-Real: Terraform
Name-Email: terraform@gitlab.com
Expire-Date: 0
Passphrase: ''
$ gpg --batch --gen-key gen-key-script
*/

var testEd25519GPGPubKey string = `-----BEGIN PGP PUBLIC KEY BLOCK-----

mDMEYtzLgRYJKwYBBAHaRw8BAQdArtnkLgGRJyUI0QiBidLC6tZ++xAQ0ofQ0sxR
+lZbIsO0IFRlcnJhZm9ybSA8dGVycmFmb3JtQGdpdGxhYi5jb20+iJAEExYIADgW
IQTP8v16wSfKRS9E/hWPJsCZR6QoUQUCYtzLgQIbAwULCQgHAgYVCgkICwIEFgID
AQIeAQIXgAAKCRCPJsCZR6QoUeiaAQDEqci/HGBKHgGy3G2wJmvdywL3SEwgsSfu
j8betD1cKQD/Q1FRHOT6nsshhgHLxuuM6Nc7EaTjQyaXxFvEJ8YOXwS4MwRi3MuB
FgkrBgEEAdpHDwEBB0BneBgfymn64IZIxsaJQNpRAehycYE9VyNMLbfPzHFUm4jv
BBgWCAAgFiEEz/L9esEnykUvRP4VjybAmUekKFEFAmLcy4ECGwIAgQkQjybAmUek
KFF2IAQZFggAHRYhBPKpgJFOiwH7W8NFVckaYlcKYj8ZBQJi3MuBAAoJEMkaYlcK
Yj8ZGcoA/RUR48JzgKB0QtUsVHbxH5HxGljJzX+fojk9amOfly+aAQCPB6//nRA2
RLwa75kqGYGYS6cpY/ZSaqi19342XjiJD+NtAQDq+XWq7qVf1DZ6VPBpsZJfQ/ws
1piPsvmKI/2koOa1wQD/ebxTge120L4AVDMpGhDjpqt+B4qN4SiEKynbSJcWiAg=
=q0aF
-----END PGP PUBLIC KEY BLOCK-----`
var updatedEd25519GPGPubKey string = `-----BEGIN PGP PUBLIC KEY BLOCK-----

mDMEYtzQ6RYJKwYBBAHaRw8BAQdASvO1H8QsiJSN+qmEmBwtgeNi61lXzCmCfUF3
/5e1qgC0IFRlcnJhZm9ybSA8dGVycmFmb3JtQGdpdGxhYi5jb20+iJAEExYIADgW
IQTqH8ino1RNUdUDG7H+AGMyH5oPPwUCYtzQ6QIbAwULCQgHAgYVCgkICwIEFgID
AQIeAQIXgAAKCRD+AGMyH5oPP+pbAP0YgRmwvpipmaBuK4/gDWnE3gO2lR4x+vzL
ciPmYxshZQD+KAK69/eUwsdbW+uapbtSUNDF4l6PWyKES6qYSUlpswi4MwRi3NDp
FgkrBgEEAdpHDwEBB0A2d2HUVSFIudz5xegw3HjmH5tDoodbyoURpD7NZhslpojv
BBgWCAAgFiEE6h/Ip6NUTVHVAxux/gBjMh+aDz8FAmLc0OkCGwIAgQkQ/gBjMh+a
Dz92IAQZFggAHRYhBEA7LZIX/tW2nvyiCfF9dLzyMph4BQJi3NDpAAoJEPF9dLzy
Mph4TYEBAPSRfGNlQESyYfUmqV795iFkIgCM5nVzqHHfw0mcH50EAP4qZK9Iobvs
/yG4eD5jp4nGkRfDlXA+ZsIjmChtIRT2DVy/AQDRFeywZIsvuje6DF1AP7qeqbs+
dZcWp5qRyFu5zodW2gEAgw7OxfYmlQbz7wwKr3IYnY9MyVp+JwSxyAxZ4X+odAM=
=XSo3
-----END PGP PUBLIC KEY BLOCK-----`
var testEd25519GPGPubKeyForCurrentUser string = `-----BEGIN PGP PUBLIC KEY BLOCK-----

mDMEYtzl7hYJKwYBBAHaRw8BAQdAF8RE2x2wE3w/QJAbB8uVv+9HqYQZJGyNZeLt
eKyyh2i0IFRlcnJhZm9ybSA8dGVycmFmb3JtQGdpdGxhYi5jb20+iJAEExYIADgW
IQShzcoy44uIRzOHNnGXiWgn/RdiBQUCYtzl7gIbAwULCQgHAgYVCgkICwIEFgID
AQIeAQIXgAAKCRCXiWgn/RdiBUDpAQDpcNH5rxVlJS4ATpkFPxbyMuZBdr6eaw2G
3o8ZPykfbAEAt62eRBhyep+rjPyAHW8YyX4+e+SdCZw/KH0NDbZCGAW4MwRi3OXu
FgkrBgEEAdpHDwEBB0DE69f1NA/05px8W0ZUkIl0kheMosQ3+UFbt/Jba0DQiojv
BBgWCAAgFiEEoc3KMuOLiEczhzZxl4loJ/0XYgUFAmLc5e4CGwIAgQkQl4loJ/0X
YgV2IAQZFggAHRYhBCB2JV+gsGWeld0PcwehRlcCKxpcBQJi3OXuAAoJEAehRlcC
KxpcSmcBALG/ZILOrVHratA+cZchDQEtOM2rymXdw6AgnllhMj48AQCjWmESbs+C
GcNNwjo6hAu59BiCvU5+W2of9fpSxBM5CJQ4AQDjyZwDcf6kMu5+bYY9aNcz9skX
BEMqhn2i2EtNNfH6eAEA1m21vnE8pseBCYHtl9/XJGkB5JH0gUqqRrCf5CAqsgk=
=/bV6
-----END PGP PUBLIC KEY BLOCK-----`
var testEd25519GPGPubKeyWithTrailingNewline string = `-----BEGIN PGP PUBLIC KEY BLOCK-----

mDMEYtzSYhYJKwYBBAHaRw8BAQdAaSePzlrqUWT/hBRqO/oIdfBPh5m57cPwpmjL
dgUJGcW0IFRlcnJhZm9ybSA8dGVycmFmb3JtQGdpdGxhYi5jb20+iJAEExYIADgW
IQTFohelsWBqQdLTMuS4BhhIqbpAJQUCYtzSYgIbAwULCQgHAgYVCgkICwIEFgID
AQIeAQIXgAAKCRC4BhhIqbpAJX8hAP9XCLVIMzO8pleoU82XV+yEVxgKA44uSvat
zcsIlQrUpwEAkFllbHnquiMxOr+UTXFXIj2L48V2oysbwbhadfIeGwO4MwRi3NJi
FgkrBgEEAdpHDwEBB0AzJRFF4Ufc8MFl42TDhkMNvMYxVdHogLyOXwK5zoQIo4jv
BBgWCAAgFiEExaIXpbFgakHS0zLkuAYYSKm6QCUFAmLc0mICGwIAgQkQuAYYSKm6
QCV2IAQZFggAHRYhBKkMHhbuzwzdKeMiMw7wawPgy8RgBQJi3NJiAAoJEA7wawPg
y8RgobsA/1uo4Z1ybThL//tNXFeV8J4vr5Kj+1mUO92v4E8oMhL0AQDsWH2H6kn2
yXEq4d7IEMN2bZOgnIK8q9wcQm7ouYXkDZ2rAP96iz2c2HYWouFZXf0vqzrgJWVB
ubvgn0HwJkdP6VIUKwEArIF90cFLtKmephxMlwp7h+djEWpVIQ/T31pLYSAH4QA=
=5+sB
-----END PGP PUBLIC KEY BLOCK-----
`

func TestAccGitlabUserGPGKey_basic(t *testing.T) {
	testUser := testAccCreateUsers(t, 1)[0]

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabUserGPGKeyDestroy,
		Steps: []resource.TestStep{
			// Create a user + gpgkey
			{
				Config: fmt.Sprintf(`
					resource "gitlab_user_gpgkey" "foo_key" {
					  key = <<EOF
					  %s
					  EOF
					  user_id = %d
					}
				`, testEd25519GPGPubKey, testUser.ID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gitlab_user_gpgkey.foo_key", "created_at"),
				),
			},
			{
				ResourceName:      "gitlab_user_gpgkey.foo_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update the key
			{
				Config: fmt.Sprintf(`
					resource "gitlab_user_gpgkey" "foo_key" {
					  key = <<EOF
					  %s
					  EOF
					  user_id = %d
					}
				`, updatedEd25519GPGPubKey, testUser.ID),
			},
			{
				ResourceName:      "gitlab_user_gpgkey.foo_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGitlabUserGPGKey_currentuser(t *testing.T) {
	testUser := testAccCreateUsers(t, 1)[0]

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabUserGPGKeyDestroy,
		Steps: []resource.TestStep{
			// Create a gpgkey
			{
				Config: fmt.Sprintf(`
					resource "gitlab_user_gpgkey" "foo_key" {
					  key = <<EOF
					  %s
					  EOF
					}
				`, testEd25519GPGPubKeyForCurrentUser),
			},
			{
				ResourceName:      "gitlab_user_gpgkey.foo_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Switch to a specific user
			{
				Config: fmt.Sprintf(`
					resource "gitlab_user_gpgkey" "foo_key" {
					  key = <<EOF
					  %s
					  EOF
					  user_id = %d
					}
				`, testEd25519GPGPubKeyForCurrentUser, testUser.ID),
			},
			{
				ResourceName:      "gitlab_user_gpgkey.foo_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Switch back to the current user
			{
				Config: fmt.Sprintf(`
					resource "gitlab_user_gpgkey" "foo_key" {
					  key = <<EOF
					  %s
					  EOF
					}
				`, testEd25519GPGPubKeyForCurrentUser),
			},
			{
				ResourceName:      "gitlab_user_gpgkey.foo_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGitlabUserGPGKey_ignoreTrailingWhitespaces(t *testing.T) {
	testUser := testAccCreateUsers(t, 1)[0]

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabUserGPGKeyDestroy,
		Steps: []resource.TestStep{
			// Create a user + gpgkey
			{
				Config: fmt.Sprintf(`
					resource "gitlab_user_gpgkey" "foo_key" {
					  key = <<EOF
					  %s
					  EOF
					  user_id = %d
					}
				`, testEd25519GPGPubKeyWithTrailingNewline, testUser.ID),
			},
			// Check for no-op plan
			{
				Config: fmt.Sprintf(`
					resource "gitlab_user_gpgkey" "foo_key" {
					  key = <<EOF
					  %s
					  EOF
					  user_id = %d
					}
				`, testEd25519GPGPubKeyWithTrailingNewline, testUser.ID),
				PlanOnly: true,
			},
			{
				ResourceName:      "gitlab_user_gpgkey.foo_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGitlabUserGPGKeyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_user_gpgkey" {
			continue
		}

		var err error
		userID, keyID, err := resourceGitlabUserGPGKeyParseID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to parse user GPG key resource ID: %s", err)
		}

		var key *gitlab.GPGKey
		if userID != 0 {
			key, _, err = testGitlabClient.Users.GetGPGKeyForUser(userID, keyID)
		} else {
			key, _, err = testGitlabClient.Users.GetGPGKey(keyID)
		}
		if err != nil && !is404(err) {
			return err
		}

		if key != nil {
			return fmt.Errorf("GPG Key still exists")
		}

		return nil
	}
	return nil
}
