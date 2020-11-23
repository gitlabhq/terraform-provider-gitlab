package gitlab

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/logging"
	"github.com/xanzy/go-gitlab"
)

// Config is per-provider, specifies where to connect to gitlab
type Config struct {
	Token      string
	BaseURL    string
	Insecure   bool
	CACertFile string
	ClientCert string
	ClientKey  string
}

// Client returns a *gitlab.Client to interact with the configured gitlab instance
func (c *Config) Client() (*gitlab.Client, error) {
	// Configure TLS/SSL
	tlsConfig := &tls.Config{}

	// If a CACertFile has been specified, use that for cert validation
	if c.CACertFile != "" {
		caCert, err := ioutil.ReadFile(c.CACertFile)
		if err != nil {
			return nil, err
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
	}

	// If configured as insecure, turn off SSL verification
	if c.Insecure {
		tlsConfig.InsecureSkipVerify = true
	}

	// add client cert and key to connection
	if c.ClientCert != "" && c.ClientKey != "" {
		clientPair, err := tls.LoadX509KeyPair(c.ClientCert, c.ClientKey)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{clientPair}
	}

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.TLSClientConfig = tlsConfig
	t.MaxIdleConnsPerHost = 100

	opts := []gitlab.ClientOptionFunc{
		gitlab.WithHTTPClient(
			&http.Client{
				Transport: logging.NewTransport("GitLab", t),
			},
		),
	}

	if c.BaseURL != "" {
		opts = append(opts, gitlab.WithBaseURL(c.BaseURL))
	}

	client, err := gitlab.NewClient(c.Token, opts...)
	if err != nil {
		return nil, err
	}

	// Test the credentials by checking we can get information about the authenticated user.
	_, _, err = client.Users.CurrentUser()

	return client, err
}
