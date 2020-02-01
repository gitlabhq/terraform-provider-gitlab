package tls

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"net"
	"net/url"
)

const pemCertReqType = "CERTIFICATE REQUEST"

func resourceCertRequest() *schema.Resource {
	return &schema.Resource{
		Create: CreateCertRequest,
		Delete: DeleteCertRequest,
		Read:   ReadCertRequest,

		Schema: map[string]*schema.Schema{

			"dns_names": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of DNS names to use as subjects of the certificate",
				ForceNew:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"ip_addresses": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of IP addresses to use as subjects of the certificate",
				ForceNew:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"uris": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of URIs to use as subjects of the certificate",
				ForceNew:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"key_algorithm": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the algorithm to use to generate the certificate's private key",
				ForceNew:    true,
			},

			"private_key_pem": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PEM-encoded private key that the certificate will belong to",
				ForceNew:    true,
				Sensitive:   true,
				StateFunc: func(v interface{}) string {
					return hashForState(v.(string))
				},
			},

			"subject": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     nameSchema,
				ForceNew: true,
			},

			"cert_request_pem": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func CreateCertRequest(d *schema.ResourceData, meta interface{}) error {
	key, err := parsePrivateKey(d, "private_key_pem", "key_algorithm")
	if err != nil {
		return err
	}

	subjectConfs := d.Get("subject").([]interface{})
	if len(subjectConfs) != 1 {
		return fmt.Errorf("must have exactly one 'subject' block")
	}
	subjectConf, ok := subjectConfs[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("subject block cannot be empty")
	}
	subject, err := nameFromResourceData(subjectConf)
	if err != nil {
		return fmt.Errorf("invalid subject block: %s", err)
	}

	certReq := x509.CertificateRequest{
		Subject: *subject,
	}

	dnsNamesI := d.Get("dns_names").([]interface{})
	for _, nameI := range dnsNamesI {
		certReq.DNSNames = append(certReq.DNSNames, nameI.(string))
	}
	ipAddressesI := d.Get("ip_addresses").([]interface{})
	for _, ipStrI := range ipAddressesI {
		ip := net.ParseIP(ipStrI.(string))
		if ip == nil {
			return fmt.Errorf("invalid IP address %#v", ipStrI.(string))
		}
		certReq.IPAddresses = append(certReq.IPAddresses, ip)
	}
	urisI := d.Get("uris").([]interface{})
	for _, uriI := range urisI {
		uri, err := url.Parse(uriI.(string))
		if err != nil {
			return fmt.Errorf("invalid URI %#v", uriI.(string))
		}
		certReq.URIs = append(certReq.URIs, uri)
	}

	certReqBytes, err := x509.CreateCertificateRequest(rand.Reader, &certReq, key)
	if err != nil {
		return fmt.Errorf("Error creating certificate request: %s", err)
	}
	certReqPem := string(pem.EncodeToMemory(&pem.Block{Type: pemCertReqType, Bytes: certReqBytes}))

	d.SetId(hashForState(string(certReqBytes)))
	d.Set("cert_request_pem", certReqPem)

	return nil
}

func DeleteCertRequest(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func ReadCertRequest(d *schema.ResourceData, meta interface{}) error {
	return nil
}
