package gitlab

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"

	"github.com/xanzy/go-gitlab"
)

// Config is per-provider, specifies where to connect to gitlab
type Config struct {
	Token      string
	BaseURL    string
	Insecure   bool
	CACertFile string
}

// Client returns a *gitlab.Client to interact with the configured gitlab instance
func (c *Config) Client() (interface{}, error) {
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

	transport := &http.Transport{TLSClientConfig: tlsConfig}

	httpClient := &http.Client{Transport: transport}

	client := gitlab.NewClient(httpClient, c.Token)
	if c.BaseURL != "" {
		err := client.SetBaseURL(c.BaseURL)
		if err != nil {
			// The BaseURL supplied wasn't valid, bail.
			return nil, err
		}
	}

	// Test the credentials by checking we can get information about the authenticated user.
	_, _, err := client.Users.CurrentUser()
	if err != nil {
		return nil, err
	}

	return client, nil
}
