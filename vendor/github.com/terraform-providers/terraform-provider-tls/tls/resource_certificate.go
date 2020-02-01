package tls

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const pemCertType = "CERTIFICATE"

var keyUsages map[string]x509.KeyUsage = map[string]x509.KeyUsage{
	"digital_signature":  x509.KeyUsageDigitalSignature,
	"content_commitment": x509.KeyUsageContentCommitment,
	"key_encipherment":   x509.KeyUsageKeyEncipherment,
	"data_encipherment":  x509.KeyUsageDataEncipherment,
	"key_agreement":      x509.KeyUsageKeyAgreement,
	"cert_signing":       x509.KeyUsageCertSign,
	"crl_signing":        x509.KeyUsageCRLSign,
	"encipher_only":      x509.KeyUsageEncipherOnly,
	"decipher_only":      x509.KeyUsageDecipherOnly,
}

var extKeyUsages map[string]x509.ExtKeyUsage = map[string]x509.ExtKeyUsage{
	"any_extended":                  x509.ExtKeyUsageAny,
	"server_auth":                   x509.ExtKeyUsageServerAuth,
	"client_auth":                   x509.ExtKeyUsageClientAuth,
	"code_signing":                  x509.ExtKeyUsageCodeSigning,
	"email_protection":              x509.ExtKeyUsageEmailProtection,
	"ipsec_end_system":              x509.ExtKeyUsageIPSECEndSystem,
	"ipsec_tunnel":                  x509.ExtKeyUsageIPSECTunnel,
	"ipsec_user":                    x509.ExtKeyUsageIPSECUser,
	"timestamping":                  x509.ExtKeyUsageTimeStamping,
	"ocsp_signing":                  x509.ExtKeyUsageOCSPSigning,
	"microsoft_server_gated_crypto": x509.ExtKeyUsageMicrosoftServerGatedCrypto,
	"netscape_server_gated_crypto":  x509.ExtKeyUsageNetscapeServerGatedCrypto,
}

// rsaPublicKey reflects the ASN.1 structure of a PKCS#1 public key.
type rsaPublicKey struct {
	N *big.Int
	E int
}

var now = func() time.Time {
	return time.Now()
}

// generateSubjectKeyID generates a SHA-1 hash of the subject public key.
func generateSubjectKeyID(pub crypto.PublicKey) ([]byte, error) {
	var publicKeyBytes []byte
	var err error

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		publicKeyBytes, err = asn1.Marshal(rsaPublicKey{N: pub.N, E: pub.E})
		if err != nil {
			return nil, err
		}
	case *ecdsa.PublicKey:
		publicKeyBytes = elliptic.Marshal(pub.Curve, pub.X, pub.Y)
	default:
		return nil, errors.New("only RSA and ECDSA public keys supported")
	}

	hash := sha1.Sum(publicKeyBytes)
	return hash[:], nil
}

func resourceCertificateCommonSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"validity_period_hours": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Number of hours that the certificate will remain valid for",
			ForceNew:    true,
		},

		"early_renewal_hours": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "Number of hours before the certificates expiry when a new certificate will be generated",
		},

		"is_ca_certificate": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether the generated certificate will be usable as a CA certificate",
			ForceNew:    true,
		},

		"allowed_uses": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Uses that are allowed for the certificate",
			ForceNew:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"cert_pem": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"ready_for_renewal": {
			Type:     schema.TypeBool,
			Computed: true,
		},

		"validity_start_time": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"validity_end_time": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"set_subject_key_id": &schema.Schema{
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "If true, the generated certificate will include a subject key identifier.",
			ForceNew:    true,
		},
	}
}

func createCertificate(d *schema.ResourceData, template, parent *x509.Certificate, pub crypto.PublicKey, priv interface{}) error {
	var err error

	template.NotBefore = now()
	template.NotAfter = template.NotBefore.Add(time.Duration(d.Get("validity_period_hours").(int)) * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	template.SerialNumber, err = rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %s", err)
	}

	keyUsesI := d.Get("allowed_uses").([]interface{})
	for _, keyUseI := range keyUsesI {
		keyUse := keyUseI.(string)
		if usage, ok := keyUsages[keyUse]; ok {
			template.KeyUsage |= usage
		}
		if usage, ok := extKeyUsages[keyUse]; ok {
			template.ExtKeyUsage = append(template.ExtKeyUsage, usage)
		}
	}

	if d.Get("is_ca_certificate").(bool) {
		template.IsCA = true

		template.SubjectKeyId, err = generateSubjectKeyID(pub)
		if err != nil {
			return fmt.Errorf("failed to set subject key identifier: %s", err)
		}
	}

	if d.Get("set_subject_key_id").(bool) {
		template.SubjectKeyId, err = generateSubjectKeyID(pub)
		if err != nil {
			return fmt.Errorf("failed to set subject key identifier: %s", err)
		}
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, parent, pub, priv)
	if err != nil {
		return fmt.Errorf("error creating certificate: %s", err)
	}
	certPem := string(pem.EncodeToMemory(&pem.Block{Type: pemCertType, Bytes: certBytes}))

	validFromBytes, err := template.NotBefore.MarshalText()
	if err != nil {
		return fmt.Errorf("error serializing validity_start_time: %s", err)
	}
	validToBytes, err := template.NotAfter.MarshalText()
	if err != nil {
		return fmt.Errorf("error serializing validity_end_time: %s", err)
	}

	d.SetId(template.SerialNumber.String())
	d.Set("cert_pem", certPem)
	d.Set("ready_for_renewal", false)
	d.Set("validity_start_time", string(validFromBytes))
	d.Set("validity_end_time", string(validToBytes))

	return nil
}

func DeleteCertificate(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func ReadCertificate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func CustomizeCertificateDiff(d *schema.ResourceDiff, meta interface{}) error {
	var readyForRenewal bool

	endTimeStr := d.Get("validity_end_time").(string)
	endTime := now()
	err := endTime.UnmarshalText([]byte(endTimeStr))
	if err != nil {
		// If end time is invalid then we'll treat it as being at the time for renewal.
		readyForRenewal = true
	} else {
		earlyRenewalPeriod := time.Duration(-d.Get("early_renewal_hours").(int)) * time.Hour
		endTime = endTime.Add(earlyRenewalPeriod)

		currentTime := now()
		timeToRenewal := endTime.Sub(currentTime)
		if timeToRenewal <= 0 {
			readyForRenewal = true
		}
	}

	if readyForRenewal {
		err = d.SetNew("ready_for_renewal", true)
		if err != nil {
			return err
		}
		err = d.ForceNew("ready_for_renewal")
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateCertificate(d *schema.ResourceData, meta interface{}) error {
	return nil
}
