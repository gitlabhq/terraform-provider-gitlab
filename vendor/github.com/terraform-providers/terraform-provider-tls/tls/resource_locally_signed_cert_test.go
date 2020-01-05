package tls

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	r "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestLocallySignedCert(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: locallySignedCertConfig(1, 0),
				Check: func(s *terraform.State) error {
					gotUntyped := s.RootModule().Outputs["cert_pem"].Value
					got, ok := gotUntyped.(string)
					if !ok {
						return fmt.Errorf("output for \"cert_pem\" is not a string")
					}
					if !strings.HasPrefix(got, "-----BEGIN CERTIFICATE----") {
						return fmt.Errorf("key is missing cert PEM preamble")
					}
					block, _ := pem.Decode([]byte(got))
					cert, err := x509.ParseCertificate(block.Bytes)
					if err != nil {
						return fmt.Errorf("error parsing cert: %s", err)
					}
					if expected, got := "2", cert.Subject.SerialNumber; got != expected {
						return fmt.Errorf("incorrect subject serial number: expected %v, got %v", expected, got)
					}
					if expected, got := "example.com", cert.Subject.CommonName; got != expected {
						return fmt.Errorf("incorrect subject common name: expected %v, got %v", expected, got)
					}
					if expected, got := "Example, Inc", cert.Subject.Organization[0]; got != expected {
						return fmt.Errorf("incorrect subject organization: expected %v, got %v", expected, got)
					}
					if expected, got := "Department of Terraform Testing", cert.Subject.OrganizationalUnit[0]; got != expected {
						return fmt.Errorf("incorrect subject organizational unit: expected %v, got %v", expected, got)
					}
					if expected, got := "5879 Cotton Link", cert.Subject.StreetAddress[0]; got != expected {
						return fmt.Errorf("incorrect subject street address: expected %v, got %v", expected, got)
					}
					if expected, got := "Pirate Harbor", cert.Subject.Locality[0]; got != expected {
						return fmt.Errorf("incorrect subject locality: expected %v, got %v", expected, got)
					}
					if expected, got := "CA", cert.Subject.Province[0]; got != expected {
						return fmt.Errorf("incorrect subject province: expected %v, got %v", expected, got)
					}
					if expected, got := "US", cert.Subject.Country[0]; got != expected {
						return fmt.Errorf("incorrect subject country: expected %v, got %v", expected, got)
					}
					if expected, got := "95559-1227", cert.Subject.PostalCode[0]; got != expected {
						return fmt.Errorf("incorrect subject postal code: expected %v, got %v", expected, got)
					}

					if expected, got := 2, len(cert.DNSNames); got != expected {
						return fmt.Errorf("incorrect number of DNS names: expected %v, got %v", expected, got)
					}
					if expected, got := "example.com", cert.DNSNames[0]; got != expected {
						return fmt.Errorf("incorrect DNS name 0: expected %v, got %v", expected, got)
					}
					if expected, got := "example.net", cert.DNSNames[1]; got != expected {
						return fmt.Errorf("incorrect DNS name 0: expected %v, got %v", expected, got)
					}

					if expected, got := 2, len(cert.IPAddresses); got != expected {
						return fmt.Errorf("incorrect number of IP addresses: expected %v, got %v", expected, got)
					}
					if expected, got := "127.0.0.1", cert.IPAddresses[0].String(); got != expected {
						return fmt.Errorf("incorrect IP address 0: expected %v, got %v", expected, got)
					}
					if expected, got := "127.0.0.2", cert.IPAddresses[1].String(); got != expected {
						return fmt.Errorf("incorrect IP address 0: expected %v, got %v", expected, got)
					}
					if expected, got := 2, len(cert.URIs); got != expected {
						return fmt.Errorf("incorrect number of URIs: expected %v, got %v", expected, got)
					}
					if expected, got := "spiffe://example-trust-domain/workload", cert.URIs[0].String(); got != expected {
						return fmt.Errorf("incorrect URI 0: expected %v, got %v", expected, got)
					}
					if expected, got := "spiffe://example-trust-domain/workload2", cert.URIs[1].String(); got != expected {
						return fmt.Errorf("incorrect URI 1: expected %v, got %v", expected, got)
					}

					if expected, got := 2, len(cert.ExtKeyUsage); got != expected {
						return fmt.Errorf("incorrect number of ExtKeyUsage: expected %v, got %v", expected, got)
					}
					if expected, got := x509.ExtKeyUsageServerAuth, cert.ExtKeyUsage[0]; got != expected {
						return fmt.Errorf("incorrect ExtKeyUsage[0]: expected %v, got %v", expected, got)
					}
					if expected, got := x509.ExtKeyUsageClientAuth, cert.ExtKeyUsage[1]; got != expected {
						return fmt.Errorf("incorrect ExtKeyUsage[1]: expected %v, got %v", expected, got)
					}

					if expected, got := x509.KeyUsageKeyEncipherment|x509.KeyUsageDigitalSignature, cert.KeyUsage; got != expected {
						return fmt.Errorf("incorrect KeyUsage: expected %v, got %v", expected, got)
					}

					// This time checking is a bit sloppy to avoid inconsistent test results
					// depending on the power of the machine running the tests.
					now := time.Now()
					if cert.NotBefore.After(now) {
						return fmt.Errorf("certificate validity begins in the future")
					}
					if now.Sub(cert.NotBefore) > (2 * time.Minute) {
						return fmt.Errorf("certificate validity begins more than two minutes in the past")
					}
					if cert.NotAfter.Sub(cert.NotBefore) != time.Hour {
						return fmt.Errorf("certificate validity is not one hour")
					}

					caBlock, _ := pem.Decode([]byte(testCACert))
					caCert, err := x509.ParseCertificate(caBlock.Bytes)
					if err != nil {
						return fmt.Errorf("error parsing ca cert: %s", err)
					}
					certPool := x509.NewCertPool()

					// Verify certificate
					_, err = cert.Verify(x509.VerifyOptions{Roots: certPool})
					if err == nil {
						return errors.New("incorrectly verified certificate")
					} else if _, ok := err.(x509.UnknownAuthorityError); !ok {
						return fmt.Errorf("incorrect verify error: expected UnknownAuthorityError, got %v", err)
					}
					certPool.AddCert(caCert)
					if _, err = cert.Verify(x509.VerifyOptions{Roots: certPool}); err != nil {
						return fmt.Errorf("verify failed: %s", err)
					}

					return nil
				},
			},
		},
	})
}

