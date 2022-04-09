package provider

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/xanzy/go-gitlab"
)

// Config is per-provider, specifies where to connect to gitlab
type Config struct {
	Token         string
	ForcePAToken  bool
	BaseURL       string
	Insecure      bool
	CACertFile    string
	ClientCert    string
	ClientKey     string
	EarlyAuthFail bool
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

	var client *gitlab.Client
	var err error

	// If a personal access token has been generated programmatically, as described at
	// https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html#create-a-personal-access-token-programmatically,
	// it will not be OAuth-compatible. Thus, personal accesstoken authentication may be enforced.
	if c.ForcePAToken {
		client, err = gitlab.NewClient(c.Token, opts...)
	} else {
		// The OAuth method is also compatible with project/group/personal access and job tokens (if generated via the GitLab web UI) because they are all usable as Bearer tokens.
		// Although the job token API access is very limited.
		// see https://docs.gitlab.com/ee/api#authentication
		client, err = gitlab.NewOAuthClient(c.Token, opts...)
	}
	if err != nil {
		return nil, err
	}

	// Test the credentials by checking we can get information about the authenticated user.
	if c.EarlyAuthFail {
		_, _, err = client.Users.CurrentUser()
	}

	return client, err
}