func TestAccLocallySignedCertRecreatesAfterExpired(t *testing.T) {
	oldNow := now
	var previousCert string
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		PreCheck:  setTimeForTest("2019-06-14T12:00:00Z"),
		Steps: []r.TestStep{
			{
				Config: locallySignedCertConfig(10, 2),
				Check: func(s *terraform.State) error {
					gotUntyped := s.RootModule().Outputs["cert_pem"].Value
					got, ok := gotUntyped.(string)
					if !ok {
						return fmt.Errorf("output for \"cert_pem\" is not a string")
					}
					previousCert = got
					return nil
				},
			},
			{
				Config: locallySignedCertConfig(10, 2),
				Check: func(s *terraform.State) error {
					gotUntyped := s.RootModule().Outputs["cert_pem"].Value
					got, ok := gotUntyped.(string)
					if !ok {
						return fmt.Errorf("output for \"cert_pem\" is not a string")
					}

					if got != previousCert {
						return fmt.Errorf("certificate updated even though no time has passed")
					}

					previousCert = got
					return nil
				},
			},
			{
				PreConfig: setTimeForTest("2019-06-14T19:00:00Z"),
				Config:    locallySignedCertConfig(10, 2),
				Check: func(s *terraform.State) error {
					gotUntyped := s.RootModule().Outputs["cert_pem"].Value
					got, ok := gotUntyped.(string)
					if !ok {
						return fmt.Errorf("output for \"cert_pem\" is not a string")
					}

					if got != previousCert {
						return fmt.Errorf("certificate updated even though not enough time has passed")
					}

					previousCert = got
					return nil
				},
			},
			{
				PreConfig: setTimeForTest("2019-06-14T21:00:00Z"),
				Config:    locallySignedCertConfig(10, 2),
				Check: func(s *terraform.State) error {
					gotUntyped := s.RootModule().Outputs["cert_pem"].Value
					got, ok := gotUntyped.(string)
					if !ok {
						return fmt.Errorf("output for \"cert_pem\" is not a string")
					}

					if got == previousCert {
						return fmt.Errorf("certificate not updated even though passed early renewal")
					}

					previousCert = got
					return nil
				},
			},
		},
	})
	now = oldNow
}

func TestAccLocallySignedCertNotRecreatedForEarlyRenewalUpdateInFuture(t *testing.T) {
	oldNow := now
	var previousCert string
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		PreCheck:  setTimeForTest("2019-06-14T12:00:00Z"),
		Steps: []r.TestStep{
			{
				Config: locallySignedCertConfig(10, 2),
				Check: func(s *terraform.State) error {
					gotUntyped := s.RootModule().Outputs["cert_pem"].Value
					got, ok := gotUntyped.(string)
					if !ok {
						return fmt.Errorf("output for \"cert_pem\" is not a string")
					}
					previousCert = got
					return nil
				},
			},
			{
				Config: locallySignedCertConfig(10, 3),
				Check: func(s *terraform.State) error {
					gotUntyped := s.RootModule().Outputs["cert_pem"].Value
					got, ok := gotUntyped.(string)
					if !ok {
						return fmt.Errorf("output for \"cert_pem\" is not a string")
					}

					if got != previousCert {
						return fmt.Errorf("certificate updated even though still time until early renewal")
					}

					previousCert = got
					return nil
				},
			},
			{
				PreConfig: setTimeForTest("2019-06-14T16:00:00Z"),
				Config:    locallySignedCertConfig(10, 3),
				Check: func(s *terraform.State) error {
					gotUntyped := s.RootModule().Outputs["cert_pem"].Value
					got, ok := gotUntyped.(string)
					if !ok {
						return fmt.Errorf("output for \"cert_pem\" is not a string")
					}

					if got != previousCert {
						return fmt.Errorf("certificate updated even though still time until early renewal")
					}

					previousCert = got
					return nil
				},
			},
			{
				PreConfig: setTimeForTest("2019-06-14T16:00:00Z"),
				Config:    locallySignedCertConfig(10, 9),
				Check: func(s *terraform.State) error {
					gotUntyped := s.RootModule().Outputs["cert_pem"].Value
					got, ok := gotUntyped.(string)
					if !ok {
						return fmt.Errorf("output for \"cert_pem\" is not a string")
					}

					if got == previousCert {
						return fmt.Errorf("certificate not updated even though early renewal time has passed")
					}

					previousCert = got
					return nil
				},
			},
		},
	})
	now = oldNow
}

func locallySignedCertConfig(validity uint32, earlyRenewal uint32) string {
	return fmt.Sprintf(`
                    resource "tls_locally_signed_cert" "test" {
                        cert_request_pem = <<EOT
%s
EOT

                        validity_period_hours = %d
                        early_renewal_hours = %d

                        allowed_uses = [
                            "key_encipherment",
                            "digital_signature",
                            "server_auth",
                            "client_auth",
                        ]

                        ca_key_algorithm = "RSA"
                        ca_cert_pem = <<EOT
%s
EOT
                        ca_private_key_pem = <<EOT
%s
EOT
                    }
                    output "cert_pem" {
                        value = "${tls_locally_signed_cert.test.cert_pem}"
                    }
                `, testCertRequest, validity, earlyRenewal, testCACert, testCAPrivateKey)
}
